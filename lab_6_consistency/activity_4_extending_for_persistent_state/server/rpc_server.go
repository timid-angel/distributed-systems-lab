package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

// arguments for multiplication
type Args struct {
	A, B int
}

type Calculator struct {
	lastResult int
	mu         sync.Mutex
}

func (c *Calculator) Add(args *Args, reply *int) error {
	*reply = args.A + args.B
	c.SetLastResult(*reply)
	return nil
}

func (c *Calculator) Subtract(args *Args, reply *int) error {
	*reply = args.A - args.B
	c.SetLastResult(*reply)
	return nil
}

func (c *Calculator) Multiply(args *Args, reply *int) error {
	if args.A == 0 || args.B == 0 {
		return errors.New("multiplication by zero is not allowed")
	}

	*reply = args.A * args.B
	c.SetLastResult(*reply)
	return nil
}

func (c *Calculator) Divide(args *Args, reply *int) error {
	if args.B == 0 {
		return errors.New("can not divide by zero")
	}

	*reply = args.A / args.B
	c.SetLastResult(*reply)
	return nil
}

func (c *Calculator) GetLastResult(args *struct{}, reply *int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	*reply = c.lastResult
	return nil
}

func (c *Calculator) SetLastResult(oldResult int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastResult = oldResult
}

func main() {
	calc := new(Calculator)
	rpc.Register(calc)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error starting RPC server:", err)
		return
	}

	fmt.Println("RPC server is listening on port 1234...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go rpc.ServeConn(conn)
	}
}
