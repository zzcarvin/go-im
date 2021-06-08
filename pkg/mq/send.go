package mq

import (
	"log"

	"github.com/streadway/amqp"
)

func Send(route string, body string) {
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

	err = ch.Publish(
		"msg_hub_direct", // exchange
		route,            // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(body),
			DeliveryMode: 2,
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}
