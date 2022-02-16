package simulation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/entities"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/obstacles"
)

type Settings struct {
	TickLength float64 `json:"tickLength"`
	Noise      float64 `json:"noise"`
	FPS        uint64  `json:"fps"`
}

type Configuration struct {
	Path         string
	Settings     Settings
	Participants []*entities.ParticipantSetupInformation `json:"participants"`
	Obstacles    []*obstacles.Obstacle                   `json:"obstacles"`
}

// Parse the simulation and print error messages, if the path or the format of the file are wrong
func ParseConfigurationFrom(path string) *Configuration {
	configuration := &Configuration{}
	// try to read file
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Not a valid path to a configuration file: '%v'\n", path)
		os.Exit(-1)
	}
	// try to parse file
	err = json.Unmarshal(data, &configuration)
	if err != nil {
		fmt.Printf("Configuration file does not have a valid format!\n")
		os.Exit(-1)
	}

	configuration.Path = path

	for _, p := range configuration.Participants {
		p.Script = filepath.Join(filepath.Dir(path), p.Script)
	}

	if configuration.Settings.TickLength <= 0 {
		configuration.Settings.TickLength = 1
	}

	if configuration.Settings.FPS <= 0 {
		configuration.Settings.FPS = 1
	}
	if len(configuration.Obstacles) == 0 {
		configuration.Obstacles = []*obstacles.Obstacle{}
	}
	return configuration
}
