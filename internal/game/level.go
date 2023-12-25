package game

import (
	"bytes"
	"dungeon/internal/gfx"
	"dungeon/internal/numerics"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"go.uber.org/zap"
	"image"
	"image/color"
)

// Globally defined tiles
var (
	// Grass tile
	Grass      *Tile
	GrassLevel *Room
)

const (
	// TileSize is the size of a tile in pixels
	TileSize = 16

	// TileCount is the number of tiles in the tileset
	TileCount = 25
)

func init() {
	Grass = NewTileFromImage(images.Tiles_png, 243)

	// Calculate the number of layers needed to fill a 1000x500 room

	position := numerics.NewVec2(gfx.ScreenWidth/2-500, gfx.ScreenHeight/2-250)
	dimensions := numerics.NewVec2(adjustToTileSize(1000), adjustToTileSize(500))
	tilesNeeded := nTilesNeeded(int(dimensions.X() * dimensions.Y()))

	for i := 0; i < tilesNeeded; i++ {
		grassLayers = append(grassLayers, *Grass)
	}

	// Thisll be generated later
	GrassLevel = &Room{
		Position:   position,
		Dimensions: dimensions,
		IsBossRoom: false,
	}
	GrassLevel.Layers = append(GrassLevel.Layers, grassLayers)
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

	Index int
}

func NewTileFromImage(imgBytes []byte, index int) *Tile {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		zap.L().Fatal("Failed to decode image", zap.Error(err))
	}

	startX := (index % TileCount) * TileSize
	startY := (index / TileCount) * TileSize

	si := ebiten.NewImageFromImage(img).SubImage(image.Rect(startX, startY, startX+TileSize, startY+TileSize)).(*ebiten.Image)
	return &Tile{
		Image: si,
		Index: index,
	}
}

func (t *Tile) Render(screen *ebiten.Image, cameraTransform *ebiten.GeoM, pos numerics.Vec2) {
	op := &ebiten.DrawImageOptions{
		GeoM: *cameraTransform,
	}
	op.GeoM.Translate(pos.X(), pos.Y())

	screen.DrawImage(t.Image, op)
}

type Room struct {
	Layers [][]Tile

	// IsBossRoom just determines if this room needs to load a boss.
	IsBossRoom bool

	// Position is the position of the top-level corner of the rectangle.
	Position numerics.Vec2

	// Dimensions is the width and height of the rectangle.
	Dimensions numerics.Vec2
}

func NewRoom() *Room {
	position := numerics.NewVec2(gfx.ScreenWidth/2-500, gfx.ScreenHeight/2-250)
	dimensions := numerics.NewVec2(adjustToTileSize(1000), adjustToTileSize(500))
	tilesNeeded := nTilesNeeded(int(dimensions.X() * dimensions.Y()))

	grassLayers := make([]Tile, 0)
	for i := 0; i < tilesNeeded; i++ {
		grassLayers = append(grassLayers, *Grass)
	}

	// Thisll be generated later
	GrassLevel = &Room{
		Position:   position,
		Dimensions: dimensions,
		IsBossRoom: false,
	}
	GrassLevel.Layers = append(GrassLevel.Layers, grassLayers)
	return &Room{}
}

func (r *Room) Render(screen *ebiten.Image, cameraTransform *ebiten.GeoM) {
	worldSizeX := int(r.Dimensions.X() / TileSize)
	worldSizeY := int(r.Dimensions.Y() / TileSize)

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
