#!/usr/bin/env bash
#author : junun
# SERVER -- alarm_center 钉钉api地址
# DING_CHANNEL -- 配置文件中钉钉机器人的名字
# msgtype -- markdown类型， link类型，其他为text类型
# content -- 消息内容

SERVER="http://127.0.0.1:1338/api/v1/dt/"
DING_CHANNEL="sops"

function sendDingTalk() {
    curl --request POST \
        --url ${SERVER}\
        --header 'content-type: multipart/form-data' \
        --form msgtype=1 \
        --form name=${DING_CHANNEL} \
        --form content="$1"
}

sendDingTalk "$@"