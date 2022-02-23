package util

import (
	"math"
	"math/rand"
	"time"
)

// Utility for storing and manipulating 2D information
type Vec2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (v *Vec2D) Copy() *Vec2D {
	return &Vec2D{
		X: v.X,
		Y: v.Y,
	}
}

func (v *Vec2D) AddI(other *Vec2D) *Vec2D {
	v.X += other.X
	v.Y += other.Y
	return v
}

func (v *Vec2D) Add(other *Vec2D) *Vec2D {
	return v.Copy().AddI(other)
}

func (v *Vec2D) ScaleI(scale float64) *Vec2D {
	v.X *= scale
	v.Y *= scale
	return v
}

func (v *Vec2D) Scale(scale float64) *Vec2D {
	return v.Copy().ScaleI(scale)
}

func (v *Vec2D) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vec2D) NormalizeI() *Vec2D {
	v.X /= v.Length()
	v.Y /= v.Length()
	return v
}

func (v *Vec2D) Normalize() *Vec2D {
	return v.Copy().NormalizeI()
}

func (v *Vec2D) NoiseI(stddev float64) *Vec2D {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	v.X += r.NormFloat64() * stddev
	v.Y += r.NormFloat64() * stddev
	return v
}

func (v *Vec2D) Noise(stddev float64) *Vec2D {
	return v.Copy().NoiseI(stddev)
}
