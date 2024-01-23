package numerics

import (
	"golang.org/x/image/math/f64"
	"math"
)

type Vec2 struct {
	f64.Vec2
}

func NewVec2(a, b float64) Vec2 {
	return Vec2{f64.Vec2{a, b}}
}

func ZeroVec2() Vec2 {
	return NewVec2(0, 0)
}

func OneVec2() Vec2 {
	return NewVec2(1, 1)
}

func (v Vec2) X() float64 {
	return v.Vec2[0]
}

func (v Vec2) Y() float64 {
	return v.Vec2[1]
}

func (v Vec2) IsZero() bool {
	return v.X() == 0 && v.Y() == 0
}

func (v Vec2) Add(b Vec2) Vec2 {
	return NewVec2(v.X()+b.X(), v.Y()+b.Y())
}

func (v Vec2) Sub(b Vec2) Vec2 {
	return NewVec2(v.X()-b.X(), v.Y()-b.Y())
}

func (v Vec2) Mul(b Vec2) Vec2 {
	return NewVec2(v.X()*b.X(), v.Y()*b.Y())
}

func (v Vec2) Div(b Vec2) Vec2 {
	return NewVec2(v.X()/b.X(), v.Y()/b.Y())
}

func (v Vec2) AddScalar(b float64) Vec2 {
	return NewVec2(v.X()+b, v.Y()+b)
}

func (v Vec2) SubScalar(b float64) Vec2 {
	return NewVec2(v.X()-b, v.Y()-b)
}

func (v Vec2) MulScalar(b float64) Vec2 {
	return NewVec2(v.X()*b, v.Y()*b)
}

func (v Vec2) DivScalar(b float64) Vec2 {
	return NewVec2(v.X()/b, v.Y()/b)
}

func (v Vec2) Dot(b Vec2) float64 {
	return v.X()*b.X() + v.Y()*b.Y()
}

func (v Vec2) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vec2) LengthSquared() float64 {
	l := v.Length()
	return l * l
}

func (v Vec2) Normalized() Vec2 {
	l := v.Length()
	return v.Div(NewVec2(l, l))
}
