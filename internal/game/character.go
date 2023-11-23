package game

import (
	"dungeon/internal/animation"
	"dungeon/internal/physics"
	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"
	"golang.org/x/image/math/f64"
	"math"
)

// PlayerCharacter is a player character
type PlayerCharacter struct {
	*physics.Object
}

func NewPlayerCharacter(screenWidth, screenHeight int) *PlayerCharacter {
	zap.L().Info("Loading player character")
	pc := physics.NewObjectFromImage(animation.Runner)
	//pc.Position = f64.Vec2{float64(screenWidth / 4), float64(screenHeight / 4)}
	pc.UpdatePosition(float64(screenWidth/2), float64(screenHeight/2))
	return &PlayerCharacter{pc}

	//return &PlayerCharacter{
	//	Count:         0,
	//	Position:      f64.Vec2{float64(screenWidth / 4), float64(screenHeight / 4)},
	//	Image:         ebiten.NewImageFromImage(img),
	//	movementSpeed: 2.0,
	//}
}

func (c *PlayerCharacter) Move(camera *Camera) {
	// Handle the movement of the player with the keys
	dx, dy := c.handleKeyPress()
	if dx != 0 || dy != 0 {
		c.Count++
		//c.Position[0] += dx
		//c.Position[1] += dy
		c.UpdatePosition(dx, dy)
	} else {
		c.Count = 0
	}

	// Handle the rotation of the player to face the direction of the mouse pointer
	mouseX, mouseY := ebiten.CursorPosition()
	mx, my := camera.ScreenToWorld(mouseX, mouseY)
	normal := f64.Vec2{
		mx - c.Position[0],
		my - c.Position[1],
	}
	c.Rotation = c.calculateRotationFromVector(normal)
}

func (c *PlayerCharacter) handleCollision(collisionObjects []physics.Collidable) {
}

func (c *PlayerCharacter) handleKeyPress() (float64, float64) {
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

	return dx * c.Velocity[0], dy * c.Velocity[1]
}

// calculateRotationFromVector calculates the roation of the player to face in the direction
// of the vector. So we compute the angle between the vector and the x-axis, then rotate.
func (c *PlayerCharacter) calculateRotationFromVector(vec f64.Vec2) float64 {
	// Rotate the player
	return math.Atan2(vec[1], vec[0])
}
