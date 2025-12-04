package main

import "fmt"

type Circuit struct {
	debug      bool
	maxDefers  int
	terminals  []Component
	components []Component
	meters     []Component
}

func NewCircuit(components []Component, maxDefers int, debug bool) *Circuit {
	circuit := &Circuit{
		debug:     debug,
		maxDefers: maxDefers,
	}
	circuit.AddComponents(append(BaseComponents, components...)...)
	return circuit
}

func (c *Circuit) addComponent(component Component) {
	switch component.(type) {
	case *Terminal:
		fmt.Println("adding terminal ", component.Debug())
		c.terminals = append(c.terminals, component)
	case *Meter:
		fmt.Println("adding meter ", component.Debug())
		c.meters = append(c.meters, component)
	default:
		fmt.Println("adding other ", component.Debug())
		c.components = append(c.components, component)
	}
}

func (c *Circuit) AddComponents(components ...Component) {
	for _, component := range components {
		c.addComponent(component)
	}
}

func (c *Circuit) Tick() (err error) {
	for _, terminal := range c.terminals {
		if err := terminal.Act(); err != nil {
			return err
		}
	}
	if err = ActComponents(c.components, c.maxDefers); err != nil {
		return err
	}
	for _, meter := range c.meters {
		if err := meter.Act(); err != nil {
			return err
		}
	}
	return nil
}
