package main

import (
	"dungeon/internal/game"
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
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Dungeon")

	playerCharacter := game.NewPlayerCharacter(game.ScreenWidth, game.ScreenHeight)

	zap.L().Info("Starting game")
	if err := ebiten.RunGame(&game.Game{
		PlayerCharacter: playerCharacter,
		Camera:          &game.Camera{ViewPort: f64.Vec2{game.ScreenWidth, game.ScreenHeight}},
	}); err != nil {
		log.Fatal(err)
	}
}
