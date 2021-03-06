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
	configurationPath := flag.String("c", "", "configuration file for the simulation")
	noise := flag.Float64("n", 0, "noise for the data send to participants")
	consensus := flag.Bool("consensus", false, "flag for indicating consensus between all participants regardning the noises values")
	output := flag.String("o", "", "path to folder for outputting summary of simulation")
	tau := flag.Float64("t", 0, "tau for the simulation")
	flag.Parse()

	if *usage {
		flag.PrintDefaults()
		return
	}
	configuration := simulation.ParseConfigurationFrom(*configurationPath)
	if *noise != 0 {
		configuration.Settings.Noise = *noise
	}
	if *tau != 0 {
		configuration.Settings.TAU = *tau
	}
	configuration.Settings.Consensus = *consensus

	views := []simulation.View{}

	if !*noGui {
		win, err := window.New("Swarm Simulation", 1024, !*noGrid)
		if err != nil {
			panic(err)
		}
		views = append(views, win)
	}

	sim := simulation.New(configuration, views)
	if err := sim.Start(); err != nil {
		panic(err)
	}
	defer sim.GenerateSummary(*output)

	// detach simulation loop in background so it does not freeze the window
	go sim.Loop()

	// actually draw something
draw_loop:
	for !sim.IsFinished() {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				break draw_loop
			}
		}
		sim.Draw()
	}
}
