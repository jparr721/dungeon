package main

import (
	"dungeon/internal/game"
	"dungeon/internal/gfx"
	"dungeon/internal/numerics"

	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"
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
	objects := make([]*game.Object, 0)
	objects = append(
		objects,
		playerCharacter.Object,
	)

	if err := ebiten.RunGame(&game.Game{
		PlayerCharacter: playerCharacter,
		Camera:          &game.Camera{ViewPort: numerics.NewVec2(gfx.ScreenWidth, gfx.ScreenHeight)},
		CurrentLevel:    game.NewLevel(),
		Objects:         objects,
	}); err != nil {
		log.Fatal(err)
	}
}
