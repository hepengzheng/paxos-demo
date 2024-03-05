package internal

import (
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
	nodes     map[string]Node // all nodes
	acceptors map[string]Node // only acceptors
	learners  map[string]Node // only learners
}

var _ Messenger = (*MessengerImpl)(nil)

func NewMessenger() Messenger {
	return &MessengerImpl{
		nodes:     make(map[string]Node),
		acceptors: make(map[string]Node),
		learners:  make(map[string]Node),
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
			timeoutRequest()
			acceptor.Receive(from, message)
		}()
	}
}

func (m *MessengerImpl) BroadCastToAcceptors(from string, acceptors []string, message Message) {
	for _, acceptorID := range acceptors {
		if acceptor, ok := m.acceptors[acceptorID]; ok {
			go func() {
				<-time.After(time.Duration(rand.Intn(30)+10) * time.Millisecond)
				acceptor.Receive(from, message)
			}()
		}
	}
}

func (m *MessengerImpl) BroadCastToLearners(from string, message DecideMessage) {
	for _, l := range m.learners {
		learner := l
		go func() {
			<-time.After(time.Duration(rand.Intn(30)+10) * time.Millisecond)
			learner.Receive(from, message)
		}()
	}
}

func (m *MessengerImpl) UniCast(from, to string, message Message) {
	if node, ok := m.nodes[to]; ok {
		go func() {
			timeoutRequest()
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

func timeoutRequest() time.Time {
	return <-time.After(time.Duration(rand.Intn(30)+10) * time.Millisecond)
}
