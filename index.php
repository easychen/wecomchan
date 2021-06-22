<?php

define('SENDKEY', 'set_a_sendkey');
define('WECOM_CID', '企业微信公司ID');
define('WECOM_SECRET', '企业微信应用Secret');
define('WECOM_AID', '企业微信应用ID');
define('WECOM_TOUID', '@all');

if (strlen(@$_REQUEST['sendkey'])  < 1
    || strlen(@$_REQUEST['text'])  < 1 || @$_REQUEST['sendkey'] != SENDKEY
) {
    die('bad params');
}

header("Content-Type: application/json; charset=UTF-8");
echo send_to_wecom(@$_REQUEST['text'], WECOM_CID, WECOM_SECRET, WECOM_AID, WECOM_TOUID);


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
