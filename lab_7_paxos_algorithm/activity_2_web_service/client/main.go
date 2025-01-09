package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func sendProposal(server string, proposalNumber int, value string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()
	body, _ := json.Marshal(map[string]interface{}{
		"proposal_number": proposalNumber,
		"value":           value,
	})

	resp, err := http.Post(fmt.Sprintf("http://%s/propose", server), "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error sending proposal to %s: %v\n", server, err)
		results <- ""
		return
	}
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)
	results <- string(responseBody)
}

func main() {
	servers := []string{"192.168.100.143:8081", "192.168.100.143:8082", "192.168.100.143:8083"}
	proposalNumber := 2
	value := "Distributed Systems"

	var wg sync.WaitGroup
	results := make(chan string, len(servers))

	for _, server := range servers {
		wg.Add(1)
		go sendProposal(server, proposalNumber, value, &wg, results)
	}

	wg.Wait()
	close(results)

	for res := range results {
		fmt.Print(res)
	}
}
