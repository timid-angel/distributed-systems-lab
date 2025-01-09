package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"paxos-lab/paxos"
	"sync"
)

var (
	acceptors = []*paxos.Acceptor{{}, {}, {}}
	mu        sync.Mutex
)

func proposeHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ProposalNumber int    `json:"proposal_number"`
		Value          string `json:"value"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err)
		return
	}

	proposer := paxos.Proposer{ProposalNumber: body.ProposalNumber, Value: body.Value}
	mu.Lock()
	value := proposer.Propose(body.Value, acceptors)
	mu.Unlock()

	if value != nil {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Consensus reached: %s\n", value)
	} else {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "Consensus not reached\n")
	}
}

func main() {
	http.HandleFunc("/propose", proposeHandler)
	port := ":8080"
	if len(os.Args) > 1 {
		port = ":" + os.Args[1]
	}

	log.Printf("Starting server on port %s\n", port)
	http.ListenAndServe(port, nil)
}
