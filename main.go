package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	MAX_DEFERS = 10
	DEBUG      = false
)

var (
	width  = int32(1280)
	height = int32(720)
)

var gridSize = int32(20)

var (
	componentsToolboxRatio = int32(10)
	componentPadding       = int32(20)
	componentSize          = width / componentsToolboxRatio
	componentImageSize     = componentSize - componentPadding
	componentNameFontSize  = int32(19)
	componentBoxSize       = componentSize + componentNameFontSize + componentPadding
)

type DrawableComponent struct {
	Name     string
	Resource rl.Texture2D
	x        int32
	y        int32
}

func loadTextureWithSize(resourcePath string, width, height int32) rl.Texture2D {
	image := rl.LoadImage(resourcePath)
	rl.ImageResize(image, width, height)
	return rl.LoadTextureFromImage(image)
}

func drawComponentsToolbox(drawableComponents []DrawableComponent) {
	rl.DrawRectangle(0, 0, componentSize, height, rl.NewColor(48, 48, 48, 255))
	for i, component := range drawableComponents {
		rl.DrawTexture(component.Resource, componentPadding/2, int32(i)*componentBoxSize+componentPadding/2, rl.White)
		rl.DrawText(component.Name, componentPadding/2, int32(i)*componentBoxSize+componentSize, componentNameFontSize, rl.White)
	}
}

func drawGridLines() {
	for i := componentSize + gridSize; i < width; i += gridSize {
		rl.DrawLine(i, 0, i, height, rl.NewColor(48, 48, 48, 255))
	}
	for i := int32(0); i < height; i += gridSize {
		rl.DrawLine(componentSize, i, width, i, rl.NewColor(48, 48, 48, 255))
	}
}

func snapToGrid(pos rl.Vector2) (int32, int32) {
	x := int32(pos.X) - componentImageSize/2 - componentSize
	y := int32(pos.Y) - componentImageSize/2
	return (x / gridSize * gridSize) + componentSize, y / gridSize * gridSize
}

func main() {
	// output := NewNode("middle")
	// input := NewNode("input")
	// parent := "Circuit"
	// components := []Component{
	// 	NewResistor(parent, SharedSourceNode, output),
	// 	NewTransistor(parent, output, input, SharedGroundNode),
	// 	NewMultimeter(output),
	// 	NewGround(input),
	// }

	// circuit := NewCircuit(components, MAX_DEFERS, DEBUG)
	// circuit.Tick()

	rl.InitWindow(width, height, "copooter")
	defer rl.CloseWindow()

	drawableComponents := []DrawableComponent{
		{
			Name:     "Resistor",
			Resource: loadTextureWithSize("./resources/resistor.png", componentImageSize, componentImageSize),
		},
		{
			Name:     "Transistor",
			Resource: loadTextureWithSize("./resources/transistor.jpg", componentImageSize, componentImageSize),
		},
	}
	for _, component := range drawableComponents {
		defer rl.UnloadTexture(component.Resource)
	}

	var components []DrawableComponent

	var draggingComponent *DrawableComponent = nil
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.NewColor(12, 12, 12, 255))

		drawComponentsToolbox(drawableComponents)
		drawGridLines()

		mousePos := rl.GetMousePosition()

		// Handle drag components from toolbox
		if rl.IsMouseButtonDown(rl.MouseButtonLeft) && draggingComponent == nil {
			if int32(mousePos.X) < componentSize {
				var selectedComponent DrawableComponent
				componentIndex := int32(mousePos.Y) / componentBoxSize
				if componentIndex < int32(len(drawableComponents)) {
					selectedComponent = drawableComponents[componentIndex]
				}
				draggingComponent = &selectedComponent
			}
		}

		if rl.IsMouseButtonUp(rl.MouseButtonLeft) && draggingComponent != nil {
			x, y := snapToGrid(mousePos)
			components = append(components, DrawableComponent{
				Name:     draggingComponent.Name,
				Resource: draggingComponent.Resource,
				x:        x,
				y:        y,
			})
			draggingComponent = nil
		}

		if draggingComponent != nil {
			rl.DrawTexture(draggingComponent.Resource, int32(mousePos.X)-componentImageSize/2, int32(mousePos.Y)-componentImageSize/2, rl.White)
		}

		for _, component := range components {
			if component.x != 0 && component.y != 0 {
				rl.DrawTexture(component.Resource, component.x, component.y, rl.White)
			}
		}

		rl.EndDrawing()
	}
}
