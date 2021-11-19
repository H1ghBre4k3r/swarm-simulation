package simulation

import (
	"github.com/H1ghBre4k3r/swarm-simulation/internal/collision"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/entities"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/window"
	"github.com/veandco/go-sdl2/sdl"
)

type Simulation struct {
	window   *window.Window
	entities *entities.EntityManger
	spatial  *collision.SpatialHashmap
}

func New() *Simulation {
	return &Simulation{
		entities: entities.Manager(),
		spatial:  collision.New(64),
	}
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

	insert := func(entity *entities.Entity) {
		s.spatial.Insert(entity)
	}

	remove := func(entity *entities.Entity) {
		s.spatial.Remove(entity)
	}

	for i := int32(0); i < 1; i++ {
		entity := entities.Create("1", entities.Position{
			R: 5,
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
