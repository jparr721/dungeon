package main

import (
	"dungeon/internal/game"
	"dungeon/internal/gfx"
	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"
	"golang.org/x/image/math/f64"
	_ "image/png"
	"log"
	"os"
)

func init() {
	logger := zap.Must(zap.NewDevelopment())
	if os.Getenv("APP_ENV") == "release" {
		logger = zap.Must(zap.NewProduction())
	}
	zap.ReplaceGlobals(logger)
}

func main() {
	ebiten.SetWindowSize(gfx.ScreenWidth, gfx.ScreenHeight)
	ebiten.SetWindowTitle("Dungeon")

	playerCharacter := game.NewPlayerCharacter(gfx.ScreenWidth, gfx.ScreenHeight)

	zap.L().Info("Starting game")
	if err := ebiten.RunGame(&game.Game{
		PlayerCharacter: playerCharacter,
		Camera:          &game.Camera{ViewPort: f64.Vec2{gfx.ScreenWidth, gfx.ScreenHeight}},
		CurrentLevel:    game.GrassLevel,
	}); err != nil {
		log.Fatal(err)
	}
}
