package entities

type SimulationMessage struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

type MovementPayload struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type MovementMessage struct {
	Action  string          `json:"action"`
	Payload MovementPayload `json:"payload"`
}
