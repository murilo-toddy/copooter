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

	var moving Component
	var wireStart *rl.Vector2

	wires := []Wire{}
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.White)

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			mouse := rl.GetMousePosition()
			for _, component := range components {
				if component.Drawable().Clicked(mouse) {
					moving = component
				}
			}
			if moving == nil {
				wireStart = &mouse
			}
		}

		if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
			if moving != nil {
				moving = nil
			}
			if wireStart != nil {
				wires = append(wires, Wire{start: *wireStart, end: rl.GetMousePosition()})
				wireStart = nil
			}
		}

		if moving != nil {
			moving.Drawable().Move(rl.GetMouseDelta())
		}

		if wireStart != nil {
			rl.DrawLineEx(*wireStart, rl.GetMousePosition(), 15, rl.Black)
		}

		for _, component := range components {
			component.Drawable().Draw()
		}

		for _, wire := range wires {
			wire.Draw()
		}

		rl.EndDrawing()
	}
}
