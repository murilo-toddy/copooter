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

// TODO: abstract rendering part
type Node struct {
	ID          string
	State       NodeState
	connections []*Node
}

// TODO: determine ID automatically
func NewNode(id string) *Node {
	return &Node{
		ID:          id,
		State:       Undefined,
		connections: []*Node{},
	}
}

func (n *Node) Change(newState NodeState) error {
	if n.State != Undefined && n.State != newState {
		return fmt.Errorf("conflicting values for node %s", n.ID)
	}
	n.State = newState
	for _, node := range n.connections {
		if node.State == newState {
			continue
		}
		node.Change(newState)
	}
	return nil
}

func (n *Node) Debug() string {
	return fmt.Sprintf("%s=<state: %s>", n.ID, n.State)
}

func (n *Node) Connect(n1 *Node) *Node {
	if n1 != nil {
		n.connections = append(n.connections, n1)
		n1.connections = append(n1.connections, n)
	}
	return n
}

var SharedSourceNode = NewNode("SharedSource")
var SharedGroundNode = NewNode("SharedGround")
