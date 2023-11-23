package character

import (
	"bytes"
	"dungeon/internal/physics"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"go.uber.org/zap"
	"golang.org/x/image/math/f64"
)

type Character struct {
	Count    int
	Position f64.Vec2
	Image    *ebiten.Image

	movementSpeed float64
}

func NewCharacter(screenWidth, screenHeight int) *Character {
	zap.L().Info("Loading character")
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}

	return &Character{
		Count:         0,
		Position:      f64.Vec2{float64(screenWidth / 4), float64(screenHeight / 4)},
		Image:         ebiten.NewImageFromImage(img),
		movementSpeed: 2.0,
	}
}

func (c *Character) Move() {
	dx, dy := c.handleKeyPress()

	if dx != 0 || dy != 0 {
		c.Count++
		c.Position[0] += dx
		c.Position[1] += dy
	} else {
		c.Count = 0
	}
}

func (c *Character) handleCollision(collisionObjects []physics.Collidable) {
}

func (c *Character) handleKeyPress() (float64, float64) {
	var dx, dy float64

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		dx = -1
	} else if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		dx = 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		dy = -1
	} else if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		dy = 1
	}

	return dx * c.movementSpeed, dy * c.movementSpeed
}
