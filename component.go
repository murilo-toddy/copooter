// TODO: share NewX and NewXFromNodes
package main

import (
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ComponentType int

const (
	TypeTerminal ComponentType = iota
	TypeResistor
	TypeTransistor
	TypeMeter
	TypeInput
)

type RenderData struct {
	ID   string
	Name string
	// TODO: this is not a good way of storing resources for components
	// that may change based on state, such as input or meter
	idleResource         rl.Texture2D
	idleResourceName     string
	selectedResource     rl.Texture2D
	selectedResourceName string
	X                    int32
	Y                    int32
	Component
}

type Component interface {
	// resets component to default state
	Reset()
	// checks if a component is ready to be executed
	Ready() bool
	// propagates component input to its outputs, should only be called if c.Ready() returns true
	Act() error

	GetRenderData() *RenderData

	// TODO: make components own a []*Node list to avoid multiple instantiations per render cycle
	Nodes() []*Node

	// TODO: make clone empty node connections
	Clone() Component

	Debug() string
}

// TODO: make terminalType an enum
type Terminal struct {
	Node         *Node
	state        NodeState
	terminalType string
	RenderData
}

func NewTerminal(
	name string,
	node *Node,
	state NodeState,
	terminalType string,
) *Terminal {
	t := &Terminal{
		Node:         NewNode(fmt.Sprintf("%s-Node", name)).Connect(node),
		state:        state,
		terminalType: terminalType,
	}
	t.Node.Parent = t
	return t
}

func NewTerminalFromNodes(
	node *Node,
	state NodeState,
	terminalType string,
) *Terminal {
	t := &Terminal{Node: node, state: state, terminalType: terminalType}
	t.Node.Parent = t
	return t
}

func NewSource(name string, node *Node) *Terminal {
	return NewTerminal(name, node, On, "Source")
}

func NewSourceFromNodes(node *Node) *Terminal {
	return NewTerminalFromNodes(node, On, "Source")
}

func NewGround(name string, node *Node) *Terminal {
	return NewTerminal(name, node, Off, "Ground")
}

func NewGroundFromNodes(node *Node) *Terminal {
	return NewTerminalFromNodes(node, Off, "Ground")
}

func NewInput(name string, node *Node, state NodeState) *Terminal {
	return NewTerminal(name, node, state, "Input")
}

func NewInputFromNodes(node *Node, state NodeState) *Terminal {
	return NewTerminalFromNodes(node, state, "Input")
}

func (t *Terminal) Reset() {
	t.Node.State = Undefined
}

func (t *Terminal) Ready() bool {
	return true
}

func (t *Terminal) Act() error {
	return t.Node.Change(t.state)
}

func (t *Terminal) GetRenderData() *RenderData {
	return &t.RenderData
}

func (t *Terminal) Nodes() []*Node {
	return []*Node{t.Node}
}

func (t Terminal) Clone() Component {
	nodeCopy := *t.Node

	newTerminal := t
	newTerminal.Node = &nodeCopy
	newTerminal.Node.Parent = &newTerminal

	return &newTerminal
}

func (t *Terminal) Debug() string {
	return fmt.Sprintf("%s<node: %s>", t.terminalType, t.Node.Debug())
}

type Meter struct {
	Node *Node
	RenderData
}

func NewMultimeter(name string, node *Node) *Meter {
	m := &Meter{
		Node: NewNode(fmt.Sprintf("%s-Node", name)).Connect(node),
	}
	m.Node.Parent = m
	return m
}

func NewMultimeterFromNodes(node *Node) *Meter {
	m := &Meter{Node: node}
	m.Node.Parent = m
	return m
}

func (m *Meter) Reset() {
	m.Node.State = Undefined
}

func (m *Meter) Ready() bool {
	return true
}

func (m *Meter) Act() error {
	if m.Node.State == Undefined {
		fmt.Println("WARN: acting on ", m.Debug(), " in undefined state")
	}
	fmt.Println(m.Debug())
	return nil
}

func (m *Meter) Nodes() []*Node {
	return []*Node{m.Node}
}

func (m Meter) Clone() Component {
	nodeCopy := *m.Node

	newMeter := m
	newMeter.Node = &nodeCopy
	newMeter.Node.Parent = &newMeter
	return &newMeter
}

func (m *Meter) Debug() string {
	return fmt.Sprintf("Multimeter<node=%s, state=%s>", m.Node.ID, m.Node.State)
}

func (m *Meter) GetRenderData() *RenderData {
	return &m.RenderData
}

type Resistor struct {
	Node1 *Node
	Node2 *Node
	RenderData
}

func NewResistor(name string, node1, node2 *Node) *Resistor {
	r := &Resistor{
		Node1: NewNode(fmt.Sprintf("%s-Node1", name)).Connect(node1),
		Node2: NewNode(fmt.Sprintf("%s-Node2", name)).Connect(node2),
	}
	r.Node1.Parent = r
	r.Node2.Parent = r
	return r
}

func NewResistorFromNodes(node1, node2 *Node) *Resistor {
	r := &Resistor{
		Node1: node1,
		Node2: node2,
	}
	r.Node1.Parent = r
	r.Node2.Parent = r
	return r
}

func (r *Resistor) Reset() {
	r.Node1.State = Undefined
	r.Node2.State = Undefined
}

func (r *Resistor) Ready() bool {
	return r.Node1.State != Undefined || r.Node2.State != Undefined
}

func (r *Resistor) Act() error {
	if !r.Ready() {
		return fmt.Errorf("component %s was executed before it was ready", r.Debug())
	}
	if r.Node1.State == Undefined {
		return r.Node1.Change(r.Node2.State)
	}
	if r.Node2.State == Undefined {
		return r.Node2.Change(r.Node1.State)
	}
	return nil
}

func (r *Resistor) Nodes() []*Node {
	return []*Node{r.Node1, r.Node2}
}

func (r Resistor) Clone() Component {
	node1Copy := *r.Node1
	node2Copy := *r.Node2

	newResistor := r
	newResistor.Node1 = &node1Copy
	newResistor.Node2 = &node2Copy
	newResistor.Node1.Parent = &newResistor
	newResistor.Node2.Parent = &newResistor
	return &newResistor
}

func (r *Resistor) GetRenderData() *RenderData {
	return &r.RenderData
}

func (r *Resistor) Debug() string {
	return fmt.Sprintf("Resistor<node1: %s, node2: %s>", r.Node1.Debug(), r.Node2.Debug())
}

type Transistor struct {
	Source *Node
	Drain  *Node
	Gate   *Node
	RenderData
}

func NewTransistor(name string, source, gate, drain *Node) *Transistor {
	t := &Transistor{
		Source: NewNode(fmt.Sprintf("%s-Source", name)).Connect(source),
		Drain:  NewNode(fmt.Sprintf("%s-Drain", name)).Connect(drain),
		Gate:   NewNode(fmt.Sprintf("%s-Gate", name)).Connect(gate),
	}
	t.Source.Parent = t
	t.Drain.Parent = t
	t.Gate.Parent = t
	return t
}

func NewTransistorFromNodes(source, gate, drain *Node) *Transistor {
	t := &Transistor{Source: source, Drain: drain, Gate: gate}
	t.Source.Parent = t
	t.Drain.Parent = t
	t.Gate.Parent = t
	return t
}

func (t *Transistor) Reset() {
	t.Source.State = Undefined
	t.Drain.State = Undefined
	t.Gate.State = Undefined
}

func (t *Transistor) Ready() bool {
	return t.Gate.State != Undefined && (t.Source.State != Undefined || t.Drain.State != Undefined)
}

// Transistor Act
// the transistor will short-circuit source and drain if gate is on and isolate
// them otherwise
func (t *Transistor) Act() error {
	if !t.Ready() {
		return fmt.Errorf("component %s was executed before it was ready", t.Debug())
	}
	if t.Gate.State != On {
		return nil
	}
	if t.Source.State == On {
		if err := t.Drain.Change(On); err != nil {
			return err
		}
	} else if t.Drain.State == Off {
		if err := t.Source.Change(Off); err != nil {
			return err
		}
	}
	return nil
}

func (t *Transistor) Nodes() []*Node {
	return []*Node{t.Source, t.Gate, t.Drain}
}

func (t Transistor) Clone() Component {
	sourceCopy := *t.Source
	gateCopy := *t.Gate
	drainCopy := *t.Drain

	newTransistor := t
	newTransistor.Source = &sourceCopy
	newTransistor.Gate = &gateCopy
	newTransistor.Drain = &drainCopy

	newTransistor.Source.Parent = &newTransistor
	newTransistor.Gate.Parent = &newTransistor
	newTransistor.Drain.Parent = &newTransistor

	return &newTransistor
}

func (t *Transistor) GetRenderData() *RenderData {
	return &t.RenderData
}

func (t *Transistor) Debug() string {
	return fmt.Sprintf("Transistor<source=%s, gate=%s, drain=%s> (x: %d, y: %d)",
		t.Source.Debug(), t.Gate.Debug(), t.Drain.Debug(), t.RenderData.X, t.RenderData.Y)
}

type CustomComponent struct {
	ComponentType string
	Subcomponents []Component
	Inputs        []*Node
	maxDefers     int
}

func NewCustomComponent(componentType string, subcomponents []Component, inputs []*Node) *CustomComponent {
	return &CustomComponent{
		ComponentType: componentType,
		Subcomponents: subcomponents,
		Inputs:        inputs,
		maxDefers:     4,
	}
}

func (c *CustomComponent) Reset() {
	for _, s := range c.Subcomponents {
		s.Reset()
	}
}

func (c *CustomComponent) Ready() bool {
	for _, input := range c.Inputs {
		if input.State == Undefined {
			return false
		}
	}
	return true
}

func (c *CustomComponent) Nodes() []*Node {
	// TODO
	return []*Node{}
}

func (c *CustomComponent) GetRenderData() *RenderData {
	return &RenderData{}
}

func (c *CustomComponent) Clone() Component {
	// TODO
	return c
}

func (c *CustomComponent) runSubcomponents(subcomponents []Component) (notExecutedComponents []Component, err error) {
	for _, subcomponent := range subcomponents {
		if subcomponent.Ready() {
			if err = subcomponent.Act(); err != nil {
				return
			}
		} else {
			notExecutedComponents = append(notExecutedComponents, subcomponent)
		}
	}
	return
}

// Returns as []Component even for transistors and resistors since go cannot cast
// []*Transistor to []Component for some reason (or maybe I'm just dumb)
func SplitComponents(components []Component) (transistors []Component, resistors []Component, others []Component) {
	for _, component := range components {
		switch component := component.(type) {
		case *Transistor:
			transistors = append(transistors, component)
		case *Resistor:
			resistors = append(resistors, component)
		default:
			others = append(others, component)
		}
	}
	return
}

func tick(components []Component) (deferred []Component, err error) {
	debug := true
	for _, component := range components {
		if component.Ready() {
			if debug {
				fmt.Println("component ready\n", "before: ", component.Debug())
			}
			if err = component.Act(); err != nil {
				return
			}
			if debug {
				fmt.Println("after: ", component.Debug())
			}
		} else {
			if debug {
				fmt.Println(component.Debug(), "component not ready")
			}
			deferred = append(deferred, component)
		}
	}
	return
}

func ActComponents(components []Component, maxDefers int) (err error) {
	deferredComponents := components
	for range maxDefers {
		transistors, resistors, others := SplitComponents(deferredComponents)
		transistorsLen := len(transistors)
		for {
			transistors, err = tick(transistors)
			if err != nil {
				return
			}
			if len(transistors) == transistorsLen {
				break
			}
			transistorsLen = len(transistors)
		}
		others, err = tick(others)
		if err != nil {
			return
		}
		resistors, err = tick(resistors)
		if err != nil {
			return
		}
		deferredComponents = append(transistors, others...)
		deferredComponents = append(deferredComponents, resistors...)
		if len(deferredComponents) == 0 {
			break
		}
	}
	return
}

func (c *CustomComponent) Act() error {
	return ActComponents(c.Subcomponents, c.maxDefers)
}

func (c *CustomComponent) Debug() string {
	var builder strings.Builder
	builder.WriteString(c.ComponentType + "\n")
	for _, subcomponent := range c.Subcomponents {
		builder.WriteString(subcomponent.Debug() + "\n")
	}
	return builder.String()
}

var BaseComponents = []Component{
	NewSource("SharedSource", SharedSourceNode),
	NewGround("SharedGround", SharedGroundNode),
}
