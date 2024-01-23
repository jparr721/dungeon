package game

import (
	"dungeon/internal/numerics"
	"github.com/hajimehoshi/ebiten/v2"
)

func MousePosition() numerics.Vec2 {
	x, y := ebiten.CursorPosition()
	return numerics.NewVec2(float64(x), float64(y))
}
