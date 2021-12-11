package entities

import "github.com/H1ghBre4k3r/swarm-simulation/internal/model/util"

// Message to send when starting the entities
type SetupMessage struct {
	Position util.Vec2D `json:"position"`
	Radius   float64    `json:"radius"`
	Vmax     float64    `json:"vmax"`
	Target   util.Vec2D `json:"target"`
}

type SimulationMessage struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

// Message containing the velocity information of the entity/process
type VectorPayload struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Message to receive from the process as a response to InformationMessage
type MovementMessage struct {
	Action  string        `json:"action"`
	Payload VectorPayload `json:"payload"`
}

// Message to send to the process about the current information about the system
type InformationMessage struct {
	Position     util.Vec2D               `json:"position"`
	Participants []ParticipantInformation `json:"participants"`
}

type ParticipantInformation struct {
	Position util.Vec2D `json:"position"`
	Velocity util.Vec2D `json:"velocity"`
	Distance float64    `json:"distance"`
	Radius   float64    `json:"radius"`
}
