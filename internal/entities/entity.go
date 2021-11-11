package entities

import (
	"math"
	"time"
)

type Position struct {
	X int32
	Y int32
	R int32
}

type Obstacle interface {
	GetX() int32
	GetY() int32
	GetR() int32
}

type GetObstacles func(int) []Obstacle

// type manager interface {
// 	Get() []Obstacle
// }

// Basic entity type which can be renderred in SDL
type Entity struct {
	id        string
	pos       Position
	color     uint32
	obstacles GetObstacles
	running   bool
}

func Create(id string, position Position, color uint32, obstacles GetObstacles) *Entity {
	return &Entity{
		id:        id,
		color:     color,
		pos:       position,
		obstacles: obstacles,
		running:   false,
	}
}

func (e *Entity) Id() string {
	return e.id
}

func (e *Entity) GetX() int32 {
	return e.pos.X
}

func (e *Entity) SetX(x int32) {
	e.pos.X = x
}

func (e *Entity) GetY() int32 {
	return e.pos.Y
}

func (e *Entity) SetY(y int32) {
	e.pos.Y = y
}

func (e *Entity) GetR() int32 {
	return e.pos.R
}

func (e *Entity) SetR(r int32) {
	e.pos.R = r
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
	i := float64(e.pos.X)
	for ; e.running; <-ticker.C {
		e.pos.X = int32(math.Sin(float64(i*(math.Pi/180)))*300) + 512
		e.pos.Y = int32(math.Cos(float64(i*(math.Pi/180)))*300) + 512
		i += 0.1
	}
}

func (e *Entity) Stop() {
	e.running = false
}
