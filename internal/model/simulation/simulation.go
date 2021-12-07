package simulation

import (
	"time"

	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/collision"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/entities"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/util"
)

type Drawable interface {
	GetX() float64
	GetY() float64
	GetR() float64
	GetColor() uint32
	GetVelocity() entities.Velocity
}

type View interface {
	Render([]Drawable)
}

type Simulation struct {
	views    []View
	entities *entities.EntityManger
	spatial  *collision.SpatialHashmap
	barrier  *util.Barrier
}

func New(views []View) *Simulation {
	return &Simulation{
		views:    views,
		entities: entities.Manager(),
		spatial:  collision.New(64),
		barrier:  util.NewBarrier(),
	}
}

// Start the simulation.
// This also includes initialization of all entities etc.
func (s *Simulation) Start() error {
	s.init()
	return nil
}

func (s *Simulation) init() {

	insert := func(entity *entities.Entity) {
		s.spatial.Insert(entity)
	}

	remove := func(entity *entities.Entity) {
		s.spatial.Remove(entity)
	}

	s.addEntity(entities.Create("1", entities.Position{
		X: 0.1,
		Y: 0.5,
		R: 0.005,
	}, 0.005, entities.Position{
		X: 0.9,
		Y: 0.5,
	}, 0xffff0000, insert, remove, "./test.py", s.barrier))

	s.addEntity(entities.Create("2", entities.Position{
		X: 0.9,
		Y: 0.5,
		R: 0.005,
	}, 0.005, entities.Position{
		X: 0.1,
		Y: 0.5,
	}, 0xffff0000, insert, remove, "./test.py", s.barrier))

	for _, e := range s.entities.Get() {
		err := e.Start()
		if err != nil {
			panic(err)
		}
	}
}

func (s *Simulation) addEntity(entity *entities.Entity) {
	if entity != nil {
		s.entities.Add(entity)
		s.spatial.Insert(entity)
	}
}

// Main loop for the simulation.
func (s *Simulation) Loop() {
	// create a new ticker which ticks every X milliseconds
	ticker := time.NewTicker(50 * time.Millisecond)
	for ; ; <-ticker.C {
		ents := s.entities.Get()
		s.barrier.Tick(len(ents))
	}
}

func (s *Simulation) Draw() {
	// draw all entities to the screen
	ents := s.entities.Get()
	ns := make([]Drawable, 0, len(ents))
	for _, e := range ents {
		ns = append(ns, e)
	}
	for _, v := range s.views {
		v.Render(ns)
	}
}

func (s *Simulation) Stop() {
	for _, e := range s.entities.Get() {
		e.Stop()
	}
}
