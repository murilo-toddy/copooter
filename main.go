package main

import (
	"fmt"
	"os"
)

func main() {
	components := []Component{}
	components = append(components, BaseComponents...)

	input1 := NewNode("input1")
	input2 := NewNode("input2")
	inputComponents := []Component{
		NewInput(input1, Off),
		NewInput(input2, Off),
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
	andOutput, andGate := NewAndGate(input1, input2)
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
	debug := false
	circuit := NewCircuit(components, maxDefers, debug)
	err := circuit.Simulate()
	if err != nil {
		fmt.Printf("unable to run simulation: %s\n", err.Error())
		os.Exit(1)
	}
}
