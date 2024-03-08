package internal

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type Messenger interface {
	AddNodes(nodes ...Node)
	BroadCastToAllAcceptors(from string, message Message)
	BroadCastToAcceptors(from string, acceptors []string, message Message)
	BroadCastToLearners(from string, message DecideMessage)
	UniCast(from, to string, message Message)
	QueryFinished() <-chan struct{}
}

type MessengerImpl struct {
	sync.Mutex
	nodes                    map[string]Node // all nodes
	acceptors                map[string]Node // only acceptors
	learners                 map[string]Node // only learners
	possibilityToDropMessage float32
}

var _ Messenger = (*MessengerImpl)(nil)

func NewMessenger(possibilityToDropMessage float32) Messenger {
	if possibilityToDropMessage > 1 || possibilityToDropMessage < 0 {
		log.Fatalf("invalid possibilityToDropMessage: %f\n", possibilityToDropMessage)
	}
	return &MessengerImpl{
		nodes:                    make(map[string]Node),
		acceptors:                make(map[string]Node),
		learners:                 make(map[string]Node),
		possibilityToDropMessage: possibilityToDropMessage,
	}
}

func (m *MessengerImpl) AddNodes(nodes ...Node) {
	m.Lock()
	defer m.Unlock()

	for _, n := range nodes {
		node := n // compatible for Go versions < 1.22
		m.nodes[node.NodeID()] = node
		if node.Role() == RoleAcceptor {
			m.acceptors[node.NodeID()] = node
		}
		if node.Role() == RoleLearner {
			m.learners[node.NodeID()] = node
		}
	}
}

func (m *MessengerImpl) BroadCastToAllAcceptors(from string, message Message) {
	for _, a := range m.acceptors {
		acceptor := a
		go func() {
			if m.mapDropOrDelayMessage() {
				log.Printf("message %v from %s to acceptor %s was dropped\n", message, from, acceptor.NodeID())
				return
			}
			acceptor.Receive(from, message)
		}()
	}
}

func (m *MessengerImpl) BroadCastToAcceptors(from string, acceptors []string, message Message) {
	for _, acceptorID := range acceptors {
		if acceptor, ok := m.acceptors[acceptorID]; ok {
			go func() {
				if m.mapDropOrDelayMessage() {
					log.Printf("message %v from %s to acceptor %s was dropped\n", message, from, acceptor.NodeID())
				}
				acceptor.Receive(from, message)
			}()
		}
	}
}

func (m *MessengerImpl) BroadCastToLearners(from string, message DecideMessage) {
	for _, l := range m.learners {
		learner := l
		go func() {
			if m.mapDropOrDelayMessage() {
				log.Printf("message %v from %s to %s was dropped\n", message, from, learner.NodeID())
				return
			}
			learner.Receive(from, message)
		}()
	}
}

func (m *MessengerImpl) UniCast(from, to string, message Message) {
	if node, ok := m.nodes[to]; ok {
		go func() {
			if m.mapDropOrDelayMessage() {
				log.Printf("message %v from %s to %s was dropped\n", message, from, to)
				return
			}
			node.Receive(from, message)
		}()
	}
}

func (m *MessengerImpl) QueryFinished() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		for _, l := range m.learners {
			learner := l.(*Learner)
			<-learner.done
		}
	}()
	return ch
}

func (m *MessengerImpl) mapDropOrDelayMessage() bool {
	r := rand.Intn(100)
	<-time.After(time.Duration(r) * time.Millisecond)
	return m.possibilityToDropMessage*100 > float32(r)
}
