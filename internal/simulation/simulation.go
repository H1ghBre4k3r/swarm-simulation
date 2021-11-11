package simulation

import (
	"github.com/H1ghBre4k3r/swarm-simulation/internal/entities"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/window"
	"github.com/veandco/go-sdl2/sdl"
)

type Simulation struct {
	window *window.Window
}

func New() *Simulation {
	return &Simulation{}
}

func (s *Simulation) Start() error {
	win, err := window.New("Hello, World!", 1024, 1024)
	if err != nil {
		return err
	}
	s.window = win

	return nil
}

func (s *Simulation) Loop() {
	ents := []entities.Entity{*entities.New("randomId")}

main_loop:
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				break main_loop
			}
		}
		// draw all entities to the screen
		ns := make([]window.Drawable, 0, len(ents))
		for _, e := range ents {
			ns = append(ns, &e)
		}
		s.window.Render(ns)
		println("Render")
	}
}

func (s *Simulation) Stop() {
	s.window.Destroy()
}
