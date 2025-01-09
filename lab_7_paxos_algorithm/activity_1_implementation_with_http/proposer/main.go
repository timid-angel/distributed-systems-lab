package main

import (
	"fmt"
	"paxos-lab/paxos"
)

func main() {
	acceptorURLs := []string{
		"http://192.168.100.143:8081",
		"http://192.168.100.143:8082",
		"http://192.168.100.143:8083",
		"http://192.168.100.143:8084",
		"http://192.168.100.143:8085",
	}

	proposer := paxos.Proposer{
		ProposalNumber: 3,
		Value:          "Distributed Systems",
		AcceptorURLs:   acceptorURLs,
	}

	consensusValue, err := proposer.Propose(proposer.Value)
	if err != nil {
		fmt.Println("Consensus not reached:", err)
	} else {
		fmt.Printf("Consensus reached on value: %s\n", consensusValue)
	}
}
