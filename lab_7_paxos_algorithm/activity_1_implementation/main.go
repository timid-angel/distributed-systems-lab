package main

import (
	"fmt"
	"paxos-lab/paxos"
)

func main() {
	acceptors := []*paxos.Acceptor{
		&paxos.Acceptor{Is_working: true},
		&paxos.Acceptor{Is_working: true},
		&paxos.Acceptor{Is_working: false},
		&paxos.Acceptor{Is_working: true},
		&paxos.Acceptor{Is_working: false},
	}

	proposer := paxos.Proposer{ProposalNumber: 1, Value: "Distributed Systems"}
	value := proposer.Propose("Distributed Systems", acceptors)

	if value != nil {
		fmt.Printf("Consensus reached on value: %s\n", value)
	} else {
		fmt.Println("Consensus not reached")
	}
}
