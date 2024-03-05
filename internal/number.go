package internal

import "time"

type ProposalID interface {
	Next(n int64) int64
}

type TimestampProposalID struct{}

func (t *TimestampProposalID) Next(n int64) int64 {
	return time.Now().UnixMilli()
}

var _ ProposalID = (*TimestampProposalID)(nil)

func NewProposalID() ProposalID {
	return &TimestampProposalID{}
}
