package entities

import (
	"encoding/json"
	"fmt"

	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/process"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/util"
)

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
	target  Position
	vmax    float64
	vel     Velocity
	color   uint32
	insert  UpdateFn
	remove  UpdateFn
	running bool
	process *process.Process
	barrier *util.Barrier
}

func Create(id string, position Position, vmax float64, target Position, color uint32, insertFn UpdateFn, removeFn UpdateFn, script string, barrier *util.Barrier) *Entity {
	p, err := process.Spawn(script)
	if err != nil {
		fmt.Printf("Cannot start process for entity '%v': %v\n", id, err.Error())
		return nil
	}
	return &Entity{
		id:      id,
		color:   color,
		pos:     position,
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

func (e *Entity) Move() {
	e.remove(e)
	e.pos.Move(&e.vel)
	e.insert(e)
}

func (e *Entity) GetVelocity() Velocity {
	return e.vel
}

func (e *Entity) SetVelocity(vel *Velocity) {
	e.vel = *vel
}

func (e *Entity) sendSetupMessage() {
	setupInformation := SetupMessage{}
	setupInformation.Position = e.pos
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
		information.Position = e.pos
		information.Participants = []ParticipantInformation{}
		outMsg, err := json.Marshal(&information)
		if err != nil {
			panic(err)
		}
		// TODO lome: this panics, if underlying process terminates
		e.process.In <- string(outMsg)
		fmt.Println(string(outMsg))

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
			e.SetVelocity(&Velocity{
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
