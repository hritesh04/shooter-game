package new

import (
	"encoding/json"
	"image"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hritesh04/shooter-game/entities/player"
	"github.com/hritesh04/shooter-game/maps/common"
	"github.com/hritesh04/shooter-game/types"
	"github.com/solarlune/resolv"
)

type DefaultMap struct {
	Game          types.Game
	MapJson       *common.TiledMapJSON
	Space         *resolv.Space
	Obstacles     []*resolv.Object
	Players       []*player.Player
	Width, Height int
	Scale         float64
	TileImage     *ebiten.Image
}

func NewDefaultMap(game types.Game) types.IMap {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	mapData, err := os.ReadFile(pwd + "/maps/default/map.json")
	if err != nil {
		log.Fatal(err)
	}

	tileMap := &DefaultMap{
		Game:    game,
		MapJson: &common.TiledMapJSON{},
	}

	if err := json.Unmarshal(mapData, tileMap.MapJson); err != nil {
		log.Fatal(err)
	}

	tileMap.Init()
	return tileMap
}

func (m *DefaultMap) Init() {
	tileImage, _, err := ebitenutil.NewImageFromFile("assets/dungeonSheet.png")
	if err != nil {
		log.Fatal(err)
	}
	m.TileImage = tileImage

	gw, gh := m.Game.GetSize()
	cellSize := 16

	m.Space = resolv.NewSpace(int(gw), int(gh), cellSize, cellSize)

	for i := 0; i < 2; i++ {
		player := player.NewPlayer(i, m.Space)
		player.Init()
		m.Players = append(m.Players, player)
	}

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
func (m *DefaultMap) Update() error {
	for _, player := range m.Players {
		player.Update()
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
