package main

import (
	"github.com/H1ghBre4k3r/swarm-simulation/internal/window"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	window, err := window.New("Hello, World!", 1024, 1024)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	rect := &sdl.Rect{
		X: 0,
		Y: 0,
		W: 200,
		H: 200,
	}

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break

			default:
				window.Render(rect)
				println("Render")
			}
		}
	}
}
