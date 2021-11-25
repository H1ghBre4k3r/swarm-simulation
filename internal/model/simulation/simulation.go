package simulation

import (
	"math"

	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/collision"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/entities"
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
}

func New(views []View) *Simulation {
	return &Simulation{
		views:    views,
		entities: entities.Manager(),
		spatial:  collision.New(64),
	}
}

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

	for i := int32(0); i < 1; i++ {
		entity := entities.Create("1", entities.Position{
			X: math.Sin(0)*0.3 + 0.5,
			Y: math.Cos(0)*0.3 + 0.5,
			R: 0.005,
		}, 0xffff0000, insert, remove, "./test.py")

		if entity != nil {
			s.entities.Add(entity)
			s.spatial.Insert(entity)
		}
	}

	for _, e := range s.entities.Get() {
		err := e.Start()
		if err != nil {
			panic(err)
		}
	}
}

func (s *Simulation) Tick() {
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
