package entities

type Position struct {
	X float64
	Y float64
	R float64
}

type Velocity struct {
	X float64
	Y float64
}

func (pos *Position) Move(velocity *Velocity) {
	pos.X += velocity.X
	pos.Y += velocity.Y
}
