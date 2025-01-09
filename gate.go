package main

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
	return outputNode, NewCustomComponent(
		"NotGate",
		[]Component{
			NewTransistor(outputNode, input, SharedGroundNode),
			NewResistor(SharedSourceNode, outputNode),
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
	intermediateNode := NewNode("AndIntermediate")
	outputNode := NewNode("AndOutput")
	return outputNode, NewCustomComponent(
		"AndGate",
		[]Component{
			NewTransistor(SharedSourceNode, input1, intermediateNode),
			NewTransistor(intermediateNode, input2, outputNode),
			NewResistor(outputNode, SharedGroundNode),
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
	outputNode := NewNode("OrOutput")
	return outputNode, NewCustomComponent(
		"OrGate",
		[]Component{
			NewTransistor(SharedSourceNode, input1, outputNode),
			NewTransistor(SharedSourceNode, input2, outputNode),
			NewResistor(outputNode, SharedGroundNode),
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
	intermediateNode := NewNode("NandIntermediate")
	outputNode := NewNode("NandOutput")
	return outputNode, NewCustomComponent(
		"NandGate",
		[]Component{
			NewTransistor(outputNode, input1, intermediateNode),
			NewTransistor(intermediateNode, input2, SharedGroundNode),
			NewResistor(SharedSourceNode, outputNode),
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
