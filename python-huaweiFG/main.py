# coding=utf-8
import logging
import json
import requests
import base64

# 修改为公司ID
WECOM_ID = "修改为自己的公司ID"


def send_to_wecom(text, wecom_cid, wecom_aid, wecom_secret, wecom_touid):
    get_token_url = f"https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid={wecom_cid}&corpsecret={wecom_secret}"
    response = requests.get(get_token_url).content
    access_token = json.loads(response).get('access_token')
    if access_token and len(access_token) > 0:
        send_msg_url = f'https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token={access_token}'
        data = {
            "touser": wecom_touid,
            "agentid": wecom_aid,
            "msgtype": "text",
            "text": {
                "content": text
            },
            "duplicate_check_interval": 600
        }
        response = requests.post(send_msg_url, data=json.dumps(data)).content
        return hook_return(str(response))
    else:
        return hook_return(None)


def send_to_wecom_markdown(text, wecom_cid, wecom_aid, wecom_secret, wecom_touid):
    get_token_url = f"https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid={wecom_cid}&corpsecret={wecom_secret}"
    response = requests.get(get_token_url).content
    access_token = json.loads(response).get('access_token')
    if access_token and len(access_token) > 0:
        send_msg_url = f'https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token={access_token}'
        data = {
            "touser": wecom_touid,
            "agentid": wecom_aid,
            "msgtype": "markdown",
            "markdown": {
                "content": text
            },
            "duplicate_check_interval": 600
        }
        response = requests.post(send_msg_url, data=json.dumps(data)).content
        return hook_return(str(response))
    else:
        return hook_return(None)


def send_to_wecom_pic(base64_content, wecom_cid, wecom_aid, wecom_secret, wecom_touid):
    get_token_url = f"https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid={wecom_cid}&corpsecret={wecom_secret}"
    response = requests.get(get_token_url).content
    access_token = json.loads(response).get('access_token')
    if access_token and len(access_token) > 0:
        upload_url = f'https://qyapi.weixin.qq.com/cgi-bin/media/upload?access_token={access_token}&type=image'
        upload_response = requests.post(upload_url, files={
            "picture": base64.b64decode(base64_content)
        }).json()

        logging.info('upload response: ' + str(upload_response))

        media_id = upload_response['media_id']

        send_msg_url = f'https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token={access_token}'
        data = {
            "touser": wecom_touid,
            "agentid": wecom_aid,
            "msgtype": "image",
            "image": {
                "media_id": media_id
            },
            "duplicate_check_interval": 600
        }
        response = requests.post(send_msg_url, data=json.dumps(data)).content
        return hook_return(str(response))
    else:
        return hook_return(None)


def send_to_wecom_file(base64_content, file_name, wecom_cid, wecom_aid, wecom_secret, wecom_touid):
    get_token_url = f"https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid={wecom_cid}&corpsecret={wecom_secret}"
    response = requests.get(get_token_url).content
    access_token = json.loads(response).get('access_token')
    if access_token and len(access_token) > 0:
        upload_url = f'https://qyapi.weixin.qq.com/cgi-bin/media/upload?access_token={access_token}&type=file'
        upload_response = requests.post(upload_url + "&debug=1", files={
            "media": (file_name, base64.b64decode(base64_content))  # 此处上传中文文件名文件旧版本 urllib 有 bug.
        }).json()

        logging.info('upload response: ' + str(upload_response))

        media_id = upload_response['media_id']

        send_msg_url = f'https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token={access_token}'
        data = {
            "touser": wecom_touid,
            "agentid": wecom_aid,
            "msgtype": "file",
            "file": {
                "media_id": media_id
            },
            "duplicate_check_interval": 600
        }
        response = requests.post(send_msg_url, data=json.dumps(data)).content
        return hook_return(str(response))
    else:
        return hook_return(None)


def handler(event, context):
    request_body = json.dumps(event.get('queryStringParameters', ''))

    print(request_body)
    # get path info
    try:
        path_info = event.get('path')
    except ValueError:
        path_info = "ERROR"

    logging.info("request body: " + request_body)
    logging.info("path_info: " + str(path_info))

    # path_info 修改为触发器对应的值,上面设定的 / 这里也写 / 
    # 应用1
    if path_info == "/WecomPush":
        send_key = '修改为应用1的'
        wecom_agentid = "修改为应用1的"
        wecom_secret = "修改为应用1的"

    #    # 应用2
    #    elif path_info == "/app2":
    #        send_key = "修改为应用2的"
    #        wecom_agentid = "修改为应用2的"
    #        wecom_secret = "修改为应用2的"
    #    # 应用3
    #    elif path_info == "/app3":
    #        send_key = "修改为应用3的"
    #        wecom_agentid = "修改为应用3的"
    #        wecom_secret = "修改为应用3的"
    #   # 如果有多个应用，模拟上面添加多个elif即可。path_info可自定义，但需要添加和path_info一致的HTTP触发器

    else:
        response = '{"code": -6, "msg": "invalid path info"}'
        return hook_return(str(response))
    # input_json = None
    try:
        input_json = json.loads(request_body)
        if input_json['key'] != send_key:
            status = '403 Forbidden'
            response_headers = [('Content-type', 'text/json')]
            response = '{"code": -2, "msg": "invalid send key"}'
            return hook_return(str(response))
    except Exception as e:
        logging.exception(e)
        status = '403 Forbidden'
        response_headers = [('Content-type', 'text/json')]
        response = '{"code": -1, "msg": "invalid json input"}'
        return hook_return(str(response))
    # 获取发送的用户
    wecom_touid = input_json.get('uid', '@all')

    logging.info("wecom_touid: " + str(wecom_touid))

    code = 0
    msg = "ok"
    status = '200 OK'

    try:
        if 'type' not in input_json or input_json['type'] == 'text':
            result = send_to_wecom(input_json['msg'], WECOM_ID, wecom_agentid, wecom_secret, wecom_touid)
        elif input_json['type'] == 'image':
            result = send_to_wecom_pic(input_json['msg'], WECOM_ID, wecom_agentid, wecom_secret, wecom_touid)
        elif input_json['type'] == 'markdown':
            result = send_to_wecom_markdown(input_json['msg'], WECOM_ID, wecom_agentid, wecom_secret, wecom_touid)
        elif input_json['type'] == 'file':
            if 'filename' in input_json:
                result = send_to_wecom_file(input_json['msg'], input_json['filename'], WECOM_ID, wecom_agentid,
                                            wecom_secret, wecom_touid)
            else:
                result = send_to_wecom_file(input_json['msg'], "Wepush推送", WECOM_ID, wecom_agentid, wecom_secret,
                                            wecom_touid)
                msg = "filename not found. using default."
        else:
            code = -5
            msg = "invalid msg type. type should be text(default), image, markdown or file."
            status = "500 Internal Server Error"
            result = ""

        logging.info('wechat api response: ' + str(result))
        if result is None:
            status = "500 Internal Server Error"
            code = -4
            msg = "wechat api error: wrong config?"
    except Exception as e:
        status = "500 Internal Server Error"
        code = -3
        msg = "unexpected error: " + str(e)
        logging.exception(e)

    response_headers = [('Content-type', 'text/json')]
    response = json.dumps({"code": code, "msg": msg})
    return hook_return(str(response))


def hook_return(string):
    return {
        "statusCode": 200,
        "isBase64Encoded": False,
        "body": string,
        "headers": {
            "Content-Type": "application/json"
        }
    }
