package main

import "fmt"

type NodeState int

const (
	Off = iota
	On
	Undefined
)

func (n NodeState) String() string {
	switch n {
	case Off:
		return "off"
	case On:
		return "on"
	case Undefined:
		return "undefined"
	default:
		return "unknown"
	}
}

type Node struct {
	ID    string
	State NodeState
}

// TODO: determine ID automatically
func NewNode(id string) *Node {
	return &Node{
		ID:    id,
		State: Undefined,
	}
}

func (n *Node) Change(newState NodeState) error {
	if n.State != Undefined && n.State != newState {
		return fmt.Errorf("conflicting values for node %s", n.ID)
	}
	n.State = newState
	return nil
}

func (n *Node) Debug() string {
	return fmt.Sprintf("%s=<state: %s>", n.ID, n.State)
}

var SharedSourceNode = NewNode("SharedSource")
var SharedGroundNode = NewNode("SharedGround")
