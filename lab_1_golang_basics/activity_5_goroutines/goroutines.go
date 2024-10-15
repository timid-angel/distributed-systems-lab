package main

import (
	"fmt"
	"time"
)

func printNumbers() {
	for i := 0; i <= 5; i++ {
		fmt.Println(i)
		time.Sleep(1 * time.Second)
	}
}

func printLetters() {
	for i := 'A'; i <= 'E'; i++ {
		fmt.Printf("%c\n", i)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	// run printNumbers concurrently
	go printNumbers()

	// run printLetters concurrently
	go printLetters()

	// sleep main function to allow goroutines to finish
	time.Sleep(6 * time.Second)
	fmt.Println("Main function finished")
}
