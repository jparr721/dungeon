package game

import (
	"fmt"
	imgui "github.com/gabstv/cimgui-go"
	ebimgui "github.com/gabstv/ebiten-imgui/v3"
)

func Image(tid imgui.TextureID, size imgui.Vec2) {
	uv0 := imgui.NewVec2(0, 0)
	uv1 := imgui.NewVec2(1, 1)
	borderCol := imgui.NewVec4(1, 1, 1, 1)
	tintCol := imgui.NewVec4(1, 1, 1, 1)
	imgui.ImageV(tid, size, uv0, uv1, tintCol, borderCol)
}

func ImageTile(textureIDRef *int, tile *Tile, size imgui.Vec2) {
	ebimgui.GlobalManager().Cache.SetTexture(imgui.TextureID(textureIDRef), tile.Image)
	uv0 := imgui.NewVec2(0, 0)
	uv1 := imgui.NewVec2(1, 1)
	tid := imgui.TextureID(textureIDRef)
	borderCol := imgui.NewVec4(1, 0, 0, 1)
	tintCol := imgui.NewVec4(1, 1, 1, 1)
	imgui.ImageV(tid, size, uv0, uv1, tintCol, borderCol)
}

func TableRow(rowName string, rowValue any) {
	imgui.TableNextColumn()
	imgui.Text(rowName)
	imgui.TableNextColumn()
	imgui.Text(fmt.Sprint(rowValue))
}
