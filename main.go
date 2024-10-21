package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	maps "github.com/hritesh04/shooter-game/maps"
	"github.com/hritesh04/shooter-game/types"
)

const (
	windowWidth  = 1280
	windowHeight = 704
)

type Game struct {
	Width, Height int
	Scale         float64
	World         types.IMap
}

func (g *Game) GetSize() (float64, float64) {
	return float64(g.Width), float64(g.Height)
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	g.World.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.World.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if ebiten.IsFullscreen() {
		return ebiten.WindowSize()
	}
	return windowWidth, windowHeight
}

func main() {
	g := &Game{
		Width:  windowWidth,
		Height: windowHeight,
		Scale:  1.8,
	}
	g.World = maps.NewMap(maps.NewDefMap, g)
	ebiten.SetWindowTitle("Shooter")
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetFullscreen(true)
	ebiten.SetTPS(60)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
