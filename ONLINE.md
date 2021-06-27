# 在线服务搭建指南(PHP版)

## 安装条件

- PHP7.4+ 
- JSON &&CURL 模块
- 可访问外部网络的运行环境

## 安装说明

1. 用编辑器打开 `index.php`，按提示修改头部 define 的值（ sendkey自己随意写，其他参见企业微信配置文档 ）
1. 将 `index.php` 上传运行环境
1. 通过 `http://指向运行环境的域名/?sendkey=你设定的sendkey&text=你要发送的内容` 即可发送内容
