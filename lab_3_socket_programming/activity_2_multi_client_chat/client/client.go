package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}

	defer conn.Close()

	go receiveMessages(conn)

	for {
		fmt.Printf("Send message to server: ")
		message, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		conn.Write([]byte(message))
	}
}

func receiveMessages(conn net.Conn) {
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("\033[92m\n\n\tMessage from server: \033[0m", message)
		fmt.Print("Send message to server: ")
	}
}
