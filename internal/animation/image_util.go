package animation

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"
	"image"
)

func LoadImage(imgBytes []byte) (*ebiten.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		zap.L().Fatal("Failed to decode image", zap.Error(err))
		return nil, err
	}

	ei := ebiten.NewImageFromImage(img)

	return ei, nil
}
