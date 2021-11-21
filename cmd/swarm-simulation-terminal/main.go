package main

import (
	"flag"

	"github.com/H1ghBre4k3r/swarm-simulation/internal/simulation"
)

func main() {

	usage := flag.Bool("h", false, "Show this information")

	flag.Parse()

	if *usage {
		flag.PrintDefaults()
		return
	}

	views := []simulation.View{}

	sim := simulation.New(views)
	if err := sim.Start(); err != nil {
		panic(err)
	}
	defer sim.Stop()

	for {
		sim.Tick()
	}
}
