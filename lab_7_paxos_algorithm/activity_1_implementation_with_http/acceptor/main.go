package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"paxos-lab/paxos"
)

type Handler struct {
	Acceptor *paxos.Acceptor
}

func NewHandler(acceptor *paxos.Acceptor) *Handler {
	return &Handler{Acceptor: acceptor}
}

func (h *Handler) handlePrepare(w http.ResponseWriter, r *http.Request) {
	log.Println("Received prepare request on port:", r.Host)
	var prepare paxos.Prepare
	err := json.NewDecoder(r.Body).Decode(&prepare)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	promise := h.Acceptor.HandlePrepare(prepare)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(promise)
}

func (h *Handler) handleAccept(w http.ResponseWriter, r *http.Request) {
	log.Println("Received accept request on port:", r.Host)
	var accept paxos.Accept
	err := json.NewDecoder(r.Body).Decode(&accept)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	accepted := h.Acceptor.HandleAccept(accept)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accepted)
}

func initServer(port string) {
	acceptor := &paxos.Acceptor{Is_working: true}
	handler := NewHandler(acceptor)
	http.HandleFunc("/prepare", handler.handlePrepare)
	http.HandleFunc("/accept", handler.handleAccept)
	log.Println("Server listening on port: " + port)
	http.ListenAndServe(":"+port, nil)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Port not provided as an argument")
	}
	port := os.Args[1]
	initServer(port)
}
