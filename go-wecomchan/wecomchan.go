package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

var SENDKEY string = "set_a_sendkey"
var WECOM_CID string = "企业微信公司ID"
var WECOM_SECRET string = "企业微信应用Secret"
var WECOM_AID string = "企业微信应用ID"
var WECOM_TOUID string = "@all"

var REDIS_ADDR string = "localhost:6379"
var ctx = context.Background()
var redis_stat string = "OFF"

func praser_json(json_str string) map[string]interface{} {
	var wecom_response map[string]interface{}
	if string(json_str) != "" {
		err := json.Unmarshal([]byte(string(json_str)), &wecom_response)
		if err != nil {
			log.Println("生成json字符串错误")
		}
	}
	return wecom_response
}

func get_token(corpid, app_secret string) string {
	resp, err := http.Get("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + corpid + "&corpsecret=" + app_secret)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	resp_data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	token_response := praser_json(string(resp_data))
	return token_response["access_token"].(string)
}

func redis_client() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     REDIS_ADDR,
		Password: "", // no password set
		DB:       8,  // use default DB
	})
	return rdb
}

func post_msg(text_msg, msg_type, post_url string) string {
	type msg struct {
		Content string `json:"content"`
	}
	type JsonData struct {
		Touser                   string `json:"touser"`
		Agentid                  string `json:"agentid"`
		Msgtype                  string `json:"msgtype"`
		Text                     msg    `json:"text"`
		Duplicate_check_interval int    `json:"duplicate_check_interval"`
	}
	post_data := JsonData{
		Touser:                   WECOM_TOUID,
		Agentid:                  WECOM_AID,
		Msgtype:                  msg_type,
		Duplicate_check_interval: 600,
		Text:                     msg{Content: text_msg},
	}

	post_json, _ := json.Marshal(post_data)
	log.Println(string(post_json))
	msg_req, err := http.NewRequest("POST", post_url, bytes.NewBuffer(post_json))
	if err != nil {
		log.Println(err)
	}
	msg_req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(msg_req)
	if err != nil {
		panic(err)
	}
	defer msg_req.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func main() {
	var access_token string
	if redis_stat == "ON" {
		rdb := redis_client()
		vals, err := rdb.Get(ctx, "access_token").Result()
		if err != redis.Nil {
			log.Println(err)
		}
		log.Println("从redis获取token")
		access_token = string(vals)
	}
	if access_token == "" {
		access_token = get_token(WECOM_CID, WECOM_SECRET)
	}
	wecom_chan := func(res http.ResponseWriter, req *http.Request) {
		post_url := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + access_token
		req.ParseForm()
		sendkey := req.FormValue("sendkey")
		if sendkey != SENDKEY {
			log.Panicln("sendkey 错误，请检查")
		}
		msg := req.FormValue("msg")
		msg_type := req.FormValue("msg_type")
		post_status := post_msg(msg, msg_type, post_url)
		log.Println(post_status)
		post_response := praser_json(string(post_status))
		log.Println(post_response)
		if post_response["errmsg"] == "ok" && redis_stat == "ON" {
			log.Println("pre to set redis key")
			rdb := redis_client()
			set, err := rdb.SetEX(ctx, "access_token", access_token, 7000*time.Second).Result()
			log.Println(set)
			if err != redis.Nil {
				log.Println(err)
			}
		}
	}
	http.HandleFunc("/wecomchan", wecom_chan)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
