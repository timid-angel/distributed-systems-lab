package main

import "fmt"

func sendData(ch chan int) {
	for i := 0; i < 5; i++ {
		ch <- i
	}

	close(ch)
}

func main() {
	ch := make(chan int)

	go sendData(ch)

	for val := range ch {
		fmt.Println("Received:", val)
	}

	fmt.Println("Channel closed, program finished.")
}
