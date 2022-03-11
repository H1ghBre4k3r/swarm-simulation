package entities

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/AvraamMavridis/randomcolor"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/process"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/util"
)

// Basic entity type which can be renderred in SDL
type Entity struct {
	id           string
	shape        Shape
	target       util.Vec2D
	vmax         float64
	vel          util.Vec2D
	color        randomcolor.RGBColor
	portal       Portal
	running      bool
	process      *process.Process
	barrier      *util.Barrier
	collisions   uint64
	mutex        sync.Mutex
	ignoreFinish bool
	r            *rand.Rand
	stddev       float64
}

func Create(id string, configuration *ParticipantSetupInformation, portal Portal, barrier *util.Barrier) *Entity {
	p, err := process.Spawn(configuration.Script)
	if err != nil {
		fmt.Printf("Cannot start process for entity '%v': %v\n", id, err.Error())
		return nil
	}
	s := rand.NewSource(time.Now().UnixNano())
	return &Entity{
		id:    id,
		color: randomcolor.GetRandomColorInRgb(),
		shape: Shape{
			Position: configuration.Start,
			Radius:   configuration.Radius,
		},
		target:       configuration.Target,
		vmax:         configuration.VMax,
		portal:       portal,
		running:      false,
		process:      p,
		barrier:      barrier,
		collisions:   0,
		ignoreFinish: configuration.IgnoreFinish,
		r:            rand.New(s),
	}
}

func (e *Entity) Copy() *Entity {
	return &Entity{
		id:         e.id,
		shape:      *e.shape.Copy(),
		target:     *e.target.Copy(),
		vmax:       e.vmax,
		vel:        *e.vel.Copy(),
		color:      e.color,
		portal:     e.portal,
		running:    e.running,
		process:    e.process,
		barrier:    e.barrier,
		collisions: e.collisions,
		r:          e.r,
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

func (e *Entity) GetCollisions() uint64 {
	return e.collisions
}

func (e *Entity) Move() {
	e.portal.Remove(e)
	e.shape.Position.AddI(&e.vel)
	e.portal.Insert(e)
}

func (e *Entity) sendSetupMessage() {
	setupInformation := SetupMessage{
		Position: e.shape.Position,
		Radius:   e.shape.Radius,
		Target:   e.target,
		Vmax:     e.vmax,
		TAU:      e.portal.TAU(),
	}
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
		defer func() {
			// we need to recover in case of a timeout, where barrier.Resolve will panic
			recover()
		}()
		if r := recover(); r != nil {
			fmt.Printf("Process for entity '%v' ended unexpectetly!\n", e.id)
			e.Stop()
			e.barrier.Resolve()
		}
	}()

	for e.IsRunning() {
		// wait for barrier to drop
		e.barrier.Wait()
		e.tick()
		e.barrier.Resolve()
	}
}

func (e *Entity) tick() {
	// perform movement with current velocity
	e.Move()

	// check, if this entity is close enough to target
	if !e.ignoreFinish && e.shape.Position.Scale(-1).Add(&e.target).Length() < 1e-4 {
		e.Stop()
		return
	}

	e.sendInformationMessage()
	e.evaluteProcessMessage()
}

// Send information about current system state to this participant
func (e *Entity) sendInformationMessage() {
	// send sample message to process
	stddev := math.Abs(e.r.NormFloat64() * e.portal.Noise())
	information := InformationMessage{
		Participants: []ParticipantInformation{},
		Obstacles:    e.portal.Obstacles(),
	}

	participants := e.portal.Participants()
	for _, x := range participants {
		// check if other participant is "in sight"
		otherIsInSight := e.shape.Position.Add(x.shape.Position.Scale(-1)).Length()-(x.shape.Radius+e.shape.Radius) < (e.vmax+x.vmax)*e.portal.TAU()
		// get information about the "this" participant (needed like this, since only the portal has the 'consensus' of noised position)
		if e.id == x.id {
			if e.portal.Consensus() {
				information.Position = *x.shape.NoisedPoisition.Copy()
				information.StdDev = x.stddev
			} else {
				// if there is no consesus, we simply noise it ourself
				information.Position = *x.shape.Position.Noise(stddev)
				information.StdDev = stddev
			}
		} else if otherIsInSight {
			stddev = e.r.NormFloat64() * e.portal.Noise()
			// check for collision with other participant
			if x.shape.Position.Add(e.shape.Position.Scale(-1)).Length() < e.shape.Radius+x.shape.Radius {
				fmt.Printf("%v collides with %v\n", e.id, x.id)
				e.collisions++
			}

			// create noised information about other participant
			participantInformation := ParticipantInformation{
				// TODO lome: noise velocity aswell?
				Velocity: x.vel,
				Distance: e.shape.Position.Add(x.shape.Position.Scale(-1)).Length(),
				Radius:   x.shape.Radius,
			}

			if e.portal.Consensus() {
				participantInformation.Position = *x.shape.NoisedPoisition.Copy()
				participantInformation.StdDev = x.stddev
			} else {
				participantInformation.Position = *x.shape.Position.Noise(stddev)
				participantInformation.StdDev = stddev
			}

			information.Participants = append(information.Participants, participantInformation)
		}
	}

	// send information to participant
	outMsg, err := json.Marshal(&information)
	if err != nil {
		panic(err)
	}
	e.process.In <- string(outMsg)
}

// Decode and evaluate the message received from process
func (e *Entity) evaluteProcessMessage() {
	// TODO lome: maybe use some kind of general message handler for the process, which decodes messages, adds them to different queues and, on demand, terminates the process
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
			e.SetVelocity(vel)
		case "stop":
			// underlying process finished computation
			e.Stop()
		}
	}
}

// Stop this entity and the underlying process.
func (e *Entity) Stop() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.running = false
	e.process.Stop()
	e.SetVelocity(&util.Vec2D{X: 0, Y: 0})
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

func (e *Entity) NoisePosition(stddev float64) {
	stddev = math.Abs(stddev)
	e.shape.NoisedPoisition = *e.shape.Position.Noise(stddev)
	e.stddev = stddev
}
