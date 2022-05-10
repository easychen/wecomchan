# Wecom酱

通过企业微信向微信推送消息的解决方案。包括：

1. 配置说明（本页下方）
2. 推送函数（支持多种语言，见本页下方）
3. 自行搭建的在线服务源码或使用本地库
   1. [PHP版搭建说明](ONLINE.md) 
   2. [Go版说明](go-wecomchan/README.md)
   3. [腾讯云云函数搭建说明](go-scf/) 
   4. [阿里云云函数搭建说明](python-aliyunfc/)
   5. [无需服务器的 Python 库](https://pypi.org/project/wechat-push/)

## 🎈 本项目属于方糖推送生态。该生态包含项目如下：

- [Server酱Turbo](https://sct.ftqq.com)：支持企业微信、微信服务号、钉钉、飞书群机器人等多通道的在线服务，无需搭建直接使用，每天有免费额度
- [Wecom酱](https://github.com/easychen/wecomchan)：通过企业微信推送消息到微信的消息推送函数和在线服务方案，开源免费，可自己搭建。支持多语言。
- [PushDeer](https://github.com/easychen/pushdeer)：可自行搭建的、无需安装APP的开源推送方案。同时也提供安装APP的降级方案给低版本/没有快应用的系统。支持作为Server酱的通道进行推送，所有支持Server酱的软件和插件都能直接整合PushDeer。

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

应用名称填入「Server酱」，应用logo到[这里](./20210208142819.png)下载，可见范围选择公司名。


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

PS：如果出现`接口请求正常，企业微信接受消息正常，个人微信无法收到消息`的情况：

1. 进入「我的企业」 → 「[微信插件](https://work.weixin.qq.com/wework_admin/frame#profile/wxPlugin)」，拉到最下方，勾选 “允许成员在微信插件中接收和回复聊天消息”
![](https://img.ams1.imgbed.xyz/2021/06/01/HPIRU.jpg)

2. 在企业微信客户端 「我」 → 「设置」  → 「新消息通知」中关闭 “仅在企业微信中接受消息” 限制条件
![](https://img.ams1.imgbed.xyz/2021/06/01/HPKPX.jpg)

#### 第五步，通过以下函数发送消息：

PS：为使用方便，以下函数没有对 `access_token` 进行缓存。对于个人低频调用已经够用。带缓存的实现可查看 `index.php` 中的示例代码（依赖Redis实现）。

PHP版：

```php
function send_to_wecom($text, $wecom_cid, $wecom_aid, $wecom_secret,  $wecom_touid = '@all')
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

PYTHON版:

```python
import json,requests,base64
def send_to_wecom(text,wecom_cid,wecom_aid,wecom_secret,wecom_touid='@all'):
    get_token_url = f"https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid={wecom_cid}&corpsecret={wecom_secret}"
    response = requests.get(get_token_url).content
    access_token = json.loads(response).get('access_token')
    if access_token and len(access_token) > 0:
        send_msg_url = f'https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token={access_token}'
        data = {
            "touser":wecom_touid,
            "agentid":wecom_aid,
            "msgtype":"text",
            "text":{
                "content":text
            },
            "duplicate_check_interval":600
        }
        response = requests.post(send_msg_url,data=json.dumps(data)).content
        return response
    else:
        return False

def send_to_wecom_image(base64_content,wecom_cid,wecom_aid,wecom_secret,wecom_touid='@all'):
    get_token_url = f"https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid={wecom_cid}&corpsecret={wecom_secret}"
    response = requests.get(get_token_url).content
    access_token = json.loads(response).get('access_token')
    if access_token and len(access_token) > 0:
        upload_url = f'https://qyapi.weixin.qq.com/cgi-bin/media/upload?access_token={access_token}&type=image'
        upload_response = requests.post(upload_url, files={
            "picture": base64.b64decode(base64_content)
        }).json()
        if "media_id" in upload_response:
            media_id = upload_response['media_id']
        else:
            return False

        send_msg_url = f'https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token={access_token}'
        data = {
            "touser":wecom_touid,
            "agentid":wecom_aid,
            "msgtype":"image",
            "image":{
                "media_id": media_id
            },
            "duplicate_check_interval":600
        }
        response = requests.post(send_msg_url,data=json.dumps(data)).content
        return response
    else:
        return False

def send_to_wecom_markdown(text,wecom_cid,wecom_aid,wecom_secret,wecom_touid='@all'):
    get_token_url = f"https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid={wecom_cid}&corpsecret={wecom_secret}"
    response = requests.get(get_token_url).content
    access_token = json.loads(response).get('access_token')
    if access_token and len(access_token) > 0:
        send_msg_url = f'https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token={access_token}'
        data = {
            "touser":wecom_touid,
            "agentid":wecom_aid,
            "msgtype":"markdown",
            "markdown":{
                "content":text
            },
            "duplicate_check_interval":600
        }
        response = requests.post(send_msg_url,data=json.dumps(data)).content
        return response
    else:
        return False
```

使用实例：

```python
ret = send_to_wecom("推送测试\r\n测试换行", "企业ID③", "应用ID①", "应用secret②");
print( ret );
ret = send_to_wecom('<a href="https://www.github.com/">文本中支持超链接</a>', "企业ID③", "应用ID①", "应用secret②");
print( ret );
ret = send_to_wecom_image("此处填写图片Base64", "企业ID③", "应用ID①", "应用secret②");
print( ret );
ret = send_to_wecom_markdown("**Markdown 内容**", "企业ID③", "应用ID①", "应用secret②");
print( ret );
```

TypeScript 版:

```typescript
import request from 'superagent'

async function sendToWecom(body: {
  text: string
  wecomCId: string
  wecomSecret: string
  wecomAgentId: string
  wecomTouid?: string
}): Promise<{ errcode: number; errmsg: string; invaliduser: string }> {
  body.wecomTouid = body.wecomTouid ?? '@all'
  const getTokenUrl = `https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=${body.wecomCId}&corpsecret=${body.wecomSecret}`
  const getTokenRes = await request.get(getTokenUrl)
  const accessToken = getTokenRes.body.access_token
  if (accessToken?.length <= 0) {
    throw new Error('获取 accessToken 失败')
  }
  const sendMsgUrl = `https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=${accessToken}`
  const sendMsgRes = await request.post(sendMsgUrl).send({
    touser: body.wecomTouid,
    agentid: body.wecomAgentId,
    msgtype: 'text',
    text: {
      content: body.text,
    },
    duplicate_check_interval: 600,
  })
  return sendMsgRes.body
}
```

使用实例：

```typescript
sendToWecom({
  text: '推送测试\r\n测试换行',
  wecomAgentId: '应用ID①',
  wecomSecret: '应用secret②',
  wecomCId: '企业ID③',
})
  .then((res) => {
    console.log(res)
  })
  .catch((err) => {
    console.log(err)
  })
```

.NET Core 版:

```C#
using System;
using RestSharp;
using Newtonsoft.Json;
namespace WeCom.Demo
{
    class WeCom
    {   
        public  string SendToWeCom(
            string text,// 推送消息
            string weComCId,// 企业Id①
            string weComSecret,// 应用secret②
            string weComAId,// 应用ID③
            string weComTouId = "@all")
        {
            // 获取Token
            string getTokenUrl = $"https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid={weComCId}&corpsecret={weComSecret}";
            string token = JsonConvert
            .DeserializeObject<dynamic>(new RestClient(getTokenUrl)
            .Get(new RestRequest()).Content).access_token;
            System.Console.WriteLine(token);
            if (!String.IsNullOrWhiteSpace(token))
            {
                var request = new RestRequest();
                var client = new RestClient($"https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token={token}");
                var data = new
                {
                    touser = weComTouId,
                    agentid = weComAId,
                    msgtype = "text",
                    text = new
                    {
                        content = text
                    },
                    duplicate_check_interval = 600
                };
                string serJson = JsonConvert.SerializeObject(data);
                System.Console.WriteLine(serJson);
                request.Method = Method.POST;
                request.AddHeader("Accept", "application/json");
                request.Parameters.Clear();
                request.AddParameter("application/json", serJson, ParameterType.RequestBody);
                return client.Execute(request).Content;
            }
            return "-1";
        }
}


```
使用实例:
```C#
   static void Main(string[] args)
        {   // 测试
            Console.Write(new WeCom().SendToWeCom(
            "msginfo",
            "企业Id①"
            , "应用secret②",
            "应用ID③"
            ));
        }

    }
```

- [纯Bash版本参考](https://gitee.com/Hemingway2003/pushservice/blob/master/wecom.sh)

其他版本的函数可参照上边的逻辑自行编写，欢迎PR。

发送图片、卡片、文件或 Markdown 消息的高级用法见 [企业微信API](https://work.weixin.qq.com/api/doc/90000/90135/90236)。



