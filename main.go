package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
    MAX_DEFERS = 10
    DEBUG = false
)

func main() {
	output := NewNode("middle")
    input := NewNode("input")
	parent := "Circuit"
	components := []Component{
		NewResistor(parent, SharedSourceNode, output),
		NewTransistor(parent, output, input, SharedGroundNode),
        NewMultimeter(output),
        NewGround(input),
	}

	circuit := NewCircuit(components, MAX_DEFERS, DEBUG)
	circuit.Tick()

    width, height := int32(1280), int32(720)
    rl.InitWindow(width, height, "copooter")
}
