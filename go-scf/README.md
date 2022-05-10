# 腾讯云云函数部署Server酱📣

## 注意：腾讯云函数从2022年5月23日起不再有免费额度，最低需要9.9元/月才能使用

本项目是对 [Wecom酱](https://github.com/easychen/wecomchan) 进行的扩展，可以通过企业微信 OpenAPI 向微信推送消息，实现微信消息提醒。

利用 [腾讯云云函数](https://cloud.tencent.com/product/scf)  ServerLess 的能力，以极低的费用（按量付费，且有大量免费额度）来完成部署

优点：

- 便宜：说是免费也不过分
- 简单：不需要购买vps, 也不需要备案, 腾讯云速度有保障.
- 易搭建：一个可执行二进制文件，直接上传至腾讯云函数控制面板即可，虽然使用 Golang 编写，但是搭建无需 Golang 环境
- Serverless：无服务器，函数调用完资源会释放

## 🖐️ 简单介绍

我们要实现的目标是把消息推送到微信上，此处借助了使用 企业微信，可以创建机器人，利用微信的 OpenAPI 来实现消息推送，本项目做了一个简单的封装。

欢迎PR代码。

> 老用户注意：
>
> 自 2.0 版本之后，不再需要 `config.yaml` 文件，配置改为从云函数的环境变量中读取，请直接下载 `main.zip` 上传至云函数并且设置环境变量即可。

## 👋 使用方法

### 1. 注册企业 & 创建机器人 & 获取相关配置信息

此处不再赘述，项目主页有完整的操作方法，见：https://github.com/riba2534/wecomchan

### 2. 下载编译好的二进制文件

下载文件 [版本发布页面](https://github.com/riba2534/wecomchan/releases)：

- [main.zip](https://github.com/riba2534/wecomchan/releases/download/2.1/main.zip) ：云函数可执行二进制文件，不用改动，等会直接上传即可。

### 3. 在腾讯云中创建云函数 & 配置环境变量

打开云函数控制台：https://console.cloud.tencent.com/scf/list

点击新建：

![image-20210705014652334](https://image-1252109614.cos.ap-beijing.myqcloud.com/img/image-20210705014652334.png)

如图所示选择

1. 自定义创建，函数类型为 `事件函数`
2. 填 `wecomchan`
3. 运行环境选择 Go1
4. 函数代码选择本地上传ZIP包，直接上传刚才下载的 `main.zip`
5. 在 `高级配置` 中配置环境变量，6 个环境变量，**缺一不可**，（后续想改环境变量，直接在创建好的函数中编辑即可）

环境变量配置说明

|      key       |                                     value                                      | 备注                                                                                                                                                                                |
| :------------: | :----------------------------------------------------------------------------: | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
|  `FUNC_NAME`   |                            填 `wecomchan`                            |                                                                                                                                                                                     |
|   `SEND_KEY`   | 最终调用HTTP接口时校验是否是本人调用的密钥，随意设置，最终发起HTTP请求携带即可 |                                                                                                                                                                                     |
|  `WECOM_CID`   |                                 企业微信公司ID                                 |                                                                                                                                                                                     |
| `WECOM_SECRET` |                               企业微信应用Secret                               |                                                                                                                                                                                     |
|  `WECOM_AID`   |                                 企业微信应用ID                                 |                                                                                                                                                                                     |
| `WECOM_TOUID`  |                                     `@all`                                     | 此处指推送消息的默认发送对象，填 `@all`，则代表向该企业的全部成员推送消息（如果是个人用的话，一个企业中只有你自己，直接填 `@all` 即可），如果想指定具体发送的人，后面会说明怎么发。 |

6. 在 `触发器配置` 中，新增 `API网关触发`，保持默认配置即可。
7. 点击完成

![基础配置](https://image-1252109614.cos.ap-beijing.myqcloud.com/img/image-20210707204518173.png)

![高级配置](https://image-1252109614.cos.ap-beijing.myqcloud.com/img/image-20210707204936310.png)

![触发器配置](https://image-1252109614.cos.ap-beijing.myqcloud.com/img/image-20210707205811630.png)

稍等一会，进入你创建的函数：

![image-20210705015301810](https://image-1252109614.cos.ap-beijing.myqcloud.com/img/image-20210705015301810.png)

图中所示的访问路径就是函数的请求路径，至此，所有的配置完成。

## 👌 发起HTTP请求测试是否成功

现已支持 `GET`、`POST` 方法进行请求。

>  当发送的文本中存在有换行符或其他字符时，请把 msg 参数进行 url 编码（使用 GET 方法注意，POST不需要）

### 简单使用：

在你刚才获得的路径之后拼几个GET参数，在后面加上：`?sendkey=你配置的sendkey&msg_type=text&msg=hello`

![image-20210705015727720](https://image-1252109614.cos.ap-beijing.myqcloud.com/img/image-20210705015727720.png)

可以看见返回 success 字样。

观察手机推送，也可以收到消息：

![image-20210705015804023](https://image-1252109614.cos.ap-beijing.myqcloud.com/img/image-20210705015804023.png)

之后，想怎么用就是你的事了，想给自己的微信推送，只需要给这个 URL 发一条 HTTP 请求即可。

### 给指定成员推送消息：

如果你的需求是给企业微信中的指定成员发送消息而不是所有成员，则在 GET 请求中多加一个参数 `to_user`，值为 成员ID列表，如果想指定多个成员，则多个成员ID之间用 `|` 隔开。如请求：`https://xxxxx/wecomchan?sendkey=123456&msg_type=text&msg=测试消息&to_user=User1|User2` ，也能收到消息。

![image-20210707211125345](https://image-1252109614.cos.ap-beijing.myqcloud.com/img/image-20210707211125345.png)

> 成员的 ID 在企业微信后台，`通讯录`，点开指定成员资料，有个 `账号` 字段，该字段即为该成员的ID.

### 使用 `POST` 进行请求

大部分情况下，`GET` 请求已经可以很好的满足发送一些短消息的需求，但是当消息体过长时，云函数可能报参数过长错误，故在 `V2.1` 版本加入 `POST` 请求支持。

与 `GET` 请求不同的是，`POST` 请求不从 [Query string](https://en.wikipedia.org/wiki/Query_string) 获取参数，所有参数改为从 [HTTP message body](https://en.wikipedia.org/wiki/HTTP_message_body) 中获取，这里要求 Body 中必须是 `JSON` 格式，参数字段名称仍与 `GET` 请求的名称保持一致，且 `json` 的 `key` 和 `value` 必须是 `string` 类型，Body 格式例如：

```json
{
    "sendkey": "123456",
    "msg_type": "text",
    "msg": "这是一条POST消息",
    "to_user": "User1|User2"
}
```

### 参数说明：

下表为请求的参数说明（`GET` 与 `POST` 字段名相同）：

| 参数名称   | 说明                                                                                                            | 是否可选 |
| ---------- | --------------------------------------------------------------------------------------------------------------- | -------- |
| `sendkey`  | 校验是否是本人调用的密钥，随意设置，最终发起HTTP请求携带即可                                                    | 必须     |
| `msg_type` | 消息类型，目前只有纯文本一种类型，值为 `text`                                                                   | 必须     |
| `msg`      | 消息内容，支持多行和UTF8字符，在程序中构建字符串时加上**换行符**即可，如果有特殊符号，记得使用 `urlencode` 编码 | 必须     |
| `to_user`  | 如果需要给企业内指定成员发消息，可在此参数中指定成员。如果不传本参数，默认所有成员。                            | 可选     |

👇👇👇

---

如果发现bug，或者对本项目有任何建议，欢迎联系 `riba2534@qq.com` 或者直接提 [Issue](https://github.com/riba2534/wecomchan/issues).

