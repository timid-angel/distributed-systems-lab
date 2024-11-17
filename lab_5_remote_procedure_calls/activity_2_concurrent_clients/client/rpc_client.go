package main

import (
	"fmt"
	"log"
	"net/rpc"
)

// arguments for multiplication
type Args struct {
	A, B int
}

func callRpc(client *rpc.Client, operation string, operationSign string, args Args) {
	var reply int
	err := client.Call(operation, &args, &reply)
	if err != nil {
		log.Fatalln("Error calling RPC:", err)
	}

	fmt.Printf("Result of %d %v %d = %d\n", args.A, operationSign, args.B, reply)
}

func callDivide(client *rpc.Client, operation string, operationSign string, args Args) {
	var reply float32
	err := client.Call(operation, &args, &reply)
	if err != nil {
		log.Fatalln("Error calling RPC:", err)
	}

	fmt.Printf("Result of %d %v %d = %f\n", args.A, operationSign, args.B, reply)
}

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatalln("Error connecting to RPC server:", err)
	}

	args := Args{A: 3, B: 5}
	callRpc(client, "Calculator.Add", "+", args)
	callRpc(client, "Calculator.Subtract", "-", args)
	callRpc(client, "Calculator.Multiply", "*", args)
	callDivide(client, "Calculator.Divide", "/", args)
}
