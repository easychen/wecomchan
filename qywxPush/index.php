<?php
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


class qywxPush {
  
    function __construct( $corpid, $corpsecret, $agentid) {
        $this->corpid = $corpid;
        $this->corpsecret = $corpsecret;
        $this->agentid = $agentid;
        // 兼容SAE，SAE中用kvdb存token，普通环境写到文件
        if(defined('SAE_ACCESSKEY')){
            $this->tokenFilename = 'saekv://qywxToken_'.$this->corpid.'_'.$this->agentid;
        }else{
            $this->tokenFilename = '.\qywxToken_'.$this->corpid.'_'.$this->agentid;
        }
        $this->access_token = @file_get_contents($this->tokenFilename);
    }
    
    function updateToken(){
        $url = 'https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid='.$this->corpid.'&corpsecret='.$this->corpsecret;
        $json = file_get_contents($url);
        $data = json_decode($json,true);
        $this->access_token = $data['access_token'];
        file_put_contents($this->tokenFilename, $data['access_token']);
        return $data;
    }

    function sendMsg($touser, $content){
        $message = array(
            'touser' => $touser,
            'msgtype' => 'text',
            'agentid' => $this->agentid,
            'text' => array(
                'content' => $content
            )
        );
        $url = 'https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token='.$this->access_token;
        $return = $this->post($url, json_encode($message));
        $data = json_decode($return,true);
        if($data['errcode']>0){
            if($data['errcode']==40014 || $data['errcode']==41001 || $data['errcode']==42001){
                $this->updateToken();
                return $this->sendMsg($touser, $content);
            }else{
                $data['message'] = $message;
                return json_encode($data, 512);
            }
        }else{
            return $return;
        }
    }

    function post($url, $postdata){
        if(is_array($postdata))
            $postdata = http_build_query($postdata);
        $opts = array(
            'http' =>array(
                'method' => 'POST', 
                'header' => 'Content-type: application/x-www-form-urlencoded', 
                'content' => $postdata 
            )
        );
        $context = stream_context_create($opts);
        $result = file_get_contents($url, false, $context);
        return $result;
    }

}