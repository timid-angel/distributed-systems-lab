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
		fmt.Println(err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Send message to server: ")
	data, _ := reader.ReadString('\n')

	_, err = conn.Write([]byte(data + "\n"))
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("Received from server: ", string(res))
	conn.Close()
}
