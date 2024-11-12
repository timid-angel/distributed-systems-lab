package main

import (
	"fmt"
	"time"
)

func producer(channel chan<- string, quit <-chan bool) {
	for i := 1; ; i++ {
		select {
		case <-quit:
			fmt.Println("Producer shutting down")

			return
		case channel <- fmt.Sprintf("Message %d", i):
			fmt.Printf("Produced: Message %d\n", i)
			time.Sleep(time.Second)
		}
	}

}

func consumer(channel <-chan string, quit chan<- bool) {
	for i := 0; i <= 10; i++ {
		fmt.Println("Consumed: ", <-channel)
	}

	quit <- true
}

func main() {
	channel := make(chan string)
	quit := make(chan bool)

	go producer(channel, quit)
	go consumer(channel, quit)

	<-quit
	fmt.Println("Main shutting down")
}
