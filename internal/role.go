package internal

type Role int

const (
	RoleProposer Role = iota
	RoleAcceptor
	RoleLearner
)
