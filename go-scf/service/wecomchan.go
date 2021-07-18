package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/riba2534/wecomchan/go-scf/consts"
	"github.com/riba2534/wecomchan/go-scf/dal"
	"github.com/riba2534/wecomchan/go-scf/model"
	"github.com/riba2534/wecomchan/go-scf/utils"
	"github.com/tencentyun/scf-go-lib/events"
)

func WeComChanService(ctx context.Context, event events.APIGatewayRequest) map[string]interface{} {
	sendKey := getQuery("sendkey", event)
	msgType := getQuery("msg_type", event)
	msg := getQuery("msg", event)
	if msgType == "" || msg == "" {
		return utils.MakeResp(-1, "param error")
	}
	if sendKey != consts.SEND_KEY {
		return utils.MakeResp(-1, "sendkey error")
	}
	toUser := getQuery("to_user", event)
	if toUser == "" {
		toUser = consts.WECOM_TOUID
	}
	if err := postWechatMsg(dal.AccessToken, msg, msgType, toUser); err != nil {
		return utils.MakeResp(0, err.Error())
	}
	return utils.MakeResp(0, "success")
}

func postWechatMsg(accessToken, msg, msgType, toUser string) error {
	content := &model.WechatMsg{
		ToUser:                 toUser,
		AgentId:                consts.WECOM_AID,
		MsgType:                msgType,
		DuplicateCheckInterval: 600,
		Text: &model.MsgText{
			Content: msg,
		},
	}
	b, _ := jsoniter.Marshal(content)
	client := http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("POST", fmt.Sprintf(consts.WeComMsgSendURL, accessToken), bytes.NewBuffer(b))
	req.Header.Set("Content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[postWechatMsg] failed, err=", err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("postWechatMsg statusCode is not 200")
		return errors.New("statusCode is not 200")
	}
	respBodyBytes, _ := ioutil.ReadAll(resp.Body)
	postResp := &model.PostResp{}
	if err := jsoniter.Unmarshal(respBodyBytes, postResp); err != nil {
		fmt.Println("postWechatMsg json Unmarshal failed, err=", err)
		return err
	}
	if postResp.Errcode != 0 {
		fmt.Println("postWechatMsg postResp.Errcode != 0, err=", postResp.Errmsg)
		return errors.New(postResp.Errmsg)
	}
	return nil
}

func getQuery(key string, event events.APIGatewayRequest) string {
	switch event.Method {
	case "GET":
		value := event.QueryString[key]
		if len(value) > 0 && value[0] != "" {
			return value[0]
		}
		return ""
	case "POST":
		return jsoniter.Get([]byte(event.Body), key).ToString()
	default:
		return ""
	}
}
