package entities

import (
	"encoding/json"
	"fmt"

	"github.com/AvraamMavridis/randomcolor"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/process"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/util"
)

type UpdateFn func(*Entity)

// Basic entity type which can be renderred in SDL
type Entity struct {
	id      string
	shape   Shape
	target  util.Vec2D
	vmax    float64
	vel     util.Vec2D
	color   randomcolor.RGBColor
	insert  UpdateFn
	remove  UpdateFn
	running bool
	process *process.Process
	barrier *util.Barrier
}

type Shape struct {
	Position util.Vec2D `json:"position"`
	Radius   float64    `json:"radius"`
}

func Create(id string, shape Shape, vmax float64, target util.Vec2D, insertFn UpdateFn, removeFn UpdateFn, script string, barrier *util.Barrier) *Entity {
	p, err := process.Spawn(script)
	if err != nil {
		fmt.Printf("Cannot start process for entity '%v': %v\n", id, err.Error())
		return nil
	}
	return &Entity{
		id:      id,
		color:   randomcolor.GetRandomColorInRgb(),
		shape:   shape,
		target:  target,
		vmax:    vmax,
		insert:  insertFn,
		remove:  removeFn,
		running: false,
		process: p,
		barrier: barrier,
	}
}

func (e *Entity) Id() string {
	return e.id
}

func (e *Entity) GetX() float64 {
	return e.shape.Position.X
}

func (e *Entity) SetX(x float64) {
	e.shape.Position.X = x
}

func (e *Entity) GetY() float64 {
	return e.shape.Position.Y
}

func (e *Entity) SetY(y float64) {
	e.shape.Position.Y = y
}

func (e *Entity) GetR() float64 {
	return e.shape.Radius
}

func (e *Entity) GetColor() randomcolor.RGBColor {
	return e.color
}

func (e *Entity) Move() {
	e.remove(e)
	e.shape.Position.AddI(&e.vel)
	e.insert(e)
}

func (e *Entity) GetVelocity() util.Vec2D {
	return e.vel
}

func (e *Entity) SetVelocity(vel *util.Vec2D) {
	e.vel = *vel
}

func (e *Entity) sendSetupMessage() {
	setupInformation := SetupMessage{}
	setupInformation.Position = e.shape.Position
	setupInformation.Radius = e.shape.Radius
	setupInformation.Target = e.target
	setupInformation.Vmax = e.vmax
	setupMessage, err := json.Marshal(&setupInformation)
	if err != nil {
		panic(err)
	}
	e.process.In <- string(setupMessage)
}

func (e *Entity) Start() error {
	err := e.process.Start()
	if err != nil {
		return err
	}
	e.sendSetupMessage()
	e.running = true
	go e.loop()
	return nil
}

func (e *Entity) loop() {
	for e.running {
		// wait for barrier to drop
		e.barrier.Wait()
		// perform movement with current velocity
		e.Move()

		// send sample message to process
		information := InformationMessage{}
		information.Position = e.shape.Position
		information.Participants = []ParticipantInformation{}
		outMsg, err := json.Marshal(&information)
		if err != nil {
			panic(err)
		}
		// TODO lome: this panics, if underlying process terminates
		e.process.In <- string(outMsg)

		// receive answer message from process
		msg := <-e.process.Out
		parsed := SimulationMessage{}
		if err := json.Unmarshal([]byte(msg), &parsed); err != nil {
			panic(err)
		}

		switch parsed.Action {
		case "move":
			message := MovementMessage{}
			if err := json.Unmarshal([]byte(msg), &message); err != nil {
				panic(err)
			}
			// update current velocity
			e.SetVelocity(&util.Vec2D{
				X: message.Payload.X,
				Y: message.Payload.Y,
			})
		}
		e.barrier.Resolve()
	}
}

func (e *Entity) Stop() {
	e.running = false
}
