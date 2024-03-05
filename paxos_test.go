package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunPaxos(t *testing.T) {
	for i := 0; i < 10; i++ {
		assert.Nil(t, RunPaxos())
	}
}
