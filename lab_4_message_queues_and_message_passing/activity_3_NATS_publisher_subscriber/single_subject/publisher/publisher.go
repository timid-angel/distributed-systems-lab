package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalln(err)
	}

	defer nc.Close()

	subject := "updates"
	message := "Hello, NATS!"
	if err := nc.Publish(subject, []byte(message)); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Sent:", message)
}
