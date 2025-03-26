package client

import (
	"embed"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hritesh04/shooter-game/maps"
	pb "github.com/hritesh04/shooter-game/proto"
	screen "github.com/hritesh04/shooter-game/scene"
	"github.com/hritesh04/shooter-game/types"
	"github.com/joho/godotenv"
)

const (
	windowWidth  = 1280
	windowHeight = 704
)

//go:embed assets/*
var assets embed.FS

const (
	Onboarding = iota
	Winner
	World
	Exit
)

type Game struct {
	Width, Height int
	Scale         float64
	World         types.IMap
	Device        types.Device
	Filesys       embed.FS
	SceneIndex    int
	ShowPopup     bool
	RoodID        string
	Screens       []types.IScreen
	Client        *pb.MovementEmitterClient
	Popup         chan bool
	Address       string
	Type          string
}

func init() {
	if env := os.Getenv("ENV"); env == "DEV" {
		err := godotenv.Load(".env")
		if err != nil {
			panic(err.Error())
		}
	}
}

func NewGame(device types.Device) *Game {
	g := &Game{
		Width:      windowWidth,
		Height:     windowHeight,
		Scale:      1.8,
		Device:     device,
		Filesys:    assets,
		SceneIndex: Onboarding,
		ShowPopup:  true,
		Popup:      make(chan bool),
		// Client:      NewGrpcClient(),
	}
	g.Screens = []types.IScreen{screen.NewOnBoardingScreen(Onboarding, World, g, assets, g.Popup), screen.NewWinnerScreen(g)}
	g.World = maps.NewMap(maps.NewDefMap, g)
	return g
}

func (g *Game) GetSize() (float64, float64) {
	return float64(g.Width), float64(g.Height)
}

func (g *Game) GetClient() *pb.MovementEmitterClient {
	return g.Client
}

func (g *Game) SetServerInfo(ID, address, connType string) {
	g.RoodID = ID
	g.Address = address
	g.Type = connType
	if err := g.World.JoinRoom(address, ID, connType); err != nil {
		log.Println("Failed to join room")
		return
	}
	go g.World.ListenCommand(address, ID)
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) && g.Device == types.Desktop {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	if g.SceneIndex == World {
		g.World.Update()
	} else {
		g.SceneIndex = g.Screens[g.SceneIndex].Update()
	}
	return nil
}

func (g *Game) TogglePopUp(flag bool) {
	g.ShowPopup = flag
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.World.Draw(screen)
	if g.ShowPopup {
		g.Screens[g.SceneIndex].Draw(screen)
	}
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
