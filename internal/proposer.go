package internal

import (
	"log"
	"math/rand"
	"time"
)

type Proposer struct {
	Role_   Role
	NodeID_ string

	quorum     int
	messenger  Messenger
	proposerID ProposalID

	value                   string
	number                  int64
	promisedAcceptors       []string // node ID list
	highestAcceptedProposal *Proposal
}

func (p *Proposer) Role() Role {
	return p.Role_
}

func (p *Proposer) NodeID() string {
	return p.NodeID_
}

var _ Node = (*Proposer)(nil)

func NewProposer(
	nodeID string,
	quorum int,
	messenger Messenger,
	proposerID ProposalID,
) *Proposer {
	return &Proposer{
		NodeID_:    nodeID,
		Role_:      RoleProposer,
		quorum:     quorum,
		messenger:  messenger,
		proposerID: proposerID,
	}
}

func (p *Proposer) Propose(value string, timeoutMilli int64) {
	p.value = value

	ticker := time.NewTicker(time.Duration(timeoutMilli) * time.Millisecond)
	defer ticker.Stop()
	done := p.messenger.QueryFinished()
	for {
		select {
		case <-ticker.C:
			r := rand.Intn(100) + 100
			<-time.After(time.Duration(r) * time.Millisecond)
			log.Printf("%s start new proposal", p.NodeID_)
			p.doPropose()
		case <-done:
			log.Printf("%s finished", p.NodeID_)
			return
		}
	}
}

func (p *Proposer) Receive(from string, message Message) {

	switch msg := message.(type) {
	case PromiseMessage:
		log.Printf("%s received promise message %v from %s\n", p.NodeID_, msg, from)
		if p.number != msg.Number {
			return
		}

		p.promisedAcceptors = append(p.promisedAcceptors, from)
		if msg.HighestAccepted != nil {
			if p.highestAcceptedProposal == nil {
				p.highestAcceptedProposal = msg.HighestAccepted
			} else if p.highestAcceptedProposal.Number < msg.HighestAccepted.Number {
				p.highestAcceptedProposal = msg.HighestAccepted
			}
		}

		if len(p.promisedAcceptors) == p.quorum {
			valueToAccept := p.value
			if p.highestAcceptedProposal != nil {
				valueToAccept = p.highestAcceptedProposal.Value
			}
			am := AcceptMessage{Number: p.number, Value: valueToAccept}
			p.messenger.BroadCastToAcceptors(p.NodeID_, p.promisedAcceptors, am)
		}

	default:

	}

}

func (p *Proposer) doPropose() {
	p.number = 0
	p.promisedAcceptors = make([]string, 0)
	p.highestAcceptedProposal = nil

	number_ := p.proposerID.Next(p.number)
	p.number = number_
	prepareMsg := PrepareMessage{Number: number_}
	p.messenger.BroadCastToAllAcceptors(p.NodeID_, prepareMsg)
}
