package main

import (
	"flag"

	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/simulation"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/view/window"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {

	usage := flag.Bool("h", false, "Show this information")
	noGui := flag.Bool("no-gui", false, "hide gui")
	noGrid := flag.Bool("no-grid", false, "hide grid")

	flag.Parse()

	if *usage {
		flag.PrintDefaults()
		return
	}

	views := []simulation.View{}

	if !*noGui {
		win, err := window.New("Swarm Simulation", 1024, !*noGrid)
		if err != nil {
			panic(err)
		}
		views = append(views, win)
	}

	sim := simulation.New(views)
	if err := sim.Start(); err != nil {
		panic(err)
	}
	defer sim.Stop()

	// detach simulation loop in background so it does not freeze the window
	go sim.Loop()

	// actually draw something
draw_loop:
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				break draw_loop
			}
		}
		sim.Draw()
	}
}
