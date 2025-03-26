package new

import (
	"embed"
	"encoding/json"
	"image"
	"io"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hritesh04/shooter-game/conn"
	"github.com/hritesh04/shooter-game/entities/player"
	"github.com/hritesh04/shooter-game/maps/common"
	pb "github.com/hritesh04/shooter-game/proto"
	"github.com/hritesh04/shooter-game/types"
	"github.com/solarlune/resolv"
)

type DefaultMap struct {
	Game          types.Game
	MapJson       *common.TiledMapJSON
	Space         *resolv.Space
	Obstacles     []*resolv.Object
	Players       map[string]*player.Player
	Width, Height int
	Scale         float64
	TileImage     *ebiten.Image
	Device        types.Device
	Assets        embed.FS
	name          string
	Conn          types.IConnection
	Rec           chan *pb.Data
	Sender        chan *pb.Data
}

func NewDefaultMap(game types.Game) types.IMap {
	fs := game.GetFS()
	file, err := fs.Open("assets/dungeonSheet.png")
	if err != nil {
		log.Fatal(err)
	}
	tileImage, _, err := ebitenutil.NewImageFromReader(file)
	if err != nil {
		log.Fatal(err)
	}
	mapFile, err := fs.ReadFile("assets/map/map.json")
	if err != nil {
		log.Fatal(err)
	}
	tileMap := &DefaultMap{
		Game:      game,
		MapJson:   &common.TiledMapJSON{},
		Device:    game.GetDevice(),
		TileImage: tileImage,
		Assets:    fs,
		Players:   make(map[string]*player.Player),
		Rec:       make(chan *pb.Data),
		Sender:    make(chan *pb.Data),
	}
	if err := json.Unmarshal(mapFile, tileMap.MapJson); err != nil {
		log.Fatal(err)
	}

	tileMap.Init()
	return tileMap
}

func (m *DefaultMap) Init() {

	gw, gh := m.Game.GetSize()
	cellSize := 16

	m.Space = resolv.NewSpace(int(gw), int(gh), cellSize, cellSize)

	m.Width = m.MapJson.Layers[0].Width * 16
	m.Height = len(m.MapJson.Layers[0].Data) / m.MapJson.Layers[0].Width * 16

	scaleX := gw / float64(m.Width)
	scaleY := gh / float64(m.Height)
	m.Scale = math.Min(scaleX, scaleY)
	var obsLayerIndex int
	var layerW int
	for index, layer := range m.MapJson.Layers {
		if layer.Name == "obstacle" {
			obsLayerIndex = index
			layerW = layer.Width
		}
	}

	for index, id := range m.MapJson.Layers[obsLayerIndex].Data {
		if id == 0 {
			continue
		}

		x := float64((index % layerW) * 16)
		y := float64((index / layerW) * 16)
		obstacle := resolv.NewObject(x*m.Scale, y*m.Scale, 32, 32, "obstacle")
		m.Obstacles = append(m.Obstacles, obstacle)
		m.Space.Add(obstacle)
	}
}

func (m *DefaultMap) JoinRoom(address, ID, connType string) error {
	var join *pb.Room
	var err error
	if connType == "grpc" {
		log.Println("Making Grpc connection")
		m.Conn = conn.NewGrpcClient(address, m.Device)
		join, err = m.Conn.JoinRoom(ID)
		if err != nil {
			log.Println("Error joining room with grpc :%s", err.Error())
		}
	} else {
		log.Println("Making websocket connection : ")
		m.Conn = conn.NewRestClient(address, m.Device)
		join, err = m.Conn.JoinRoom(ID)
		if err != nil {
			log.Println("Error joining room with http connection :%s", err.Error())
		}
	}
	for _, p := range join.GetPlayer() {
		m.name = p.GetName()
		player := player.NewPlayer(p.GetName(), float64(p.GetX()), float64(p.GetY()), 0, m.Space, m.Device, m.Assets, m.Conn, ID)
		player.Init()
		player.AddStream()
		m.Players[p.GetName()] = player
	}
	return nil
}

func (m *DefaultMap) ListenCommand(address, ID string) {
	conn := m.Conn.GetEventConn()
	log.Printf("Backend URL : %s for room %s", address, ID)
	go func() {
		for {
			resp, err := conn.Recv()

			log.Printf("received data from move :%v", resp)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("receiving Streaming message: %v", err)
			}
			switch resp.Type {
			case pb.Action_Join, pb.Action_Info:
				log.Printf("New player joined : %+v", resp.Player)
				players := resp.GetPlayer()
				for _, p := range players {
					player := player.NewPlayer(p.GetName(), float64(p.GetX()), float64(p.GetY()), 0, m.Space, m.Device, m.Assets, nil, ID)
					player.Init()
					m.Players[p.GetName()] = player
				}
			case pb.Action_Movement:
				players := resp.GetPlayer()
				for _, p := range players {
					player := m.Players[p.GetName()]
					switch resp.GetData() {
					case pb.Direction_RIGHT:
						player.Dir = types.Right
					case pb.Direction_LEFT:
						player.Dir = types.Left
					}
					player.Src.Position.X = float64(p.GetX())
					player.Src.Position.Y = float64(p.GetY())
				}
			case pb.Action_Fire:
				players := resp.GetPlayer()
				for _, p := range players {
					player := m.Players[p.GetName()]
					switch resp.GetData() {
					case pb.Direction_RIGHT:
						player.Weapon.Fire(resolv.Vector{X: float64(p.GetX()), Y: float64(p.GetY())}, types.Right, p.GetName())
					case pb.Direction_LEFT:
						player.Weapon.Fire(resolv.Vector{X: float64(p.GetX()), Y: float64(p.GetY())}, types.Left, p.GetName())
					}
				}
			}
		}
	}()
}

func (m *DefaultMap) Update() error {
	for name, player := range m.Players {
		if name == m.name {
			player.Update()
		} else {
			player.Simulate()
		}
	}
	return nil
}

func (m *DefaultMap) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	for _, layer := range m.MapJson.Layers {
		for index, id := range layer.Data {
			if id == 0 {
				continue
			}

			x := float64((index % layer.Width) * 16)
			y := float64((index / layer.Width) * 16)

			srcX := ((id - 1) % 12) * 16
			srcY := ((id - 1) / 24) * 16

			opts.GeoM.Reset()
			opts.GeoM.Scale(m.Scale, m.Scale)
			opts.GeoM.Translate(x*m.Scale, y*m.Scale)
			screen.DrawImage(m.TileImage.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image), &opts)
		}
	}
	for _, player := range m.Players {
		player.Draw(screen)
	}

	// debug code
	// for _, obj := range m.Obstacles {
	// 	if obj.HasTags("obstacle") {
	// 		vector.DrawFilledRect(screen, float32(obj.Position.X), float32(obj.Position.Y), float32(obj.Size.X), float32(obj.Size.Y), color.RGBA{0, 0, 255, 128}, true)
	// 	}
	// }
	// 	if obj.HasTags("player") {
	// 		vector.DrawFilledRect(screen, float32(obj.Position.X), float32(obj.Position.Y), float32(obj.Size.X), float32(obj.Size.Y), color.RGBA{0, 0, 255, 128}, true)
	// 	}
	// }
}
