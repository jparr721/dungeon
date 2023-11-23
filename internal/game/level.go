package game

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"go.uber.org/zap"
	"image"
)

// Globally defined tiles
var (
	// Grass tile
	Grass *Tile
)

type Tile struct {
	*ebiten.Image
}

func NewTileFromImage(imgBytes []byte) *Tile {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		zap.L().Fatal("Failed to decode image", zap.Error(err))
	}

	return &Tile{
		Image: ebiten.NewImageFromImage(img),
	}
}

type Level struct {
	Layers [][]int
}

func init() {
	Grass = NewTileFromImage(images.Tiles_png)
}
