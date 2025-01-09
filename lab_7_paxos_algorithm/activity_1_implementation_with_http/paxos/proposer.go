package paxos

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Proposer struct {
	ProposalNumber int
	Value          interface{}
	AcceptorURLs   []string
}

func (p *Proposer) sendPrepare(url string, prepare Prepare) (Promise, error) {
	data, _ := json.Marshal(prepare)
	resp, err := http.Post(url+"/prepare", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return Promise{}, err
	}

	defer resp.Body.Close()
	var promise Promise
	err = json.NewDecoder(resp.Body).Decode(&promise)
	return promise, err
}

func (p *Proposer) sendAccept(url string, accept Accept) (Accepted, error) {
	data, _ := json.Marshal(accept)
	resp, err := http.Post(url+"/accept", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return Accepted{}, err
	}

	defer resp.Body.Close()
	var accepted Accepted
	err = json.NewDecoder(resp.Body).Decode(&accepted)
	return accepted, err
}

func (p *Proposer) Propose(value interface{}) (interface{}, error) {
	promises := 0
	for _, url := range p.AcceptorURLs {
		promise, err := p.sendPrepare(url, Prepare{ProposalNumber: p.ProposalNumber})
		if err == nil && promise.ProposalNumber == p.ProposalNumber {
			promises++
		}
	}

	if promises > len(p.AcceptorURLs)/2 {
		accepted := 0
		for _, url := range p.AcceptorURLs {
			ack, err := p.sendAccept(url, Accept{ProposalNumber: p.ProposalNumber, Value: value})
			if err == nil && ack.ProposalNumber == p.ProposalNumber {
				accepted++
			}

			if accepted > len(p.AcceptorURLs)/2 {
				return value, nil
			}
		}
	}

	return nil, errors.New("consensus not reached")
}
