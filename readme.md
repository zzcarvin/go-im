## go-im

> go 消息中心模型

### Run

> Docker 运行 rabbitmq 服务

* cd docker/rabbitmq && docker-compose up

> 运行 websocket 服务

* go run main.go

> 浏览 websocket 网页客户端

* 127.0.0.1:5000

> 登录 rabbitmq 管理界面向 123456 队列 publish 消息

* 127.0.0.1:15672
* rabbitmq : rabbitmq

### Undo

* mongodb 做消息持久化
* redis 实现定时消息