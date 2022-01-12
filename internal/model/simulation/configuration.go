package simulation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/util"
)

type ParticipantSetupInformation struct {
	Start  util.Vec2D `json:"start"`
	Radius float64    `json:"radius"`
	VMax   float64    `json:"vmax"`
	Target util.Vec2D `json:"target"`
	Script string     `json:"script"`
}

type Settings struct {
	TickLength int64   `json:"tickLength"`
	Noise      float64 `json:"noise"`
	FPS        uint64  `json:"fps"`
}

type Configuration struct {
	Settings     Settings
	Participants []*ParticipantSetupInformation `json:"participants"`
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
	for _, p := range configuration.Participants {
		p.Script = filepath.Join(filepath.Dir(path), p.Script)
	}

	if configuration.Settings.FPS == 0 {
		configuration.Settings.FPS = 1
	}
	return configuration
}
