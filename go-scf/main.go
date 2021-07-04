package main

import (
	"context"
	"strings"

	"github.com/riba2534/wecomchan/go-scf/service"
	"github.com/riba2534/wecomchan/go-scf/utils"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"github.com/tencentyun/scf-go-lib/events"
)

func HTTPHandler(ctx context.Context, event events.APIGatewayRequest) (events.APIGatewayResponse, error) {
	path := event.Path
	var result interface{}
	if strings.HasPrefix(path, "/wecomchan") {
		result = service.WeComChanService
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
