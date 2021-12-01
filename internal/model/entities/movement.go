package entities

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	R float64 `json:"r"`
}

type Velocity struct {
	X float64
	Y float64
}

func (pos *Position) Move(velocity *Velocity) {
	pos.X += velocity.X
	pos.Y += velocity.Y
}
