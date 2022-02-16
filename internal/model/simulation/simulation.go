package simulation

import (
	"fmt"
	"log"
	"time"

	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/collision"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/entities"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/util"
)

type Simulation struct {
	configuration *Configuration
	views         []View
	entities      *entities.EntityManger
	spatial       *collision.SpatialHashmap
	barrier       *util.Barrier
	portal        *SimulationPortal
	running       bool
	finished      bool
	duration      time.Duration
	ticks         uint64
}

func New(configuration *Configuration, views []View) *Simulation {
	return &Simulation{
		configuration: configuration,
		views:         views,
		entities:      entities.Manager(),
		spatial:       collision.New(64),
		barrier:       util.NewBarrier(),
	}
}

// Start the simulation.
// This also includes initialization of all entities etc.
func (s *Simulation) Start() error {
	s.init()
	return nil
}

func (s *Simulation) init() {
	s.portal = CreateSimulationPortal(s.spatial, s.entities, s.configuration.Settings.Noise, s.configuration.Settings.FPS, s.configuration.Obstacles)

	// initialize all participants mentioned in the configuration
	for i, p := range s.configuration.Participants {
		s.addEntity(entities.Create(fmt.Sprintf("id_%v", i), p, s.portal, s.barrier))
	}

	for _, e := range s.entities.Get() {
		err := e.Start()
		if err != nil {
			panic(err)
		}
	}
	s.running = true
}

func (s *Simulation) addEntity(entity *entities.Entity) {
	if entity != nil {
		s.entities.Add(entity)
		s.spatial.Insert(entity)
	}
}

// Main loop for the simulation.
func (s *Simulation) Loop() {
	start := time.Now()
	// create a new ticker which ticks every X milliseconds
	ticker := time.NewTicker(time.Duration(s.configuration.Settings.TickLength * float64(time.Millisecond)))
	for ; s.running; <-ticker.C {
		s.tick()
	}
	s.duration = time.Since(start)
	s.finished = true
}

func (s *Simulation) tick() {
	s.portal.Update()
	start := time.Now()
	active, obstacles := s.entities.GetRunning()
	if len(active) == 0 {
		s.running = false
	}
	s.barrier.Tick(len(active) + len(obstacles))
	elapsed := time.Since(start)
	log.Printf("Tick took %s\n", elapsed)
	s.ticks++
}

// Draw current state of simulation
func (s *Simulation) Draw() {
	// draw all entities to the screen
	ents := s.entities.Get()
	ns := make([]Drawable, 0, len(ents))
	for _, e := range ents {
		ns = append(ns, e)
	}
	for _, v := range s.views {
		v.Render(ns, s.configuration.Obstacles)
	}
}

func (s *Simulation) Stop() {
	for _, e := range s.entities.Get() {
		e.Stop()
	}
}

// Print the summary about the simulation
func (s *Simulation) PrintSummary() {
	most := int64(0)
	least := int64(0)
	total := int64(0)
	participants := s.entities.Get()
	for i, e := range participants {
		collisions := e.GetCollisions()
		total += collisions
		if i == 0 {
			least = collisions
			most = collisions
		} else {
			if collisions > most {
				most = collisions
			} else if collisions < least {
				least = collisions
			}
		}
	}

	// TODO lome: Add total runtime & avg tick length to summary
	fmt.Printf("\n-----------------------------\n")
	fmt.Printf("           Summary           \n")
	fmt.Printf("-----------------------------\n\n")
	fmt.Printf("N° of Participants: \t%v\n", len(participants))
	fmt.Printf("Total Collisions: \t%v\n", total)
	fmt.Printf("Most Collisions: \t%v\n", most)
	fmt.Printf("Least Collisions: \t%v\n", least)
	fmt.Printf("Avg Collisions: \t%v\n\n", total/int64(len(participants)))

	fmt.Printf("Total ticks: \t\t%v\n", s.ticks)
	fmt.Printf("Duration: \t\t%v\n", s.duration)
	fmt.Printf("Avg Tick Length: \t%v\n\n", time.Duration(s.duration/time.Duration(s.ticks)))
}

func (s *Simulation) IsRunning() bool {
	return s.running
}

func (s *Simulation) IsFinished() bool {
	return s.finished
}
