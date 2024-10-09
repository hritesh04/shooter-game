package main

import (
	"image"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hritesh04/shooter-game/entities/player"
	"github.com/hritesh04/shooter-game/maps"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/resolv"
)

const (
	windowWidth  = 860
	windowHeight = 860
)

var (
	playerImage *ebiten.Image
	tilesImage  *ebiten.Image
)

type Game struct {
	count         int
	player        *player.Player
	inputSystem   input.System
	scale         float64
	tiledMap      *maps.TiledMapJSON
	tileMapImage  *ebiten.Image
	obstacleSpace *resolv.Space
}

func init() {
	var err error
	tilesImage, _, err = ebitenutil.NewImageFromFile("assets/mapSheet.png")
	if err != nil {
		log.Fatal(err)
	}

	playerImage, _, err = ebitenutil.NewImageFromFile("assets/runner.png")
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	playerObj := g.player.Src
	moved := false

	if g.player.Input.ActionIsPressed(player.ActionMoveLeft) {
		if collision := playerObj.Check(-2, 0, "obstacle"); collision == nil {
			playerObj.Position.X -= 2
			moved = true
		}
	}
	if g.player.Input.ActionIsPressed(player.ActionMoveRight) {
		if collision := playerObj.Check(2, 0, "obstacle"); collision == nil {
			playerObj.Position.X += 2
			moved = true
		}
	}
	if g.player.Input.ActionIsPressed(player.ActionMoveUp) {
		if collision := playerObj.Check(0, -2, "obstacle"); collision == nil {
			playerObj.Position.Y -= 2
			moved = true
		}
	}
	if g.player.Input.ActionIsPressed(player.ActionMoveDown) {
		if collision := playerObj.Check(0, 2, "obstacle"); collision == nil {
			playerObj.Position.Y += 2
			moved = true
		}
	}

	if moved {
		playerObj.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}

	for _, layer := range g.tiledMap.Layers {
		for index, id := range layer.Data {
			if id == 0 {
				continue // Skip empty tiles
			}

			x := float64((index % layer.Width) * 16)
			y := float64((index / layer.Width) * 16)

			srcX := ((id - 1) % 12) * 16
			srcY := ((id - 1) / 24) * 16

			opts.GeoM.Reset()
			opts.GeoM.Scale(g.scale, g.scale)
			opts.GeoM.Translate(x*g.scale, y*g.scale)
			screen.DrawImage(g.tileMapImage.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image), &opts)

			if layer.Name == "obstacle" {
				obstacle := resolv.NewObject(x*g.scale, y*g.scale, 10*g.scale, 10*g.scale, "obstacle")
				g.obstacleSpace.Add(obstacle)
			}
		}
	}

	g.player.Draw(screen)

	objs := g.obstacleSpace.Objects()
	for _, obj := range objs {
		if obj.HasTags("obstacle") {
			vector.DrawFilledRect(screen, float32(obj.Position.X), float32(obj.Position.Y), float32(obj.Size.X), float32(obj.Size.Y), color.RGBA{0, 0, 255, 128}, true)
		}
		if obj.HasTags("player") {
			vector.DrawFilledRect(screen, float32(obj.Position.X), float32(obj.Position.Y), float32(obj.Size.X), float32(obj.Size.Y), color.RGBA{0, 0, 255, 128}, true)
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowWidth, windowHeight
}

func main() {
	gameMap, err := maps.NewMap(maps.DefaultMap)
	if err != nil {
		log.Fatal(err)
	}
	mapImage, _, err := ebitenutil.NewImageFromFile("assets/dungeonSheet.png")
	if err != nil {
		log.Fatal(err)
	}
	g := &Game{
		scale:         1.8,
		tiledMap:      gameMap,
		tileMapImage:  mapImage,
		obstacleSpace: resolv.NewSpace(windowWidth, windowHeight, 16, 16),
	}
	g.player = player.NewPlayer(g.obstacleSpace, playerImage)
	g.inputSystem.Init(input.SystemConfig{DevicesEnabled: input.AnyDevice})
	g.player.Input = g.inputSystem.NewHandler(0, player.Keymap)
	ebiten.SetWindowTitle("Shooter")
	ebiten.SetWindowSize(windowWidth, windowHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
