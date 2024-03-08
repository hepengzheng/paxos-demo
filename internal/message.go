package internal

import "fmt"

// There are four types of messages:
// (1). Prepare message: proposer -> acceptor
// (2). Promise message: acceptor -> proposer
// (3). Accept message: proposer -> acceptor
// (4). Decide message: acceptor -> learner

type Message interface{}

type PrepareMessage struct {
	Number int64
}

func (p *PrepareMessage) String() string {
	return fmt.Sprintf("PrepareMessage(%d)", p.Number)
}

type PromiseMessage struct {
	Number          int64
	HighestAccepted *Proposal
}

func (p *PromiseMessage) String() string {
	if p.HighestAccepted == nil {
		return fmt.Sprintf("PromiseMessage(%d)", p.Number)
	}
	return fmt.Sprintf("PromiseMessage(%d, %s)", p.Number, p.HighestAccepted)
}

type AcceptMessage struct {
	Number int64
	Value  string
}

func (a *AcceptMessage) String() string {
	return fmt.Sprintf("AcceptMessage(%d, %s)", a.Number, a.Value)
}

type DecideMessage struct {
	Number int64
	Value  string
}

func (d *DecideMessage) String() string {
	return fmt.Sprintf("DecideMessage(%d, %s)", d.Number, d.Value)
}

var _ Message = (*PrepareMessage)(nil)
var _ Message = (*PromiseMessage)(nil)
var _ Message = (*AcceptMessage)(nil)
var _ Message = (*DecideMessage)(nil)

type Proposal struct {
	Number int64
	Value  string
}

func (p *Proposal) String() string {
	if p == nil {
		return "nil"
	}
	return fmt.Sprintf("Proposal(%d, %s)", p.Number, p.Value)
}
