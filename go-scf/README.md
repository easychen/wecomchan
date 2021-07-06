# 腾讯云云函数部署Server酱

本项目是对 [Wecom酱](https://github.com/easychen/wecomchan) 进行的扩展。

利用 [腾讯云云函数](https://cloud.tencent.com/product/scf)  ServerLess 的能力，以极低的费用（按量付费，且有大量免费额度）来完成部署，优点：

- 便宜：说是免费也不过分
- 易搭建：一个可执行二进制文件加一个配置文件，直接上传至腾讯云函数控制面板即可，虽然使用 Golang 编写，但是搭建无需 Golang 环境
- Serverless：无服务器，不用自己购买VPS

## 简单介绍

我们要实现的目标是把消息推送到微信上，此处借助了使用 企业微信，可以创建机器人，利用微信的 OpenAPI 来实现消息推送，本项目做了一个简单的封装。

欢迎PR代码。

## 使用方法

### 1. 注册企业 & 创建机器人 & 获取相关配置信息

此处不再赘述，项目主页有完整的操作方法，见：https://github.com/easychen/wecomchan

### 2. 下载编译好的二进制文件，填写配置

下载文件：

- [main](https://github.com/riba2534/wecomchan/releases/download/1.0/main) ：腾讯云云函数，可执行二进制文件
- [config.yaml.example](https://github.com/riba2534/wecomchan/releases/download/1.0/config.yaml.example) ： 示例的配置文件

修改 `config.yaml.example` 中的内容，先把文件名改为 `config.yaml`，在对应位置填入相关配置，企业微信相关配置，请参考项目主页获取。

配置中的 `FUNC_NAME` 指下一步在腾讯云创建云函数时的函数名称，推荐填 `wecomchan`

配置中的 `SEND_KEY` 指最终调用HTTP接口时校验是否是本人调用的密钥，随意设置，最终发起HTTP请求携带即可。

```yaml
config:
  FUNC_NAME: 函数名称
  SEND_KEY: 这里写你设置的send_key
  WECOM_CID: 企业微信公司ID
  WECOM_SECRET: 企业微信应用Secret
  WECOM_AID: 企业微信应用ID
  WECOM_TOUID: "@all" # 别改
```

修改完配置完成后，把 `main` 和 `config.yaml` 这两个文件打包在一个 ZIP 文件中（注意，压缩包中不要有其他目录，应该只包含这两个文件），打包zip的方法就不赘述了，你甚至可以用360压缩

![image-20210705014521569](https://image-1252109614.cos.ap-beijing.myqcloud.com/img/image-20210705014521569.png)

### 3. 在腾讯云中创建云函数

打开云函数控制台：https://console.cloud.tencent.com/scf/list

点击新建：

![image-20210705014652334](https://image-1252109614.cos.ap-beijing.myqcloud.com/img/image-20210705014652334.png)

如图所示选择

1. 自定义创建
2. 函数名称填你在 `config.yaml` 中配置的那个名称
3. 运行环境选择 Go1
4. 函数代码选择本地上传ZIP包
5. 触发器配置里面新增一个 API网关触发，保持默认配置就行
6. 点击完成

![image-20210705015022784](https://image-1252109614.cos.ap-beijing.myqcloud.com/img/image-20210705015022784.png)

稍等一会，进入你创建的函数：

![image-20210705015301810](https://image-1252109614.cos.ap-beijing.myqcloud.com/img/image-20210705015301810.png)

图中所示的访问路径就是函数的请求路径。

## 测试是否成功

在你刚才获得的路径之后拼几个GET参数，在后面加上：`?sendkey=你配置的sendkey&msg_type=text&msg=hello`

![image-20210705015727720](https://image-1252109614.cos.ap-beijing.myqcloud.com/img/image-20210705015727720.png)

可以看见返回 success 字样。

观察手机推送，也可以收到消息：

![image-20210705015804023](https://image-1252109614.cos.ap-beijing.myqcloud.com/img/image-20210705015804023.png)

之后，想怎么用就是你的事了，想给自己的微信推送，只需要给这个 URL 发一条 HTTP 请求即可。

---

如果发现bug，或者对本文档有任何建议，欢迎联系 `riba2534@qq.com`

