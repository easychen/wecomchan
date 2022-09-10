# 发送应用消息

应用支持推送文本、图片、视频、文件、卡片、图文等类型。

**请求方式：**POST（**HTTPS**）
**请求地址：** https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=ACCESS_TOKEN

# 全局说明

| 参数                     | 是否必须 | 说明                                                         |
| ------------------------ | ---------------- | ------------------------------------------------------------ |
| touser                   | <font color='cornflowerblue'>否</font>               | 指定接收消息的成员，成员ID列表（多个接收者用‘\|’分隔，最多支持1000个）。<br>:warning:<font color='red'>特殊情况：指定为"@all"，则向该企业应用的全部成员发送</font> |
| toparty                  | <font color='cornflowerblue'>否</font>               | 指定接收消息的部门，部门ID列表，多个接收者用‘\|’分隔，最多支持100个。<br/><font color='gray'>当touser为"@all"时忽略本参数</font> |
| totag                    | <font color='cornflowerblue'>否</font>               | 指定接收消息的标签，标签ID列表，多个接收者用‘\|’分隔，最多支持100个。<br/><font color='gray'>当touser为"@all"时忽略本参数</font> |
| agentid | <font color='red'>是 </font> | 企业应用的id，整型。登录 **网页版企业微信** 在配置的应用中进行查看。 |

touser toparty totag 三个不能同时为空！



# 文本信息

~~~json
{
   "touser" : "UserID1|UserID2|UserID3",
   "toparty" : "PartyID1|PartyID2",
   "totag" : "TagID1 | TagID2",
   "agentid" : {你的agentid},
    
   "msgtype" : "text",
   "text" : {
       "content" : "你的快递已到，请携带工卡前往邮件中心领取。\n出发前可查看<a href=\"http://work.weixin.qq.com\">邮件中心视频实况</a>，聪明避开排队。"
   },
}
~~~

| 参数    | 是否必须                    | 说明                                                         |
| ------- | --------------------------- | ------------------------------------------------------------ |
| msgtype | <font color='red'>是</font> | 消息类型，此时固定为：**text**                               |
| text    | <font color='red'>是</font> | text参数的content字段可以支持换行、以及A标签，即可打开自定义的网页（可参考以上示例代码）<br/>:warning:注意：换行符等符号请转义`\n`，引号用`\"` |

# 媒体信息

获取的media_id 的有效期只有三天，谨慎使用！

## 图片信息

~~~json
{
   "touser" : "UserID1|UserID2|UserID3",
   "toparty" : "PartyID1|PartyID2",
   "totag" : "TagID1 | TagID2",
   "agentid" : {你的agentid},
    
   "msgtype" : "image",
   "image" : {
        "media_id" : "MEDIA_ID"
   },
~~~

| 参数    | 是否必须                    | 说明                                                         |
| ------- | --------------------------- | ------------------------------------------------------------ |
| msgtype  | <font color='red'>是</font> | 消息类型，此时固定为：**image**                              |
| media_id | <font color='red'>是</font> | 图片媒体文件id，可以调用上传临时素材接口获取                 |



## 视频信息

~~~json
{
   "touser" : "UserID1|UserID2|UserID3",
   "toparty" : "PartyID1|PartyID2",
   "totag" : "TagID1 | TagID2",
   "agentid" : {你的agentid},
   
   "msgtype" : "video",
   "video" : {
        "media_id" : "MEDIA_ID",
        "title" : "Title",
       "description" : "Description"
   },
}
~~~

| 参数        | 是否必须                               | 说明                                                         |
| ----------- | -------------------------------------- | ------------------------------------------------------------ |
| msgtype     | <font color='red'>是</font>            | 消息类型，此时固定为：**video**                              |
| media_id    | <font color='red'>是</font>            | 视频媒体文件id，可以调用[上传临时素材](https://developer.work.weixin.qq.com/document/path/90236#10112)接口获取 |
| title       | <font color='cornflowerblue'>否</font> | 视频消息的标题，不超过128个字节，超过会自动截断              |
| description | <font color='cornflowerblue'>否</font> | 视频消息的描述，不超过512个字节，超过会自动截断              |



## 文件信息

~~~json
{
   "touser" : "UserID1|UserID2|UserID3",
   "toparty" : "PartyID1|PartyID2",
   "totag" : "TagID1 | TagID2",
   "agentid" : {你的agentid},
    
   "msgtype" : "file",
   "file" : {
        "media_id" : "1Yv-zXfHjSjU-7LH-GwtYqDGS-zz6w22KmWAT5COgP7o"
   },
}
~~~
| 参数    | 是否必须                    | 说明                                                         |
| ------- | --------------------------- | ------------------------------------------------------------ |
| msgtype  | <font color='red'>是</font> | 消息类型，此时固定为：**file**                           |
| media_id | <font color='red'>是</font> | 图片媒体文件id，可以调用上传临时素材接口获取                 |





# 卡片消息

卡片信息无法使用字体颜色!所以省略这部分内容。

图片链接pic_url是永久保存的，可以一直复用。

~~~json
{
   "touser" : "UserID1|UserID2|UserID3",
   "toparty" : "PartyID1 | PartyID2",
   "totag" : "TagID1 | TagID2",
   "agentid" : {你的agentid},
    
   "msgtype" : "textcard",
   "textcard" : {
            "title" : "领奖通知",
            "description" : 2016年9月26日\n
       恭喜你抽中iPhone 7一台，领奖码：xxxx\n
       请于2016年10月10日前联系行政同事领取",
            "url" : "URL",
            "btntxt":"更多"
   },
}
~~~

| 参数        | 是否必须                               | 说明                                                         |
| ----------- | -------------------------------------- | ------------------------------------------------------------ |
| msgtype     | <font color='red'>是</font>            | 消息类型，此时固定为：**textcard**                           |
| title       | <font color='red'>是</font>            | 标题，不超过128个字节，超过会自动截断                        |
| description | <font color='red'>是</font>            | 描述，不超过512个字节，超过会自动截断                        |
| url         | <font color='red'>是</font>            | 点击后跳转的链接。最长2048字节，请确保包含了协议头(http/https) |
| btntxt      | <font color='cornflowerblue'>否</font> | 按钮文字。 默认为“详情”， 不超过4个文字，超过自动截断。      |

![卡片消息示例](https://res.mail.qq.com/node/ww/wwopenmng/images/independent/doc/cardmsg.jpeg)



# 图文消息

图片链接pic_url是永久保存的，可以一直复用。

~~~json
{
   "touser" : "UserID1|UserID2|UserID3",
   "toparty" : "PartyID1 | PartyID2",
   "totag" : "TagID1 | TagID2",
   "agentid" : {你的agentid},
    
   "msgtype" : "news",
   "news" : {
       "articles" : [
           {
               "title" : "中秋节礼品领取",
               "description" : "今年中秋节公司有豪礼相送",
               "url" : "URL",
               "picurl" : "http://res.mail.qq.com/node/ww/wwopenmng/images/independent/doc/test_pic_msg1.png", 
			   "appid": "wx123123123123123",
        	   "pagepath": "pages/index?userid=zhangsan&orderid=123123123",
           }
                ]
   },
}
~~~

| 参数        | 是否必须                               | 说明                                                         |
| ----------- | -------------------------------------- | ------------------------------------------------------------ |
| msgtype     | <font color='red'>是</font>            | 消息类型，此时固定为：**news**                               |
| articles    | <font color='red'>是</font>            | 图文消息，一个图文消息支持1到8条图文                         |
| title       | <font color='red'>是</font>            | 标题，不超过128个字节，超过会自动截断                        |
| description | <font color='cornflowerblue'>否</font> | 描述，不超过512个字节，超过会自动截断                        |
| url         | <font color='cornflowerblue'>否</font> | 点击后跳转的链接。 最长2048字节，请确保包含了协议头(http/https)，<u>小程序或者url必须填写一个</u> |
| picurl      | <font color='cornflowerblue'>否</font> | 图文消息的图片链接，最长2048字节，支持JPG、PNG格式，较好的效果为大图 1068*455，小图150*150。 |
| appid       | <font color='cornflowerblue'>否</font> | 小程序appid，必须是与当前应用关联的小程序，appid和pagepath必须同时填写，填写后会忽略url字段。没有就删掉吧 |
| pagepath    | <font color='cornflowerblue'>否</font> | 点击消息卡片后的小程序页面，最长128字节，仅限本小程序内的页面。appid和pagepath必须同时填写，填写后会忽略url字段 |



![图文消息示例](https://p.qpic.cn/pic_wework/3478722865/7a7c92760b2bd396e3b856a660f43c8b7db11271bddb3f34/0)



这里分成两部分好了，一个是url跳转，另外一个是小程序跳转



测试代码：

```
{
  "touser": "LiuWenLong",
  "toparty": "PartyID1 | PartyID2",
  "totag": "TagID1 | TagID2",
  "msgtype": "textcard",
  "agentid": 1000002,
  "textcard": {
    "title": "领奖通知",
    "description": "<div class=\"gray\">2016年9月26日</div><div class=\"normal\">恭喜你抽中iPhone 7一台，领奖码：xxxx</div><div class=\"highlight\">请于2016年10月10日前联系行政同事领取</div>",
    "url": "URL",
    "btntxt": "更多"
  },
  "enable_id_trans": 0,
  "enable_duplicate_check": 0,
  "duplicate_check_interval": 1800
}
```

~~~
{
  "touser": "LiuWenLong",
  "toparty": "PartyID1 | PartyID2",
  "totag": "TagID1 | TagID2",
  "msgtype": "news",
  "agentid": 1000002,
  "news": {
    "articles": [{
      "title": "中秋节礼品领取",
      "description": "今年中秋节公司有豪礼相送",
      "url": "URL",
      "picurl": "https://wework.qpic.cn/wwpic/886088_r6OOQ_SPSb26Ptf_1662734751/0",
      "appid": "wx123123123123123",
      "pagepath": "pages/index?userid=zhangsan&orderid=123123123"
    }]
  },
  "enable_id_trans": 0,
  "enable_duplicate_check": 0,
  "duplicate_check_interval": 1800
}
~~~

