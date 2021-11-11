package main

import (
	"github.com/H1ghBre4k3r/swarm-simulation/internal/simulation"
)

func main() {
	sim := simulation.New()
	if err := sim.Start(); err != nil {
		panic(err)
	}
	defer sim.Stop()

	sim.Loop()
}
