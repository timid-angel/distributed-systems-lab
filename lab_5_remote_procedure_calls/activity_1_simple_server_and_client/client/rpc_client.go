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

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatalln("Error connecting to RPC server:", err)
	}

	args := Args{A: 3, B: 5}
	var reply int

	err = client.Call("Calculator.Multiply", &args, &reply)
	if err != nil {
		log.Fatalln("Error calling RPC:", err)
	}

	fmt.Printf("Result of %d * %d = %d\n", args.A, args.B, reply)
}
