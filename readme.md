
# beego 聊天室 demo

## 说明
主要使用 go 和 beego 实现 web 聊天室功能，网络请求使用 http 配合 websocket 实现；

聊天信息存储在 mysql 内

session 信息存储在 redis

实现的功能：





## 项目 demo



## 安装
> go version go1.13.10 darwin/amd64
>
> bee version v1.11.0

`git clone https://github.com/Elibool/beego-chat.git`

1. 将 `config/database.conf.exp` 修改为 `config/database.conf`,
且将里面内容修改为自己的配置

2. `config/session.conf.exp` 可更改亦可不更改，如不更改`beego` 默认使用临时文件存储  

## 使用
