package game

import (
	"dungeon/internal/animation"
	"dungeon/internal/numerics"
	"dungeon/internal/physics"
	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"
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
		physics.Back:  animation.WizardFront,
		physics.Left:  animation.WizardSide,
		physics.Right: animation.WizardSide,
	})
	pc.UpdatePosition(numerics.NewVec2(float64(screenWidth/2), float64(screenHeight/2)))
	return &PlayerCharacter{pc}
}

func (c *PlayerCharacter) Move(camera *Camera) {
	// Handle the movement of the player with the keys
	diff := c.handleKeyPress()

	// Only increment the count when the player is moving, otherwise reset to the start frame.
	if diff.IsZero() {
		c.Count = 0
	} else {
		c.Count++
	}

	c.UpdatePosition(diff)

	c.handleMouseMovement(camera)
}

func (c *PlayerCharacter) handleCollision(collisionObjects []physics.Collidable) {
}

func (c *PlayerCharacter) handleMouseMovement(camera *Camera) {
	// Handle the rotation of the player to face the direction of the mouse pointer
	mouseX, mouseY := ebiten.CursorPosition()
	mx, my := camera.ScreenToWorld(mouseX, mouseY)
	normal := numerics.NewVec2(
		mx-c.Position.X(),
		my-c.Position.Y(),
	)
	c.Rotation = c.calculateXAxisAngleFromVec(normal.Normalized())

	rotDeg := c.Rotation * 180 / math.Pi

	if rotDeg >= -45 && rotDeg <= 45 {
		c.Orientation = physics.Right
	} else if rotDeg >= 45 && rotDeg <= 135 {
		c.Orientation = physics.Front
	} else if rotDeg >= 135 || rotDeg <= -135 {
		c.Orientation = physics.Left
	} else if rotDeg >= -135 && rotDeg <= -45 {
		c.Orientation = physics.Back
	}
}

func (c *PlayerCharacter) handleKeyPress() numerics.Vec2 {
	diff := numerics.Zero()

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		diff = diff.Add(numerics.NewVec2(0, -1))
	} else if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		diff = diff.Add(numerics.NewVec2(0, 1))
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		diff = diff.Add(numerics.NewVec2(-1, 0))
	} else if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		diff = diff.Add(numerics.NewVec2(1, 0))
	}

	if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		c.Velocity = numerics.One().MulScalar(2)
	} else {
		c.Velocity = numerics.One()
	}

	return diff.Mul(c.Velocity)
}

// calculateXAxisAngleFromVec calculates the angle of the vector with respect to the x-axis. Assumes that the input
// vector is normalized.
func (c *PlayerCharacter) calculateXAxisAngleFromVec(vec numerics.Vec2) float64 {
	return math.Atan2(vec.Y(), vec.X())
}
