package main

import (
	"fmt"
	"time"
)

func producer(channel chan<- string) {
	for i := 1; i <= 5; i++ {
		channel <- fmt.Sprintf("Message %d", i)
		time.Sleep(2 * time.Second)
	}

	close(channel)
}

func consumer(channel <-chan string) {
	for message := range channel {
		fmt.Println("Received: ", message)
	}
}

func main() {
	channel := make(chan string, 3)
	go producer(channel)
	consumer(channel)
}
