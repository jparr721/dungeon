package game

import (
	"dungeon/internal/animation"
	"dungeon/internal/numerics"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type Collidable interface {
	// IsCollidingInternal checks if an object, which exists WITHIN a bounding volume, is coming near
	// the edge of the shape where it would break out. This is for things like levels.
	IsCollidingInternal(b Collidable) bool

	// IsCollidingExternal check if an object, which exists OUTSIDE a bounding volume, is going
	// to clip inside the edge of the shape. This is for basically everything else.
	IsCollidingExternal(b Collidable) bool

	BoundingBox() *AABB
}

// Orientation enum representing the orientation of an object (Front, Left, Right)
type Orientation int

const (
	Front Orientation = iota
	Back
	Left
	Right
	All
)

func (o Orientation) String() string {
	switch o {
	case Front:
		return "Front"
	case Back:
		return "Back"
	case Left:
		return "Left"
	case Right:
		return "Right"
	case All:
		return "All"
	default:
		return "Unknown"
	}
}

// Object represents any game object which can be interactive
type Object struct {
	// The images representing the currently loaded object and its various orientations
	Image map[Orientation]*animation.Image

	// The image options for this object
	Op *ebiten.DrawImageOptions

	// The current position of the object
	Position numerics.Vec2

	// The center of the object in world space
	Center numerics.Vec2

	// The current velocity of the object
	Velocity numerics.Vec2

	// The current rotation of the object
	Rotation float64

	// The count of the animation frame
	Count int

	// The current orientation of the image
	Orientation Orientation

	*AABB
}

// NewObjectFromImage creates a new object from an image instance.
func NewObjectFromImage(image *animation.Image) *Object {
	return NewObjectFromImages(map[Orientation]*animation.Image{All: image})
}

func NewObjectFromImages(images map[Orientation]*animation.Image) *Object {
	var aabb *AABB
	var center numerics.Vec2
	var orientation Orientation
	if _, ok := images[All]; !ok {
		aabb = NewAABB(numerics.ZeroVec2(), images[Front])
		center = numerics.NewVec2(float64(images[Front].FrameWidth/2), float64(images[Front].FrameHeight/2))
		orientation = Front
	} else {
		aabb = NewAABB(numerics.ZeroVec2(), images[All])
		center = numerics.NewVec2(float64(images[All].FrameWidth/2), float64(images[All].FrameHeight/2))
		orientation = All
	}

	return &Object{
		Image:       images,
		Op:          &ebiten.DrawImageOptions{},
		Center:      center,
		Velocity:    numerics.OneVec2(),
		AABB:        aabb,
		Orientation: orientation,
	}
}

func (o *Object) UpdatePosition(diff numerics.Vec2, objects []*Object, room *Room) {
	// Check for object collisions
	for _, obj := range objects {
		if obj == o {
			continue
		}

		// If any collision is occurring, do not let the player move in the direction of the collision
		if o.IsCollidingInternal(obj) || o.IsCollidingExternal(obj) {
			diff = numerics.ZeroVec2()
		}
	}

	// Check for room collisions
	if room != nil && room.CheckCollisionAndUpdatePosition(o, diff) {
		diff = numerics.ZeroVec2()
	}

	o.Position = o.Position.Add(diff)
	o.Center = o.Center.Add(diff)

	o.AABB.UpdatePosition(diff)

}

func (o *Object) Render(screen *ebiten.Image, cameraTransform *ebiten.GeoM) {
	img := o.Image[o.Orientation]

	// First, quick check if an image for "All" is set, if it is, always use that
	if img == nil {
		img = o.Image[All]
	}

	// First, rotate BEFORE any translation has occurred, we MUST create a new geom every time.
	o.Op.GeoM = ebiten.GeoM{}

	if o.Orientation == Left {
		// Left to right flips over the y axis
		o.Op.GeoM.Scale(1.0, -1.0)
		o.Op.GeoM.Translate(0, float64(img.FrameHeight))
	}

	if o.Orientation == Back {
		// TODO: This doesn't really fix the issue
		o.Op.GeoM.Scale(-1.0, 1.0)
		o.Op.GeoM.Translate(float64(img.FrameWidth), 0)
	}

	// Translate to the center of the object
	o.Op.GeoM.Translate(-float64(img.FrameWidth)/2, -float64(img.FrameHeight)/2)
	// Apply rotation
	o.Op.GeoM.Rotate(o.Rotation)
	// Translate back to the original position
	o.Op.GeoM.Translate(float64(img.FrameWidth)/2, float64(img.FrameHeight)/2)

	// Now, apply the camera transformation to this
	o.Op.GeoM.Concat(*cameraTransform)

	// Move to the object position including any camera offset
	o.Op.GeoM.Translate(o.Position.X(), o.Position.Y())

	// This just chooses the character frame from the sprite sheet. We divide by 5 so that way the transition
	// between animation frames is less intense.
	i := (o.Count / 10) % img.FrameCount
	//sx, sy := img.FrameOX+i*img.FrameWidth, img.FrameOY
	sx, sy := img.FrameOX, img.FrameOY+i*img.FrameHeight

	screen.DrawImage(img.SubImage(image.Rect(sx, sy, sx+img.FrameWidth, sy+img.FrameHeight)).(*ebiten.Image), o.Op)

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

func (o *Object) BoundingBox() *AABB {
	return o.AABB
}
