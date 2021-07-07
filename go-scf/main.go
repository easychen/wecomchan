package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/riba2534/wecomchan/go-scf/consts"
	"github.com/riba2534/wecomchan/go-scf/dal"
	"github.com/riba2534/wecomchan/go-scf/service"
	"github.com/riba2534/wecomchan/go-scf/utils"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"github.com/tencentyun/scf-go-lib/events"
)

func init() {
	consts.FUNC_NAME = utils.GetEnvDefault("FUNC_NAME", "")
	consts.SEND_KEY = utils.GetEnvDefault("SEND_KEY", "")
	consts.WECOM_CID = utils.GetEnvDefault("WECOM_CID", "")
	consts.WECOM_SECRET = utils.GetEnvDefault("WECOM_SECRET", "")
	consts.WECOM_AID = utils.GetEnvDefault("WECOM_AID", "")
	consts.WECOM_TOUID = utils.GetEnvDefault("WECOM_TOUID", "@all")
	if consts.FUNC_NAME == "" || consts.SEND_KEY == "" || consts.WECOM_CID == "" ||
		consts.WECOM_SECRET == "" || consts.WECOM_AID == "" || consts.WECOM_TOUID == "" {
		fmt.Printf("os.env load Fail, please check your os env.\nFUNC_NAME=%s\nSEND_KEY=%s\nWECOM_CID=%s\nWECOM_SECRET=%s\nWECOM_AID=%s\nWECOM_TOUID=%s\n", consts.FUNC_NAME, consts.SEND_KEY, consts.WECOM_CID, consts.WECOM_SECRET, consts.WECOM_AID, consts.WECOM_TOUID)
		panic("os.env param error")
	}
	fmt.Println("os.env load success!")
}

func HTTPHandler(ctx context.Context, event events.APIGatewayRequest) (events.APIGatewayResponse, error) {
	path := event.Path
	fmt.Println("req->", utils.MarshalToStringParam(event))
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
	dal.Init()
	cloudfunction.Start(HTTPHandler)
}
