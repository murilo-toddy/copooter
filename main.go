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
	width  = int32(1920)
	height = int32(1080)
)

var (
	gridCellSize                = int32(10)
	gridComponentImageSize      = gridCellSize * 10
	gridComponentFontSize       = int32(8)
	gridComponentTerminalRadius = float32(5.0)
	gridWireWidth               = int32(4)
)

var (
	toolkitSidebarRatio          = int32(15)
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
	StateNodeSelected
	StateSimulating
)

type DrawingState struct {
	state             State
	toolkitComponents []DrawableComponent
	components        []DrawableComponent
	nextComponentID   int

	draggingComponent *DrawableComponent
	selectedComponent *DrawableComponent
	selectedNode      *int
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
	case StateNodeSelected:
		state = "node-selected"
	}
	logMessage += fmt.Sprintf("Current state: %s", state)

	if d.draggingComponent != nil {
		logMessage += fmt.Sprintf(" | dragging component %s", d.draggingComponent.Name)
	}
	if d.selectedComponent != nil {
		logMessage += fmt.Sprintf(" | selected component %s", d.selectedComponent.Name)
	}
	if d.selectedNode != nil {
		if d.selectedComponent == nil {
			logMessage += fmt.Sprintf(" | selected node %d but component nil", *d.selectedNode)
		} else {
			logMessage += fmt.Sprintf(" | selected node %d of %s", *d.selectedNode, d.selectedComponent.Name)
		}
	}
	fmt.Println(logMessage)
}

type DrawableComponent struct {
	ID   string
	Name string
	// TODO: this is not a good way of storing resources for components
	// that may change based on state, such as input or meter
	idleResource         rl.Texture2D
	idleResourceName     string
	selectedResource     rl.Texture2D
	selectedResourceName string
	X                    int32
	Y                    int32
	Component
}

func createComponent(c *DrawableComponent) {
	// TODO: select component based on type enum instead of name
	// At this point, name is already updated to contain index,
	// this should also be a field inside DrawableComponent struct
	// type DrawableComponent struct {
	//     componentType Enum
	//     index         int
	// }

	// TODO: ideally we would need to set only the component
	// and the terminals would be derived automatically

	// TODO: remove extra mumbo jumbo once DrawableComponent and Component are integrated
	if strings.HasPrefix(c.Name, "Resistor") {
		resistor := NewResistor(c.Name, nil, nil)
		resistor.RenderData.X = c.X
		resistor.RenderData.Y = c.Y

		c.Component = resistor
		return
	} else if strings.HasPrefix(c.Name, "Transistor") {
		transistor := NewTransistor(c.Name, nil, nil, nil)
		transistor.RenderData.X = c.X
		transistor.RenderData.Y = c.Y
		c.Component = transistor
		return
	} else if strings.HasPrefix(c.Name, "Multimeter") {
		meter := NewMultimeter(c.Name, nil)
		meter.RenderData.X = c.X
		meter.RenderData.Y = c.Y
		c.Component = meter
		return
	} else if strings.HasPrefix(c.Name, "Ground") {
		ground := NewGround(c.Name, nil)
		ground.RenderData.X = c.X
		ground.RenderData.Y = c.Y
		c.Component = ground
		return
	} else if strings.HasPrefix(c.Name, "Source") {
		source := NewSource(c.Name, nil)
		source.RenderData.X = c.X
		source.RenderData.Y = c.Y
		c.Component = source
		return
	} else if strings.HasPrefix(c.Name, "Input") {
		input := NewInput(c.Name, nil, Off)
		input.RenderData.X = c.X
		input.RenderData.Y = c.Y
		c.Component = input
		return
	}
	fmt.Println("Component type not found for ", c.Name)
}

func addComponent(s *DrawingState, c DrawableComponent) {
	n := 1
	for _, existingComponent := range s.components {
		if strings.HasPrefix(existingComponent.Name, c.Name) {
			n++
		}
	}

	c.Name = fmt.Sprintf("%s %d", c.Name, n)

	newComponent := c.Component.Clone()

	// set x,y coordinates of new component
	newR := newComponent.GetRenderData()
	newR.X = c.X
	newR.Y = c.Y

	c.Component = newComponent

	setComponentID(s, &c)
	s.components = append(s.components, c)
}

// Select component from toolbox
func checkToolkitComponentSelected(s *DrawingState, pos rl.Vector2) {
	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		if int32(pos.X) < toolkitSidebarSize {
			componentIndex := int32(pos.Y) / toolkitComponentBoxSize
			if componentIndex < int32(len(s.toolkitComponents)) {
				selectedComponent := s.toolkitComponents[componentIndex]
				selectedComponent.idleResource = loadTextureWithSize(selectedComponent.idleResourceName, gridComponentImageSize, gridComponentImageSize)
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
		s.selectedNode = nil
		s.state = StateIdle
	}
}

func checkNodeSelected(s *DrawingState, pos rl.Vector2) {
	if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		if s.selectedComponent != nil && isInsideComponent(pos, *s.selectedComponent) {
			for nodeIndex, node := range s.selectedComponent.Nodes() {
				if isInsideNode(pos, *s.selectedComponent, node) {
					s.selectedNode = &nodeIndex
					s.state = StateNodeSelected
					return
				}
			}
		}
	}
}

func checkChangeInputComponentState(s *DrawingState) {
	if rl.IsKeyPressed(rl.KeyEnter) {
		terminal, ok := s.selectedComponent.Component.(*Terminal)
		if !ok {
			return
		}

		var newState NodeState
		switch terminal.Node.State {
		case Off:
			newState = On
		case On:
			newState = Off
		}
		terminal.state = newState
	}
}

func checkConnectNodes(s *DrawingState, pos rl.Vector2) {
	if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
		if s.selectedComponent == nil {
			fmt.Println("Attempting to connect terminal but selected component is nil")
			return
		}
		selectedTerminal := s.selectedComponent.Nodes()[*s.selectedNode]
		for _, component := range s.components {
			for _, term := range component.Nodes() {
				if isInsideNode(pos, component, term) {
					selectedTerminal.Connect(term)
					return
				}
			}
		}
		s.selectedComponent = nil
		s.selectedNode = nil
		s.state = StateIdle
	}
}

func checkRemoveConnections(s *DrawingState) {
	if rl.IsKeyPressed(rl.KeyD) {
		s.selectedComponent.Nodes()[*s.selectedNode].DisconnectAll()
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

func loadTextureWithSize(resourcePath string, width, height int32) (t rl.Texture2D) {
	if resourcePath != "" {
		image := rl.LoadImage(resourcePath)
		rl.ImageResize(image, width, height)
		t = rl.LoadTextureFromImage(image)
	}
	return
}

func drawComponentsToolbox(drawableComponents []DrawableComponent) {
	rl.DrawRectangle(0, 0, toolkitSidebarSize, height, rl.NewColor(48, 48, 48, 255))
	for i, component := range drawableComponents {
		rl.DrawTexture(component.idleResource, toolkitComponentPadding/2, int32(i)*toolkitComponentBoxSize+toolkitComponentPadding/2, rl.White)
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

func getTerminalCoordinates(n *Node) (float32, float32) {
	r := n.Parent.GetRenderData()
	return float32(r.X) + float32(gridComponentImageSize)*n.OffsetX, float32(r.Y) + float32(gridComponentImageSize)*n.OffsetY
}

func minAndDist(v1, v2 int32) (int32, int32) {
	if v1 < v2 {
		return v1, v2 - v1
	}
	return v2, v1 - v2
}

func drawWire(fromX, fromY, toX, toY int32, color rl.Color) {
	if fromX == toX {
		startY, dy := minAndDist(fromY, toY)
		rl.DrawRectangle(fromX-gridWireWidth/2, startY-gridWireWidth/2, gridWireWidth, dy+gridWireWidth, color)
	} else if fromY == toY {
		startX, dx := minAndDist(fromX, toX)
		rl.DrawRectangle(startX-gridWireWidth/2, fromY-gridWireWidth/2, dx+gridWireWidth, gridWireWidth, color)
	}
}

func drawSchematicComponents(components []DrawableComponent) {
	for _, component := range components {
		if component.X == 0 && component.Y == 0 {
			continue
		}
		rl.DrawTexture(component.idleResource, component.X, component.Y, rl.White)
		rl.DrawText(component.Name, component.X, component.Y+gridComponentImageSize, gridComponentFontSize, rl.White)
		for _, term := range component.Nodes() {
			termX, termY := getTerminalCoordinates(term)
			var color rl.Color
			switch term.State {
			case Off:
				color = rl.White
			case On:
				color = rl.Yellow
			case Undefined:
				color = rl.White
			default:
				panic("unreachable state")
			}
			for _, conn := range term.connections {
				connX, connY := getTerminalCoordinates(conn)
				if termX > connX {
					drawWire(int32(termX), int32(termY), int32(termX), int32(connY), color)
					drawWire(int32(termX), int32(connY), int32(connX), int32(connY), color)
				} else {
					drawWire(int32(termX), int32(termY), int32(connX), int32(termY), color)
					drawWire(int32(connX), int32(termY), int32(connX), int32(connY), color)
				}
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

func isInsideNode(pos rl.Vector2, c DrawableComponent, term *Node) bool {
	termCenterX, termCenterY := getTerminalCoordinates(term)
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

func drawTerminal(c DrawableComponent, n *Node, color rl.Color) {
	rl.DrawCircle(
		c.X+int32(float32(gridComponentImageSize)*n.OffsetX),
		c.Y+int32(float32(gridComponentImageSize)*n.OffsetY),
		gridComponentTerminalRadius,
		color,
	)
}

func drawPlayButton() {
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

func setComponentID(s *DrawingState, c *DrawableComponent) {
	c.ID = fmt.Sprintf("%d", s.nextComponentID)
	s.nextComponentID += 1
}

func NewDrawableComponent(name string, idleResourceName string, selectedResourceName string, component Component) DrawableComponent {
	return DrawableComponent{
		Name:                 name,
		idleResourceName:     idleResourceName,
		idleResource:         loadTextureWithSize(idleResourceName, toolkitComponentImageSize, toolkitComponentImageSize),
		selectedResource:     loadTextureWithSize(selectedResourceName, toolkitComponentImageSize, toolkitComponentImageSize),
		selectedResourceName: selectedResourceName,
		Component:            component,
	}
}

func main() {
	rl.InitWindow(width, height, "copooter")
	defer rl.CloseWindow()

	s := DrawingState{
		state: StateIdle,
		toolkitComponents: []DrawableComponent{
			NewDrawableComponent(
				"Resistor",
				"./resources/resistor.png", "./resources/resistor-selected.png",
				NewResistorFromNodes(&Node{OffsetX: 0.0, OffsetY: 0.5}, &Node{OffsetX: 1.0, OffsetY: 0.5}),
			),
			NewDrawableComponent(
				"Transistor",
				"./resources/transistor.jpg", "",
				NewTransistorFromNodes(
					&Node{OffsetX: 0.6, OffsetY: 0.05},
					&Node{OffsetX: 0.05, OffsetY: 0.5},
					&Node{OffsetX: 0.6, OffsetY: 0.95},
				),
			),
			NewDrawableComponent(
				"Source",
				"./resources/source.png", "",
				NewSourceFromNodes(&Node{OffsetX: 0.5, OffsetY: 0.05}),
			),
			NewDrawableComponent(
				"Ground",
				"./resources/ground.png", "",
				NewGroundFromNodes(&Node{OffsetX: 0.5, OffsetY: 0.05}),
			),
			NewDrawableComponent(
				"Multimeter",
				"./resources/meter.jpg", "",
				NewMultimeterFromNodes(&Node{OffsetX: 0.3, OffsetY: 0.5}),
			),
			NewDrawableComponent(
				"Input",
				"./resources/input.jpg", "",
				NewInputFromNodes(&Node{OffsetX: 0.7, OffsetY: 0.5}, Off),
			),
		},
	}

	for _, component := range s.toolkitComponents {
		defer rl.UnloadTexture(component.idleResource)
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
			checkNodeSelected(&s, mousePos)
			checkChangeInputComponentState(&s)
		case StateNodeSelected:
			checkConnectNodes(&s, mousePos)
			checkRemoveConnections(&s)
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
			rl.DrawTexture(s.draggingComponent.idleResource, x, y, rl.White)
		case StateComponentSelected:
			if s.selectedComponent.selectedResourceName != "" {
				if s.selectedComponent.selectedResource.ID <= 0 {
					s.selectedComponent.selectedResource = loadTextureWithSize(
						s.selectedComponent.selectedResourceName,
						gridComponentImageSize,
						gridComponentImageSize,
					)
				}
				rl.DrawTexture(s.selectedComponent.selectedResource, s.selectedComponent.X, s.selectedComponent.Y, rl.White)
			} else {
				drawComponentOutline(*s.selectedComponent, rl.Yellow)
				rl.DrawRectangleLines(s.selectedComponent.X, s.selectedComponent.Y, gridComponentImageSize, gridComponentImageSize, rl.Yellow)
			}
			for _, term := range s.selectedComponent.Nodes() {
				drawTerminal(*s.selectedComponent, term, rl.Red)
			}
		case StateNodeSelected:
			for _, component := range s.components {
				for termIndex, term := range component.Nodes() {
					var color rl.Color
					if component.ID == s.selectedComponent.ID && termIndex == *s.selectedNode {
						color = rl.Blue
					} else {
						color = rl.Red
					}
					drawTerminal(component, term, color)
				}
			}
		case StateSimulating:
			components := make([]Component, len(s.components))
			for i, component := range s.components {
				component.Component.Reset()
				components[i] = component.Component
			}
			circuit := NewCircuit(components, 50, true)
			if err := circuit.Tick(); err != nil {
				fmt.Println("Failed to run circuit: ", err.Error())
			}
			s.state = StateIdle
		}
		drawPlayButton()
		rl.EndDrawing()
	}
}
