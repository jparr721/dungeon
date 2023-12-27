package game

import (
	"bytes"
	"dungeon/internal/gfx"
	"dungeon/internal/numerics"
	imgui "github.com/gabstv/cimgui-go"
	ebimgui "github.com/gabstv/ebiten-imgui/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"go.uber.org/zap"
	"image"
	"image/color"
	"math/rand"
)

// Globally defined tiles
var (
	Debug       *Tile
	DebugIndex  int32
	GrassEmpty  *Tile
	GrassPlant1 *Tile
	GrassPlant2 *Tile

	Cobblestone1 *Tile
	Cobblestone2 *Tile
	Cobblestone3 *Tile
	Cobblestone4 *Tile
	Cobblestone5 *Tile

	GrassLevel *Room
)

const (
	// TileSize is the size of a tile in pixels
	TileSize = 16

	// TileCount is the number of tiles in the tileset
	TileCount = 25
)

func init() {
	Debug = NewTileFromImage(images.Tiles_png, TileCount, TileSize, int(DebugIndex))
	GrassEmpty = NewTileFromImage(images.Tiles_png, TileCount, TileSize, 243)
	GrassPlant1 = NewTileFromImage(images.Tiles_png, TileCount, TileSize, 218)
	GrassPlant2 = NewTileFromImage(images.Tiles_png, TileCount, TileSize, 219)
	GrassLevel = NewRoom(
		numerics.NewVec2(gfx.ScreenWidth/2-500, gfx.ScreenHeight/2-250),
		numerics.NewVec2(adjustToTileSize(1000), adjustToTileSize(500)))
}

func adjustToTileSize(dimension int) float64 {
	if dimension%TileSize != 0 {
		return float64((dimension / TileSize) * TileSize)
	}
	return float64(dimension)
}

func nTilesNeeded(area int) int {
	return area / (TileSize * TileSize)
}

type Tile struct {
	*ebiten.Image
	// The index into the image
	Index int
}

func NewTileFromImage(imgBytes []byte, tileCount, tileSize, index int) *Tile {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		zap.L().Fatal("Failed to decode image", zap.Error(err))
	}

	startX := (index % tileCount) * tileSize
	startY := (index / tileCount) * tileSize

	si := ebiten.NewImageFromImage(img).SubImage(image.Rect(
		startX,
		startY,
		startX+tileSize,
		startY+tileSize,
	)).(*ebiten.Image)

	return &Tile{
		Image: si,
		Index: index,
	}
}

func (t *Tile) Render(screen *ebiten.Image, cameraTransform *ebiten.GeoM, pos numerics.Vec2) {
	if cameraTransform == nil {
		cameraTransform = &ebiten.GeoM{}
	}

	op := &ebiten.DrawImageOptions{
		GeoM: *cameraTransform,
	}
	op.GeoM.Translate(pos.X(), pos.Y())

	screen.DrawImage(t.Image, op)
}

type Room struct {
	Layers [][]*Tile

	// IsBossRoom just determines if this room needs to load a boss.
	IsBossRoom bool

	// Position is the position of the top-level corner of the rectangle.
	Position numerics.Vec2

	// Dimensions is the width and height of the rectangle.
	Dimensions numerics.Vec2
}

func NewRoom(position, dimensions numerics.Vec2) *Room {
	// The tiles needed to cover the floor
	tilesNeeded := nTilesNeeded(int(dimensions.X() * dimensions.Y()))

	layer := make([]*Tile, 0)
	for i := 0; i < tilesNeeded; i++ {
		if rand.Float64() < 0.01 {
			if rand.Float64() < 0.5 {
				layer = append(layer, GrassPlant1)
			} else {
				layer = append(layer, GrassPlant2)
			}
		} else {
			layer = append(layer, Debug)
		}

	}

	// This'll be generated later
	room := &Room{
		Position:   position,
		Dimensions: dimensions,
		IsBossRoom: false,
	}
	room.Layers = append(room.Layers, layer)
	return room
}

func (r *Room) Render(screen *ebiten.Image, cameraTransform *ebiten.GeoM) {
	ebimgui.Update(1.0 / 60.0)
	ebimgui.BeginFrame()
	defer ebimgui.EndFrame()
	imgui.InputIntV(
		"Index",
		&DebugIndex,
		1,
		5,
		0,
	)

	worldSizeX := int(r.Dimensions.X() / TileSize)
	worldSizeY := int(r.Dimensions.Y() / TileSize)

	newDebug := NewTileFromImage(images.Tiles_png, TileCount, TileSize, int(DebugIndex))
	// Update all the tiles
	for x := 0; x < worldSizeX; x++ {
		for y := 0; y < worldSizeY; y++ {
			r.Layers[0][x+y*worldSizeX] = newDebug
		}
	}

	bg := ebiten.NewImage(int(r.Dimensions.X()), int(r.Dimensions.Y()))
	bg.Fill(color.White)
	op := &ebiten.DrawImageOptions{
		GeoM: *cameraTransform,
	}
	op.GeoM.Translate(r.Position.X(), r.Position.Y())
	screen.DrawImage(bg, op)

	for x := 0; x < worldSizeX; x++ {
		for y := 0; y < worldSizeY; y++ {
			t := r.Layers[0][x+y*worldSizeX]
			drawPos := numerics.NewVec2(
				r.Position.X()+float64(x*TileSize),
				r.Position.Y()+float64(y*TileSize),
			)
			t.Render(screen,
				cameraTransform,
				drawPos,
			)
		}
	}
}
