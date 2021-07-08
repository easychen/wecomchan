package consts

var (
	FUNC_NAME    string
	SEND_KEY     string
	WECOM_CID    string
	WECOM_SECRET string
	WECOM_AID    string
	WECOM_TOUID  string
)

// 微信发消息API
const (
	// https://work.weixin.qq.com/api/doc/90000/90135/90236
	WeComMsgSendURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"
	// https://work.weixin.qq.com/api/doc/90000/90135/91039
	WeComAccessTokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
)
