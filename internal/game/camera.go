package game

import (
	"dungeon/internal/numerics"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type Camera struct {
	ViewPort     numerics.Vec2
	Position     numerics.Vec2
	ZoomFactor   float64
	ZoomFactorTo float64
	Rotation     float64
}

func (c *Camera) String() string {
	return fmt.Sprintf(
		"T: %.1f, R: %.2f, S: %.2f",
		c.Position, c.Rotation, c.ZoomFactor,
	)
}

func (c *Camera) viewportCenter() numerics.Vec2 {
	return numerics.NewVec2(
		c.ViewPort.X()*0.5,
		c.ViewPort.Y()*0.5,
	)
}

func (c *Camera) worldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-c.Position.X(), -c.Position.Y())
	// We want to scale and rotate around center of image / screen
	m.Translate(-c.viewportCenter().X(), -c.viewportCenter().Y())
	m.Scale(
		math.Pow(1.01, c.ZoomFactor),
		math.Pow(1.01, c.ZoomFactor),
	)
	m.Rotate(float64(c.Rotation) * 2 * math.Pi / 360)
	m.Translate(c.viewportCenter().X(), c.viewportCenter().Y())
	return m
}

func (c *Camera) ScreenToWorld(posX, posY int) (float64, float64) {
	inverseMatrix := c.worldMatrix()
	if inverseMatrix.IsInvertible() {
		inverseMatrix.Invert()
		return inverseMatrix.Apply(float64(posX), float64(posY))
	} else {
		// When scaling it can happened that matrix is not invertable
		return math.NaN(), math.NaN()
	}
}
