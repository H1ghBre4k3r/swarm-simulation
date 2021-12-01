package entities

type SimulationMessage struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

// Message containing the velocity information of the entity/process
type MovementPayload struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Message to receive from the process as a response to InformationMessage
type MovementMessage struct {
	Action  string          `json:"action"`
	Payload MovementPayload `json:"payload"`
}

// Message to send to the process about the current information about the system
type InformationMessage struct {
	Position Position `json:"position"`
}
