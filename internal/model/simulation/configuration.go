package simulation

import "github.com/H1ghBre4k3r/swarm-simulation/internal/model/util"

type ParticipantSetupInformation struct {
	Start  util.Vec2D `json:"start"`
	Radius float64    `json:"radius"`
	VMax   float64    `json:"vmax"`
	Target util.Vec2D `json:"target"`
	Script string     `json:"script"`
}

type Configuration struct {
	Participants []*ParticipantSetupInformation `json:"participants"`
}
