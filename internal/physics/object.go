package physics

import (
	"dungeon/internal/animation"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/math/f64"
	"image"
)

type Collidable interface {
	// IsCollidingInternal checks if an object, which exists WITHIN a bounding volume, is coming near
	// the edge of the shape where it would break out. This is for things like levels.
	IsCollidingInternal(b *Collidable) bool

	// IsCollidingExternal check if an object, which exists OUTSIDE a bounding volume, is going
	// to clip inside the edge of the shape. This is for basically everything else.
	IsCollidingExternal(b *Collidable) bool

	BoundingBox() *AABB
}

// Object represents any game object which can be interactive
type Object struct {
	// The image representing the currently loaded object
	Image *animation.Image

	// The image options for this object
	Op *ebiten.DrawImageOptions

	// The current position of the object
	Position f64.Vec2

	// The center of the object in world space
	Center f64.Vec2

	// The current velocity of the object
	Velocity f64.Vec2

	// The current rotation of the object
	Rotation float64

	// The count of the animation frame
	Count int

	*AABB
}

// NewObjectFromImage creates a new object from an image instance.
func NewObjectFromImage(image *animation.Image) *Object {
	return &Object{
		Image:    image,
		Op:       &ebiten.DrawImageOptions{},
		Center:   f64.Vec2{float64(image.FrameWidth / 2), float64(image.FrameHeight / 2)},
		Velocity: f64.Vec2{1.0, 1.0},
		AABB:     NewAABB(f64.Vec2{0, 0}, image),
	}
}

func (o *Object) UpdatePosition(dx, dy float64) {
	o.Position[0] += dx
	o.Position[1] += dy
	o.Center[0] += dx
	o.Center[1] += dy
	o.AABB.UpdatePosition(dx, dy)
}

func (o *Object) Render(screen *ebiten.Image, cameraTransform *ebiten.GeoM) {
	// First, rotate BEFORE any translation has occurred, we MUST create a new geom every time.
	o.Op.GeoM = ebiten.GeoM{}

	// Translate to the center of the object
	o.Op.GeoM.Translate(-float64(o.Image.FrameWidth)/2, -float64(o.Image.FrameHeight)/2)
	// Apply rotation
	o.Op.GeoM.Rotate(o.Rotation)
	// Translate back to the original position
	o.Op.GeoM.Translate(float64(o.Image.FrameWidth)/2, float64(o.Image.FrameHeight)/2)

	// Now, apply the camera transformation to this
	o.Op.GeoM.Concat(*cameraTransform)

	// Move to the object position including any camera offset
	o.Op.GeoM.Translate(o.Position[0], o.Position[1])

	// This just chooses the character frame from the sprite sheet. We divide by 5 so that way the transition
	// between animation frames is less intense (basically going at 5 frames per second).
	i := (o.Count / 5) % o.Image.FrameCount
	sx, sy := o.Image.FrameOX+i*o.Image.FrameWidth, o.Image.FrameOY

	screen.DrawImage(o.Image.SubImage(image.Rect(sx, sy, sx+o.Image.FrameWidth, sy+o.Image.FrameHeight)).(*ebiten.Image), o.Op)

	// Draw the player's bounding box
	o.AABB.Render(screen, &o.Op.GeoM)
}

// IsCollidingInternal implements the collidable interface for the Object
func (o *Object) IsCollidingInternal(b Collidable) bool {
	return o.IsInternallyColliding2D(b.BoundingBox())
}

// IsCollidingExternal implements the collidable interface for the Object
func (o *Object) IsCollidingExternal(b Collidable) bool {
	return o.IsExternallyColliding2D(b.BoundingBox())
}
