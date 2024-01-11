package game

import (
	"dungeon/internal/gfx"
	"dungeon/internal/numerics"
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
)

const (
	// TileSize is the size of a tile in pixels
	TileSize = 16
)

func init() {
	//GrassEmpty = NewTileFromImage(images.Tiles_png, 288, 139, TileSize)
	//GrassLevel = NewRoom(
	//	numerics.NewVec2(gfx.ScreenWidth/2-500, gfx.ScreenHeight/2-250),
	//	numerics.NewVec2(adjustToTileSize(1000), adjustToTileSize(500)))
}

type Level struct {
	rooms []*Room

	currentRoom int
}

func NewLevel() *Level {
	// Generate a random number of rooms between 10-20
	nRooms := 10 + rand.Intn(10)
	rooms := make([]*Room, nRooms)

	for i := 0; i < nRooms; i++ {
		// Random number between 1000-2000
		roomWidth := 500 + rand.Intn(1000)
		roomHeight := 500 + rand.Intn(1000)

		position := numerics.NewVec2(
			gfx.ScreenWidth/2-float64(roomWidth)/2,
			gfx.ScreenHeight/2-float64(roomHeight)/2,
		)
		dimensions := numerics.NewVec2(
			adjustToTileSize(roomWidth),
			adjustToTileSize(roomHeight),
		)

		rooms[i] = NewRoom(position, dimensions)
	}

	// Every room has at least one door, and up to 2 more
	for i := 0; i < nRooms-1; i++ {
		// Door boundary position

		//rooms[i].Doors = append(rooms[i].Doors, NewDoor())
	}

	return &Level{
		rooms:       rooms,
		currentRoom: 0,
	}
}

func (l *Level) CurrentRoom() *Room {
	return l.rooms[l.currentRoom]
}

func (l *Level) Render(screen *ebiten.Image, cameraTransform *ebiten.GeoM) {
	l.CurrentRoom().Render(screen, cameraTransform)
}
