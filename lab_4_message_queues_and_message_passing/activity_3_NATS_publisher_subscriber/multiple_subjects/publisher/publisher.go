package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func sendMessage(subject string, message string, nc *nats.Conn) {
	if err := nc.Publish(subject, []byte(message)); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Sent: '%v' on subject: '%v'\n", message, subject)
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalln(err)
	}

	defer nc.Close()

	sendMessage("updates.info", "This is an informative message", nc)
	sendMessage("updates.error", "This is an error message", nc)
}
