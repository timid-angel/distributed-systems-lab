package main

import (
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Replica struct {
	value float64
	mu    sync.Mutex
	peers []string
}

func (r *Replica) Update(newValue, delta float64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if math.Abs(newValue-r.value) <= delta {
		fmt.Printf("> [SUCCESS] Replica value updated from %.2f to %.2f\n", r.value, newValue)
		r.value = newValue
		return true
	}

	fmt.Printf("> [FAILURE] Failed to update replica value: difference exceeds delta value of %.2f\n", delta)
	return false
}

func (r *Replica) propagateUpdates() {
	for _, peer := range r.peers {
		go func(peer string) {
			conn, err := net.Dial("tcp", peer)
			if err != nil {
				fmt.Println("Error connecting to peer:", peer, err)
				return
			}

			defer conn.Close()

			r.mu.Lock()
			message := fmt.Sprintf("%.2f\n", r.value)
			r.mu.Unlock()

			conn.Write([]byte(message))
		}(peer)
	}
}

func handleConnection(conn net.Conn, replica *Replica, delta float64) {
	defer conn.Close()
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			break
		}

		newValue := strings.TrimSpace(string(buffer[:n]))
		var value float64
		fmt.Sscanf(newValue, "%f", &value)
		replica.Update(value, delta)
	}
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run replica_numerical.go <delta> <machine_ip:port> <peer1_ip:port> [<peer2_ip:port>...]")
		return
	}

	delta, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		panic(err)
	}

	machineAddress := os.Args[2]
	peers := os.Args[3:]
	replica := &Replica{
		value: rand.Float64() * 20,
		peers: peers,
	}

	listener, err := net.Listen("tcp", machineAddress)
	if err != nil {
		panic(err)
	}

	defer listener.Close()
	fmt.Println("Replica listening on " + machineAddress)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}

			go handleConnection(conn, replica, delta)
		}
	}()

	time.Sleep(2 * time.Second) // wait for replicas to start up
	replica.value = rand.Float64() * 20
	replica.propagateUpdates()
	fmt.Printf("Replica Value: %.2f\n", replica.value)
	time.Sleep(3 * time.Second) // wait for replicas to finish propagating
}
