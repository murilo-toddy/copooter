package main

import "testing"

func TestSimpleAdder(t *testing.T) {
	tt := []struct {
		input1        NodeState
		input2        NodeState
		expectedOut   NodeState
		expectedCarry NodeState
	}{
		{input1: Off, input2: Off, expectedOut: Off, expectedCarry: Off},
		{input1: On, input2: Off, expectedOut: On, expectedCarry: Off},
		{input1: Off, input2: On, expectedOut: On, expectedCarry: Off},
		{input1: On, input2: On, expectedOut: Off, expectedCarry: On},
	}
	for _, tc := range tt {
		components := []Component{}
		input1 := NewNode("Input1")
		input2 := NewNode("Input2")
		inputComponents := []Component{
			NewInput(input1, tc.input1),
			NewInput(input2, tc.input2),
		}
		components = append(components, inputComponents...)

		adderOut, adderCarry, adder := NewSimpleAdder(input1, input2)
		components = append(components, adder)

		c := NewCircuit(components, 4, false)
		if err := c.Simulate(); err != nil {
			t.Errorf(err.Error())
		}
		if adderOut.State != tc.expectedOut || adderCarry.State != tc.expectedCarry {
			t.Errorf("Inputs<input1: %s, input2: %s> generated output state <out: %s, carry: %s> instead of <out: %s, carry: %s>",
				tc.input1, tc.input2, adderOut.State, adderCarry.State, tc.expectedOut, tc.expectedCarry)
		}
	}
}

func TestFullAdder(t *testing.T) {
	tt := []struct {
		input1        NodeState
		input2        NodeState
		carryIn       NodeState
		expectedOut   NodeState
		expectedCarry NodeState
	}{
		{input1: Off, input2: Off, carryIn: Off, expectedOut: Off, expectedCarry: Off},
		{input1: On, input2: Off, carryIn: Off, expectedOut: On, expectedCarry: Off},
		{input1: Off, input2: On, carryIn: Off, expectedOut: On, expectedCarry: Off},
		{input1: On, input2: On, carryIn: Off, expectedOut: Off, expectedCarry: On},
		{input1: Off, input2: Off, carryIn: On, expectedOut: On, expectedCarry: Off},
		{input1: On, input2: Off, carryIn: On, expectedOut: Off, expectedCarry: On},
		{input1: Off, input2: On, carryIn: On, expectedOut: Off, expectedCarry: On},
		{input1: On, input2: On, carryIn: On, expectedOut: On, expectedCarry: On},
	}
	for _, tc := range tt {
		components := []Component{}
		input1 := NewNode("Input1")
		input2 := NewNode("Input2")
		carryIn := NewNode("CarryIn")
		inputComponents := []Component{
			NewInput(input1, tc.input1),
			NewInput(input2, tc.input2),
			NewInput(carryIn, tc.carryIn),
		}
		components = append(components, inputComponents...)

		adderOut, adderCarry, adder := NewFullAdder(input1, input2, carryIn)
		components = append(components, adder)

		c := NewCircuit(components, 4, false)
		if err := c.Simulate(); err != nil {
			t.Errorf(err.Error())
		}
		if adderOut.State != tc.expectedOut || adderCarry.State != tc.expectedCarry {
			t.Errorf("Inputs<input1: %s, input2: %s> generated output state <out: %s, carry: %s> instead of <out: %s, carry: %s>",
				tc.input1, tc.input2, adderOut.State, adderCarry.State, tc.expectedOut, tc.expectedCarry)
		}
	}
}

func TestAdderSubtractor(t *testing.T) {
	tt := []struct {
		input1        NodeState
		input2        NodeState
		carryIn       NodeState
		operation     NodeState
		expectedOut   NodeState
		expectedCarry NodeState
	}{
		// addition
		{input1: Off, input2: Off, carryIn: Off, operation: Off, expectedOut: Off, expectedCarry: Off},
		{input1: On, input2: Off, carryIn: Off, operation: Off, expectedOut: On, expectedCarry: Off},
		{input1: Off, input2: On, carryIn: Off, operation: Off, expectedOut: On, expectedCarry: Off},
		{input1: On, input2: On, carryIn: Off, operation: Off, expectedOut: Off, expectedCarry: On},
		{input1: Off, input2: Off, carryIn: On, operation: Off, expectedOut: On, expectedCarry: Off},
		{input1: On, input2: Off, carryIn: On, operation: Off, expectedOut: Off, expectedCarry: On},
		{input1: Off, input2: On, carryIn: On, operation: Off, expectedOut: Off, expectedCarry: On},
		{input1: On, input2: On, carryIn: On, operation: Off, expectedOut: On, expectedCarry: On},
		// subtraction
		{input1: Off, input2: Off, carryIn: Off, operation: On, expectedOut: On, expectedCarry: Off},
		{input1: On, input2: Off, carryIn: Off, operation: On, expectedOut: Off, expectedCarry: On},
		{input1: Off, input2: On, carryIn: Off, operation: On, expectedOut: Off, expectedCarry: Off},
		{input1: On, input2: On, carryIn: Off, operation: On, expectedOut: On, expectedCarry: Off},
		{input1: Off, input2: Off, carryIn: On, operation: On, expectedOut: Off, expectedCarry: On},
		{input1: On, input2: Off, carryIn: On, operation: On, expectedOut: On, expectedCarry: On},
		{input1: Off, input2: On, carryIn: On, operation: On, expectedOut: On, expectedCarry: Off},
		{input1: On, input2: On, carryIn: On, operation: On, expectedOut: Off, expectedCarry: On},
	}
	for _, tc := range tt {
		components := []Component{}
		input1 := NewNode("Input1")
		input2 := NewNode("Input2")
		carryIn := NewNode("CarryIn")
		operation := NewNode("Operation")
		inputComponents := []Component{
			NewInput(input1, tc.input1),
			NewInput(input2, tc.input2),
			NewInput(carryIn, tc.carryIn),
			NewInput(operation, tc.operation),
		}
		components = append(components, inputComponents...)

		adderOut, adderCarry, adder := NewAdderSubtractor(input1, input2, carryIn, operation)
		components = append(components, adder)

		c := NewCircuit(components, 4, false)
		if err := c.Simulate(); err != nil {
			t.Errorf(err.Error())
		}
		if adderOut.State != tc.expectedOut || adderCarry.State != tc.expectedCarry {
			t.Errorf("Inputs<input1: %s, input2: %s, carryIn: %s, operation: %s> generated output state <out: %s, carry: %s> instead of <out: %s, carry: %s>",
				tc.input1, tc.input2, carryIn.State, operation.State, adderOut.State, adderCarry.State, tc.expectedOut, tc.expectedCarry)
		}
	}
}
