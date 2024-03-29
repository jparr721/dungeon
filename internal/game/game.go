package game

import (
	"dungeon/internal/gfx"
	"dungeon/internal/numerics"
	"fmt"
	ebimgui "github.com/gabstv/ebiten-imgui/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"math"
)

type Game struct {
	PlayerCharacter *PlayerCharacter
	Camera          *Camera
	CurrentLevel    *Level
	Objects         []*Object
}

func (g *Game) Update() error {
	ebimgui.Update(1.0 / 60.0)
	ebimgui.BeginFrame()
	defer ebimgui.EndFrame()

	for _, a := range g.Objects {
		a.ResetCollisionState()
	}

	for _, a := range g.Objects {
		for _, b := range g.Objects {
			if a == b {
				continue
			}

			a.IsExternallyColliding2D(b.AABB)
		}
	}

	g.PlayerCharacter.Move(g.Camera, g.Objects, g.CurrentLevel.CurrentRoom())

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.PlayerCharacter.FireProjectile()
	}

	for _, proj := range g.PlayerCharacter.Projectiles {
		proj.Step()
	}

	// Camera is always centered on the main PlayerCharacter
	g.Camera.Position = numerics.NewVec2(
		g.PlayerCharacter.Position.X()-gfx.ScreenWidth/2,
		g.PlayerCharacter.Position.Y()-gfx.ScreenHeight/2,
	)

	return nil
}

// Draw is the main draw function for the game. It handles drawing all Object types to the screen.
func (g *Game) Draw(screen *ebiten.Image) {
	// Get the Camera matrix transform
	cameraTransform := g.Camera.worldMatrix()

	// Render the level before the character otherwise it'll draw overtop of it.
	g.CurrentLevel.Render(screen, &cameraTransform)

	// Draw the PlayerCharacter and translate them to whatever their current position is
	g.PlayerCharacter.Render(screen, &cameraTransform)

	// Draw the projectiles
	for _, proj := range g.PlayerCharacter.Projectiles {
		proj.Render(screen, &cameraTransform)
	}

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("TPS: %0.2f, FPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()),
	)

	mx, my := ebiten.CursorPosition()
	ex, ey := g.Camera.ScreenToWorld(mx, my)
	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf(
			"Pos x: %.2f, y: %.2f; Center x: %.2f, y: %.2f; Mouse x: %.2f, y: %.2f",
			g.PlayerCharacter.Position.X(),
			g.PlayerCharacter.Position.Y(),
			g.PlayerCharacter.Center.X(),
			g.PlayerCharacter.Center.Y(),
			ex,
			ey,
		),
		0, gfx.ScreenHeight-32,
	)

	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("Camera %s", g.Camera.String()),
		0, gfx.ScreenHeight-64,
	)

	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("Player Rotation (Degrees) %.2f (Radians) %.2f", g.PlayerCharacter.Rotation*180/math.Pi, g.PlayerCharacter.Rotation),
		0, gfx.ScreenHeight-96,
	)

	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("Bounding Box %s", g.PlayerCharacter.AABB.String()),
		0, gfx.ScreenHeight-108,
	)

	ebimgui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	ebimgui.SetDisplaySize(float32(gfx.ScreenWidth), float32(gfx.ScreenHeight))
	return gfx.ScreenWidth, gfx.ScreenHeight
}
