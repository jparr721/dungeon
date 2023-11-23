package physics

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/math/f64"
)

type AABB struct {
	Min f64.Vec2
	Max f64.Vec2
}

func NewAABB(img *ebiten.Image) *AABB {
	bounds := img.Bounds()
	minBounds := bounds.Min
	maxBounds := bounds.Max
	return &AABB{
		Min: f64.Vec2{float64(minBounds.X), float64(minBounds.Y)},
		Max: f64.Vec2{float64(maxBounds.X), float64(maxBounds.Y)},
	}
}

func (a *AABB) IsColliding2D(b *AABB) bool {
	if a.Max[0] < b.Min[0] || a.Min[0] > b.Max[0] {
		return false
	}

	if a.Max[1] < b.Min[1] || a.Min[1] > b.Max[1] {
		return false
	}

	return true
}

type Collidable interface {
	// IsCollidingInternal checks if an object, which exists WITHIN a bounding volume, is coming near
	// the edge of the shape where it would break out. This is for things like levels.
	IsCollidingInternal(b *Collidable) bool

	// IsCollidingExternal check if an ojbect, which exists OUTSIDE of a bounding volume, is going
	// to clip inside of the edge of the shape. This is for basically everything else.
	IsCollidingExternal(b *Collidable) bool
}
