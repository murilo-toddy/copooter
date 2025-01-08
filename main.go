package main

import (
	"fmt"
	"os"
)

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

type ComponentType int

type Completed bool

type Component interface {
	// component actions can be deferred if they are in a state where they
	// do not have sufficient data to perform an operation.
	// the force flag removes the possibility of postponing an action and the
	// function instead returns an error if an undefined state is found
	Act(force bool) (Completed, error)
	Debug() string
}

type Source struct {
	Node *Node
}

func NewSource(node *Node) *Source {
	return &Source{
		Node: node,
	}
}

func (s *Source) Act(force bool) (Completed, error) {
	return true, s.Node.Change(On)
}

func (s *Source) Debug() string {
	return fmt.Sprintf("Source<node: %s>", s.Node.Debug())
}

type Ground struct {
	Node *Node
}

func NewGround(node *Node) *Ground {
	return &Ground{
		Node: node,
	}
}

func (g *Ground) Act(force bool) (Completed, error) {
	return true, g.Node.Change(Off)
}

func (g *Ground) Debug() string {
	return fmt.Sprintf("Ground<node: %s>", g.Node.Debug())
}

type Resistor struct {
	Node1 *Node
	Node2 *Node
}

func NewResistor(node1, node2 *Node) *Resistor {
	return &Resistor{
		Node1: node1,
		Node2: node2,
	}
}

func (r *Resistor) Act(force bool) (Completed, error) {
	if r.Node1.State == Undefined && r.Node2.State == Undefined {
		if !force {
			return false, nil
		} else {
			return true, fmt.Errorf("resistor is hanging in undefined with terminals %s and %s",
				r.Node1.ID, r.Node2.ID)
		}
	}
	if r.Node1.State == Undefined {
		if !force {
			return false, nil
		}
		return true, r.Node1.Change(r.Node2.State)
	}
	if r.Node2.State == Undefined {
		if !force {
			return false, nil
		}
		return true, r.Node2.Change(r.Node1.State)
	}
	return true, nil
}

func (r *Resistor) Debug() string {
	return fmt.Sprintf("Resistor<node1: %s, node2: %s>", r.Node1.Debug(), r.Node2.Debug())
}

type Transistor struct {
	Source *Node
	Drain  *Node
	Gate   *Node
}

func NewTransistor(source, gate, drain *Node) *Transistor {
	return &Transistor{
		Source: source,
		Drain:  drain,
		Gate:   gate,
	}
}

// Transistor Act
// the transistor will short-circuit source and drain if gate is on and isolate
// them otherwise
func (t *Transistor) Act(force bool) (Completed, error) {
	if t.Gate.State == Undefined || (t.Source.State == Undefined && t.Drain.State == Undefined) {
		if !force {
			return false, nil
		} else {
			return true, fmt.Errorf("cannot perform action because transistor is in inconsistent state: %s", t.Debug())
		}
	}
	if t.Gate.State != On {
		return true, nil
	}
	if t.Source.State == On {
		if err := t.Drain.Change(On); err != nil {
			return true, err
		}
	} else if t.Drain.State == Off {
		if err := t.Source.Change(Off); err != nil {
			return true, err
		}
	}
	return true, nil
}

func (t *Transistor) Debug() string {
	return fmt.Sprintf("Transistor<source=%s, gate=%s, drain=%s>",
		t.Source.Debug(), t.Gate.Debug(), t.Drain.Debug())
}

type Multimeter struct {
	Node *Node
}

func NewMultimeter(node *Node) *Multimeter {
	return &Multimeter{
		Node: node,
	}
}

func (m *Multimeter) Act(force bool) (Completed, error) {
	if m.Node.State != Undefined || force {
		fmt.Println(m.Debug())
	}
	return true, nil
}

func (m *Multimeter) Debug() string {
	return fmt.Sprintf("Multimeter<node=%s, state=%s>", m.Node.ID, m.Node.State)
}

type Input struct {
	Node  *Node
	Value NodeState
}

func NewInput(node *Node, value NodeState) *Input {
	return &Input{
		Node:  node,
		Value: value,
	}
}

func (i *Input) Act(force bool) (Completed, error) {
	return true, i.Node.Change(i.Value)
}

func (i *Input) Debug() string {
	return fmt.Sprintf("Input<node:%s>", i.Node.Debug())
}

type CustomComponent struct {
	ComponentType string
	Subcomponents []Component
}

func NewCustomComponent(componentType string, subcomponents []Component) *CustomComponent {
	return &CustomComponent{
		ComponentType: componentType,
		Subcomponents: subcomponents,
	}
}

// TODO: add execution for custom component, will be needed for chaining executions
func (c *CustomComponent) Act(force bool) (Completed, error) {
	return true, nil
}

func (c *CustomComponent) Debug() string {
	return ""
}

type Circuit struct {
	debug     bool
	maxDefers int
	// components that cannot be deferred (sources, grounds and inputs)
	terminals   []Component
	resistors   []Component
	transistors []Component
	meters      []Component
}

func NewCircuit(components []Component, maxDefers int, debug bool) *Circuit {
	circuit := &Circuit{
		debug:     debug,
		maxDefers: maxDefers,
	}
	circuit.AddComponents(components...)
	return circuit
}

func (c *Circuit) addComponent(component Component) {
	switch component.(type) {
	case *Source, *Ground, *Input:
		c.terminals = append(c.terminals, component)
	case *Resistor:
		c.resistors = append(c.resistors, component)
	case *Transistor:
		c.transistors = append(c.transistors, component)
	case *Multimeter:
		c.meters = append(c.meters, component)
	case *CustomComponent:
		for _, subcomponent := range component.(*CustomComponent).Subcomponents {
			c.addComponent(subcomponent)
		}
	default:
		panic("component type not found")
	}
}

func (c *Circuit) AddComponents(components ...Component) {
	for _, component := range components {
		c.addComponent(component)
	}
}

// executes a simulation tick and returns the
// list of deferred components and potential errors
func (c *Circuit) tick(components []Component, force bool) ([]Component, error) {
	deferred := []Component{}
	for _, component := range components {
		if c.debug {
			fmt.Println("before:", component.Debug())
		}
		completed, err := component.Act(force)
		if c.debug {
			fmt.Println("after: ", component.Debug())
		}
		if err != nil {
			return []Component{}, err
		}
		if !completed {
			deferred = append(deferred, component)
		}
	}
	return deferred, nil
}

func (c *Circuit) simulateTicks(components []Component, force bool, maxDefers int) ([]Component, error) {
	deferred := components
	var err error
	for i := range maxDefers {
		fmt.Println("----- running tick", i, "-------")
		deferred, err = c.tick(deferred, force)
		if err != nil {
			return []Component{}, err
		}
		if len(deferred) == 0 {
			return deferred, nil
		}
	}
	return deferred, nil
}

func (c *Circuit) resistorHailMary(components []Component) (deferred []Component, err error) {
	fmt.Println("----- running resistor hail mary -------")
	for _, component := range components {
		switch component.(type) {
		case *Resistor:
			if c.debug {
				fmt.Println("before:", component.Debug())
			}
			_, err := component.Act(true)
			if c.debug {
				fmt.Println("after: ", component.Debug())
			}
			if err != nil {
				return []Component{}, err
			}
		default:
			deferred = append(deferred, component)
		}
	}
	return deferred, nil
}

// TODO: update function to propagate in queue format instead of running all components at once
func (c *Circuit) Simulate() (err error) {
	for _, terminal := range c.terminals {
		if _, err := terminal.Act(true); err != nil {
			return err
		}
	}
	components := []Component{}
	components = append(components, c.transistors...)
	components = append(components, c.resistors...)
	components, err = c.simulateTicks(components, false, c.maxDefers/2)
	if err != nil {
		return err
	}
	components, err = c.resistorHailMary(components)
	if err != nil {
		return err
	}
	components, err = c.simulateTicks(components, false, c.maxDefers/2)
	if err != nil {
		return err
	}
	components, err = c.tick(components, true)
	if err != nil {
		return err
	}
	if len(components) > 0 {
		// TODO: improve error message
		return fmt.Errorf("there were components that didn't get to a stable position")
	}

	for _, meter := range c.meters {
		if _, err := meter.Act(true); err != nil {
			return err
		}
	}
	return nil
}

var SharedSourceNode = NewNode("SharedSource")
var SharedGroundNode = NewNode("SharedGround")
var BaseComponents = []Component{
	NewSource(SharedSourceNode),
	NewGround(SharedGroundNode),
}

// performs NOT logic for input
//
//	         Vcc
//	         ───
//	          │
//	          >
//	          >
//	          >
//	          ├───o output
//	        ┌─┘
//	input o─│
//	        └─┐
//	        ──┴──
//	         GND
func NewNotGate(input *Node) (*Node, *CustomComponent) {
	outputNode := NewNode("NotOutput")
	return outputNode, &CustomComponent{
		ComponentType: "NotGate",
		Subcomponents: []Component{
			NewTransistor(outputNode, input, SharedGroundNode),
			NewResistor(SharedSourceNode, outputNode),
		},
	}
}

// performs AND logic for input1 and input2
//
//	          Vcc
//	          ───
//	           │
//	         ┌─┘
//	input1 o─│
//	         └─┐
//	         ┌─┘
//	input2 o─│
//	         └─┐
//	           ├───o output
//	           >
//	           >
//	           >
//	           │
//	         ──┴──
//	          GND
func NewAndGate(input1, input2 *Node) (*Node, *CustomComponent) {
	intermediateNode := NewNode("AndIntermediate")
	outputNode := NewNode("AndOutput")
	return outputNode, &CustomComponent{
		ComponentType: "AndGate",
		Subcomponents: []Component{
			NewTransistor(SharedSourceNode, input1, intermediateNode),
			NewTransistor(intermediateNode, input2, outputNode),
			NewResistor(outputNode, SharedGroundNode),
		},
	}
}

// performs OR logic for input1 and input2
//
//	          Vcc
//	          ───
//	           │
//	         ┌─┘─┐
//	input1 o─│   │─o input2
//	         └─┐─┘
//	           │
//	           ├───o output
//	           >
//	           >
//	           >
//	           │
//	         ──┴──
//	          GND
func NewOrGate(input1, input2 *Node) (*Node, *CustomComponent) {
	outputNode := NewNode("OrOutput")
	return outputNode, &CustomComponent{
		ComponentType: "AndGate",
		Subcomponents: []Component{
			NewTransistor(SharedSourceNode, input1, outputNode),
			NewTransistor(SharedSourceNode, input2, outputNode),
			NewResistor(outputNode, SharedGroundNode),
		},
	}
}

// performs NAND logic for input1 and input2
//
//	          Vcc
//	          ───
//	           │
//	           >
//	           >
//	           >
//	           │
//	           ├───o output
//	         ┌─┘
//	input1 o─│
//	         └─┐
//	         ┌─┘
//	input2 o─│
//	         └─┐
//	         ──┴──
//	          GND
func NewNandGate(input1, input2 *Node) (*Node, *CustomComponent) {
	intermediateNode := NewNode("NandIntermediate")
	outputNode := NewNode("NandOutput")
	return outputNode, &CustomComponent{
		ComponentType: "NandGate",
		Subcomponents: []Component{
			NewResistor(SharedSourceNode, outputNode),
			NewTransistor(outputNode, input1, intermediateNode),
			NewTransistor(intermediateNode, input2, SharedGroundNode),
		},
	}
}

func NewXorGate(input1, input2 *Node) (*Node, *CustomComponent) {
	orOut, orComponent := NewOrGate(input1, input2)
	nandOut, nandComponent := NewNandGate(input1, input2)
	outputNode, andComponent := NewAndGate(orOut, nandOut)
	return outputNode, &CustomComponent{
		ComponentType: "XorGate",
		Subcomponents: []Component{orComponent, nandComponent, andComponent},
	}
}

func main() {
	components := []Component{}
	components = append(components, BaseComponents...)

	input1 := NewNode("input1")
	input2 := NewNode("input2")
	inputComponents := []Component{
		NewInput(input1, Off),
		NewInput(input2, On),
	}
	components = append(components, inputComponents...)

	// NOT gate
	notOutput, notGate := NewNotGate(input1)
	notGateComponents := []Component{
		notGate,
		NewMultimeter(notOutput),
	}
	components = append(components, notGateComponents...)

	// AND gate
	andOutput, andGate := NewAndGate(notOutput, input2)
	andGateComponents := []Component{
		andGate,
		NewMultimeter(andOutput),
	}
	components = append(components, andGateComponents...)

	// OR gate
	orOutput, orGate := NewOrGate(input1, input2)
	orGateComponents := []Component{
		orGate,
		NewMultimeter(orOutput),
	}
	components = append(components, orGateComponents...)

	// NAND gate
	nandOutput, nandGate := NewNandGate(input1, input2)
	nandGateComponents := []Component{
		nandGate,
		NewMultimeter(nandOutput),
	}
	components = append(components, nandGateComponents...)

	// XOR gate
	xorOutput, xorGate := NewXorGate(input1, input2)
	xorGateComponents := []Component{
		xorGate,
		NewMultimeter(xorOutput),
	}
	components = append(components, xorGateComponents...)

	maxDefers := 10
	debug := true
	circuit := NewCircuit(components, maxDefers, debug)
	err := circuit.Simulate()
	if err != nil {
		fmt.Printf("unable to run simulation: %s\n", err.Error())
		os.Exit(1)
	}
}
