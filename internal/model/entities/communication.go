package entities

// Message to send when starting the entities
type SetupMessage struct {
	Position Position `json:"position"`
	Vmax     float64  `json:"vmax"`
	Target   Position `json:"target"`
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
	Position     Position                 `json:"position"`
	Participants []ParticipantInformation `json:"participants"`
}

type ParticipantInformation struct {
	Position Position `json:"position"`
	Velocity Velocity `json:"velocity"`
}
