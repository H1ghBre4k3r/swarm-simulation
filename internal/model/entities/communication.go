package entities

import (
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/obstacles"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/util"
)

// Message to send when starting the entities
type SetupMessage struct {
	Position util.Vec2D `json:"position"`
	Radius   float64    `json:"radius"`
	SafeZone float64    `json:"safezone"`
	Vmax     float64    `json:"vmax"`
	Target   util.Vec2D `json:"target"`
	TAU      float64    `json:"tau"`
}

type SimulationMessage struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

// Message to receive from the process as a response to InformationMessage
type MovementMessage struct {
	Action  string     `json:"action"`
	Payload util.Vec2D `json:"payload"`
}

// Message to send to the process about the current information about the system
type InformationMessage struct {
	Position     util.Vec2D               `json:"position"`
	Participants []ParticipantInformation `json:"participants"`
	Obstacles    []*obstacles.Obstacle    `json:"obstacles"`
}

type ParticipantInformation struct {
	Position util.Vec2D `json:"position"`
	Velocity util.Vec2D `json:"velocity"`
	Distance float64    `json:"distance"`
	Radius   float64    `json:"radius"`
	SafeZone float64    `json:"safezone"`
}
