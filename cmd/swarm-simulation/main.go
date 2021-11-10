package main

import (
	"github.com/H1ghBre4k3r/swarm-simulation/internal/entities"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/window"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	w, err := window.New("Hello, World!", 1024, 1024)
	if err != nil {
		panic(err)
	}
	defer w.Destroy()

	ents := []entities.Entity{*entities.New("randomId")}

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break

			default:
				ns := make([]window.Drawable, 0, len(ents))
				for _, e := range ents {
					ns = append(ns, &e)
				}
				w.Render(ns)
				println("Render")
			}
		}
	}
}
