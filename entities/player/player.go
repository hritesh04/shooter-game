package player

import (
	"embed"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hritesh04/shooter-game/entities/weapon"
	pb "github.com/hritesh04/shooter-game/stubs"
	"github.com/hritesh04/shooter-game/types"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/resolv"
)

const (
	ActionMoveUp input.Action = iota
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight
	Fire
	Reload
)

var keymap = input.Keymap{
	ActionMoveUp:    {input.KeyGamepadUp, input.KeyUp, input.KeyW},
	ActionMoveDown:  {input.KeyGamepadBack, input.KeyDown, input.KeyS},
	ActionMoveLeft:  {input.KeyGamepadLeft, input.KeyLeft, input.KeyA},
	ActionMoveRight: {input.KeyGamepadRight, input.KeyRight, input.KeyD},
	Fire:            {input.KeyEnter, input.KeySpace},
	Reload:          {input.KeyR, input.KeyP},
}

type Player struct {
	Name   string
	X, Y   float64
	Input  *input.Handler
	Keymap input.Keymap
	Image  *ebiten.Image
	Src    *resolv.Object
	roomID string
	Weapon *weapon.Weapon
	Dir    types.Direction
	conn   types.IConnection
}

func NewPlayer(name string, w float64, h float64, index int, space *resolv.Space, device types.Device, assets embed.FS, Conn types.IConnection, roomID string) *Player {
	var playerImage *ebiten.Image
	var err error
	runner, err := assets.Open("assets/runner.png")
	if err != nil {
		log.Fatal(err)
	}
	if device == types.Desktop {
		playerImage, _, err = ebitenutil.NewImageFromReader(runner)
		if err != nil {
			log.Fatal(err)
		}
	} else if device == types.Web {
		playerImage, _, err = ebitenutil.NewImageFromReader(runner)
		if err != nil {
			log.Fatal(err)
		}
	}
	// var player *resolv.Object
	player := resolv.NewObject(w, h, 20, 28, "player")
	// if index == 0 {
	// } else {
	// 	player = resolv.NewObject(16*36*2+20, 16*19*2, 20, 28, "player")
	// }
	space.Add(player)
	return &Player{
		Name:   name,
		X:      w,
		Y:      h,
		Src:    player,
		Image:  playerImage,
		Weapon: weapon.NewWeapon(space, weapon.Pistol, device, assets),
		conn:   Conn,
		roomID: roomID,
	}
}

func (p *Player) Init() {
	inputSystem := input.System{}
	inputSystem.Init(input.SystemConfig{DevicesEnabled: input.AnyDevice})
	p.Input = inputSystem.NewHandler(0, keymap)
	// conn := p.conn.GetEventConn()
	// // do in seperate
	// conn.Send(&pb.Data{Type: pb.Action_Join, RoomID: p.roomID, Name: p.Name})
	// if err := p.conn.Send(&pb.Data{Type: pb.Action_Join, RoomID: p.roomID, Name: p.Name}); err != nil {
	// log.Println("Error setting up stream")
	// }
	p.Weapon.Init()
}

func (p *Player) AddStream() {
	conn := p.conn.GetEventConn()
	// do in seperate
	conn.Send(&pb.Data{Type: pb.Action_Join, RoomID: p.roomID, Name: p.Name})
}

func (p *Player) Update() {
	// p.Input.EmitKeyEvent(input.SimulatedKeyEvent{})
	playerObj := p.Src
	moved := false
	player := &pb.Player{Name: p.Name, X: float32(p.Src.Position.X), Y: float32(p.Src.Position.Y)}
	var conn pb.MovementEmitter_SendMoveClient
	if p.conn != nil {
		conn = p.conn.GetEventConn()
	}
	// fmt.Println("PLAYER UPDATE")
	if p.Input.ActionIsPressed(ActionMoveLeft) {
		if collision := playerObj.Check(-2, 0, "obstacle"); collision == nil {
			playerObj.Position.X -= 2
			p.Dir = types.Left
			if p.conn != nil {
				conn.Send(&pb.Data{Type: pb.Action_Movement, Data: pb.Direction_LEFT, Name: p.Name, RoomID: p.roomID, Player: []*pb.Player{player}})
			}
			// playerObj.Shape.Move(-2, 0)
			moved = true
		}
	}
	if p.Input.ActionIsPressed(ActionMoveRight) {
		if collision := playerObj.Check(2, 0, "obstacle"); collision == nil {
			playerObj.Position.X += 2
			p.Dir = types.Right
			if p.conn != nil {
				conn.Send(&pb.Data{Type: pb.Action_Movement, Data: pb.Direction_RIGHT, Name: p.Name, RoomID: p.roomID, Player: []*pb.Player{player}})
				// playerObj.Shape.Move(2, 0)
			}
			// fmt.Println(playerObj.Shape.Rotation())
			moved = true
		}
	}
	if p.Input.ActionIsPressed(ActionMoveUp) {
		if collision := playerObj.Check(0, -2, "obstacle"); collision == nil {
			playerObj.Position.Y -= 2
			if p.conn != nil {
				conn.Send(&pb.Data{Type: pb.Action_Movement, Data: pb.Direction_UP, Name: p.Name, RoomID: p.roomID, Player: []*pb.Player{player}})
			}
			// playerObj.Shape.Move(0, -2)
			moved = true
		}
	}
	if p.Input.ActionIsPressed(ActionMoveDown) {
		if collision := playerObj.Check(0, 10, "obstacle"); collision == nil {
			playerObj.Position.Y += 2
			if p.conn != nil {
				conn.Send(&pb.Data{Type: pb.Action_Movement, Data: pb.Direction_DOWN, Name: p.Name, RoomID: p.roomID, Player: []*pb.Player{player}})
				// playerObj.Shape.Move(0, 2)
			}
			moved = true
		}
	}
	if p.Input.ActionIsJustReleased(Fire) {
		p.Weapon.Fire(p.Src.Position, p.Dir, p.Name)
		if p.conn != nil {
			if p.Dir == types.Right {
				conn.Send(&pb.Data{Type: pb.Action_Fire, Data: pb.Direction_RIGHT, Name: p.Name, RoomID: p.roomID, Player: []*pb.Player{player}})
			} else {
				conn.Send(&pb.Data{Type: pb.Action_Fire, Data: pb.Direction_LEFT, Name: p.Name, RoomID: p.roomID, Player: []*pb.Player{player}})
			}
		}
	}

	if moved {
		playerObj.Update()
	}

	p.Weapon.Update()

	if collision := playerObj.Check(0, 0, "bullet"); collision != nil {
		if collision.Objects[0].Data != p.Name {
			playerObj.Position.X = p.X
			playerObj.Position.Y = p.Y
			playerObj.Update()
		}
	}

}

func (p *Player) Simulate() {
	playerObj := p.Src
	// playerObj.Update()
	p.Weapon.Update()
	if collision := playerObj.Check(0, 0, "bullet"); collision != nil {
		// log.Printf("Collision with bullet %+v", *collision.Objects[0])
		if collision.Objects[0].Data != p.Name {
			log.Printf("Collision with bullet %+v for player %s", *collision.Objects[0], p.Name)
			playerObj.Position.X = p.X
			playerObj.Position.Y = p.Y
			playerObj.Update()
		}
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	if p.Dir == types.Left {
		opts.GeoM.Scale(-1.8, 1.8)
		opts.GeoM.Translate(float64(p.Src.Position.X+21), float64(p.Src.Position.Y-10))
	} else {
		opts.GeoM.Scale(1.8, 1.8)
		opts.GeoM.Translate(float64(p.Src.Position.X), float64(p.Src.Position.Y-10))
	}
	screen.DrawImage(p.Image.SubImage(image.Rect(8, 5, 32, 32)).(*ebiten.Image), &opts)
	p.Weapon.Draw(screen, p.Src.Position, p.Dir)
	// debug code
	// vector.DrawFilledRect(screen, float32(p.Src.Position.X), float32(p.Src.Position.Y), float32(p.Src.Size.X), float32(p.Src.Size.Y), color.RGBA{0, 0, 255, 128}, true)
}
