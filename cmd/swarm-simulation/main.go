package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

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

	configuration := parseConfiguration(*configurationPath)

	if *noise != 0 {
		configuration.Settings.Noise = *noise
	}

	sim := simulation.New(configuration, views)
	if err := sim.Start(); err != nil {
		panic(err)
	}
	// defer sim.Stop()

	// detach simulation loop in background so it does not freeze the window
	go sim.Loop()

	// actually draw something
draw_loop:
	for sim.IsRunning() {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				break draw_loop
			}
		}
		sim.Draw()
	}
}

// Parse the simulation and print error messages, if the path or the format of the file are wrong
func parseConfiguration(configurationPath string) *simulation.Configuration {
	configuration := &simulation.Configuration{}
	// try to read file
	data, err := os.ReadFile(configurationPath)
	if err != nil {
		fmt.Printf("Not a valid path to a configuration file: '%v'\n", configurationPath)
		os.Exit(-1)
	}
	// try to parse file
	err = json.Unmarshal(data, &configuration)
	if err != nil {
		fmt.Printf("Configuration file does not have a valid format!\n")
		os.Exit(-1)
	}
	for _, p := range configuration.Participants {
		p.Script = filepath.Join(filepath.Dir(configurationPath), p.Script)
	}
	return configuration
}
