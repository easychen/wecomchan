package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/riba2534/wecomchan/go-scf/consts"
	"github.com/riba2534/wecomchan/go-scf/service"
	"github.com/riba2534/wecomchan/go-scf/utils"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"github.com/tencentyun/scf-go-lib/events"
)

func init() {
	config := viper.New()
	config.SetConfigFile("config.yaml")
	config.AddConfigPath(".")
	err := config.ReadInConfig()
	if err != nil {
		panic("load config.yaml failed")
	}
	consts.FUNC_NAME = cast.ToString(config.Get("config.FUNC_NAME"))
	consts.SEND_KEY = cast.ToString(config.Get("config.SEND_KEY"))
	consts.WECOM_CID = cast.ToString(config.Get("config.WECOM_CID"))
	consts.WECOM_SECRET = cast.ToString(config.Get("config.WECOM_SECRET"))
	consts.WECOM_AID = cast.ToString(config.Get("config.WECOM_AID"))
	consts.WECOM_TOUID = cast.ToString(config.Get("config.WECOM_TOUID"))
	fmt.Println("config.yaml load success!")
}

func HTTPHandler(ctx context.Context, event events.APIGatewayRequest) (events.APIGatewayResponse, error) {
	path := event.Path
	fmt.Println("req:", utils.MarshalToStringParam(event))
	var result interface{}
	if strings.HasPrefix(path, "/"+consts.FUNC_NAME) {
		result = service.WeComChanService(ctx, event)
	} else {
		// 匹配失败返回原始HTTP请求
		result = event
	}
	return events.APIGatewayResponse{
		IsBase64Encoded: false,
		StatusCode:      200,
		Headers:         map[string]string{},
		Body:            utils.MarshalToStringParam(result),
	}, nil
}

func main() {
	cloudfunction.Start(HTTPHandler)
}
