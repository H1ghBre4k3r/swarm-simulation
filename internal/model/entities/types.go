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
