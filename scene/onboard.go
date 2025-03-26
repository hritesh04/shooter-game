package screen

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hritesh04/shooter-game/entities/player"
	"github.com/hritesh04/shooter-game/maps/common"
	"github.com/hritesh04/shooter-game/scene/ui"
	"github.com/hritesh04/shooter-game/types"
	"github.com/solarlune/resolv"
)

const (
	JoinDungeon = iota
	CreateDungeon
	Onboarding
	Exit
)

type Onboard struct {
	ID                int
	Width, Height     int
	Index             int
	ParentScene       int
	collisionCoolDown int
	SceneStart        bool
	Game              types.Game
	Scale             float64
	TileImage         *ebiten.Image
	PlayerImage       *ebiten.Image
	MapJson           *common.TiledMapJSON
	Player            *player.Player
	Space             *resolv.Space
	Scene             []types.IScreen
	Assets            embed.FS
	Obstacles         []*resolv.Object
	Show              chan bool
	Done              chan bool
}

func NewOnBoardingScreen(ID, ParentScene int, game types.Game, fs embed.FS, show chan bool) *Onboard {
	file, err := fs.Open("assets/dungeonSheet.png")
	if err != nil {
		log.Fatal(err)
	}
	tileImage, _, err := ebitenutil.NewImageFromReader(file)
	if err != nil {
		log.Fatal(err)
	}
	rfile, err := fs.Open("assets/runner.png")
	if err != nil {
		log.Fatal(err)
	}
	runnerImage, _, err := ebitenutil.NewImageFromReader(rfile)
	if err != nil {
		log.Fatal(err)
	}
	mapFile, err := fs.ReadFile("assets/screen/onboard.json")
	if err != nil {
		log.Fatal(err)
	}
	scence := &Onboard{
		ID:                ID,
		ParentScene:       ParentScene,
		Game:              game,
		TileImage:         tileImage,
		PlayerImage:       runnerImage,
		MapJson:           &common.TiledMapJSON{},
		Assets:            fs,
		Index:             Onboarding,
		collisionCoolDown: 0,
		SceneStart:        false,
		Show:              show,
		Done:              make(chan bool),
	}

	if err := json.Unmarshal(mapFile, scence.MapJson); err != nil {
		log.Fatal(err)
	}
	scence.Init()
	return scence
}

func (o *Onboard) Init() {
	gw, gh := o.Game.GetSize()
	cellSize := 16
	tx := gw/2 - gw/4
	ty := gh/2 - gh/3

	o.Space = resolv.NewSpace(int(gw), int(gh), cellSize, cellSize)
	o.Player = player.NewPlayer("onboard", tx+tx*0.9, ty+ty*2.9, 0, o.Space, o.Game.GetDevice(), o.Assets, nil, "")
	o.Player.Init()
	o.Obstacles = append(o.Obstacles, o.Player.Src)

	o.Width = o.MapJson.Layers[0].Width * 16
	o.Height = len(o.MapJson.Layers[0].Data) / o.MapJson.Layers[0].Width * 16

	scaleX := gw / float64(o.Width)
	scaleY := gh / float64(o.Height)
	o.Scale = math.Min(scaleX, scaleY)
	var obsLayerIndex int
	var layerW int
	for index, layer := range o.MapJson.Layers {
		if layer.Name == "boundry" {
			obsLayerIndex = index
			layerW = layer.Width
		}
	}

	for index, id := range o.MapJson.Layers[obsLayerIndex].Data {
		if id == 0 {
			continue
		}

		x := float64((index % layerW) * 16)
		y := float64((index / layerW) * 16)
		obstacle := resolv.NewObject(x*1.8+tx, y*1.8+ty+10, 28, 20, "obstacle")
		o.Obstacles = append(o.Obstacles, obstacle)
		o.Space.Add(obstacle)
	}

	var joinLayer int
	for index, layer := range o.MapJson.Layers {
		if layer.Name == "join" {
			joinLayer = index
			layerW = layer.Width
		}
	}

	for index, id := range o.MapJson.Layers[joinLayer].Data {
		if id == 0 {
			continue
		}

		x := float64((index % layerW) * 16)
		y := float64((index / layerW) * 16)
		obstacle := resolv.NewObject(x*1.8+tx, y*1.8+ty+10, 28, 20, "join")
		o.Obstacles = append(o.Obstacles, obstacle)
		o.Space.Add(obstacle)
	}

	var createLayer int
	for index, layer := range o.MapJson.Layers {
		if layer.Name == "create" {
			createLayer = index
			layerW = layer.Width
		}
	}

	for index, id := range o.MapJson.Layers[createLayer].Data {
		if id == 0 {
			continue
		}

		x := float64((index % layerW) * 16)
		y := float64((index / layerW) * 16)
		obstacle := resolv.NewObject(x*1.8+tx, y*1.8+ty+10, 28, 20, "create")
		o.Obstacles = append(o.Obstacles, obstacle)
		o.Space.Add(obstacle)
	}

	o.Scene = append(o.Scene, ui.NewKeyBoardInput(JoinDungeon, Onboarding, "Enter the dungeon ID\n", tx+tx*0.3, ty-ty*0.3, o.JoinRoom, Exit), ui.NewKeyBoardInput(CreateDungeon, Onboarding, "Creating a new dungeon", tx+tx*1.3, ty-ty*0.3, o.CreateRoom, Exit))
}

func (o *Onboard) Update() int {
	if o.Index == Onboarding {
		o.Player.Update()
		if o.collisionCoolDown <= 0 {
			playerObj := o.Player.Src
			if isScene, scene := checkJoinRoom(playerObj); isScene {
				fmt.Println("collission detected ", scene)
				o.collisionCoolDown = 60
				o.Index = scene
				return o.ID
			}
		}
		o.collisionCoolDown--
	} else {
		screen := o.Scene[o.Index].Update()
		if screen == Exit {
			o.Game.TogglePopUp(false)
			return o.ParentScene
		} else {
			o.Index = screen
		}
	}
	return o.ID
}

func checkJoinRoom(player *resolv.Object) (bool, int) {
	if collision := player.Check(0, -2, "join"); collision != nil {
		return true, JoinDungeon
	}
	if collision := player.Check(0, -2, "create"); collision != nil {
		return true, CreateDungeon
	}
	return false, 2
}

func (o *Onboard) JoinRoom() func(string) error {
	return func(data string) error {
		out, err := json.Marshal(struct {
			RoomID string `json:"roomID"`
		}{RoomID: data})
		if err != nil {
			return fmt.Errorf("data marshaling failed")
		}
		req, err := http.NewRequest(http.MethodPost, os.Getenv("ROOT_SERVER_URL")+"/joinRoom", bytes.NewBuffer(out))
		if err != nil {
			return fmt.Errorf("error creating request")
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("error making request")
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				return fmt.Errorf("error reading response body: %w", err)
			}
			return fmt.Errorf(string(body))
		}
		decoder := json.NewDecoder(res.Body)
		var player PlayerResponse
		if err := decoder.Decode(&player); err != nil {
			return fmt.Errorf("response decoding failed")
		}
		defer res.Body.Close()
		o.Game.SetServerInfo(player.RoomID, player.Address, player.Type)
		return nil
	}
}

func (o *Onboard) CreateRoom() func(string) error {
	return func(roomID string) error {
		out, err := json.Marshal(struct {
			RoomID string `json:"roomID"`
		}{RoomID: roomID})
		if err != nil {
			return fmt.Errorf("data marshaling failed")
		}
		req, err := http.NewRequest(http.MethodGet, os.Getenv("ROOT_SERVER_URL")+"/createRoom", bytes.NewBuffer(out))
		if err != nil {
			return fmt.Errorf("error creating request")
		}
		log.Println("Making create room request to : ", os.Getenv("ROOT_SERVER_URL"))
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("error making request", err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				return fmt.Errorf("error reading response body: %w", err)
			}
			return fmt.Errorf(string(body))
		}
		decoder := json.NewDecoder(res.Body)
		var player PlayerResponse
		if err := decoder.Decode(&player); err != nil {
			return fmt.Errorf("response decoding failed %w", err)
		}
		defer res.Body.Close()
		log.Printf("Game server address received : %s", player.Address)
		o.Game.SetServerInfo(player.RoomID, player.Address, player.Type)
		return nil
	}
}

type PlayerResponse struct {
	RoomID  string `json:"roomID"`
	Address string `json:"address"`
	Type    string `json:"type"`
}

func (o *Onboard) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	w, h := o.Game.GetSize()
	for _, layer := range o.MapJson.Layers {
		if layer.Name == "join" || layer.Name == "create" {
			for index, id := range layer.Data {
				if id == 0 {
					continue
				}

				x := float64((index % layer.Width) * 16)
				y := float64((index / layer.Width) * 16)

				tx := w/2 - w/4
				ty := h/2 - h/3
				opts.GeoM.Reset()
				opts.GeoM.Scale(1.8, 1.8)
				opts.GeoM.Translate(x*1.8+tx, y+ty)
				screen.DrawImage(o.PlayerImage.SubImage(image.Rect(0, 0, 32, 32)).(*ebiten.Image), &opts)
				continue
			}
		}
		for index, id := range layer.Data {
			if id == 0 {
				continue
			}

			x := float64((index % layer.Width) * 16)
			y := float64((index / layer.Width) * 16)

			tx := w/2 - w/4
			ty := h/2 - h/3
			srcX := ((id - 1) % 12) * 16
			srcY := ((id - 1) / 24) * 16

			opts.GeoM.Reset()
			opts.GeoM.Scale(1.8, 1.8)
			opts.GeoM.Translate(x*1.8+tx, y*1.8+ty)
			screen.DrawImage(o.TileImage.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image), &opts)
		}
	}
	o.Player.Draw(screen)
	if o.Index != Onboarding {
		o.Scene[o.Index].Draw(screen)
	}
	// for _, obj := range o.Obstacles {
	// 	if obj.HasTags("obstacle") {
	// 		vector.DrawFilledRect(screen, float32(obj.Position.X), float32(obj.Position.Y), float32(obj.Size.X), float32(obj.Size.Y), color.RGBA{0, 0, 255, 128}, true)
	// 	}
	// 	if obj.HasTags("join") {
	// 		vector.DrawFilledRect(screen, float32(obj.Position.X), float32(obj.Position.Y), float32(obj.Size.X), float32(obj.Size.Y), color.RGBA{0, 0, 255, 128}, true)
	// 	}
	// 	if obj.HasTags("create") {
	// 		vector.DrawFilledRect(screen, float32(obj.Position.X), float32(obj.Position.Y), float32(obj.Size.X), float32(obj.Size.Y), color.RGBA{0, 0, 255, 128}, true)
	// 	}
	// }
}
