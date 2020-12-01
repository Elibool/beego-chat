
# beego 聊天室 demo

## 说明
主要使用 go 和 beego 实现 web 聊天室功能，网络请求使用 http 配合 websocket 实现；

聊天信息存储在 mysql 内

session 信息存储在 redis

#### 实现的功能：
1. 多人群聊天室
2. 单对单聊天 
3. 用户上下线提醒 
4. 陌生用户发送消息提醒 
5. 消息未读提醒  
6. 消息已读标记 
7. 修改用户头像和暱称 ; 现为随机默认 （待完成 ）
8. 聊天数据保存时间为 7 天，超时自动清除 

## 项目 demo



## 安装
> 本地开发版本：
>
> go version go1.13.10 darwin/amd64
>
> bee version v1.11.0

1. 克隆代码： `git clone https://github.com/Elibool/beego-chat.git`

2. 将 `config/database.conf.exp` 修改为 `config/database.conf`,
且将里面内容修改为自己的配置

3. `config/session.conf.exp` 可更改亦可不更改，如不更改`beego` 默认使用临时文件存储  

## 使用
