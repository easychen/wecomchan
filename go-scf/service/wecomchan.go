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
	"github.com/riba2534/wecomchan/go-scf/model"
	"github.com/riba2534/wecomchan/go-scf/utils"
	"github.com/tencentyun/scf-go-lib/events"
)

func WeComChanService(ctx context.Context, event events.APIGatewayRequest) interface{} {
	accessToken, err := getAccessToken()
	if err != nil {
		return utils.MakeResp(-1, "get accessToken error")
	}
	sendKey := getQuery("sendkey", event.QueryString)
	msgType := getQuery("msg_type", event.QueryString)
	msg := getQuery("msg", event.QueryString)
	if accessToken == "" || msgType == "" || msg == "" {
		return utils.MakeResp(-1, "param error")
	}
	if sendKey != consts.SENDKEY {
		return utils.MakeResp(-1, "sendkey error")
	}
	if err := postWechatMsg(accessToken, msg, msgType); err != nil {
		return utils.MakeResp(0, err.Error())
	}
	return utils.MakeResp(0, "success")
}

func getAccessToken() (string, error) {
	client := http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", fmt.Sprintf(consts.WeComAccessTokenURL, consts.WECOM_CID, consts.WECOM_SECRET), nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("getAccessToken err=", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("getAccessToken statusCode is not 200")
	}
	respBodyBytes, _ := ioutil.ReadAll(resp.Body)
	assesTokenResp := &model.AssesTokenResp{}
	if err := jsoniter.Unmarshal(respBodyBytes, assesTokenResp); err != nil {
		fmt.Println("getAccessToken json Unmarshal failed, err=", err)
		return "", err
	}
	if assesTokenResp.Errcode != 0 {
		fmt.Println("getAccessToken assesTokenResp.Errcode != 0, err=", assesTokenResp.Errmsg)
		return "", errors.New(assesTokenResp.Errmsg)
	}
	return assesTokenResp.AccessToken, nil
}

func postWechatMsg(accessToken, msg, msgType string) error {
	content := &model.WechatMsg{
		ToUser:                 consts.WECOM_TOUID,
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

func getQuery(key string, query events.APIGatewayQueryString) string {
	value := query[key]
	if len(value) > 0 && value[0] != "" {
		return value[0]
	}
	return ""
}
