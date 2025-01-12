package main

import "fmt"

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
	parent := "NotGate"
	outputNode := NewNode(fmt.Sprintf("%s-Output", parent))
	return outputNode, NewCustomComponent(
		"NotGate",
		[]Component{
			NewResistor(parent, SharedSourceNode, outputNode),
			NewTransistor(parent, outputNode, input, SharedGroundNode),
		},
		[]*Node{input},
	)
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
	parent := "AndGate"
	intermediateNode := NewNode(fmt.Sprintf("%s-AndIntermediate", parent))
	outputNode := NewNode(fmt.Sprintf("%s-AndOutput", parent))
	return outputNode, NewCustomComponent(
		"AndGate",
		[]Component{
			NewTransistor(parent, SharedSourceNode, input1, intermediateNode),
			NewTransistor(parent, intermediateNode, input2, outputNode),
			NewResistor(parent, outputNode, SharedGroundNode),
		},
		[]*Node{input1, input2},
	)
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
	parent := "OrGate"
	outputNode := NewNode(fmt.Sprintf("%s-OrOutput", parent))
	return outputNode, NewCustomComponent(
		"OrGate",
		[]Component{
			NewTransistor(parent, SharedSourceNode, input1, outputNode),
			NewTransistor(parent, SharedSourceNode, input2, outputNode),
			NewResistor(parent, outputNode, SharedGroundNode),
		},
		[]*Node{input1, input2},
	)
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
	parent := "NandGate"
	intermediateNode := NewNode("NandIntermediate")
	outputNode := NewNode("NandOutput")
	return outputNode, NewCustomComponent(
		"NandGate",
		[]Component{
			NewTransistor(parent, outputNode, input1, intermediateNode),
			NewTransistor(parent, intermediateNode, input2, SharedGroundNode),
			NewResistor(parent, SharedSourceNode, outputNode),
		},
		[]*Node{input1, input2},
	)
}

// performs XOR logic for input1 and input2
// input1 XOR input2 = (input1 OR input2) AND (input1 NAND input2)
func NewXorGate(input1, input2 *Node) (*Node, *CustomComponent) {
	orOut, orComponent := NewOrGate(input1, input2)
	nandOut, nandComponent := NewNandGate(input1, input2)
	outputNode, andComponent := NewAndGate(orOut, nandOut)
	return outputNode, NewCustomComponent(
		"XorGate",
		[]Component{orComponent, nandComponent, andComponent},
		[]*Node{input1, input2},
	)
}
