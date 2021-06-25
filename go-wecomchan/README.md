# go-wecomchan 

## 配置说明

直接使用和构建二进制文件使用需要golang环境，并且网络可以安装依赖。
docker构建镜像使用，需要安装docker，不依赖golang以及网络。

## 直接使用

`go run .`

## build命令构建二进制文件使用
1. 构建命令
`go build`

2. 启动
`./wecomchan`

## 构建docker镜像使用（推荐，不依赖golang，不依赖网络）

1. 构建镜像
`docker build -t go-wecomchan .`

2. 启动镜像
`docker run -dit -p 8080:8080 go-wecomchan`

## 调用方式

访问 `http://localhost:8080/wecomchan?sendkey=你配置的sendkey&&msg=需要发送的消息&&msg_type=text`

## 后续预计添加

~~- Dockerfile 打包镜像(不依赖网络环境)~~
- docker-compose redis + go-wecomchan 一键部署
- 通过环境变量传递企业微信id，secret等，镜像一次构建多次使用
