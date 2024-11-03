package screen

import (
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hritesh04/shooter-game/types"
)

type WinnerScreen struct {
	Index     int
	TileImage *ebiten.Image
	Scene     []types.IScreen
	Assets    *embed.FS
}

func NewWinnerScreen(game types.Game) *WinnerScreen {
	return &WinnerScreen{}
}

func (w *WinnerScreen) Init() {

}

func (w *WinnerScreen) Update() error {
	return nil
}
func (w *WinnerScreen) Draw(*ebiten.Image) {

}
