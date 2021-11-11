package simulation

import (
	"fmt"
	"math"

	"github.com/H1ghBre4k3r/swarm-simulation/internal/entities"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/window"
	"github.com/veandco/go-sdl2/sdl"
)

type Simulation struct {
	window   *window.Window
	entities *entities.EntityManger
}

func New() *Simulation {
	return &Simulation{}
}

func (s *Simulation) Start() error {
	s.init()
	win, err := window.New("Swarm Simulation", 1024, 1024)
	if err != nil {
		return err
	}
	s.window = win

	return nil
}

func (s *Simulation) init() {
	s.entities = entities.Manager()

	for i := int32(0); i < 100; i++ {
		s.entities.Add(entities.Create(fmt.Sprintf("%v", i), entities.Position{
			X: int32((float64(360) / float64(100)) * float64(i)),
			Y: 0,
			R: 5,
		}, 0xffff0000, func(entity *entities.Entity, distance int) []*entities.Entity {
			// TODO lome: we neeeeed (!) performance improvements
			ents := s.entities.Get()
			in_dist := []*entities.Entity{}

			for _, target := range ents {
				// calculate, if this entity is within range of our querying entity
				if target.Id() != entity.Id() && math.Sqrt(math.Pow(float64(target.GetX())-float64(entity.GetX()), 2)+math.Pow(float64(target.GetY())-float64(entity.GetY()), 2)) < float64(distance) {
					in_dist = append(in_dist, target)
				}
			}

			return in_dist
		}))
	}

	for _, e := range s.entities.Get() {
		e.Start()
	}
}

func (s *Simulation) Loop() {

main_loop:
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				break main_loop
			}
		}
		// draw all entities to the screen
		ents := s.entities.Get()
		ns := make([]window.Drawable, 0, len(ents))
		for _, e := range ents {
			ns = append(ns, e)
		}
		s.window.Render(ns)
	}
}

func (s *Simulation) Stop() {
	for _, e := range s.entities.Get() {
		e.Stop()
	}
	s.window.Destroy()
}
