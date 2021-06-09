## go-im

> go 消息中心模型

### Run

> Docker 运行 rabbitmq 服务

* cd docker/rabbitmq && docker-compose up

> 运行 websocket 服务

* go run main.go

> 浏览 websocket 网页客户端

* 127.0.0.1:5000
* 输入手机号 (12345678901) 回车开始接受消息

> 推送个人

* POST 127.0.0.1:5000/push/user
* msg : 锄禾日当午
* user : 12345678901

> 推送群组

* POST 127.0.0.1:5000/push/group
* msg : 汗滴禾下土
* group : default_group

### Refine

* mongodb 做消息持久化
* redis 实现定时消息