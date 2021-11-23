package entities

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/process"
)

type Position struct {
	X float64
	Y float64
	R float64
}

type Obstacle interface {
	GetX() int32
	GetY() int32
	GetR() int32
}

type GetObstacles func(*Entity, int) []*Entity

type UpdateFn func(*Entity)

// Basic entity type which can be renderred in SDL
type Entity struct {
	id      string
	pos     Position
	color   uint32
	insert  UpdateFn
	remove  UpdateFn
	running bool
	process *process.Process
}

func Create(id string, position Position, color uint32, insertFn UpdateFn, removeFn UpdateFn, script string) *Entity {
	p, err := process.Spawn(script)
	if err != nil {
		fmt.Printf("Cannot start process for entity '%v': %v\n", id, err.Error())
		return nil
	}
	return &Entity{
		id:      id,
		color:   color,
		pos:     position,
		insert:  insertFn,
		remove:  removeFn,
		running: false,
		process: p,
	}
}

func (e *Entity) Id() string {
	return e.id
}

func (e *Entity) GetX() float64 {
	return e.pos.X
}

func (e *Entity) SetX(x float64) {
	e.pos.X = x
}

func (e *Entity) GetY() float64 {
	return e.pos.Y
}

func (e *Entity) SetY(y float64) {
	e.pos.Y = y
}

func (e *Entity) GetR() float64 {
	return e.pos.R
}

func (e *Entity) SetR(r float64) {
	e.pos.R = r
}

func (e *Entity) GetColor() uint32 {
	return e.color
}

func (e *Entity) SetColor(color uint32) {
	e.color = color
}

func (e *Entity) Start() error {
	err := e.process.Start()
	if err != nil {
		return err
	}
	e.running = true
	go e.loop()
	return nil
}

func (e *Entity) loop() {
	for e.running {
		msg := <-e.process.Out
		coords := strings.Split(msg, " ")
		x, err := strconv.ParseFloat(coords[0], 64)
		if err != nil {
			continue
		}
		y, err := strconv.ParseFloat(coords[1], 64)
		if err != nil {
			continue
		}
		e.remove(e)
		e.pos.X += x
		e.pos.Y += y
		e.insert(e)
	}
}

func (e *Entity) Stop() {
	e.running = false
}
