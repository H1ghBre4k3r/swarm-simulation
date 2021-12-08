package entities

import "github.com/H1ghBre4k3r/swarm-simulation/internal/model/util"

type UpdateFn func(*Entity)

type Portal interface {
	Insert(*Entity)
	Remove(*Entity)
}

type Shape struct {
	Position util.Vec2D `json:"position"`
	Radius   float64    `json:"radius"`
}
