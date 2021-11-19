package main

import (
	"flag"

	"github.com/H1ghBre4k3r/swarm-simulation/internal/simulation"
)

func main() {

	usage := flag.Bool("h", false, "Show this information")
	// noGui := flag.Bool("no-gui", false, "hide gui")

	flag.Parse()

	if *usage {
		flag.PrintDefaults()
		return
	}

	sim := simulation.New()
	if err := sim.Start(); err != nil {
		panic(err)
	}
	defer sim.Stop()

	sim.Loop()

	// 	p, err := process.Spawn("../../test.py")

	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	p.In <- "'Hello, World!'\n"
	// 	// time.Sleep(30 * time.Second)
	// 	println(<-p.Out)
	// 	p.Stop()
}
