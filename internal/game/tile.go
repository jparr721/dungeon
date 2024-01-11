package game

import (
	"dungeon/internal/animation"
	"dungeon/internal/numerics"
	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"
	"image"
)

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

func NewTileFromImage(imgBytes []byte, startX, startY, tileSize int) *Tile {
	egImg, err := animation.LoadImage(imgBytes)
	if err != nil {
		zap.L().Fatal("Failed to load image", zap.Error(err))
		return nil
	}

	egImg = egImg.SubImage(image.Rect(
		startX,
		startY,
		startX+tileSize,
		startY+tileSize,
	)).(*ebiten.Image)

	return &Tile{
		Image: egImg,
		Index: -1,
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
