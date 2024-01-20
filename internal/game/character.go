package game

import (
	"dungeon/internal/animation"
	"dungeon/internal/numerics"

	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"
	"math"
)

// PlayerCharacter is a player character
type PlayerCharacter struct {
	*Object
}

func NewPlayerCharacter(screenWidth, screenHeight int) *PlayerCharacter {
	zap.L().Info("Loading player character")
	pc := NewObjectFromImages(map[Orientation]*animation.Image{
		Front: animation.WizardFront,
		Back:  animation.WizardFront,
		Left:  animation.WizardSide,
		Right: animation.WizardSide,
	})
	pc.UpdatePosition(numerics.NewVec2(float64(screenWidth/2), float64(screenHeight/2)))
	return &PlayerCharacter{pc}
}

func (c *PlayerCharacter) Move(camera *Camera, objects []*Object, room *Room) {
	// Handle the movement of the player with the keys
	diff := c.handleKeyPress()

	// Check if the player is colliding with the boundary of the room.
	diff = room.CheckCollisionAndUpdatePosition(c.Object, diff)

	// Depending on the collision axis, prevent movement in diff
	if c.IsColliding {
		// Keep track of the previous diff so we can undo it
		oldDiff := diff

		// First, apply the diff to the position
		c.UpdatePosition(oldDiff)

		// Does this move relieve the collision?
		anyCollision := false
		for _, a := range objects {
			if a == c.Object {
				continue
			}

			// Check for a bounding box collision
			if c.IsExternallyColliding2D(a.AABB) {
				anyCollision = true
			}
		}

		// If any of these are colliding, restrict the motion along the collision axis
		if anyCollision {
			if c.CollisionDirection.X {
				diff = numerics.NewVec2(0, diff.Y())
			}

			if c.CollisionDirection.Y {
				diff = numerics.NewVec2(diff.X(), 0)
			}
		}

		// Undo the position update.
		c.UpdatePosition(oldDiff.MulScalar(-1))
	}

	// Only increment the count when the player is moving, otherwise reset to the start frame.
	if diff.IsZero() {
		c.Count = 0
	} else {
		c.Count++
	}

	c.UpdatePosition(diff)
	c.handleMouseMovement(camera)
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
		c.Orientation = Right
	} else if rotDeg >= 45 && rotDeg <= 135 {
		c.Orientation = Front
	} else if rotDeg >= 135 || rotDeg <= -135 {
		c.Orientation = Left
	} else if rotDeg >= -135 && rotDeg <= -45 {
		c.Orientation = Back
	}
}

func (c *PlayerCharacter) handleKeyPress() numerics.Vec2 {
	diff := numerics.ZeroVec2()

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
		c.Velocity = numerics.OneVec2().MulScalar(2)
	} else {
		c.Velocity = numerics.OneVec2()
	}

	return diff.Mul(c.Velocity)
}

// calculateXAxisAngleFromVec calculates the angle of the vector with respect to the x-axis. Assumes that the input
// vector is normalized.
func (c *PlayerCharacter) calculateXAxisAngleFromVec(vec numerics.Vec2) float64 {
	return math.Atan2(vec.Y(), vec.X())
}
