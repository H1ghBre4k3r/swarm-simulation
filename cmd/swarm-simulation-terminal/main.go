package main

import (
	"flag"

	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/simulation"
)

func main() {

	usage := flag.Bool("h", false, "Show this information")
	configurationPath := flag.String("c", "", "configuration file for the simulation")
	noise := flag.Float64("n", 0, "noise for the data send to participants")
	output := flag.String("o", "", "path to folder for outputting summary of simulation")
	flag.Parse()

	if *usage {
		flag.PrintDefaults()
		return
	}
	configuration := simulation.ParseConfigurationFrom(*configurationPath)
	if *noise != 0 {
		configuration.Settings.Noise = *noise
	}

	views := []simulation.View{}

	sim := simulation.New(configuration, views)
	if err := sim.Start(); err != nil {
		panic(err)
	}
	defer sim.GenerateSummary(*output)

	sim.Loop()
}
