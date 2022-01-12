package entities

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/AvraamMavridis/randomcolor"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/process"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/util"
)

// Basic entity type which can be renderred in SDL
type Entity struct {
	id         string
	shape      Shape
	target     util.Vec2D
	vmax       float64
	vel        util.Vec2D
	color      randomcolor.RGBColor
	portal     Portal
	running    bool
	process    *process.Process
	barrier    *util.Barrier
	collisions int64
	mutex      sync.Mutex
}

func Create(id string, shape Shape, vmax float64, target util.Vec2D, portal Portal, script string, barrier *util.Barrier) *Entity {
	p, err := process.Spawn(script)
	if err != nil {
		fmt.Printf("Cannot start process for entity '%v': %v\n", id, err.Error())
		return nil
	}
	return &Entity{
		id:         id,
		color:      randomcolor.GetRandomColorInRgb(),
		shape:      shape,
		target:     target,
		vmax:       vmax,
		portal:     portal,
		running:    false,
		process:    p,
		barrier:    barrier,
		collisions: 0,
	}
}

func (e *Entity) Copy() *Entity {
	return &Entity{
		id:         e.id,
		shape:      e.shape,
		target:     e.target,
		vmax:       e.vmax,
		vel:        e.vel,
		color:      e.color,
		portal:     e.portal,
		running:    e.running,
		process:    e.process,
		barrier:    e.barrier,
		collisions: e.collisions,
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

func (e *Entity) GetTarget() util.Vec2D {
	return e.target
}

func (e *Entity) GetVelocity() util.Vec2D {
	return e.vel
}

func (e *Entity) SetVelocity(vel *util.Vec2D) {
	e.vel = *vel
}

func (e *Entity) GetCollisions() int64 {
	return e.collisions
}

func (e *Entity) Move() {
	e.portal.Remove(e)
	e.shape.Position.AddI(&e.vel)
	e.portal.Insert(e)
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

	// in case of an unpredicted exit of underlying process, we still want the barrier to move on
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Process for entity '%v' ended unexpectetly!\n", e.id)
			e.Stop()
			e.barrier.Resolve()
		}
	}()

	for e.IsRunning() {
		// wait for barrier to drop
		e.barrier.Wait()
		// perform movement with current velocity
		e.Move()

		// send sample message to process
		information := InformationMessage{}
		information.Position = e.shape.Position
		information.Participants = []ParticipantInformation{}

		participants := e.portal.Participants()
		for _, x := range participants {
			if e.id != x.id {
				// check for collision with other participant
				if x.shape.Position.Add(e.shape.Position.Scale(-1)).Length() < e.shape.Radius+x.shape.Radius {
					fmt.Printf("%v collides with %v\n", e.id, x.id)
					e.collisions++
				}

				// create noised information about other participant
				information.Participants = append(information.Participants, ParticipantInformation{
					Position: *x.shape.Position.Noise(e.portal.Noise()),
					Velocity: *x.vel.Noise(e.portal.Noise()),
					Distance: e.shape.Position.Add(x.shape.Position.Scale(-1)).Length(),
					Radius:   x.shape.Radius,
				})
			}
		}

		// send information to participant
		outMsg, err := json.Marshal(&information)
		if err != nil {
			panic(err)
		}
		e.process.In <- string(outMsg)

		// receive answer message from process
		msg := <-e.process.Out
		parsed := SimulationMessage{}
		if err := json.Unmarshal([]byte(msg), &parsed); err != nil {
			e.Stop()
		} else {
			switch parsed.Action {
			case "move":
				// a "simple" move action
				message := MovementMessage{}
				if err := json.Unmarshal([]byte(msg), &message); err != nil {
					panic(err)
				}
				vel := &util.Vec2D{
					X: message.Payload.X,
					Y: message.Payload.Y,
				}
				// update current velocity
				e.SetVelocity(vel.NoiseI(e.portal.Noise()))
			case "stop":
				// underlying process finished computation
				e.Stop()
			}
		}
		e.barrier.Resolve()
	}
}

// Stop this entity and the underlying process.
func (e *Entity) Stop() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.running = false
	e.process.Stop()
}

// Check, if this entity AND the underlying process are still running.
// If the underlying process terminated, terminate this entity aswell.
func (e *Entity) IsRunning() bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.process.IsRunning() {
		return e.running
	}
	e.running = false
	return false
}
