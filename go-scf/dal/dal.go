package dal

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/riba2534/wecomchan/go-scf/consts"
	"github.com/riba2534/wecomchan/go-scf/model"
)

var AccessToken string

func loadAccessToken() {
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
		panic(err)
	}
	if assesTokenResp.Errcode != 0 {
		fmt.Println("getAccessToken assesTokenResp.Errcode != 0, err=", assesTokenResp.Errmsg)
		panic(err)
	}
	AccessToken = assesTokenResp.AccessToken
}

func Init() {
	loadAccessToken()
	fmt.Printf("[Init] accessToken load success, time=%s, token=%s\n", time.Now().Format("2006-01-02 15:04:05"), AccessToken)
	go func() {
		for {
			time.Sleep(30 * time.Minute)
			loadAccessToken()
			fmt.Printf("[Goroutine] accessToken load success, time=%s, token=%s\n", time.Now().Format("2006-01-02 15:04:05"), AccessToken)
		}
	}()
}
