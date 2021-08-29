#!/usr/bin/env bash
#author : junun
# SERVER -- alarm_center email api 地址
# mail_to -- 发送给谁
# subject -- 主题
# mail_type -- 0: 异步发送； 1：实时发送
# content -- 消息内容

SERVER="http://127.0.0.1:1338/api/v1/email/"
MAIL_FROM="470499989@qq.com"


function sendEmail() {
    curl --request POST \
        --url ${SERVER}\
        --header 'content-type: multipart/form-data' \
        --form mail_from=${MAIL_FROM} \
        --form mail_to="${args[0]}" \
        --form subject="${args[1]}" \
        --form content="${args[2]}" \
        --form mail_type="${args[3]}"
}

args=("$@")
sendEmail