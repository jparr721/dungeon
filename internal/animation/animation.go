package animation

import (
	"bytes"
	assets "dungeon/assets/images"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"
	"image"
	_ "image/png"
	"os"
)

// Globally defined images
var (
	WizardFront *Image
	WizardSide  *Image
)

func init() {
	WizardFront = NewImageFromImage(assets.WizardSheet, 3, 0, 24, 24, 24)
	WizardSide = NewImageFromImage(assets.WizardSheet, 3, 24, 24, 24, 24)
}

type Image struct {
	FrameCount  int
	FrameOX     int
	FrameOY     int
	FrameWidth  int
	FrameHeight int

	*ebiten.Image
}

func NewImageFromFile(filename string, frameCount, frameOX, frameOY, frameWidth, frameHeight int) *Image {
	// Read the image filename
	zap.L().Debug("Loading image", zap.String("filename", filename))

	file, err := os.Open(filename)
	if err != nil {
		zap.L().Fatal("Failed to open image file", zap.Error(err))
	}

	img, _, err := image.Decode(file)
	if err != nil {
		zap.L().Fatal("Failed to decode image", zap.Error(err))
	}

	if err := file.Close(); err != nil {
		zap.L().Fatal("Failed to close image file", zap.Error(err))
	}

	return &Image{
		FrameCount:  frameCount,
		FrameOX:     frameOX,
		FrameOY:     frameOY,
		FrameWidth:  frameWidth,
		FrameHeight: frameHeight,
		Image:       ebiten.NewImageFromImage(img),
	}
}

func NewImageFromImage(imgBytes []byte, frameCount, frameOX, frameOY, frameWidth, frameHeight int) *Image {
	img, format, err := image.Decode(bytes.NewReader(imgBytes))

	zap.L().Debug("Image format", zap.String("format", format))

	if err != nil {
		fmt.Println(err)
		zap.L().Fatal("Failed to decode image", zap.Error(err))
	}

	return &Image{
		FrameCount:  frameCount,
		FrameOX:     frameOX,
		FrameOY:     frameOY,
		FrameWidth:  frameWidth,
		FrameHeight: frameHeight,
		Image:       ebiten.NewImageFromImage(img),
	}
}
