# Paxos demo

As the name suggests, this project is a demonstration of the Basic Paxos algorithm.

## Demo

```go
func main() {

	possibilityToDropMsg := float32(0.3)
	messenger := internal.NewMessenger(possibilityToDropMsg)
	proposalID := internal.NewProposalID()
	quorum := 3

	// three nodes as proposers
	proposer1 := internal.NewProposer("proposer-1", quorum, messenger, proposalID)
	proposer2 := internal.NewProposer("proposer-2", quorum, messenger, proposalID)
	proposer3 := internal.NewProposer("proposer-3", quorum, messenger, proposalID)

	// five nodes as acceptors
	acceptor1 := internal.NewAcceptor("acceptor-1", messenger)
	acceptor2 := internal.NewAcceptor("acceptor-2", messenger)
	acceptor3 := internal.NewAcceptor("acceptor-3", messenger)
	acceptor4 := internal.NewAcceptor("acceptor-4", messenger)
	acceptor5 := internal.NewAcceptor("acceptor-5", messenger)

	// three nodes as learners
	learner1 := internal.NewLearner("learner-1", quorum, messenger)
	learner2 := internal.NewLearner("learner-2", quorum, messenger)
	learner3 := internal.NewLearner("learner-3", quorum, messenger)

	learners := []*internal.Learner{learner1, learner2, learner3}

	messenger.AddNodes(
		proposer1, proposer2, proposer3,
		acceptor1, acceptor2, acceptor3, acceptor4, acceptor5,
		learner1, learner2, learner3,
	)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		proposer1.Propose("cat", 2000)
	}()

	go func() {
		defer wg.Done()
		proposer2.Propose("dog", 2000)
	}()

	go func() {
		defer wg.Done()
		proposer3.Propose("fish", 2000)
	}()

	wg.Wait()

	// Collect the results
	results := make([]internal.Proposal, 0, len(learners))
	numbers := make(map[int64]string)
	for _, learner := range learners {
		results = append(results, *learner.Result)
		numbers[learner.Result.Number] = learner.Result.Value
	}

	log.Printf("Results: %v\n", results)
}
```


# Run the program

```shell
go run main.go
```