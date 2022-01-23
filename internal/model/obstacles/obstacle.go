package obstacles

import "github.com/H1ghBre4k3r/swarm-simulation/internal/model/util"

type Obstacle struct {
	Start util.Vec2D `json:"start"`
	End   util.Vec2D `json:"end"`
}
