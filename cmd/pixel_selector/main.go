package main

import (
	"dungeon/internal/animation"
	"dungeon/internal/game"
	"fmt"
	imgui "github.com/gabstv/cimgui-go"
	ebimgui "github.com/gabstv/ebiten-imgui/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"go.uber.org/zap"
	"log"
	"os"
)

var (
	TileSize                  = 16
	TileCount                 = 25
	TileIndex                 = 47
	ImageIDRef                = 10
	DebugIDRef                = 11
	ViewerWidth               = 1920
	ViewerHeight              = 1080
	PictureWindowOpen         = true
	InspectorWindowOpen       = true
	PixelSelectorDropdownOpen = true
	LockedWindowFlags         = imgui.WindowFlags(
		imgui.WindowFlagsNoTitleBar |
			imgui.WindowFlagsNoResize |
			imgui.WindowFlagsNoResize,
	)
)

type PixelSelector struct {
	IsHovered bool
	MousePos  imgui.Vec2
}

type Game struct {
	LoadedImage *ebiten.Image
	ImageWidth  float32
	ImageHeight float32

	selector *PixelSelector
}

func (p *PixelSelector) DrawPixelSelector(img *ebiten.Image) {
	imgui.SetNextWindowSize(imgui.NewVec2(float32(ViewerWidth), float32(ViewerHeight/2)))
	imgui.SetNextWindowPos(imgui.NewVec2(0, 0))
	imgui.BeginV(
		"Picture",
		&PictureWindowOpen,
		LockedWindowFlags,
	)
	defer imgui.End()

	io := imgui.CurrentIO()

	screenPos := imgui.CursorScreenPos()
	imgSize := img.Bounds().Size()
	ebimgui.GlobalManager().Cache.SetTexture(imgui.TextureID(&ImageIDRef), img)
	tid := imgui.TextureID(&ImageIDRef)
	game.Image(tid, imgui.NewVec2(float32(imgSize.X), float32(imgSize.Y)))
	if imgui.BeginItemTooltip() {
		p.IsHovered = imgui.IsItemHovered()
		regionSz := float32(16.0)

		regionX := io.MousePos().X - screenPos.X - regionSz*0.5
		regionY := io.MousePos().Y - screenPos.Y - regionSz*0.5

		if regionX < 0 {
			regionX = 0
		} else if regionX > float32(imgSize.X)-regionSz {
			regionX = float32(imgSize.X) - regionSz
		}

		if regionY < 0 {
			regionY = 0
		} else if regionY > float32(imgSize.Y)-regionSz {
			regionY = float32(imgSize.Y) - regionSz
		}

		tile := game.NewTileFromImage(images.Tiles_png, int(regionX), int(regionY), int(regionSz))

		imgui.Text(fmt.Sprintf("Cursor: (%.2f, %.2f)", regionX, regionY))
		imgui.Text("Bounds:")
		imgui.Text(fmt.Sprintf("Min: (%.2f, %.2f)", regionX, regionY))
		imgui.Text(fmt.Sprintf("Max: (%.2f, %.2f)", regionX+regionSz, regionY+regionSz))

		game.ImageTile(&DebugIDRef, tile, imgui.NewVec2(regionSz*6, regionSz*6))

		imgui.EndTooltip()
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebimgui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	ebimgui.SetDisplaySize(float32(ViewerWidth), float32(ViewerHeight))
	return ViewerWidth, ViewerHeight
}

func (g *Game) Update() error {
	ebimgui.Update(1.0 / 60.0)
	ebimgui.BeginFrame()
	defer ebimgui.EndFrame()

	g.selector.DrawPixelSelector(g.LoadedImage)

	imgui.SetNextWindowSize(imgui.NewVec2(float32(ViewerWidth), float32(ViewerHeight/2)))
	imgui.SetNextWindowPos(imgui.NewVec2(0, float32(ViewerHeight/2)))

	imgui.BeginV(
		"Inspector",
		&InspectorWindowOpen,
		LockedWindowFlags,
	)
	imgui.BeginTabBar("Inspector")
	imgui.BeginTabItem("Inspector")
	if imgui.CollapsingHeaderBoolPtrV("Pixel Selector", &PixelSelectorDropdownOpen, imgui.TreeNodeFlagsDefaultOpen) {
		if imgui.BeginTable("##Image", 2) {
			game.TableRow("Global Mouse Pos", imgui.MousePos())
			game.TableRow("Image Dimensions", g.LoadedImage.Bounds().Size())
			imgui.EndTable()
		}
	}
	imgui.EndTabItem()
	imgui.EndTabBar()
	imgui.End()
	return nil
}

func init() {
	logger := zap.Must(zap.NewDevelopment())
	if os.Getenv("APP_ENV") == "release" {
		logger = zap.Must(zap.NewProduction())
	}
	zap.ReplaceGlobals(logger)
}

func main() {
	loadedImage, err := animation.LoadImage(images.Tiles_png)
	if err != nil {
		log.Fatal(err)
	}
	point := loadedImage.Bounds().Size()
	imageWidth := float32(point.X)
	imageHeight := float32(point.Y)

	ebiten.SetWindowSize(ViewerWidth, ViewerHeight)
	ebiten.SetWindowTitle("Level Debug")

	if err := ebiten.RunGame(&Game{
		LoadedImage: loadedImage,
		ImageWidth:  imageWidth,
		ImageHeight: imageHeight,
		selector:    &PixelSelector{MousePos: imgui.NewVec2(0, 0)},
	}); err != nil {
		log.Fatal(err)
	}
}
