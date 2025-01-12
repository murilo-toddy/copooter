package main

import "testing"

func TestNotGate(t *testing.T) {
	tt := []struct {
		input          NodeState
		expectedOutput NodeState
	}{
		{input: Off, expectedOutput: On},
		{input: On, expectedOutput: Off},
	}
	for _, tc := range tt {
		components := []Component{}
		input := NewNode("Input")
		inputComponents := []Component{NewInput(input, tc.input)}
		components = append(components, inputComponents...)

		notOutput, notGate := NewNotGate(input)
		components = append(components, notGate)

		c := NewCircuit(components, 4, true)
		if err := c.Simulate(); err != nil {
			t.Errorf(err.Error())
		}
		if notOutput.State != tc.expectedOutput {
			t.Errorf("input: %s generated output state %s instead of %s",
				tc.input, notOutput.State, tc.expectedOutput)
		}
	}
}

func TestAndGate(t *testing.T) {
	tt := []struct {
		input1         NodeState
		input2         NodeState
		expectedOutput NodeState
	}{
		{input1: Off, input2: Off, expectedOutput: Off},
		{input1: On, input2: Off, expectedOutput: Off},
		{input1: Off, input2: On, expectedOutput: Off},
		{input1: On, input2: On, expectedOutput: On},
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

		andOutput, andGate := NewAndGate(input1, input2)
		components = append(components, andGate)

		c := NewCircuit(components, 4, false)
		if err := c.Simulate(); err != nil {
			t.Errorf(err.Error())
		}
		if andOutput.State != tc.expectedOutput {
			t.Errorf("Inputs<input1: %s, input2: %s> generated output state %s instead of %s",
				tc.input1, tc.input2, andOutput.State, tc.expectedOutput)
		}
	}
}

func TestOrGate(t *testing.T) {
	tt := []struct {
		input1         NodeState
		input2         NodeState
		expectedOutput NodeState
	}{
		{input1: Off, input2: Off, expectedOutput: Off},
		{input1: On, input2: Off, expectedOutput: On},
		{input1: Off, input2: On, expectedOutput: On},
		{input1: On, input2: On, expectedOutput: On},
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

		orOutput, orGate := NewOrGate(input1, input2)
		components = append(components, orGate)

		c := NewCircuit(components, 4, false)
		if err := c.Simulate(); err != nil {
			t.Errorf(err.Error())
		}
		if orOutput.State != tc.expectedOutput {
			t.Errorf("Inputs<input1: %s, input2: %s> generated output state %s instead of %s",
				tc.input1, tc.input2, orOutput.State, tc.expectedOutput)
		}
	}
}

func TestNandGate(t *testing.T) {
	tt := []struct {
		input1         NodeState
		input2         NodeState
		expectedOutput NodeState
	}{
		{input1: Off, input2: Off, expectedOutput: Off},
		{input1: On, input2: Off, expectedOutput: On},
		{input1: Off, input2: On, expectedOutput: On},
		{input1: On, input2: On, expectedOutput: On},
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

		nandOutput, nandGate := NewOrGate(input1, input2)
		components = append(components, nandGate)

		c := NewCircuit(components, 4, false)
		if err := c.Simulate(); err != nil {
			t.Errorf(err.Error())
		}
		if nandOutput.State != tc.expectedOutput {
			t.Errorf("Inputs<input1: %s, input2: %s> generated output state %s instead of %s",
				tc.input1, tc.input2, nandOutput.State, tc.expectedOutput)
		}
	}
}

func TestXorGate(t *testing.T) {
	tt := []struct {
		input1         NodeState
		input2         NodeState
		expectedOutput NodeState
	}{
		{input1: Off, input2: Off, expectedOutput: Off},
		{input1: On, input2: Off, expectedOutput: On},
		{input1: Off, input2: On, expectedOutput: On},
		{input1: On, input2: On, expectedOutput: Off},
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

		xorOutput, xorGate := NewXorGate(input1, input2)
		components = append(components, xorGate)

		c := NewCircuit(components, 4, false)
		if err := c.Simulate(); err != nil {
			t.Errorf(err.Error())
		}
		if xorOutput.State != tc.expectedOutput {
			t.Errorf("Inputs<input1: %s, input2: %s> generated output state %s instead of %s",
				tc.input1, tc.input2, xorOutput.State, tc.expectedOutput)
		}
	}
}

func TestChainingGates(t *testing.T) {
	input := NewNode("Input")
	components := []Component{NewInput(input, On)}

	notOut, notGate := NewNotGate(input)
	xorOut, xorGate := NewNandGate(input, notOut)
	andOut, andGate := NewAndGate(notOut, xorOut)
	orOut, orGate := NewOrGate(input, xorOut)
	nandOut, nandGate := NewNandGate(orOut, orOut)
	components = append(components, []Component{notGate, xorGate, andGate, orGate, nandGate}...)

	c := NewCircuit(components, 10, false)
	if err := c.Simulate(); err != nil {
		t.Errorf(err.Error())
	}
	if notOut.State != Off {
		t.Errorf("expected NOT gate to output off, but got %s", notOut.State)
	}
	if xorOut.State != On {
		t.Errorf("expected XOR gate to output on, but got %s", xorOut.State)
	}
	if andOut.State != Off {
		t.Errorf("expected AND gate to output off, but got %s", andOut.State)
	}
	if orOut.State != On {
		t.Errorf("expected OR gate to output on, but got %s", orOut.State)
	}
	if nandOut.State != Off {
		t.Errorf("expected NAND gate to output off, but got %s", nandOut.State)
	}
}
