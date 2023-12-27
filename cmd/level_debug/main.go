package main

import (
	"dungeon/internal/game"
	imgui "github.com/gabstv/cimgui-go"
	ebimgui "github.com/gabstv/ebiten-imgui/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"go.uber.org/zap"
	"log"
	"os"
)

var (
	TileSize  = 16
	TileCount = 25
	TileIndex = 0
)

type LevelDebug struct {
}

func (l *LevelDebug) Draw(screen *ebiten.Image) {
	ebimgui.Update(1.0 / 60.0)
	ebimgui.BeginFrame()
	defer ebimgui.EndFrame()

	imgui.Text("here")

	//pos := numerics.NewVec2(200, 300)
	//im := ebiten.NewImage(16, 16)
	//im.Fill(color.White)
	//screen.DrawImage(im, &ebiten.DrawImageOptions{})

	tile := game.NewTileFromImage(images.Tiles_png, TileCount, TileSize, TileIndex)
	screen.DrawImage(tile.Image, &ebiten.DrawImageOptions{})

	ebimgui.Draw(screen)
}

func (l *LevelDebug) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	width := 600
	height := 600

	ebimgui.SetDisplaySize(float32(width), float32(height))
	return width, height
}

func (l *LevelDebug) Update() error {

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
	ebiten.SetWindowSize(600, 600)
	ebiten.SetWindowTitle("Level Debug")

	if err := ebiten.RunGame(&LevelDebug{}); err != nil {
		log.Fatal(err)
	}
}
