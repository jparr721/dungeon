package game

import (
	"dungeon/internal/animation"
	"dungeon/internal/numerics"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

type AABB struct {
	Min numerics.Vec2
	Max numerics.Vec2
}

// NewAABB computes the bounding box from an image and player position
func NewAABB(pos numerics.Vec2, img *animation.Image) *AABB {
	x, y := pos.X(), pos.Y()
	minBounds := numerics.NewVec2(x, y)
	maxBounds := numerics.NewVec2(x+float64(img.FrameWidth), y+float64(img.FrameHeight))

	return &AABB{
		Min: minBounds,
		Max: maxBounds,
	}
}

func (a *AABB) String() string {
	return fmt.Sprintf("Min: (%.2f, %.2f), Max: (%.2f, %.2f)", a.Min.X(), a.Min.Y(), a.Max.X(), a.Max.Y())
}

func (a *AABB) UpdatePosition(diff numerics.Vec2) {
	a.Min = a.Min.Add(diff)
	a.Max = a.Max.Add(diff)
}

func (a *AABB) Render(screen *ebiten.Image, cameraTransform *ebiten.GeoM) {
	bbox := ebiten.NewImage(int(a.Max.X()-a.Min.X()), int(a.Max.Y()-a.Min.Y()))
	vector.StrokeRect(
		bbox,
		0,
		0,
		float32(bbox.Bounds().Max.X),
		float32(bbox.Bounds().Max.Y),
		1,
		color.RGBA{R: 0xff, A: 0xff},
		true,
	)
	op := &ebiten.DrawImageOptions{
		GeoM: *cameraTransform,
	}
	screen.DrawImage(bbox, op)
}

// IsExternallyColliding2D checks whether a, which is outside b, is about to clip into b
func (a *AABB) IsExternallyColliding2D(b *AABB) bool {
	if a.Max.X() < b.Min.X() || a.Min.X() > b.Max.X() {
		return false
	}

	if a.Max.Y() < b.Min.Y() || a.Min.Y() > b.Max.Y() {
		return false
	}

	return true
}

// IsInternallyColliding2D checks whether a, which is contained within b, is about to break out of b
func (a *AABB) IsInternallyColliding2D(b *AABB) bool {
	if a.Min.X() < b.Min.X() || a.Max.X() > b.Max.X() {
		return false
	}

	if a.Min.Y() < b.Min.Y() || a.Max.Y() > b.Max.Y() {
		return false
	}

	return true
}
