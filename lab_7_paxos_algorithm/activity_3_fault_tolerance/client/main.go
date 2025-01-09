package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

func sendRequest(server string, proposalNumber int, value string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, _ := json.Marshal(map[string]interface{}{
		"proposal_number": proposalNumber,
		"value":           value,
	})

	req, err := http.NewRequestWithContext(context, "POST", fmt.Sprintf("http://%s/propose", server), bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error creating request to %s: %v\n", server, err)
		results <- ""
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending proposal to %s: %v\n", server, err)
		results <- ""
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error response from %s: %d\n", server, resp.StatusCode)
		results <- ""
		return
	}

	log.Printf("Successfully sent proposal to %s\n", server)
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)
	results <- string(responseBody)
}

func sendProposal(server string, proposalNumber int, value string, retries int, results chan<- string) {
	for i := 0; i < retries; i++ {
		var wg sync.WaitGroup
		wg.Add(1)
		retryRes := make(chan string, 1)

		go sendRequest(server, proposalNumber, value, &wg, retryRes)
		wg.Wait()
		close(retryRes)

		response := <-retryRes
		if response != "" {
			results <- response
			return
		}

		log.Printf("Retry %d for server %s failed\n", i+1, server)
		time.Sleep(2 * time.Second)
	}

	log.Printf("Failed to communicate with server %s after %d retries\n", server, retries)
	results <- ""
}

func main() {
	servers := []string{"192.168.100.143:8081", "192.168.100.143:8082", "192.168.100.143:8083"}
	proposalNumber := 4
	value := "Distributed Systems"
	retries := 2

	var wg sync.WaitGroup
	results := make(chan string, len(servers))

	for _, server := range servers {
		wg.Add(1)
		go func(server string) {
			defer wg.Done()
			sendProposal(server, proposalNumber, value, retries, results)
		}(server)
	}

	wg.Wait()
	close(results)
	for res := range results {
		fmt.Print(res)
	}
}
