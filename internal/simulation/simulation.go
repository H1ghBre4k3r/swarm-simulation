package simulation

import (
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
	win, err := window.New("Hello, World!", 1024, 1024)
	if err != nil {
		return err
	}
	s.window = win

	return nil
}

func (s *Simulation) init() {
	s.entities = entities.Manager()

	s.entities.Add(entities.Create("someId", entities.Rect{
		X:      100,
		Y:      200,
		Width:  10,
		Height: 10,
	}, 0xffff0000))
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
	s.window.Destroy()
}
