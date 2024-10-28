package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}

	defer conn.Close()

	for {
		task, _ := bufio.NewReader(conn).ReadString('\n')
		task = strings.TrimSpace(task)

		num, _ := strconv.Atoi(task)
		fmt.Printf("\033[92m Computing task with value: %d \033[0m\n", num)
		result := num * num

		conn.Write([]byte(fmt.Sprintf("%d\n", result)))
	}
}
