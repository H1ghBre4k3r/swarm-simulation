package simulation

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/collision"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/entities"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/util"
)

const (
	NO_ERROR = iota
	TIMEOUT  = iota
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
	status        int
}

func New(configuration *Configuration, views []View) *Simulation {
	return &Simulation{
		configuration: configuration,
		views:         views,
		entities:      entities.Manager(),
		spatial:       collision.New(64),
		barrier:       util.NewBarrier(),
		status:        NO_ERROR,
	}
}

// Start the simulation.
// This also includes initialization of all entities etc.
func (s *Simulation) Start() error {
	s.init()
	return nil
}

func (s *Simulation) init() {
	s.portal = CreateSimulationPortal(s.spatial, s.entities, s.configuration.Settings.Noise, s.configuration.Settings.TAU, s.configuration.Obstacles, s.configuration.Settings.Consensus)

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

	// this is our timeout handling (for crashing processes)
	go func() {
		timeout := time.NewTicker(time.Duration(5 * time.Second))
		last := s.ticks
		<-timeout.C
		for ; !s.finished; <-timeout.C {
			if last == s.ticks {
				s.Stop(TIMEOUT)
				break
			}
			last = s.ticks
		}
	}()

	go func() {
		timeout := time.NewTicker(time.Duration(5 * time.Minute))
		active, _ := s.entities.GetRunning()
		last := len(active)
		<-timeout.C
		for ; !s.finished; <-timeout.C {
			active, _ := s.entities.GetRunning()
			if last == len(active) {
				s.Stop(TIMEOUT)
				break
			}
			last = len(active)
		}
	}()

	for ; s.running; <-ticker.C {
		s.tick()
	}
	s.duration = time.Since(start)
	s.finished = true
}

func (s *Simulation) tick() {
	defer func() {
		// we need to recover in case of a timeout, where barrier.Tick will panic
		recover()
	}()
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

func (s *Simulation) Stop(status int) {
	s.status = status
	for _, e := range s.entities.Get() {
		e.Stop()
	}
}

// Print the summary about the simulation
func (s *Simulation) GenerateSummary(outputPath string) {
	// do not generate summary after timeout
	if s.status == TIMEOUT {
		return
	}
	summary := GenerateSummary(s)
	fmt.Printf("%v\n", summary)

	// check, if we actually want to save to summary to a file
	if outputPath == "" {
		return
	}

	// use "${outputPath}/${exampleName}/" as folder for summary
	// this structures the output in a reasonable way
	base := filepath.Base(s.configuration.Path)
	exampleName := fmt.Sprintf("%v-%v-%v-%v", base[:len(base)-len(filepath.Ext(base))], s.configuration.Settings.TAU, s.configuration.Settings.Noise, s.configuration.Settings.Consensus)
	outputFolder := filepath.Join(outputPath, exampleName)

	// create output folder
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
			fmt.Printf("Could not create directory '%v' to output summary!\n", outputFolder)
			os.Exit(-1)
		}
	}

	// serialize summary
	out, err := json.Marshal(summary)
	if err != nil {
		fmt.Printf("Error during serialization of summary. \nError: %v\n", err.Error())
		os.Exit(-1)
	}

	// create output file
	summaryPath := filepath.Join(outputFolder, fmt.Sprintf("%v.json", time.Now().UnixMilli()))
	file, err := os.Create(summaryPath)
	if err != nil {
		fmt.Printf("Could not create file '%v' for summary. \nError: %v\n", summaryPath, err.Error())
		os.Exit(-1)
	}
	defer file.Close()

	// write to output file
	if _, err := file.Write(out); err != nil {
		fmt.Printf("Could not write summary to file '%v'. \nError: %v", summaryPath, err.Error())
		os.Exit(-1)
	}
}

func (s *Simulation) IsRunning() bool {
	return s.running
}

func (s *Simulation) IsFinished() bool {
	return s.finished
}
