package main

func NewSimpleAdder(input1, input2 *Node) (out, carry *Node, adder *CustomComponent) {
	out, xorGate := NewXorGate(input1, input2)
	carry, andGate := NewAndGate(input1, input2)
	adder = NewCustomComponent(
		"SimpleAdder",
		[]Component{xorGate, andGate},
		[]*Node{input1, input2},
	)
	return
}

func NewFullAdder(input1, input2, carryIn *Node) (out, carry *Node, adder *CustomComponent) {
	intermediateXorOut, intermediateXorGate := NewXorGate(input1, input2)
	out, xorGate := NewXorGate(intermediateXorOut, carryIn)

	intermediateAnd1Out, intermediateAnd1Gate := NewAndGate(input1, input2)
	intermediateAnd2Out, intermediateAnd2Gate := NewAndGate(input1, carryIn)
	intermediateAnd3Out, intermediateAnd3Gate := NewAndGate(input2, carryIn)

	intermediateOrOut, intermediateOrGate := NewOrGate(intermediateAnd1Out, intermediateAnd2Out)
	carry, andGate := NewOrGate(intermediateOrOut, intermediateAnd3Out)

	adder = NewCustomComponent(
		"FullAdder",
		[]Component{
			intermediateXorGate,
			xorGate,
			intermediateAnd1Gate,
			intermediateAnd2Gate,
			intermediateAnd3Gate,
			intermediateOrGate,
			andGate,
		},
		[]*Node{input1, input2, carryIn},
	)
	return
}

func NewAdderSubtractor(input1, input2, carryIn, operation *Node) (out, carry *Node, component *CustomComponent) {
	xorOut, xorGate := NewXorGate(input2, operation)
	out, carry, adder := NewFullAdder(input1, xorOut, carryIn)
	component = NewCustomComponent(
		"AdderSubtractor",
		[]Component{xorGate, adder},
		[]*Node{input1, input2, carryIn, operation},
	)
	return
}
