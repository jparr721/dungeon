package main

import (
	"dungeon/internal/character"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"go.uber.org/zap"
	"golang.org/x/image/math/f64"
)

const (
	screenWidth  = 1920
	screenHeight = 1080

	frameOX     = 0
	frameOY     = 32
	frameWidth  = 32
	frameHeight = 32
	frameCount  = 8
)

func init() {
	logger := zap.Must(zap.NewDevelopment())
	if os.Getenv("APP_ENV") == "release" {
		logger = zap.Must(zap.NewProduction())
	}
	zap.ReplaceGlobals(logger)
}

type Game struct {
	character *character.Character
	camera    *Camera
}

type Camera struct {
	ViewPort     f64.Vec2
	Position     f64.Vec2
	ZoomFactor   float64
	ZoomFactorTo float64
	Rotation     float64
}

func (c *Camera) String() string {
	return fmt.Sprintf(
		"T: %.1f, R: %.2f, S: %.2f",
		c.Position, c.Rotation, c.ZoomFactor,
	)
}

func (c *Camera) viewportCenter() f64.Vec2 {
	return f64.Vec2{
		c.ViewPort[0] * 0.5,
		c.ViewPort[1] * 0.5,
	}
}

func (c *Camera) worldMatrix() ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-c.Position[0], -c.Position[1])
	// We want to scale and rotate around center of image / screen
	m.Translate(-c.viewportCenter()[0], -c.viewportCenter()[1])
	m.Scale(
		math.Pow(1.01, float64(c.ZoomFactor)),
		math.Pow(1.01, float64(c.ZoomFactor)),
	)
	m.Rotate(float64(c.Rotation) * 2 * math.Pi / 360)
	m.Translate(c.viewportCenter()[0], c.viewportCenter()[1])
	return m
}

func (c *Camera) ScreenToWorld(posX, posY int) (float64, float64) {
	inverseMatrix := c.worldMatrix()
	if inverseMatrix.IsInvertible() {
		inverseMatrix.Invert()
		return inverseMatrix.Apply(float64(posX), float64(posY))
	} else {
		// When scaling it can happened that matrix is not invertable
		return math.NaN(), math.NaN()
	}
}

func (g *Game) Update() error {
	g.character.Move()

	// Camera is always centered on the main character
	g.camera.Position = f64.Vec2{
		g.character.Position[0] - screenWidth/4,
		g.character.Position[1] - screenHeight/4,
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Get the camera matrix transform
	cameraTransform := g.camera.worldMatrix()

	bg := ebiten.NewImage(1000, 500)
	bg.Fill(color.White)

	op := &ebiten.DrawImageOptions{
		GeoM: cameraTransform,
	}
	op.GeoM.Translate(screenWidth/4, screenHeight/4)
	screen.DrawImage(bg, op)

	// Draw the character and translate them to whatever their current position is
	op = &ebiten.DrawImageOptions{
		GeoM: cameraTransform,
	}
	op.GeoM.Translate(g.character.Position[0], g.character.Position[1])

	// Move the character to the start of their frame
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)

	// Scale by 2 since it looks kind of small
	op.GeoM.Scale(2.0, 2.0)

	// This just chooses the character frame from the sprite sheet. We divide by 5 so that way the transition
	// between animation frames is less intense (basically going at 5 frames per second).
	i := (g.character.Count / 5) % frameCount
	sx, sy := frameOX+i*frameWidth, frameOY

	screen.DrawImage(g.character.Image.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("TPS: %0.2f, FPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()),
	)

	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("Pos x: %.2f, y: %.2f", g.character.Position[0], g.character.Position[1]),
		0, screenHeight-32,
	)

	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("Camera %s", g.camera.String()),
		0, screenHeight-64,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Dungeon")

	character := character.NewCharacter(screenWidth, screenHeight)

	zap.L().Info("Starting game")
	if err := ebiten.RunGame(&Game{
		character: character,
		camera:    &Camera{ViewPort: f64.Vec2{screenWidth, screenHeight}},
	}); err != nil {
		log.Fatal(err)
	}
}
