package main

import (
	"fmt"
	"strings"
)

type ComponentType int

const (
	resistorSprite   = "resources/resistor.png"
	transistorSprite = "resources/transistor.png"
)

// TODO: remove rendering stuff to separate interface
type Component interface {
	Ready() bool
	// propagates component input to its outputs, should only be called if c.Ready() returns true
	Act() error
	Debug() string
	Drawable() *DrawableComponent
}

type Terminal struct {
	Node         *Node
	state        NodeState
	terminalType string
	drawable     *DrawableComponent
}

func NewTerminal(node *Node, state NodeState, terminalType string) *Terminal {
	return &Terminal{
		Node:         node,
		state:        state,
		terminalType: terminalType,
		drawable:     NewDrawable(0, 0, ""),
	}
}

func NewSource(node *Node) *Terminal {
	return NewTerminal(node, On, "Source")
}

func NewGround(node *Node) *Terminal {
	return NewTerminal(node, Off, "Ground")
}

func NewInput(node *Node, state NodeState) *Terminal {
	return NewTerminal(node, state, "Input")
}

func (t *Terminal) Ready() bool {
	return true
}

func (t *Terminal) Act() error {
	return t.Node.Change(t.state)
}

func (t *Terminal) Debug() string {
	return fmt.Sprintf("%s<node: %s>", t.terminalType, t.Node.Debug())
}

func (t *Terminal) Drawable() *DrawableComponent {
	return t.drawable
}

type Meter struct {
	Node     *Node
	drawable *DrawableComponent
}

func NewMultimeter(node *Node) *Meter {
	return &Meter{
		Node:     node,
		drawable: NewDrawable(0, 0, ""),
	}
}

func (m *Meter) Ready() bool {
	return true
}

func (m *Meter) Act() error {
	if m.Node.State != Undefined {
		fmt.Println(m.Debug())
	}
	return nil
}

func (m *Meter) Debug() string {
	return fmt.Sprintf("Multimeter<node=%s, state=%s>", m.Node.ID, m.Node.State)
}

func (m *Meter) Drawable() *DrawableComponent {
	return m.drawable
}

type Resistor struct {
	Node1    *Node
	Node2    *Node
	drawable *DrawableComponent
}

func NewResistor(parent string, node1, node2 *Node) *Resistor {
	return &Resistor{
		Node1:    NewNode(fmt.Sprintf("%s-Resistor-Node1", parent)).Connect(node1),
		Node2:    NewNode(fmt.Sprintf("%s-Resistor-Node2", parent)).Connect(node2),
		drawable: NewDrawable(200, 100, resistorSprite),
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

func (r *Resistor) Drawable() *DrawableComponent {
	return r.drawable
}

type Transistor struct {
	Source   *Node
	Drain    *Node
	Gate     *Node
	drawable *DrawableComponent
}

func NewTransistor(parent string, source, gate, drain *Node) *Transistor {
	return &Transistor{
		Source:   NewNode(fmt.Sprintf("%s-Transistor-Source", parent)).Connect(source),
		Drain:    NewNode(fmt.Sprintf("%s-Transistor-Drain", parent)).Connect(drain),
		Gate:     NewNode(fmt.Sprintf("%s-Transistor-Gate", parent)).Connect(gate),
		drawable: NewDrawable(200, 400, transistorSprite),
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

func (t *Transistor) Drawable() *DrawableComponent {
	return t.drawable
}

type CustomComponent struct {
	ComponentType string
	Subcomponents []Component
	Inputs        []*Node
	maxDefers     int
	drawable      *DrawableComponent
}

func NewCustomComponent(componentType string, subcomponents []Component, inputs []*Node) *CustomComponent {
	return &CustomComponent{
		ComponentType: componentType,
		Subcomponents: subcomponents,
		Inputs:        inputs,
		maxDefers:     4,
		drawable:      NewDrawable(0, 0, ""),
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

func (c *CustomComponent) Drawable() *DrawableComponent {
	return c.drawable
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
	debug := false
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
