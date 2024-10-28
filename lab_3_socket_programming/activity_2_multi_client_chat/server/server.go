package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

var clients = make(map[net.Conn]bool)
var mu sync.Mutex // mutex to ensure safe concurrent access to client list

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	defer listener.Close()
	fmt.Println("Chat server is listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		mu.Lock()
		clients[conn] = true
		mu.Unlock()

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer func() {
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected:", err)
			return
		}

		fmt.Printf("Broadcasting message: %s", message)
		broadcastMessage(message, conn)
	}
}

func broadcastMessage(message string, sender net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	for clientConn := range clients {
		if clientConn != sender {
			clientConn.Write([]byte(message))
		}
	}
}
