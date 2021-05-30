# Wecom酱

通过企业微信向微信推送消息(配置说明&推送函数)。

## 企业微信应用消息配置说明

优点：

1. 一次配置，持续使用
1. 配置好以后，只需要微信就能收消息，不再需要安装企业微信客户端

PS：消息接口无需认证即可使用，个人用微信就可以注册

### 具体操作

#### 第一步，注册企业

用电脑打开[企业微信官网](https://work.weixin.qq.com/)，注册一个企业

#### 第二步，创建应用

注册成功后，点「管理企业」进入管理界面，选择「应用管理」 → 「自建」 →  「创建应用」

![](https://theseven.ftqq.com/20210208143228.png)

应用名称填入「Server酱」，应用logo到[这里](https://theseven.ftqq.com/20210208142819.png)下载，可见范围选择公司名。


![](https://theseven.ftqq.com/20210208143327.png)

创建完成后进入应用详情页，可以得到应用ID( `agentid` )①，应用Secret( `secret` )②。

注意：`secret`推送到手机端时，只能在`企业微信客户端`中查看。

![](https://theseven.ftqq.com/20210208143553.png)

#### 第三步，获取企业ID

进入「[我的企业](https://work.weixin.qq.com/wework_admin/frame#profile)」页面，拉到最下边，可以看到企业ID③，复制并填到上方。

推送UID直接填 `@all` ，推送给公司全员。

#### 第四步，推送消息到微信

进入「我的企业」 → 「[微信插件](https://work.weixin.qq.com/wework_admin/frame#profile/wxPlugin)」，拉到下边扫描二维码，关注以后即可收到推送的消息。

![](https://theseven.ftqq.com/20210208144808.png)

#### 第五步，通过以下函数发送消息：

PHP版：

```php
function send_to_wecom($text, $wecom_cid, $wecom_secret, $wecom_aid, $wecom_touid = '@all')
{
    $info = @json_decode(file_get_contents("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=".urlencode($wecom_cid)."&corpsecret=".urlencode($wecom_secret)), true);
                
    if ($info && isset($info['access_token']) && strlen($info['access_token']) > 0) {
        $access_token = $info['access_token'];
        $url = 'https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token='.urlencode($access_token);
        $data = new \stdClass();
        $data->touser = $wecom_touid;
        $data->agentid = $wecom_aid;
        $data->msgtype = "text";
        $data->text = ["content"=> $text];
        $data->duplicate_check_interval = 600;

        $data_json = json_encode($data);
        $ch = curl_init();
        curl_setopt($ch, CURLOPT_HTTPHEADER, ['Content-Type: application/json']);
        curl_setopt($ch, CURLOPT_URL, $url);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        @curl_setopt($ch, CURLOPT_FOLLOWLOCATION, true);
        curl_setopt($ch, CURLOPT_POST, true);
        curl_setopt($ch, CURLOPT_TIMEOUT, 5);
        curl_setopt($ch, CURLOPT_POSTFIELDS, $data_json);

        curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, false);
        curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false);
        
        $response = curl_exec($ch);
        return $response;
    }
    return false;
}

```

使用实例：

```php
$ret = send_to_wecom("推送测试\r\n测试换行", "企业ID③", "应用ID①", "应用secret②");
print_r( $ret );
```

其他版本的函数可参照上边的逻辑自行编写，欢迎PR。


