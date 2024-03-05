package internal

type Node interface {
	Role() Role
	NodeID() string
	Receive(from string, message Message)
}
