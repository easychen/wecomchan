package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
)

/*-------------------------------  环境变量配置 begin  -------------------------------*/

var Sendkey = GetEnvDefault("SENDKEY", "set_a_sendkey")
var WecomCid = GetEnvDefault("WECOM_CID", "企业微信公司ID")
var WecomSecret = GetEnvDefault("WECOM_SECRET", "企业微信应用Secret")
var WecomAid = GetEnvDefault("WECOM_AID", "企业微信应用ID")
var WecomToUid = GetEnvDefault("WECOM_TOUID", "@all")
var RedisStat = GetEnvDefault("REDIS_STAT", "OFF")
var RedisAddr = GetEnvDefault("REDIS_ADDR", "localhost:6379")
var RedisPassword = GetEnvDefault("REDIS_PASSWORD", "")
var ctx = context.Background()

/*-------------------------------  环境变量配置 end  -------------------------------*/

/*-------------------------------  企业微信服务端API begin  -------------------------------*/

var GetTokenApi = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
var SendMessageApi = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"
var UploadMediaApi = "https://qyapi.weixin.qq.com/cgi-bin/media/upload?access_token=%s&type=%s"

/*-------------------------------  企业微信服务端API end  -------------------------------*/

type Msg struct {
	Content string `json:"content"`
}
type Pic struct {
	MediaId string `json:"media_id"`
}
type JsonData struct {
	ToUser                 string `json:"touser"`
	AgentId                string `json:"agentid"`
	MsgType                string `json:"msgtype"`
	DuplicateCheckInterval int    `json:"duplicate_check_interval"`
	Text                   Msg    `json:"text"`
	Image                  Pic    `json:"image"`
}

// GetEnvDefault 获取配置信息，未获取到则取默认值
func GetEnvDefault(key, defVal string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defVal
	}
	return val
}

// ParseJson 将json字符串解析为map
func ParseJson(jsonStr string) map[string]interface{} {
	var wecomResponse map[string]interface{}
	if string(jsonStr) != "" {
		err := json.Unmarshal([]byte(string(jsonStr)), &wecomResponse)
		if err != nil {
			log.Println("生成json字符串错误")
		}
	}
	return wecomResponse
}

// GetRemoteToken 从企业微信服务端API获取access_token，存在redis服务则缓存
func GetRemoteToken(corpId, appSecret string) string {
	getTokenUrl := fmt.Sprintf(GetTokenApi, corpId, appSecret)
	log.Println("getTokenUrl==>", getTokenUrl)
	resp, err := http.Get(getTokenUrl)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	tokenResponse := ParseJson(string(respData))
	log.Println("企业微信获取access_token接口返回==>", tokenResponse)
	accessToken := tokenResponse["access_token"].(string)

	if RedisStat == "ON" {
		log.Println("prepare to set redis key")
		rdb := RedisClient()
		// access_token有效时间为7200秒(2小时)
		set, err := rdb.SetNX(ctx, "access_token", accessToken, 7000*time.Second).Result()
		log.Println(set)
		if err != nil {
			log.Println(err)
		}
	}
	return accessToken
}

// RedisClient redis客户端
func RedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: RedisPassword, // no password set
		DB:       0,             // use default DB
	})
	return rdb
}

// PostMsg 推送消息
func PostMsg(postData JsonData, postUrl string) string {
	postJson, _ := json.Marshal(postData)
	log.Println("postJson ", string(postJson))
	log.Println("postUrl ", postUrl)
	msgReq, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(postJson))
	if err != nil {
		log.Println(err)
	}
	msgReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(msgReq)
	if err != nil {
		log.Fatalln("企业微信发送应用消息接口报错==>", err)
	}
	defer msgReq.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	mediaResp := ParseJson(string(body))
	log.Println("企业微信发送应用消息接口返回==>", mediaResp)
	return string(body)
}

// CheckOrUploadMedia  核对消息类型，如果为图片则上传临时素材并返回mediaId
func CheckOrUploadMedia(msgType string, req *http.Request, accessToken string) string {
	if msgType != "image" {
		log.Println("消息类型不是图片")
		return ""
	}

	// 企业微信图片上传不能大于2M
	_ = req.ParseMultipartForm(2 << 20)
	imgFile, imgHeader, err := req.FormFile("media")
	log.Printf("文件大小==>%d字节", imgHeader.Size)
	if err != nil {
		log.Fatalln("图片文件出错==>", err)
		return ""
	}
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	if createFormFile, err := writer.CreateFormFile("media", imgHeader.Filename); err == nil {
		readAll, _ := ioutil.ReadAll(imgFile)
		createFormFile.Write(readAll)
	}
	writer.Close()

	uploadMediaUrl := fmt.Sprintf(UploadMediaApi, accessToken, msgType)
	log.Println("uploadMediaUrl==>", uploadMediaUrl)
	newRequest, _ := http.NewRequest("POST", uploadMediaUrl, buf)
	newRequest.Header.Set("Content-Type", writer.FormDataContentType())
	log.Println("Content-Type ", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(newRequest)
	respData, _ := ioutil.ReadAll(resp.Body)
	mediaResp := ParseJson(string(respData))
	log.Println("企业微信上传临时素材接口返回==>", mediaResp)
	if err != nil {
		log.Fatalln("上传临时素材出错==>", err)
		return ""
	} else {
		return mediaResp["media_id"].(string)
	}
}

// IsZero 判断企业微信服务端是否正常响应
func IsZero(v interface{}) (bool, error) {
	t := reflect.TypeOf(v)
	if !t.Comparable() {
		return false, fmt.Errorf("type is not comparable: %v", t)
	}
	return v == reflect.Zero(t).Interface(), nil
}

// 获取企业微信的access_token
func getAccessToken() string {
	accessToken := ""
	if RedisStat == "ON" {
		log.Println("尝试从redis获取token")
		rdb := RedisClient()
		value, err := rdb.Get(ctx, "access_token").Result()
		if err == redis.Nil {
			log.Println("access_token does not exist, need get it from remote API")
		}
		accessToken = value
	}
	if accessToken == "" {
		log.Println("get access_token from remote API")
		accessToken = GetRemoteToken(WecomCid, WecomSecret)
	} else {
		log.Println("get access_token from redis")
	}
	return accessToken
}

// InitJsonData 初始化Json公共部分数据
func InitJsonData(msgType string) JsonData {
	return JsonData{
		ToUser:                 WecomToUid,
		AgentId:                WecomAid,
		MsgType:                msgType,
		DuplicateCheckInterval: 600,
	}
}

// 主函数入口
func main() {
	// 设置日志内容显示文件名和行号
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	wecomChan := func(res http.ResponseWriter, req *http.Request) {
		_ = req.ParseForm()
		sendkey := req.FormValue("sendkey")
		if sendkey != Sendkey {
			log.Panicln("sendkey 错误，请检查")
		}
		msgContent := req.FormValue("msg")
		msgType := req.FormValue("msg_type")
		log.Println("mes_type=", msgType)
		// 刷新token
		accessToken := getAccessToken()
		mediaId := CheckOrUploadMedia(msgType, req, accessToken)
		log.Println("企业微信上传临时素材接口返回的media_id==>", mediaId)

		// 准备发送应用消息所需参数
		postData := InitJsonData(msgType)
		postData.Text = Msg{
			Content: msgContent,
		}
		postData.Image = Pic{
			MediaId: mediaId,
		}
		// 再次刷新token
		accessToken = getAccessToken()
		sendMessageUrl := fmt.Sprintf(SendMessageApi, accessToken)

		postStatus := PostMsg(postData, sendMessageUrl)
		postResponse := ParseJson(postStatus)
		errcode := postResponse["errcode"]
		_, err := IsZero(errcode)
		if err != nil {
			log.Printf("%v", err)
		}
		res.Header().Set("Content-type", "application/json")
		_, _ = res.Write([]byte(postStatus))
	}
	http.HandleFunc("/wecomchan", wecomChan)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
