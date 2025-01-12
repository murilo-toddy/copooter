package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	width, height := int32(1280), int32(920)
	rl.InitWindow(width, height, "Copooter Sim")

	node := NewNode("")
	parent := "Circuit"
	components := []Component{
		NewResistor(parent, node, node),
		NewTransistor(parent, node, node, node),
	}
	maxDefers := 10
	debug := false
	circuit := NewCircuit(components, maxDefers, debug)
	circuit.Simulate()

	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.White)

		for _, component := range components {
			component.Render()
		}

		rl.EndDrawing()
	}
}
