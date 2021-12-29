package simulation

import "github.com/H1ghBre4k3r/swarm-simulation/internal/model/util"

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
}

type Configuration struct {
	Settings     Settings
	Participants []*ParticipantSetupInformation `json:"participants"`
}
