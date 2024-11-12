package main

import (
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	defer conn.Close()

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	defer channel.Close()

	queue, err := channel.QueueDeclare("task_queue", true, false, false, false, nil)
	failOnError(err, "Failed to declare queue")

	body := "Hello, RabbitMQ!"
	err = channel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(body),
	})

	failOnError(err, "Failed to publish a message")
}
