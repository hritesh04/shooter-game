package client

import (
	"embed"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hritesh04/shooter-game/entities/screen"
	"github.com/hritesh04/shooter-game/maps"
	pb "github.com/hritesh04/shooter-game/stubs"
	"github.com/hritesh04/shooter-game/types"
)

const (
	windowWidth  = 1280
	windowHeight = 704
)

//go:embed assets/*
var assets embed.FS

type Game struct {
	Width, Height int
	Scale         float64
	World         types.IMap
	Device        types.Device
	Filesys       embed.FS
	ScreenIndex   int
	ShowScreen    bool
	Screens       []types.IScreen
	Client        pb.MovementEmitterClient
}

func NewGame(device types.Device) *Game {
	g := &Game{
		Width:       windowWidth,
		Height:      windowHeight,
		Scale:       1.8,
		Device:      device,
		Filesys:     assets,
		ScreenIndex: types.Onboarding,
		ShowScreen:  true,
		// Client:      NewGrpcClient(),
	}
	g.Screens = []types.IScreen{screen.NewOnBoardingScreen(g, assets), screen.NewWinnerScreen(g)}
	g.World = maps.NewMap(maps.NewDefMap, g)
	// g.BoardingUI = ui.InitBoardingUI()
	return g
}

func (g *Game) GetSize() (float64, float64) {
	return float64(g.Width), float64(g.Height)
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) && g.Device == types.Desktop {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	if !g.ShowScreen {
		g.World.Update()
	}
	// g.BoardingUI.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.World.Draw(screen)
	if g.ShowScreen {
		g.Screens[g.ScreenIndex].Draw(screen)
	}
	// g.BoardingUI.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if g.Device == types.Desktop && ebiten.IsFullscreen() {
		return ebiten.WindowSize()
	}
	return g.Width, g.Height
}

func (g *Game) GetDevice() types.Device {
	return g.Device
}

func (g *Game) GetFS() embed.FS {
	return g.Filesys
}
