package main

import . "github.com/gen2brain/raylib-go/raylib"

type Drawable interface {
	Draw()
}

type DrawableComponent struct {
	position      Vector2
	width, height float32
	spritePath    string
	texture       *Texture2D
}

func NewDrawable(x, y float32, spritePath string) *DrawableComponent {
	return &DrawableComponent{
		position:   Vector2{X: x, Y: y},
		spritePath: spritePath,
		texture:    nil, // only loaded on first render
	}
}

func (d *DrawableComponent) Clicked(pos Vector2) bool {
	return pos.X >= d.position.X &&
		pos.X <= d.position.X+d.width &&
		pos.Y >= d.position.Y &&
		pos.Y <= d.position.Y+d.height
}

func (d *DrawableComponent) Move(delta Vector2) {
	d.position = Vector2{X: d.position.X + delta.X, Y: d.position.Y + delta.Y}
}

func (d *DrawableComponent) Draw() {
	if d.texture == nil {
		texture := LoadTextureFromImage(LoadImage(d.spritePath))
		d.texture = &texture
		d.height = float32(texture.Height)
		d.width = float32(texture.Width)
	}
	DrawTexture(*d.texture, int32(d.position.X), int32(d.position.Y), White)
}

type Wire struct {
	start, end Vector2
}

func (w *Wire) Draw() {
	DrawLineEx(w.start, w.end, 15, Black)
}
