package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Replica struct {
	data  map[string]string
	mu    sync.Mutex
	peers []string // peer addresses
}

func (r *Replica) Update(key, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[key] = value
}

func (r *Replica) propagateUpdates(key, value string) {
	time.Sleep(4 * time.Second) // wait for the other replicas to start
	for _, peer := range r.peers {
		go func(peer string) {
			conn, err := net.Dial("tcp", peer)
			if err != nil {
				fmt.Println("Error connecting to peer: ", peer, err)
				return
			}

			defer conn.Close()
			sentTime := time.Now().Unix()
			time.Sleep(time.Duration(rand.Intn(4)) * time.Second) // add random delay
			fmt.Fprintf(conn, "%s:%s:%v\n", key, value, sentTime)
			fmt.Println("\t> Sent propagation message to " + peer)
		}(peer)
	}
}

func handleConnection(conn net.Conn, replica *Replica) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		parts := strings.Split(strings.TrimSpace(message), ":")
		if len(parts) == 3 {
			// get the time stamp from when the propagation request had been made
			unixTime, _ := strconv.ParseInt(parts[2], 10, 64)
			sentTimestamp := time.Unix(unixTime, 0)
			interval := time.Since(sentTimestamp)

			fmt.Printf("\t> Received update propagation message: %s:%s\n\t  Time taken for consistency: %v seconds\n", parts[0], parts[1], interval.Seconds())
			replica.Update(parts[0], parts[1])
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run replica_eventual.go <machine_ip:port> <peer1_ip:port> [<peer2_ip:port>...]")
		return
	}

	machineAddr := os.Args[1]
	peers := os.Args[2:]

	replica := &Replica{
		data:  make(map[string]string),
		peers: peers,
	}

	listener, err := net.Listen("tcp", machineAddr)
	if err != nil {
		panic(err)
	}

	defer listener.Close()
	fmt.Printf("Replica listening on %s\n", machineAddr)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}

			go handleConnection(conn, replica)
		}
	}()

	replica.Update("key1", "value1")
	replica.propagateUpdates("key1", "value1")

	time.Sleep(15 * time.Second)
	replica.mu.Lock()
	fmt.Println("Replica Data:", replica.data)
	replica.mu.Unlock()
}
