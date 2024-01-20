package game

import (
	"dungeon/internal/animation"
	"dungeon/internal/numerics"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math/rand"
)

// Door is a door that the player character can pass through
type Door struct {
	// To is the pointer that this door connects to
	To *Room

	*Object
}

func NewDoor(position numerics.Vec2, to *Room, doorImage *animation.Image) *Door {
	obj := NewObjectFromImages(map[Orientation]*animation.Image{All: doorImage})

	// Position is the new position, we need to get from where we are to that position
	diff := position
	obj.UpdatePosition(diff)
	return &Door{
		To:     to,
		Object: obj,
	}
}

type Room struct {
	Layers [][]*Tile

	// IsBossRoom just determines if this room needs to load a boss.
	IsBossRoom bool

	// Position is the position of the top-level corner of the rectangle.
	Position numerics.Vec2

	// Dimensions is the width and height of the rectangle.
	Dimensions numerics.Vec2

	// Every room has at least one door
	Doors []*Door

	// Color is the color of the boundary box of the room
	Color color.Color

	// StrokeWidth is the width of the stroke of the boundary box
	StrokeWidth float32
}

func NewRoom(position, dimensions numerics.Vec2) *Room {
	// TODO Add layers
	// The tiles needed to cover the floor
	//tilesNeeded := nTilesNeeded(int(dimensions.X() * dimensions.Y()))
	layer := make([]*Tile, 0)
	//for i := 0; i < tilesNeeded; i++ {
	//	layer = append(layer, GrassEmpty)
	//}

	room := &Room{
		IsBossRoom: false,
		Position:   position,
		Dimensions: dimensions,

		Doors: make([]*Door, 0),

		// Random fill color
		Color: color.RGBA{
			R: uint8(rand.Intn(255)),
			G: uint8(rand.Intn(255)),
			B: uint8(rand.Intn(255)),
			A: 0xff,
		},

		// Stroke width to give boundaries some texture
		StrokeWidth: 50,
	}
	room.Layers = append(room.Layers, layer)
	return room
}

func (r *Room) Bounds() (numerics.Vec2, numerics.Vec2) {
	return r.Position, r.Position.Add(r.Dimensions)
}

func (r *Room) CheckCollisionAndUpdatePosition(object *Object, diff numerics.Vec2) numerics.Vec2 {
	newPos := object.Position.Add(diff)

	startBounds, endBounds := r.Bounds()
	startBounds = startBounds.AddScalar(float64(r.StrokeWidth) / 2)
	// Account for the stroke width
	endBounds = endBounds.SubScalar(float64(r.StrokeWidth) / 2)
	endBounds = endBounds.Sub(numerics.NewVec2(
		float64(object.Image[Front].FrameWidth),
		float64(object.Image[Front].FrameHeight),
	))

	// Check if the physics object is about to exceed the extent of the room
	// Check for collision on the X-axis and update position
	if newPos.X() < startBounds.X() {
		newPos = numerics.NewVec2(startBounds.X(), newPos.Y())
	} else if newPos.X() >= endBounds.X() {
		newPos = numerics.NewVec2(endBounds.X(), newPos.Y())
	}

	// Check for collision on the Y-axis and update position
	if newPos.Y() < startBounds.Y() {
		newPos = numerics.NewVec2(newPos.X(), startBounds.Y())
	} else if newPos.Y() >= endBounds.Y() {
		newPos = numerics.NewVec2(newPos.X(), endBounds.Y())
	}

	// Update diff to reflect newPos
	return newPos.Sub(object.Position)
}

func (r *Room) Render(screen *ebiten.Image, cameraTransform *ebiten.GeoM) {
	//worldSizeX := int(r.Dimensions.X() / TileSize)
	//worldSizeY := int(r.Dimensions.Y() / TileSize)

	boundary := ebiten.NewImage(int(r.Dimensions.X()), int(r.Dimensions.Y()))
	vector.StrokeRect(
		boundary,
		0,
		0,
		float32(boundary.Bounds().Max.X),
		float32(boundary.Bounds().Max.Y),
		r.StrokeWidth,
		r.Color,
		true,
	)

	op := &ebiten.DrawImageOptions{
		GeoM: *cameraTransform,
	}
	op.GeoM.Translate(r.Position.X(), r.Position.Y())
	screen.DrawImage(boundary, op)

	// Render doors
	for _, door := range r.Doors {
		door.Render(screen, cameraTransform)
	}

	//for x := 0; x < worldSizeX; x++ {
	//	for y := 0; y < worldSizeY; y++ {
	//		t := r.Layers[0][x+y*worldSizeX]
	//		drawPos := numerics.NewVec2(
	//			r.Position.X()+float64(x*TileSize),
	//			r.Position.Y()+float64(y*TileSize),
	//		)
	//		t.Render(screen,
	//			cameraTransform,
	//			drawPos,
	//		)
	//	}
	//}
}
