package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
)

// arguments for multiplication
type Args struct {
	A, B int
}

type Calculator int

func (c *Calculator) Add(args *Args, reply *int) error {
	*reply = args.A + args.B
	return nil
}

func (c *Calculator) Subtract(args *Args, reply *int) error {
	*reply = args.A - args.B
	return nil
}

func (c *Calculator) Multiply(args *Args, reply *int) error {
	if args.A == 0 || args.B == 0 {
		return errors.New("multiplication by zero is not allowed")
	}

	*reply = args.A * args.B
	return nil
}

func (c *Calculator) Divide(args *Args, reply *float32) error {
	if args.B == 0 {
		return errors.New("can not divide by zero")
	}

	*reply = float32(args.A) / float32(args.B)
	return nil
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
