package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/math/f64"
	"image/color"
)

const (
	ScreenWidth  = 1920
	ScreenHeight = 1080
)

type Game struct {
	PlayerCharacter *PlayerCharacter
	Camera          *Camera
}

func (g *Game) Update() error {
	g.PlayerCharacter.Move(g.Camera)

	// Camera is always centered on the main PlayerCharacter
	g.Camera.Position = f64.Vec2{
		g.PlayerCharacter.Position[0] - ScreenWidth/2,
		g.PlayerCharacter.Position[1] - ScreenHeight/2,
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Get the Camera matrix transform
	cameraTransform := g.Camera.worldMatrix()

	bg := ebiten.NewImage(1000, 500)
	bg.Fill(color.White)

	op := &ebiten.DrawImageOptions{
		GeoM: cameraTransform,
	}
	op.GeoM.Translate(ScreenWidth/2-500, ScreenHeight/2-250)
	screen.DrawImage(bg, op)

	// Draw the PlayerCharacter and translate them to whatever their current position is
	g.PlayerCharacter.Render(screen, &cameraTransform)

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("TPS: %0.2f, FPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()),
	)

	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf(
			"Pos x: %.2f, y: %.2f; Center x: %.2f, y: %.2f",
			g.PlayerCharacter.Position[0],
			g.PlayerCharacter.Position[1],
			g.PlayerCharacter.Center[0],
			g.PlayerCharacter.Center[1]),
		0, ScreenHeight-32,
	)

	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("Camera %s", g.Camera.String()),
		0, ScreenHeight-64,
	)

	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("Player Rotation %.2f", g.PlayerCharacter.Rotation),
		0, ScreenHeight-96,
	)

	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("Bounding Box %s", g.PlayerCharacter.AABB.String()),
		0, ScreenHeight-108,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
