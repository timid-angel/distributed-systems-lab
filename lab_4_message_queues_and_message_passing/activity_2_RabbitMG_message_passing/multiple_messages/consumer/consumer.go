package main

import (
	"fmt"
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
	failOnError(err, "Failed to register a consumer")

	msgs, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			fmt.Printf("Received a message: %s\n", d.Body)
		}
	}()

	log.Println("Waiting for messages")
	select {}
}
