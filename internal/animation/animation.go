package animation

import (
	"bytes"
	assets "dungeon/assets/images"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"
	"image"
	_ "image/png"
	"log"
	"os"
)

// Globally defined images
var (
	WizardFront *Image
	WizardSide  *Image
)

func init() {
	WizardFront = NewImageFromImageBytes(assets.WizardSheet, 3, 0, 24, 24, 24)
	WizardSide = NewImageFromImageBytes(assets.WizardSheet, 3, 24, 24, 24, 24)
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

// NewImageFromImage creates a new image from an existing ebiten image object. This is reserved for filled images.
func NewImageFromImage(img *ebiten.Image) *Image {
	bounds := img.Bounds().Size()
	return &Image{
		FrameCount:  1,
		FrameOX:     0,
		FrameOY:     0,
		FrameWidth:  bounds.X,
		FrameHeight: bounds.Y,
		Image:       img,
	}
}

func NewImageFromImageBytes(imgBytes []byte, frameCount, frameOX, frameOY, frameWidth, frameHeight int) *Image {
	if frameCount < 1 {
		log.Fatal("Frame count cannot be < 1")
	}

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
