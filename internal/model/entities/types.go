package entities

import "github.com/H1ghBre4k3r/swarm-simulation/internal/model/util"

type UpdateFn func(*Entity)

type Portal interface {
	Insert(*Entity)
	Remove(*Entity)
	Participants() []*Entity
	Noise() float64
	FPS() uint64
}

type Shape struct {
	Position util.Vec2D `json:"position"`
	Radius   float64    `json:"radius"`
}

func (s *Shape) Copy() *Shape {
	return &Shape{
		Position: *s.Position.Copy(),
		Radius:   s.Radius,
	}
}

type ParticipantSetupInformation struct {
	Start    util.Vec2D `json:"start"`
	Radius   float64    `json:"radius"`
	SafeZone float64    `json:"safezone"`
	VMax     float64    `json:"vmax"`
	Target   util.Vec2D `json:"target"`
	Script   string     `json:"script"`
	// Flag for indicating, that this participant shall continue "executing" after reaching target position
	IgnoreFinish bool `json:"ignoreFinish"`
}
