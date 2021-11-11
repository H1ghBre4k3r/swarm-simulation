package entities

import (
	"math"
	"time"
)

type Rect struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}
type manager interface {
	Get() []*Entity
}

// Basic entity type which can be renderred in SDL
type Entity struct {
	id      string
	rect    Rect
	color   uint32
	manager manager
	running bool
}

func Create(id string, rect Rect, color uint32, manager manager) *Entity {
	return &Entity{
		id:      id,
		color:   color,
		rect:    rect,
		manager: manager,
		running: false,
	}
}

func (e *Entity) Id() string {
	return e.id
}

func (e *Entity) GetX() int32 {
	return e.rect.X
}

func (e *Entity) SetX(x int32) {
	e.rect.X = x
}

func (e *Entity) GetY() int32 {
	return e.rect.Y
}

func (e *Entity) SetY(y int32) {
	e.rect.Y = y
}

func (e *Entity) GetWidth() int32 {
	return e.rect.Width
}

func (e *Entity) SetWidth(width int32) {
	e.rect.Width = width
}

func (e *Entity) GetHeight() int32 {
	return e.rect.Height
}

func (e *Entity) SetHeight(height int32) {
	e.rect.Height = height
}

func (e *Entity) GetColor() uint32 {
	return e.color
}

func (e *Entity) SetColor(color uint32) {
	e.color = color
}

func (e *Entity) Start() {
	e.running = true
	go e.loop()
}

func (e *Entity) loop() {
	ticker := time.NewTicker(1 * time.Millisecond)
	defer ticker.Stop()
	i := float64(e.rect.X)
	for ; e.running; <-ticker.C {
		e.rect.X = int32(math.Sin(float64(i*(math.Pi/180)))*300) + 512
		e.rect.Y = int32(math.Cos(float64(i*(math.Pi/180)))*300) + 512
		i += 0.1
	}
}

func (e *Entity) Stop() {
	e.running = false
}
