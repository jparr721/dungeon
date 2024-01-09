package main

import (
	"bytes"
	imgui "github.com/gabstv/cimgui-go"
	ebimgui "github.com/gabstv/ebiten-imgui/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"go.uber.org/zap"
	"image"
	"log"
	"os"
)

var (
	TileSize     = 16
	TileCount    = 25
	TileIndex    = 47
	ImageIDRef   = 10
	ViewerWidth  = 1920
	ViewerHeight = 1080
)

type LevelDebug struct {
	LoadedImage *ebiten.Image
	ImageWidth  float32
	ImageHeight float32
}

func LoadImage(imgBytes []byte) (*ebiten.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		zap.L().Fatal("Failed to decode image", zap.Error(err))
		return nil, err

	}

	ei := ebiten.NewImageFromImage(img)

	return ei, nil
}

func Image(tid imgui.TextureID, size imgui.Vec2) {
	uv0 := imgui.NewVec2(0, 0)
	uv1 := imgui.NewVec2(1, 1)
	borderCol := imgui.NewVec4(0, 0, 0, 0)
	tintCol := imgui.NewVec4(1, 1, 1, 1)

	imgui.SetNextWindowPos(imgui.NewVec2(0, 0))
	imgui.SetNextWindowSize(size)
	imgui.ImageV(tid, size, uv0, uv1, tintCol, borderCol)
}

func (l *LevelDebug) Draw(screen *ebiten.Image) {
	ebimgui.Draw(screen)
}

func (l *LevelDebug) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	ebimgui.SetDisplaySize(float32(ViewerWidth), float32(ViewerHeight))
	return ViewerWidth, ViewerHeight
}

func (l *LevelDebug) Update() error {
	ebimgui.Update(1.0 / 60.0)
	ebimgui.BeginFrame()
	defer ebimgui.EndFrame()

	ebimgui.GlobalManager().Cache.SetTexture(imgui.TextureID(&ImageIDRef), l.LoadedImage)
	Image(imgui.TextureID(&ImageIDRef), imgui.Vec2{X: l.ImageWidth, Y: l.ImageHeight})

	return nil
}

func init() {
	logger := zap.Must(zap.NewDevelopment())
	if os.Getenv("APP_ENV") == "release" {
		logger = zap.Must(zap.NewProduction())
	}
	zap.ReplaceGlobals(logger)
}

func main() {
	LoadedImage, err := LoadImage(images.Tiles_png)
	if err != nil {
		log.Fatal(err)
	}
	point := LoadedImage.Bounds().Size()
	ImageWidth := float32(point.X)
	ImageHeight := float32(point.Y)

	ebiten.SetWindowSize(ViewerWidth, ViewerHeight)
	ebiten.SetWindowTitle("Level Debug")

	if err := ebiten.RunGame(&LevelDebug{
		LoadedImage: LoadedImage,
		ImageWidth:  ImageWidth,
		ImageHeight: ImageHeight,
	}); err != nil {
		log.Fatal(err)
	}
}
