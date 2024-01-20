package game

import (
	"dungeon/internal/animation"
	"dungeon/internal/gfx"
	"dungeon/internal/numerics"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math/rand"
	"slices"
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
	doors []*Door

	currentRoom int
}

func NewLevel() *Level {
	// Generate a random number of rooms between 10-20
	nRooms := 10 + rand.Intn(10)
	rooms := make([]*Room, nRooms)

	for i := 0; i < nRooms; i++ {
		// Random number between 1000-2000
		roomWidth := adjustToTileSize(500 + rand.Intn(1000))
		roomHeight := adjustToTileSize(500 + rand.Intn(1000))

		position := numerics.NewVec2(
			gfx.ScreenWidth/2-float64(roomWidth)/2,
			gfx.ScreenHeight/2-float64(roomHeight)/2,
		)
		dimensions := numerics.NewVec2(roomWidth, roomHeight)

		rooms[i] = NewRoom(position, dimensions)
	}

	// Every room has at least one door, and up to 2 more
	for i := 0; i < nRooms-1; i++ {
		// Door exists on the boundary of a wall, but anywhere along that wall
		// Randomly choose a wall

		// 0 - Left
		// 1 - Right
		// 2 - Top
		// 3 - Bottom
		nWalls := 1

		if rand.Float64() < 0.1 {
			nWalls = 2
		}

		usedWalls := make([]int, 0)
		for w := 0; w < nWalls; w++ {
			wall := rand.Intn(4)
			usedWalls = append(usedWalls, wall)

			// Loop until the wall is not in usedWalls
			for slices.Contains(usedWalls, wall) {
				wall = rand.Intn(4)
			}

			minX := rooms[i].Position.X()
			minY := rooms[i].Position.Y()

			endPos := rooms[i].Position.Add(rooms[i].Dimensions).SubScalar(float64(rooms[i].StrokeWidth / 2))

			maxX := endPos.X()
			maxY := endPos.Y()

			// Halfway between minX and maxX
			halfX := minX + (maxX-minX)/2
			halfY := minY + (maxY-minY)/2

			// Put the door 50% of the way along the wall
			var doorPosition numerics.Vec2
			switch wall {
			case 0: // Left
				doorPosition = numerics.NewVec2(
					minX,
					halfY,
				)
				break
			case 1: // Right
				doorPosition = numerics.NewVec2(
					maxX-10,
					halfY,
				)
				break
			case 2: // Top
				doorPosition = numerics.NewVec2(
					halfX,
					minY,
				)
				break
			case 3: // Bottom
				doorPosition = numerics.NewVec2(
					halfX,
					maxY-10,
				)
				break
			}

			width := (rooms[i].StrokeWidth / 2) * 1.5
			height := (rooms[i].StrokeWidth / 2) * 1.5

			doorImg := animation.NewImageFromImage(ebiten.NewImage(int(width), int(height)))
			doorImg.Fill(color.White)
			rooms[i].Doors = append(rooms[i].Doors, NewDoor(doorPosition, rooms[i+1], doorImg))
		}
	}

	return &Level{
		rooms:       rooms,
		currentRoom: 0,
	}
}

func (l *Level) Doors() []*Door {
	doors := make([]*Door, 0)
	for _, room := range l.rooms {
		for _, door := range room.Doors {
			if !slices.Contains(doors, door) {
				doors = append(doors, door)
			}
		}
	}

	return doors
}

func (l *Level) CurrentRoom() *Room {
	return l.rooms[l.currentRoom]
}

func (l *Level) Render(screen *ebiten.Image, cameraTransform *ebiten.GeoM) {
	l.CurrentRoom().Render(screen, cameraTransform)
}
