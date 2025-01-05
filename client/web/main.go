package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hritesh04/shooter-game/client"
	"github.com/hritesh04/shooter-game/types"
)

const (
	windowWidth  = 1280
	windowHeight = 704
)

// transfer map to new file - game.go? - use that to init game
// client/main.go -> main for desktop and webassembly
// client/mobile/main.go -> main for mobile
// client/web/ -> assets-html- for webassembly
func main() {
	g := client.NewGame(types.Web)
	ebiten.SetWindowTitle("Shooter")
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetFullscreen(true)
	ebiten.SetTPS(60)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
