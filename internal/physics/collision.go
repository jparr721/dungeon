package physics

import (
	"dungeon/internal/animation"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/math/f64"
	"image/color"
)

type AABB struct {
	Min f64.Vec2
	Max f64.Vec2
}

// NewAABB computes the bounding box from an image and player position
func NewAABB(pos f64.Vec2, img *animation.Image) *AABB {
	minBounds := f64.Vec2{pos[0], pos[1]}
	maxBounds := f64.Vec2{pos[0] + float64(img.FrameWidth), pos[1] + float64(img.FrameHeight)}

	return &AABB{
		Min: f64.Vec2{minBounds[0], minBounds[1]},
		Max: f64.Vec2{maxBounds[0], maxBounds[1]},
	}
}

func (a *AABB) String() string {
	return fmt.Sprintf("Min: (%.2f, %.2f), Max: (%.2f, %.2f)", a.Min[0], a.Min[1], a.Max[0], a.Max[1])
}

func (a *AABB) UpdatePosition(dx, dy float64) {
	a.Min[0] += dx
	a.Min[1] += dy
	a.Max[0] += dx
	a.Max[1] += dy
}

func (a *AABB) Render(screen *ebiten.Image, cameraTransform *ebiten.GeoM) {
	bbox := ebiten.NewImage(int(a.Max[0]-a.Min[0]), int(a.Max[1]-a.Min[1]))
	vector.StrokeRect(bbox, 0, 0, float32(bbox.Bounds().Max.X), float32(bbox.Bounds().Max.Y), 1, color.RGBA{R: 0xff, A: 0xff}, true)
	op := &ebiten.DrawImageOptions{
		GeoM: *cameraTransform,
	}
	//op.GeoM = *cameraTransform
	//op.GeoM.Translate(a.Min[0], a.Min[1])
	screen.DrawImage(bbox, op)
}

// IsExternallyColliding2D checks whether a, which is outside b, is about to clip into b
func (a *AABB) IsExternallyColliding2D(b *AABB) bool {
	if a.Max[0] < b.Min[0] || a.Min[0] > b.Max[0] {
		return false
	}

	if a.Max[1] < b.Min[1] || a.Min[1] > b.Max[1] {
		return false
	}

	return true
}

// IsInternallyColliding2D checks whether a, which is contained within b, is about to break out of b
func (a *AABB) IsInternallyColliding2D(b *AABB) bool {
	if a.Min[0] < b.Min[0] || a.Max[0] > b.Max[0] {
		return false
	}

	if a.Min[1] < b.Min[1] || a.Max[1] > b.Max[1] {
		return false
	}

	return true
}
