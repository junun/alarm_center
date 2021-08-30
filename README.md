# alarm_center
标签（空格分隔）： go gorm mysql redis cron ddd

---
[TOC]
##  简介
写这个项目主要有以下几个目的：

1. 为运维检查脚本或者保活脚本提供一个通用的告警接口，使得运维只需要管着业务本身逻辑实现。
2. 避免钉钉告警1分钟超过20条后禁用通知功能（10分钟）。
3. 封装钉钉通知，支持每分钟超过18条后，通知消息放到队列延迟发送。
4. 封装email，支持异步发送。
5. 封装cron，支持自定义任务名。
6. 尝试是ddd的组织思想完成代码编程。
7. 配置管理使用viper。
8. 依赖注入使用dig。

## 注意事项
1. 本项目mysql没有实际作用，只是项目的一个部分。
2. 本项目所有user文件都没有实际意义，只是为了展示怎样与mysql交互。
3. 本项目的分层设计参考了 https://github.com/daheige/inject-demo 。
4. 本项目强依赖redis；目前只实现单例redis，保留了哨兵和集群接口。
5. 注意shell脚本带有空格参数传递问题。

## 使用
### curl 请求钉钉消息接口

```shell
$ bash dingtalk.sh '测试钉钉发送by curl'
{"code":200,"data":"","message":"发送成功"}%       
```

> 发送超过限制

```
{
    "code": 500,
    "data": "",
    "message": "已经到达单分钟最大的发送能力：3"
}
```
```
等这个key的窗口期过后，计划任务（StartJobOnBoot）会从queue读取消息，继续发送。
```

### curl 请求邮件消息接口

```
$ bash email.sh '470499998@qq.com,1186108666@qq.com' 'mail by curl subject' 'mail by curl context' 1 
{"code":200,"data":"","message":"邮件发送成功"}%   
```

## License
alarm_center is licensed under the MIT license.
