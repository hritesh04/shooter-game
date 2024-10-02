package main

import (
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	windowWidth  = 1080
	windowHeight = 720

	frameOX     = 0
	frameOY     = 32
	frameWidth  = 32
	frameHeight = 32
	frameCount  = 8

	tileSize  = 128
	mapWidth  = 20
	mapHeight = 20
)

var (
	playerImage *ebiten.Image
	tilesImage  *ebiten.Image
)

type layer struct {
	tileID int
	rotate float64
}

type Game struct {
	count  int
	scale  float64
	screen [][]layer
}

func init() {
	var err error
	tilesImage, _, err = ebitenutil.NewImageFromFile("./assets/mapSheet.png")
	if err != nil {
		log.Fatal(err)
	}

	playerImage, _, err = ebitenutil.NewImageFromFile("./assets/runner.png")
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	g.count++
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		fullscreen := !ebiten.IsFullscreen()
		ebiten.SetFullscreen(fullscreen)
		g.updateScale()
	}
	return nil
}

func (g *Game) updateScale() {
	if ebiten.IsFullscreen() {
		w, _ := ebiten.Monitor().Size()
		g.scale = float64(w) / float64(mapWidth*tileSize)
	} else {
		g.scale = float64(windowWidth) / float64(mapWidth*tileSize)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 0})
	// Draw the layer
	tileXCount := tilesImage.Bounds().Dx() / tileSize
	scaledTileSize := int(float64(tileSize) * g.scale)

	for y, row := range g.screen {
		for x, tile := range row {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(g.scale, g.scale)
			op.GeoM.Rotate(tile.rotate * math.Pi / 180.0)
			op.GeoM.Translate(float64(x*scaledTileSize), float64(y*scaledTileSize))
			sx := (tile.tileID % tileXCount) * tileSize
			sy := (tile.tileID / tileXCount) * tileSize
			screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
		}
	}

	// Draw the player
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.scale*4, g.scale*4)
	op.GeoM.Translate(float64(windowWidth/2-frameWidth)*g.scale, float64(windowHeight/2-frameHeight)*g.scale)
	i := (g.count / 5) % frameCount
	sx, sy := frameOX+i*frameWidth, frameOY
	screen.DrawImage(playerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if ebiten.IsFullscreen() {
		return ebiten.Monitor().Size()
	}
	return windowWidth, windowHeight
}

func main() {
	g := &Game{
		scale: 1.0,
		screen: [][]layer{
			{{tileID: 9, rotate: 0}, {tileID: 10, rotate: 0}, {tileID: 10, rotate: 0}, {tileID: 10, rotate: 0}, {tileID: 10, rotate: 0}, {tileID: 10, rotate: 0}, {tileID: 10, rotate: 0}, {tileID: 10, rotate: 0}, {tileID: 10, rotate: 0}, {tileID: 10, rotate: 0}, {tileID: 10, rotate: 0}, {tileID: 10, rotate: 0}, {tileID: 10, rotate: 0}, {tileID: 10, rotate: 0}, {tileID: 10, rotate: 0}, {tileID: 10, rotate: 0}, {tileID: 10, rotate: 0}, {tileID: 10, rotate: 0}, {tileID: 10, rotate: 0}, {tileID: 10, rotate: 0}, {tileID: 9, rotate: 90.0}},
			{{tileID: 82, rotate: 0}},
		},
	}

	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("Shooter")
	ebiten.SetWindowSize(windowWidth, windowHeight)
	g.updateScale()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
