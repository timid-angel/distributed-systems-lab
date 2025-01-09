package paxos

import "sync"

type Acceptor struct {
	mu             sync.Mutex
	promisedNumber int
	acceptedNumber int
	acceptedValue  interface{}
	Is_working     bool
}

func (a *Acceptor) HandlePrepare(p Prepare) Promise {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.Is_working && p.ProposalNumber > a.promisedNumber {
		a.promisedNumber = p.ProposalNumber
		return Promise{ProposalNumber: p.ProposalNumber, AcceptedValue: a.acceptedValue}
	}

	return Promise{}
}

func (a *Acceptor) HandleAccept(ac Accept) Accepted {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.Is_working && ac.ProposalNumber >= a.promisedNumber {
		a.promisedNumber = ac.ProposalNumber
		a.acceptedNumber = ac.ProposalNumber
		a.acceptedValue = ac.Value
		return Accepted{ProposalNumber: ac.ProposalNumber, Value: ac.Value}
	}

	return Accepted{}
}
