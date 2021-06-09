package mq

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

func Receive(queue string, groups []string) (res []string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("receive:%v\n", err)
			return
		}
	}()
	// 获取参数
	host := "127.0.0.1"
	port := "5672"
	user := "rabbitmq"
	pass := "rabbitmq"
	conn, err := amqp.Dial("amqp://" + user + ":" + pass + "@" + host + ":" + port + "/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"msg_hub_direct", // name
		"direct",         // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// 声明自己的队列
	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 绑定自己的路由
	err = ch.QueueBind(
		q.Name,           // queue name
		queue,            // routing key
		"msg_hub_direct", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	// 绑定所属群组的路由
	for index := range groups {
		err = ch.QueueBind(
			q.Name,           // queue name
			groups[index],    // routing key
			"msg_hub_direct", // exchange
			false,
			nil)
		failOnError(err, "Failed to bind a queue")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	// 获取消息
	res = make([]string, 0)
	endtime := time.Now().Unix() + 2
	go func() {
		for {
			// 超时
			if time.Now().Unix() > endtime {
				break
			}
			d, ok := <-msgs
			if !ok {
				break
			}
			d.Ack(false)
			res = append(res, string(d.Body))
		}
	}()
	for {
		// 超时
		if time.Now().Unix() > endtime {
			panic("done")
		}
	}
}
