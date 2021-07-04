package model

type AssesTokenResp struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type MsgText struct {
	Content string `json:"content"`
}

// https://work.weixin.qq.com/api/doc/90002/90151/90854
type WechatMsg struct {
	ToUser                 string   `json:"touser"`
	AgentId                string   `json:"agentid"`
	MsgType                string   `json:"msgtype"`
	Text                   *MsgText `json:"text"`
	DuplicateCheckInterval int      `json:"duplicate_check_interval"`
}

type PostResp struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	Invaliduser string `json:"invaliduser"`
}
