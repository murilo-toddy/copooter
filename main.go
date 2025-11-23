package main

import (
	"fmt"
	"strings"

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

var (
	gridCellSize                = int32(20)
	gridComponentImageSize      = gridCellSize * 8
	gridComponentFontSize       = int32(8)
	gridComponentTerminalRadius = float32(5.0)
)

var (
	toolkitSidebarRatio          = int32(10)
	toolkitComponentPadding      = int32(20)
	toolkitSidebarSize           = width / toolkitSidebarRatio
	toolkitComponentImageSize    = toolkitSidebarSize - toolkitComponentPadding
	toolkitComponentNameFontSize = int32(19)
	toolkitComponentBoxSize      = toolkitSidebarSize + toolkitComponentNameFontSize + toolkitComponentPadding
)

var (
	actionsOffset    = int32(10)
	actionButtonSize = int32(30)
)

type State int

const (
	StateIdle State = iota
	StateDragging
	StateComponentSelected
	StateTerminalSelected
	StateSimulating
)

type DrawingState struct {
	state             State
	toolkitComponents []DrawableComponent
	components        []DrawableComponent

	draggingComponent *DrawableComponent
	selectedComponent *DrawableComponent
	selectedTerminal  *int
}

func (d *DrawingState) Log() {
	var logMessage string

	var state string
	switch d.state {
	case StateIdle:
		state = "idle"
	case StateDragging:
		state = "dragging"
	case StateComponentSelected:
		state = "component-selected"
	case StateTerminalSelected:
		state = "terminal-selected"
	}
	logMessage += fmt.Sprintf("Current state: %s", state)

	if d.draggingComponent != nil {
		logMessage += fmt.Sprintf(" | dragging component %s", d.draggingComponent.Name)
	}
	if d.selectedComponent != nil {
		logMessage += fmt.Sprintf(" | selected component %s", d.selectedComponent.Name)
	}
	if d.selectedTerminal != nil {
		if d.selectedComponent == nil {
			logMessage += fmt.Sprintf(" | selected terminal %d but component nil", *d.selectedTerminal)
		} else {
			logMessage += fmt.Sprintf(" | selected terminal %d of %s", *d.selectedTerminal, d.selectedComponent.Name)
		}
	}
	fmt.Println(logMessage)
}

type DrawableConnection struct {
	component DrawableComponent
	termIndex int
}

type DrawableTerminal struct {
	OffsetX float32
	OffsetY float32
	Node    *Node
	// TODO: integrate with node flow
	connections []DrawableConnection
}

type DrawableComponent struct {
	Name         string
	Resource     rl.Texture2D
	ResourceName string
	X            int32
	Y            int32
	terminals    []DrawableTerminal
	Component
}

func addComponent(s *DrawingState, c DrawableComponent) {
	n := 1
	for _, existingComponent := range s.components {
		if strings.HasPrefix(existingComponent.Name, c.Name) {
			n++
		}
	}
	c.Name = fmt.Sprintf("%s %d", c.Name, n)
	s.components = append(s.components, c)
}

// Select component from toolbox
func checkToolkitComponentSelected(s *DrawingState, pos rl.Vector2) {
	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		if int32(pos.X) < toolkitSidebarSize {
			componentIndex := int32(pos.Y) / toolkitComponentBoxSize
			if componentIndex < int32(len(s.toolkitComponents)) {
				selectedComponent := s.toolkitComponents[componentIndex]
				selectedComponent.Resource = loadTextureWithSize(selectedComponent.ResourceName, gridComponentImageSize, gridComponentImageSize)
				// enter dragging state
				s.draggingComponent = &selectedComponent
				s.state = StateDragging
			}
		}
	}
}

// Drop toolbox component into schematic
func checkComponentDropped(s *DrawingState, pos rl.Vector2) {
	if rl.IsMouseButtonReleased(rl.MouseButtonLeft) && s.draggingComponent != nil {
		if isInsideSchematic(pos) {
			x, y := snapToGrid(pos)
			component := *s.draggingComponent
			component.X = x
			component.Y = y
			addComponent(s, component)
		}
		// release dragging component and reset state to idle
		s.draggingComponent = nil
		s.state = StateIdle
	}
}

func checkSchematicComponentSelected(s *DrawingState, pos rl.Vector2) {
	if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		for _, component := range s.components {
			if isInsideComponent(pos, component) {
				s.selectedComponent = &component
				s.state = StateComponentSelected
			}
		}
	}
}

func checkNewComponentSelected(s *DrawingState, pos rl.Vector2) {
	if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		for _, component := range s.components {
			if isInsideComponent(pos, component) {
				s.selectedComponent = &component
				s.state = StateComponentSelected
				return
			}
		}
		s.selectedComponent = nil
		s.selectedTerminal = nil
		s.state = StateIdle
	}
}

func checkTerminalSelected(s *DrawingState, pos rl.Vector2) {
	if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		if s.selectedComponent != nil && isInsideComponent(pos, *s.selectedComponent) {
			for termIndex, term := range s.selectedComponent.terminals {
				if isInsideTerminal(pos, *s.selectedComponent, term) {
					s.selectedTerminal = &termIndex
					s.state = StateTerminalSelected
					return
				}
			}
		}
	}
}

func checkConnectTerminals(s *DrawingState, pos rl.Vector2) {
	if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		if s.selectedComponent == nil {
			return
		}

		selectedTerminal := &s.selectedComponent.terminals[*s.selectedTerminal]
		for _, component := range s.components {
			for termIndex, term := range component.terminals {
				if isInsideTerminal(pos, component, term) {
					selectedTerminal.Node.Connect(term.Node)
					// TODO: move to separate function
					selectedTerminal.connections = append(selectedTerminal.connections, DrawableConnection{component, termIndex})
					term.connections = append(term.connections, DrawableConnection{*s.selectedComponent, *s.selectedTerminal})
					return
				}
			}
		}
		s.selectedComponent = nil
		s.selectedTerminal = nil
		s.state = StateIdle
	}
}

// Select component from toolbox
func checkPlayButtonSelected(s *DrawingState, pos rl.Vector2) {
	if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		if isInsideSquare(pos, width-actionsOffset-actionButtonSize, actionsOffset, actionButtonSize, actionButtonSize) {
			s.state = StateSimulating
		}
	}
}

func loadTextureWithSize(resourcePath string, width, height int32) rl.Texture2D {
	image := rl.LoadImage(resourcePath)
	rl.ImageResize(image, width, height)
	return rl.LoadTextureFromImage(image)
}

func drawComponentsToolbox(drawableComponents []DrawableComponent) {
	rl.DrawRectangle(0, 0, toolkitSidebarSize, height, rl.NewColor(48, 48, 48, 255))
	for i, component := range drawableComponents {
		rl.DrawTexture(component.Resource, toolkitComponentPadding/2, int32(i)*toolkitComponentBoxSize+toolkitComponentPadding/2, rl.White)
		rl.DrawText(component.Name, toolkitComponentPadding/2, int32(i)*toolkitComponentBoxSize+toolkitSidebarSize, toolkitComponentNameFontSize, rl.White)
	}
}

func drawGridLines() {
	for i := toolkitSidebarSize + gridCellSize; i < width; i += gridCellSize {
		rl.DrawLine(i, 0, i, height, rl.NewColor(48, 48, 48, 255))
	}
	for i := int32(0); i < height; i += gridCellSize {
		rl.DrawLine(toolkitSidebarSize, i, width, i, rl.NewColor(48, 48, 48, 255))
	}
}

func getTerminalCoordinates(c DrawableComponent, t DrawableTerminal) (float32, float32) {
	return float32(c.X) + float32(gridComponentImageSize)*t.OffsetX, float32(c.Y) + float32(gridComponentImageSize)*t.OffsetY
}

func drawSchematicComponents(components []DrawableComponent) {
	for _, component := range components {
		if component.X == 0 && component.Y == 0 {
			continue
		}
		rl.DrawTexture(component.Resource, component.X, component.Y, rl.White)
		rl.DrawText(component.Name, component.X, component.Y+gridComponentImageSize, gridComponentFontSize, rl.White)
		for _, term := range component.terminals {
			termX, termY := getTerminalCoordinates(component, term)
			for _, conn := range term.connections {
				connX, connY := getTerminalCoordinates(conn.component, conn.component.terminals[conn.termIndex])
				rl.DrawLine(int32(termX), int32(termY), int32(termX), int32(connY), rl.White)
				rl.DrawLine(int32(termX), int32(connY), int32(connX), int32(connY), rl.White)
			}
		}
	}
}

func isInsideSquare(pos rl.Vector2, x, y, w, h int32) bool {
	posX, posY := int32(pos.X), int32(pos.Y)
	return posX >= x && posX <= x+w && posY >= y && posY <= y+h
}

func isInsideComponent(pos rl.Vector2, c DrawableComponent) bool {
	return isInsideSquare(pos, c.X, c.Y, gridComponentImageSize, gridComponentImageSize)
}

func isInsideTerminal(pos rl.Vector2, c DrawableComponent, term DrawableTerminal) bool {
	termCenterX, termCenterY := getTerminalCoordinates(c, term)
	r := gridComponentTerminalRadius
	return pos.X >= termCenterX-r && pos.X <= termCenterX+r &&
		pos.Y >= termCenterY-r && pos.Y <= termCenterY+r
}

func isInsideSchematic(pos rl.Vector2) bool {
	// TODO: allow adding components close to the toolkit sidebar
	return int32(pos.X) > toolkitSidebarSize+gridComponentImageSize/2
}

func snapToGrid(pos rl.Vector2) (int32, int32) {
	x := int32(pos.X) - toolkitComponentImageSize/2 - toolkitSidebarSize
	y := int32(pos.Y) - toolkitComponentImageSize/2
	return (x / gridCellSize * gridCellSize) + toolkitSidebarSize, y / gridCellSize * gridCellSize
}

func drawComponentOutline(c DrawableComponent, color rl.Color) {
	rl.DrawRectangleLines(c.X, c.Y, gridComponentImageSize, gridComponentImageSize, color)
}

func drawTerminal(c DrawableComponent, t DrawableTerminal, color rl.Color) {
	rl.DrawCircle(
		c.X+int32(float32(gridComponentImageSize)*t.OffsetX),
		c.Y+int32(float32(gridComponentImageSize)*t.OffsetY),
		gridComponentTerminalRadius,
		color,
	)
}

func drawActionButtons() {
	x := width - actionsOffset - actionButtonSize
	y := actionsOffset
	rl.DrawRectangle(x, y, actionButtonSize, actionButtonSize, rl.Green)
	playButtonSymbolOffset := int32(3)
	rl.DrawTriangle(
		rl.Vector2{X: float32(x + playButtonSymbolOffset), Y: float32(y + playButtonSymbolOffset)},
		rl.Vector2{X: float32(x + playButtonSymbolOffset), Y: float32(y + actionButtonSize - playButtonSymbolOffset)},
		rl.Vector2{X: float32(x + actionButtonSize - playButtonSymbolOffset), Y: float32(y + actionButtonSize/2)},
		rl.White,
	)
}

func NewDrawableResistor(resourceName string) DrawableComponent {
	n1, n2 := NewNode("r1"), NewNode("r2")
	return DrawableComponent{
		Name:         "Resistor",
		ResourceName: resourceName,
		Resource:     loadTextureWithSize(resourceName, toolkitComponentImageSize, toolkitComponentImageSize),
		terminals: []DrawableTerminal{
			{0.07, 0.5, n1, []DrawableConnection{}},
			{0.93, 0.5, n2, []DrawableConnection{}},
		},
		Component: NewResistor("", n1, n2),
	}
}

func NewDrawableTransistor(resourceName string) DrawableComponent {
	source, gate, drain := NewNode("source"), NewNode("gate"), NewNode("drain")
	return DrawableComponent{
		Name:         "Transistor",
		ResourceName: resourceName,
		Resource:     loadTextureWithSize(resourceName, toolkitComponentImageSize, toolkitComponentImageSize),
		terminals: []DrawableTerminal{
			{0.05, 0.5, source, []DrawableConnection{}},
			{0.6, 0.05, gate, []DrawableConnection{}},
			{0.6, 0.95, drain, []DrawableConnection{}},
		},
		Component: NewTransistor("", source, gate, drain),
	}
}

func NewDrawableMultimeter(resourceName string) DrawableComponent {
	node := NewNode("n")
	return DrawableComponent{
		Name:         "Multimeter",
		ResourceName: resourceName,
		Resource:     loadTextureWithSize(resourceName, toolkitComponentImageSize, toolkitComponentImageSize),
		terminals: []DrawableTerminal{
			{0.3, 0.5, node, []DrawableConnection{}},
		},
		Component: NewMultimeter(node),
	}
}

func NewDrawableSource(resourceName string) DrawableComponent {
	node := NewNode("n")
	return DrawableComponent{
		Name:         "Source",
		ResourceName: resourceName,
		Resource:     loadTextureWithSize(resourceName, toolkitComponentImageSize, toolkitComponentImageSize),
		terminals: []DrawableTerminal{
			{0.5, 0.05, node, []DrawableConnection{}},
		},
		Component: NewSource(node),
	}
}

func NewDrawableGround(resourceName string) DrawableComponent {
	node := NewNode("n")
	return DrawableComponent{
		Name:         "Ground",
		ResourceName: resourceName,
		Resource:     loadTextureWithSize(resourceName, toolkitComponentImageSize, toolkitComponentImageSize),
		terminals: []DrawableTerminal{
			{0.5, 0.1, node, []DrawableConnection{}},
		},
		Component: NewGround(node),
	}
}

func main() {
	rl.InitWindow(width, height, "copooter")
	defer rl.CloseWindow()

	s := DrawingState{
		state: StateIdle,
		toolkitComponents: []DrawableComponent{
			NewDrawableResistor("./resources/resistor.png"),
			NewDrawableTransistor("./resources/transistor.jpg"),
			NewDrawableSource("./resources/source.png"),
			NewDrawableGround("./resources/ground.png"),
			NewDrawableMultimeter("./resources/meter.jpg"),
		},
	}

	for _, component := range s.toolkitComponents {
		defer rl.UnloadTexture(component.Resource)
	}
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.NewColor(12, 12, 12, 255))

		mousePos := rl.GetMousePosition()

		// s.Log()
		switch s.state {
		case StateIdle:
			checkToolkitComponentSelected(&s, mousePos)
			checkSchematicComponentSelected(&s, mousePos)
		case StateDragging:
			checkComponentDropped(&s, mousePos)
		case StateComponentSelected:
			checkNewComponentSelected(&s, mousePos)
			checkTerminalSelected(&s, mousePos)
		case StateTerminalSelected:
			checkConnectTerminals(&s, mousePos)
			checkNewComponentSelected(&s, mousePos)
		}
		checkPlayButtonSelected(&s, mousePos)

		// Render
		drawGridLines()
		drawSchematicComponents(s.components)
		drawComponentsToolbox(s.toolkitComponents)
		// draw different things depending on current state
		switch s.state {
		case StateDragging:
			offset := toolkitComponentImageSize / 2
			x, y := int32(mousePos.X)-offset, int32(mousePos.Y)-offset
			rl.DrawTexture(s.draggingComponent.Resource, x, y, rl.White)
		case StateComponentSelected:
			drawComponentOutline(*s.selectedComponent, rl.Yellow)
			rl.DrawRectangleLines(s.selectedComponent.X, s.selectedComponent.Y, gridComponentImageSize, gridComponentImageSize, rl.Yellow)
			for _, term := range s.selectedComponent.terminals {
				drawTerminal(*s.selectedComponent, term, rl.Red)
			}
		case StateTerminalSelected:
			for _, component := range s.components {
				// TODO: update to check value and not reference
				for termIndex, term := range component.terminals {
					var color rl.Color
					if component.Name == s.selectedComponent.Name && termIndex == *s.selectedTerminal {
						color = rl.Blue
					} else {
						color = rl.Red
					}
					drawTerminal(component, term, color)
				}
			}
		case StateSimulating:
			c := make([]Component, len(s.components))
			for _, component := range s.components {
				c = append(c, component)
			}
			circuit := NewCircuit(c, 50, true)
			circuit.Tick()
		}
		drawActionButtons()
		rl.EndDrawing()
	}
}
