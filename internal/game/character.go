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
	pc := physics.NewObjectFromImages(map[physics.Orientation]*animation.Image{
		physics.Front: animation.WizardFront,
		physics.Left:  animation.WizardSide,
		physics.Right: animation.WizardSide,
	})
	pc.UpdatePosition(float64(screenWidth/2), float64(screenHeight/2))
	return &PlayerCharacter{pc}
}

func (c *PlayerCharacter) Move(camera *Camera) {
	// Handle the movement of the player with the keys
	dx, dy := c.handleKeyPress()
	if dx != 0 || dy != 0 {
		c.Count++
		c.UpdatePosition(dx, dy)
	} else {
		c.Count = 0
	}

	c.handleMouseMovement(camera)
}

func (c *PlayerCharacter) handleCollision(collisionObjects []physics.Collidable) {
}

func (c *PlayerCharacter) handleMouseMovement(camera *Camera) {
	// Handle the rotation of the player to face the direction of the mouse pointer
	mouseX, mouseY := ebiten.CursorPosition()
	mx, my := camera.ScreenToWorld(mouseX, mouseY)
	normal := f64.Vec2{
		mx - c.Position[0],
		my - c.Position[1],
	}
	c.Rotation = c.calculateXAxisAngleFromVec(normal)
}

func (c *PlayerCharacter) handleKeyPress() (float64, float64) {
	var dx, dy float64

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		dy = -1
		c.Orientation = physics.Front
	} else if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		dy = 1
		c.Orientation = physics.Front
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		dx = -1
		c.Orientation = physics.Left
	} else if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		dx = 1
		c.Orientation = physics.Right
	}

	if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		c.Velocity = f64.Vec2{2.0, 2.0}
	} else {
		c.Velocity = f64.Vec2{1.0, 1.0}
	}

	return dx * c.Velocity[0], dy * c.Velocity[1]
}

// calculateXAxisAngleFromVec calculates the roation of the player to face in the direction
// of the vector. So we compute the angle between the vector and the x-axis, then rotate.
func (c *PlayerCharacter) calculateXAxisAngleFromVec(vec f64.Vec2) float64 {
	// Rotate the player
	angle := math.Atan2(vec[1], vec[0])

	// Clip the angle to always be between -90 degrees and 90 degrees in radians if facing right
	if angle < -math.Pi/2 {
		angle = -math.Pi / 2
	}

	if angle > math.Pi/2 {
		angle = math.Pi / 2
	}

	return angle
}
