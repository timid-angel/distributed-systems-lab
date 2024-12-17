package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
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
	var reply int
	call := client.Go("Calculator.GetLastResult", &args, &reply, nil)
	select {
	case <-call.Done:
		if call.Error != nil {
			log.Println("RPC error:", call.Error)
		} else {
			fmt.Printf("Stored result before method invocation: %d\n", reply)
		}
	case <-time.After(2 * time.Second):
		log.Println("RPC call timed out")
	}

	call = client.Go("Calculator.Subtract", &args, &reply, nil)
	select {
	case <-call.Done:
		if call.Error != nil {
			log.Println("RPC error:", call.Error)
		} else {
			fmt.Printf("Computed Result: %d\n", reply)
		}
	case <-time.After(2 * time.Second):
		log.Println("RPC call timed out")
	}

	call = client.Go("Calculator.GetLastResult", &args, &reply, nil)
	select {
	case <-call.Done:
		if call.Error != nil {
			log.Println("RPC error:", call.Error)
		} else {
			fmt.Printf("Stored result after method invocation: %d\n", reply)
		}
	case <-time.After(2 * time.Second):
		log.Println("RPC call timed out")
	}
}
