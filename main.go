package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
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

type AABB struct {
	Min f64.Vec2
	Max f64.Vec2
}

func NewAABB(img *ebiten.Image) *AABB {
	bounds := img.Bounds()
	minBounds := bounds.Min
	maxBounds := bounds.Max
	return &AABB{
		Min: f64.Vec2{float64(minBounds.X), float64(minBounds.Y)},
		Max: f64.Vec2{float64(maxBounds.X), float64(maxBounds.Y)},
	}
}

func (a *AABB) IsColliding2D(b *AABB) bool {
	if a.Max[0] < b.Min[0] || a.Min[0] > b.Max[0] {
		return false
	}

	if a.Max[1] < b.Min[1] || a.Min[1] > b.Max[1] {
		return false
	}

	return true
}

type Collidable interface {
	// IsCollidingInternal checks if an object, which exists WITHIN a bounding volume, is coming near
	// the edge of the shape where it would break out. This is for things like levels.
	IsCollidingInternal(b *Collidable) bool

	// IsCollidingExternal check if an ojbect, which exists OUTSIDE of a bounding volume, is going
	// to clip inside of the edge of the shape. This is for basically everything else.
	IsCollidingExternal(b *Collidable) bool
}

type Game struct {
	character *Character
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

type Character struct {
	count    int
	Position f64.Vec2
	Image    *ebiten.Image

	movementSpeed float64
}

func NewCharacter() *Character {
	zap.L().Info("Loading character")
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}

	return &Character{
		count:         0,
		Position:      f64.Vec2{screenWidth / 4, screenHeight / 4},
		Image:         ebiten.NewImageFromImage(img),
		movementSpeed: 2.0,
	}
}

func (c *Character) Move() {
	dx, dy := c.handleKeyPress()

	if dx != 0 || dy != 0 {
		if c.Position[0]+dx < 0 || c.Position[0]+dx > screenWidth {
			dx = 0
		}

		if c.Position[1]+dy < 0 || c.Position[1]+dy > screenHeight {
			dy = 0
		}

		c.count++
		c.Position[0] += dx
		c.Position[1] += dy
	} else {
		c.count = 0
	}
}

func (c *Character) handleCollision(collisionObjects []Collidable) {
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
	i := (g.character.count / 5) % frameCount
	sx, sy := frameOX+i*frameWidth, frameOY

	screen.DrawImage(g.character.Image.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)

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

	character := NewCharacter()

	zap.L().Info("Starting game")
	if err := ebiten.RunGame(&Game{
		character: character,
		camera:    &Camera{ViewPort: f64.Vec2{screenWidth, screenHeight}},
	}); err != nil {
		log.Fatal(err)
	}
}
