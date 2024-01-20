package game

import (
	"dungeon/internal/animation"
	"dungeon/internal/numerics"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

type CollisionDirection struct {
	X, Y bool
}

type AABB struct {
	IsColliding        bool
	CollisionDirection CollisionDirection
	Min                numerics.Vec2
	Max                numerics.Vec2
}

// NewAABB computes the bounding box from an image and player position
func NewAABB(pos numerics.Vec2, img *animation.Image) *AABB {
	x, y := pos.X(), pos.Y()
	minBounds := numerics.NewVec2(x, y)
	maxBounds := numerics.NewVec2(x+float64(img.FrameWidth), y+float64(img.FrameHeight))

	return &AABB{
		IsColliding: false,
		Min:         minBounds,
		Max:         maxBounds,
	}
}

func (a *AABB) String() string {
	return fmt.Sprintf(
		"IsColliding: %t Min: (%.2f, %.2f), Max: (%.2f, %.2f)",
		a.IsColliding,
		a.Min.X(),
		a.Min.Y(),
		a.Max.X(),
		a.Max.Y(),
	)
}

func (a *AABB) Dimensions() numerics.Vec2 {
	return numerics.NewVec2(
		a.Max.X()-a.Min.X(),
		a.Max.Y()-a.Min.Y(),
	)
}

// SetPosition sets the position of the AABB such that the min value is the top-left and the max is the bottom right.
func (a *AABB) SetPosition(min numerics.Vec2, max numerics.Vec2) {
	a.Min = min
	a.Max = max
}

func (a *AABB) UpdatePosition(diff numerics.Vec2) {
	a.Min = a.Min.Add(diff)
	a.Max = a.Max.Add(diff)
}

func (a *AABB) Render(screen *ebiten.Image, cameraTransform *ebiten.GeoM) {
	bbox := ebiten.NewImage(int(a.Max.X()-a.Min.X()), int(a.Max.Y()-a.Min.Y()))

	boxColor := color.RGBA{R: 0xff, A: 0xff}
	if a.IsColliding {
		boxColor = color.RGBA{G: 0xff, A: 0xff}
	}

	vector.StrokeRect(
		bbox,
		0,
		0,
		float32(bbox.Bounds().Max.X),
		float32(bbox.Bounds().Max.Y),
		3,
		boxColor,
		true,
	)
	op := &ebiten.DrawImageOptions{
		GeoM: *cameraTransform,
	}
	screen.DrawImage(bbox, op)
}

func (a *AABB) ResetCollisionState() {
	a.IsColliding = false
	a.CollisionDirection = CollisionDirection{}
}

// IsExternallyColliding2D checks whether a, which is outside b, is about to clip into b
func (a *AABB) IsExternallyColliding2D(b *AABB) bool {
	if a.Max.X() < b.Min.X() || a.Min.X() > b.Max.X() {
		return false
	}

	if a.Max.Y() < b.Min.Y() || a.Min.Y() > b.Max.Y() {
		return false
	}

	a.IsColliding = true
	b.IsColliding = true

	// Get the direction of the collision
	overlapX := min(a.Max.X(), b.Max.X()) - max(a.Min.X(), b.Min.X())
	overlapY := min(a.Max.Y(), b.Max.Y()) - max(a.Min.Y(), b.Min.Y())

	if overlapX > 0 && overlapY > 0 {
		// Determine primary collision axis
		if overlapX > overlapY {
			// Y-axis is primary
			a.CollisionDirection.Y = true
			b.CollisionDirection.Y = true
		} else {
			// X-axis is primary
			a.CollisionDirection.X = true
			b.CollisionDirection.X = true
		}

		return true
	}

	return false
}

// IsInternallyColliding2D checks whether a, which is contained within b, is about to break out of b
func (a *AABB) IsInternallyColliding2D(b *AABB) bool {
	if a.Min.X() < b.Min.X() || a.Max.X() > b.Max.X() {
		return false
	}

	if a.Min.Y() < b.Min.Y() || a.Max.Y() > b.Max.Y() {
		return false
	}

	a.IsColliding = true
	b.IsColliding = true
	return true
}
