# qywxPush

**另一个**企业微信向微信推送消息的PHP版本。这个是一直自用，也贡献出来吧，大家自行选用。

相对于 [easychen的PHP版本](https://github.com/easychen/wecomchan/blob/main/index.php) 有以下区别：

1. 缓存了 access_token ，在SAE环境中用kvdb存token，普通环境写到文件。

> 因为 [企业开发文档](https://work.weixin.qq.com/api/doc/90000/90135/91039) 中提到：开发者需要缓存access_token，用于后续接口的调用（注意：不能频繁调用gettoken接口，否则会受到频率拦截）。当access_token失效或过期时，需要重新获取。

2. 无指定接收人时只发给自己，所有人接收需传值 `&touser=@all` ，账号在 [后台通讯录](https://work.weixin.qq.com/wework_admin/frame#contacts) 中查看

## 推送函数

使用实例：

```php
$corpid = '企业微信公司ID';
$corpsecret = '企业微信应用Secret';
$agentid = '企业微信应用ID';
$me = '自己账号';
    
if( !empty($_REQUEST['text']) ){
    $qywxPush = new qywxPush($corpid, $corpsecret, $agentid);
    // 兼容 Server酱，拼接两个参数：text、desp
	$text = empty($_REQUEST['desp']) ? $_REQUEST['text'] : "【{$_REQUEST['text']}】\n{$_REQUEST['desp']}";
    // 无指定接收人时只发给自己
    $touser = isset($_REQUEST['touser']) ? $_REQUEST['touser'] : $me;
	echo $qywxPush->sendMsg($touser, $text);
}

// class 略…
```

通过 `http://指向运行环境的域名/?text=你要发送的内容` 发送内容，可选项：`&touser=接收人&desp=【标题】内容效果` 

## 其它参考教程

1. 方案及企业微信创建说明 [参考easychen方法](https://github.com/easychen/wecomchan)
2. 自行搭建的在线服务源码 [查看搭建说明](https://github.com/easychen/wecomchan/ONLINE.md) 