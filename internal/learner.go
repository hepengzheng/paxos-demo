package internal

import (
	"log"
	"sync"
)

type Learner struct {
	Role_     Role
	NodeID_   string
	quorum    int
	messenger Messenger

	sync.Mutex
	finished  bool
	Result    *Proposal
	done      chan struct{}
	acceptors map[string]Proposal           // nodeID -> proposal
	accepted  map[int64]map[string]struct{} // number -> [nodeID]
}

var _ Node = (*Learner)(nil)

func NewLearner(
	nodeID string,
	quorum int,
	messenger Messenger,
) *Learner {
	return &Learner{
		Role_:     RoleLearner,
		NodeID_:   nodeID,
		quorum:    quorum,
		messenger: messenger,
		done:      make(chan struct{}),
		acceptors: make(map[string]Proposal),
		accepted:  make(map[int64]map[string]struct{}),
	}
}

func (l *Learner) Role() Role {
	return l.Role_
}

func (l *Learner) NodeID() string {
	return l.NodeID_
}

func (l *Learner) Receive(from string, message Message) {
	switch msg := message.(type) {
	case DecideMessage:
		l.Lock()
		defer l.Unlock()
		log.Printf("%s received decide message %#v from %s\n", l.NodeID_, msg, from)
		if l.finished {
			return
		}
		if _, ok := l.acceptors[from]; !ok {
			l.acceptors[from] = Proposal(msg) // add to acceptors
			number_ := msg.Number
			nodes, ok_ := l.accepted[number_]
			if !ok_ {
				nodes = make(map[string]struct{})
			}
			nodes[from] = struct{}{}
			l.accepted[number_] = nodes

			log.Printf("In if %s has %d acceptors accepted number %d, current acceptors: %v\n", l.NodeID_, len(l.accepted[number_]), number_, l.acceptors)
			if len(l.accepted[number_]) == l.quorum {
				l.finished = true
				l.Result = &Proposal{
					Number: number_,
					Value:  msg.Value,
				}
				log.Printf("%s finished with result %#v", l.NodeID_, l.Result)
				close(l.done) // notify all observers
			}
			return
		}

		// acceptor already in acceptors
		earlyAccepted := l.acceptors[from]
		number_ := msg.Number
		if earlyAccepted.Number < number_ {
			l.acceptors[from] = Proposal(msg)
			delete(l.accepted, earlyAccepted.Number)

			nodes, ok := l.accepted[number_]
			if !ok {
				nodes = make(map[string]struct{})
			}
			nodes[from] = struct{}{}
			l.accepted[number_] = nodes

			log.Printf("Out of if %s' acceptors accepted number %d: %v, current acceptors: %v\n", l.NodeID_, number_, l.accepted[number_], l.acceptors)
			if len(l.accepted[number_]) == l.quorum {
				l.finished = true
				l.Result = &Proposal{
					Number: number_,
					Value:  msg.Value,
				}
				log.Printf("%s finished with result %#v", l.NodeID_, l.Result)
				close(l.done)
			}
			return
		}
	default:

	}
}
