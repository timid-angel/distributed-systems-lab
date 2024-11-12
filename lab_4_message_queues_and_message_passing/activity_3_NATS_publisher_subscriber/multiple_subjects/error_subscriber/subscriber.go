package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	subject := "updates.error"
	_, err = nc.Subscribe(subject, func(msg *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(msg.Data))
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Subscribed to updates. Waiting for messages...")
	select {}
}
