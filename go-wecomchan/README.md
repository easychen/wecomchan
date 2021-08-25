# go-wecomchan 

## what's new

添加 Dockerfile.architecture 使用docker buildx支持构建多架构镜像。

关于docker buildx build 使用方式参考官方文档:

[https://docs.docker.com/engine/reference/commandline/buildx_build/](https://docs.docker.com/engine/reference/commandline/buildx_build/)

## 配置说明

直接使用和构建二进制文件使用需要golang环境，并且网络可以安装依赖。  
docker构建镜像使用，需要安装docker，不依赖golang以及网络。  

## 修改默认值

修改的sendkey，企业微信公司ID 等默认值为你的企业中的相关信息，如不设置运行时和打包后都可通过环境变量传入。

```golang
var Sendkey = GetEnvDefault("SENDKEY", "set_a_sendkey")
var WecomCid = GetEnvDefault("WECOM_CID", "企业微信公司ID")
var WecomSecret = GetEnvDefault("WECOM_SECRET", "企业微信应用Secret")
var WecomAid = GetEnvDefault("WECOM_AID", "企业微信应用ID")
var WecomToUid = GetEnvDefault("WECOM_TOUID", "@all")
var RedisStat = GetEnvDefault("REDIS_STAT", "OFF")
var RedisAddr = GetEnvDefault("REDIS_ADDR", "localhost:6379")
var RedisPassword = GetEnvDefault("REDIS_PASSWORD", "")
```

## 直接使用

如果没有添加默认值，需要先引入环境变量，以SENDKEY为例：

`export SENDKEY=set_a_sendkey`
依次引入环境变量后，执行
`go run .`

## build命令构建二进制文件使用

1. 构建命令
`go build`

2. 启动
`./wecomchan`

## 构建docker镜像使用（推荐，不依赖golang，不依赖网络）

新增打包好的镜像可以直接使用

- 推送文本or图片:`docker pull aozakiaoko/go-wecomchan`  
Docker Hub 地址为:[https://hub.docker.com/r/aozakiaoko/go-wecomchan](https://hub.docker.com/r/aozakiaoko/go-wecomchan)  

已经更新latest镜像为 @fcbhank 的最新代码，并支持arm64设备。也可通过aozakiaoko/go-wecomchan:v2 获取最新镜像。

- v2_推送文本or图片:`docker pull fcbhank/go-wecomchan`
Docker Hub 地址为:[https://hub.docker.com/r/fcbhank/go-wecomchan](https://hub.docker.com/r/fcbhank/go-wecomchan)

1. 构建镜像
`docker build -t go-wecomchan .`

2. 修改默认值后启动镜像
`docker run -dit -p 8080:8080 go-wecomchan`

3. 通过环境变量启动镜像并启用redis

```bash
docker run -dit -e SENDKEY=set_a_sendkey \
-e WECOM_CID=企业微信公司ID \
-e WECOM_SECRET=企业微信应用Secret \
-e WECOM_AID=企业微信应用ID \
-e WECOM_TOUID="@all" \
-e REDIS_STAT=ON \
-e REDIS_ADDR="localhost:6379" \
-e REDIS_PASSWORD="" \
# aozakiaoko/go-wecomchan 已经更新镜像为 @fcbhank 的最新代码，并支持arm64设备。
# v2 fcbhank/go-wecomchan
-p 8080:8080 go-wecomchan
```

如不使用redis不要传入最后三个关于redis的环境变量(REDIS_STAT|REDIS_ADDR|REDIS_PASSWORD)

4. 环境变量说明

|名称|描述|
|---|---|
|SENDKEY|发送时用来验证的key|
|WECOM_CID|企业微信公司ID|
|WECOM_SECRET|企业微信应用Secret|
|WECOM_AID|企业微信应用ID|
|WECOM_TOUID|需要发送给的人，详见[企业微信官方文档](https://work.weixin.qq.com/api/doc/90000/90135/90236#%E6%96%87%E6%9C%AC%E6%B6%88%E6%81%AF)|
|REDIS_STAT|是否启用redis换缓存token,ON-启用 OFF或空-不启用|
|REDIS_ADDR|redis服务器地址，如不启用redis缓存可不设置|
|REDIS_PASSWORD|redis的连接密码，如不启用redis缓存可不设置|

## 使用docker-compose 部署

修改docker-compose.yml 文件内上述的环境变量，之后执行

`docker-compose up -d`

## 调用方式
- v1_推送文本
访问 `http://localhost:8080/wecomchan?sendkey=你配置的sendkey&&msg=需要发送的消息&&msg_type=text`

- v2_推送文本or图片

```bash
# 推送文本消息
curl --location --request GET 'http://localhost:8080/wecomchan?sendkey={你的sendkey}&msg={你的文本消息}&msg_type=text'

# 推送图片消息
curl --location --request POST 'http://localhost:8080/wecomchan?sendkey={你的sendkey}&msg_type=image' \
--form 'media=@"test.jpg"'
```

## 后续预计添加

* [x] Dockerfile 打包镜像(不依赖网络环境)
* [x] 通过环境变量传递企业微信id，secret等，镜像一次构建多次使用
* [x] docker-compose redis + go-wecomchan 一键部署