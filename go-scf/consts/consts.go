package consts

const (
	SENDKEY      = "set_a_sendkey"
	WECOM_CID    = "企业微信公司ID"
	WECOM_SECRET = "企业微信应用Secret"
	WECOM_AID    = "企业微信应用ID"
	WECOM_TOUID  = "@all"
)

// 微信发消息API
const (
	WeComMsgSendURL     = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"
	WeComAccessTokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
)
