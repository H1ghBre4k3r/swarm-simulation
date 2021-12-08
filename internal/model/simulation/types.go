package simulation

import (
	"github.com/AvraamMavridis/randomcolor"
	"github.com/H1ghBre4k3r/swarm-simulation/internal/model/util"
)

type Drawable interface {
	GetX() float64
	GetY() float64
	GetR() float64
	GetColor() randomcolor.RGBColor
	GetVelocity() util.Vec2D
}

type View interface {
	Render([]Drawable)
}
