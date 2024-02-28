<?php
// config
// ======================================
define('SENDKEY', 'set_a_sendkey');
define('WECOM_CID', '企业微信公司ID');
define('WECOM_SECRET', '企业微信应用Secret');
define('WECOM_AID', '企业微信应用ID');
define('WECOM_TOUID', '@all');

// 以下配置需要有 redis 服务和 phpredis 扩展
define('REDIS_ON', false);
define('REDIS_HOST', '127.0.0.1');
define('REDIS_PORT', '6379');
define('REDIS_EXPIRED', '7000');
define('REDIS_PASSWORD', '');
define('REDIS_KEY', 'wecom_access_token');

// code
// ======================================

if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    if (strlen(@$_POST['sendkey']) < 1 || strlen(@$_POST['msg']) < 1 || @$_POST['sendkey'] != SENDKEY) {
        die('bad params');
    }

    header("Content-Type: application/json; charset=UTF-8");
    echo send_to_wecom(@$_POST['msg'], WECOM_CID, WECOM_SECRET, WECOM_AID, WECOM_TOUID);
} elseif ($_SERVER['REQUEST_METHOD'] === 'GET') {
    if (strlen(@$_GET['sendkey']) < 1 || strlen(@$_GET['msg']) < 1 || @$_GET['sendkey'] != SENDKEY) {
        die('bad params');
    }

    header("Content-Type: application/json; charset=UTF-8");
    echo send_to_wecom(@$_GET['msg'], WECOM_CID, WECOM_SECRET, WECOM_AID, WECOM_TOUID);
} else {
    die('Method not supported');
}


function redis()
{
    if (!isset($GLOBALS['REDIS_INSTANCE']) || !$GLOBALS['REDIS_INSTANCE']) {
        $GLOBALS['REDIS_INSTANCE'] = new Redis();
        $GLOBALS['REDIS_INSTANCE']->connect(REDIS_HOST, REDIS_PORT);
        $GLOBALS['REDIS_INSTANCE']->auth(REDIS_PASSWORD);
    }

    return $GLOBALS['REDIS_INSTANCE'];
}

function send_to_wecom($text, $wecom_cid, $wecom_secret, $wecom_aid, $wecom_touid = '@all')
{
    $access_token = false;
    // 如果启用redis作为缓存
    if (REDIS_ON) {
        $access_token = redis()->get(REDIS_KEY);
    }

    if (!$access_token) {
        $info = @json_decode(file_get_contents("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=".urlencode($wecom_cid)."&corpsecret=".urlencode($wecom_secret)), true);
                
        if ($info && isset($info['access_token']) && strlen($info['access_token']) > 0) {
            $access_token = $info['access_token'];
        }
    }
    
    if ($access_token) {
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
        if ($response !== false && REDIS_ON) {
            redis()->set(REDIS_KEY, $access_token, ['nx', 'ex'=>REDIS_EXPIRED]);
        }
        return $response;
    }
    
    
    return false;
}
