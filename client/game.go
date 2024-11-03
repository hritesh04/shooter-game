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
	ShowPopup     bool
	RoodID        string
	Screens       []types.IScreen
	Client        *pb.MovementEmitterClient
	Popup         chan bool
	Address       string
}

func NewGame(device types.Device) *Game {
	g := &Game{
		Width:       windowWidth,
		Height:      windowHeight,
		Scale:       1.8,
		Device:      device,
		Filesys:     assets,
		ScreenIndex: types.Onboarding,
		ShowPopup:   true,
		Popup:       make(chan bool),
		// Client:      NewGrpcClient(),
	}
	g.Screens = []types.IScreen{screen.NewOnBoardingScreen(g, assets, g.Popup), screen.NewWinnerScreen(g)}
	g.World = maps.NewMap(maps.NewDefMap, g)
	// g.BoardingUI = ui.InitBoardingUI()
	return g
}

func (g *Game) GetSize() (float64, float64) {
	return float64(g.Width), float64(g.Height)
}

// func (g *Game) Connect() {
// 	conn, err := grpc.NewClient(":3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Fatalf("error creating coonnection,%v", err)
// 	}
// 	defer conn.Close()

// 	g.Client := pb.NewMovementEmitterClient(conn)
// }

func (g *Game) GetClient() *pb.MovementEmitterClient {
	return g.Client
}

func (g *Game) SetServerInfo(ID, address string) {
	g.RoodID = ID
	g.Address = address
	// g.Client = NewGrpcClient(address)
	go g.World.ListenCommand(address, ID)
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) && g.Device == types.Desktop {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	g.Screens[g.ScreenIndex].Update()
	if !g.ShowPopup {
		g.World.Update()
	}
	if g.ShowPopup {
		select {
		case <-g.Popup:
			g.ShowPopup = false
		default:
		}
	}
	// g.BoardingUI.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.World.Draw(screen)
	if g.ShowPopup {
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
