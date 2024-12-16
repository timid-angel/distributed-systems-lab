package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"sync"
	"time"
)

type Replica struct {
	data    map[string]string
	mu      sync.Mutex
	peers   []string
	ackLock sync.Mutex
	acks    map[string]int // track acknowledgements
}

type Args struct {
	Key    string
	Value  string
	Source string
}

func (r *Replica) Update(args *Args, reply *bool) error {
	fmt.Println("> [UPDATE] Update request committed with args: " + args.Key + "<->" + args.Value + " sourced from: " + args.Source)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[args.Key] = args.Value
	if reply != nil {
		*reply = true
	}
	return nil
}

func (r *Replica) propagateUpdates(key, value, machineAddress string) {
	time.Sleep(4 * time.Second) // wait for the other replicas to start
	r.ackLock.Lock()
	r.acks[key] = 0
	r.ackLock.Unlock()

	for _, peer := range r.peers {
		go func(peer string) {
			client, err := rpc.Dial("tcp", peer)
			if err != nil {
				fmt.Println("Error connecting to peer:", peer, err)
				return
			}

			defer client.Close()
			args := &Args{Key: key, Value: value, Source: machineAddress}
			var reply bool = false
			err = client.Call("Replica.Update", args, &reply)
			fmt.Println("> [PROPAGATION] Update request propagated to peer " + peer)
			if err == nil && reply {
				fmt.Println("> [ACK] Peer " + peer + " acknowledged update request")
				r.ackLock.Lock()
				r.acks[key]++
				r.ackLock.Unlock()
			}
		}(peer)
	}
}

func (r *Replica) waitForAcknowledgements(key string, quorum int) {
	for {
		r.ackLock.Lock()
		if r.acks[key] >= quorum {
			r.ackLock.Unlock()
			break
		}

		r.ackLock.Unlock()
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run replica_strong.go <machine_ip:port> <peer1_ip:port> [<peer2_ip:port>...]")
		return
	}

	machineAddr := os.Args[1]
	peers := os.Args[2:]

	replica := &Replica{
		data:  make(map[string]string),
		peers: peers,
		acks:  make(map[string]int),
	}

	rpc.Register(replica)

	listener, err := net.Listen("tcp", machineAddr)
	if err != nil {
		panic(err)
	}

	defer listener.Close()
	fmt.Printf("Replica RPC server listening on %s\n", machineAddr)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}

			go rpc.ServeConn(conn)
		}
	}()

	key, value := "key1", "value1"
	fmt.Println("> [INIT] Update initialized locally")
	replica.Update(&Args{Key: key, Value: value, Source: "local"}, nil)
	replica.propagateUpdates(key, value, machineAddr)

	var Q int = (len(replica.peers) / 2)
	replica.waitForAcknowledgements(key, Q)
	fmt.Println("> [COMMITED] Update committed after receiving sufficient acknowledgments")
	time.Sleep(5 * time.Second) // wait for machine to acknowledge update requests from peers
}
