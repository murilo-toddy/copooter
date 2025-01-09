package main

import (
	"github.com/gen2brain/raylib-go/raylib"
)

func main() {
	components := []Component{}
	maxDefers := 10
	debug := false
	circuit := NewCircuit(components, maxDefers, debug)
	circuit.Simulate()

	rl.InitWindow(800, 600, "Copooter Sim")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.NewColor(18, 18, 18, 255))
		rl.EndDrawing()
	}
}
