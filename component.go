package main

import (
	"fmt"
	"strings"
)

type ComponentType int

type Component interface {
	Ready() bool
	// propagates component input to its outputs, should only be called if c.Ready() returns true
	Act() error
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

func (s *Source) Ready() bool {
	return true
}

func (s *Source) Act() error {
	return s.Node.Change(On)
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

func (g *Ground) Ready() bool {
	return true
}

func (g *Ground) Act() error {
	return g.Node.Change(Off)
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

func (r *Resistor) Ready() bool {
	return r.Node1.State != Undefined || r.Node2.State != Undefined
}

func (r *Resistor) Act() error {
	if r.Node1.State == Undefined && r.Node2.State == Undefined {
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

func (t *Transistor) Ready() bool {
	return t.Gate.State != Undefined && (t.Source.State != Undefined || t.Drain.State != Undefined)
}

// Transistor Act
// the transistor will short-circuit source and drain if gate is on and isolate
// them otherwise
func (t *Transistor) Act() error {
	if t.Gate.State == Undefined || (t.Source.State == Undefined && t.Drain.State == Undefined) {
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

func (m *Multimeter) Ready() bool {
	return true
}

func (m *Multimeter) Act() error {
	if m.Node.State != Undefined {
		fmt.Println(m.Debug())
	}
	return nil
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

func (i *Input) Ready() bool {
	return true
}

func (i *Input) Act() error {
	return i.Node.Change(i.Value)
}

func (i *Input) Debug() string {
	return fmt.Sprintf("Input<node:%s>", i.Node.Debug())
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

func (c *CustomComponent) Ready() bool {
	for _, input := range c.Inputs {
		if input.State == Undefined {
			return false
		}
	}
	return true
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
		switch component.(type) {
		case *Transistor:
			transistors = append(transistors, component.(*Transistor))
		case *Resistor:
			resistors = append(resistors, component.(*Resistor))
		default:
			others = append(others, component)
		}
	}
	return
}

func tick(components []Component) (deferred []Component, err error) {
	for _, component := range components {
		if component.Ready() {
			if err = component.Act(); err != nil {
				return
			}
		} else {
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
	NewSource(SharedSourceNode),
	NewGround(SharedGroundNode),
}
