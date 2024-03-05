package internal

import (
	"log"
	"sync"
)

type Acceptor struct {
	Role_     Role
	NodeID_   string
	messenger Messenger

	sync.Mutex
	highestPromised *Proposal
	highestAccepted *Proposal
}

func (a *Acceptor) Role() Role {
	return a.Role_
}

func (a *Acceptor) NodeID() string {
	return a.NodeID_
}

var _ Node = (*Acceptor)(nil)

func NewAcceptor(
	nodeID string,
	messenger Messenger,
) *Acceptor {
	return &Acceptor{
		Role_:     RoleAcceptor,
		NodeID_:   nodeID,
		messenger: messenger,
	}
}

func (a *Acceptor) Receive(from string, message Message) {
	switch msg := message.(type) {
	case PrepareMessage:
		a.Lock()
		defer a.Unlock()
		log.Printf("%s receivied prepare message %#v from %s\n", a.NodeID_, msg, from)
		defer func() {
			log.Printf("%s current state: highestPromised: %#v, highestAccepted: %#v\n", a.NodeID_, a.highestPromised, a.highestAccepted)
		}()
		if a.highestPromised == nil || a.highestPromised.Number < msg.Number {
			promise := PromiseMessage{
				Number: msg.Number,
			}
			if a.highestAccepted != nil {
				promise.HighestAccepted = a.highestAccepted
			}
			a.highestPromised = &Proposal{Number: msg.Number}
			a.messenger.UniCast(a.NodeID_, from, promise)
		}

		return

	case AcceptMessage:
		a.Lock()
		defer a.Unlock()
		log.Printf("%s receivied accept message %#v from %s\n", a.NodeID_, msg, from)
		defer func() {
			log.Printf("%s current state: highestPromised: %#v, highestAccepted: %#v\n", a.NodeID_, a.highestPromised, a.highestAccepted)
		}()
		if a.highestPromised == nil || a.highestPromised.Number <= msg.Number {
			a.highestAccepted = &Proposal{Number: msg.Number, Value: msg.Value}
			decide := DecideMessage(msg)
			a.messenger.BroadCastToLearners(a.NodeID_, decide)
		}

		return
	default:
	}
}
