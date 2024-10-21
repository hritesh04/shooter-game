package player

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hritesh04/shooter-game/entities/weapon"
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
	Name   int
	Input  *input.Handler
	Keymap input.Keymap
	Image  *ebiten.Image
	Src    *resolv.Object
	Weapon *weapon.Weapon
	Dir    types.Direction
}

func NewPlayer(index int, space *resolv.Space) *Player {
	playerImage, _, err := ebitenutil.NewImageFromFile("assets/runner.png")
	if err != nil {
		log.Fatal(err)
	}
	var player *resolv.Object
	if index == 0 {
		player = resolv.NewObject(60, 70, 26, 32, "player")
	} else {
		player = resolv.NewObject(16*36*2+20, 16*19*2, 32, 32, "player")
	}
	space.Add(player)
	return &Player{
		Src:    player,
		Image:  playerImage,
		Weapon: weapon.NewWeapon(space, weapon.Pistol),
	}
}

func (p *Player) Init() {
	inputSystem := input.System{}
	inputSystem.Init(input.SystemConfig{DevicesEnabled: input.AnyDevice})
	p.Input = inputSystem.NewHandler(0, keymap)
	p.Weapon.Init()
}

func (p *Player) Update() {
	playerObj := p.Src
	moved := false

	if p.Input.ActionIsPressed(ActionMoveLeft) {
		if collision := playerObj.Check(-2, 0, "obstacle"); collision == nil {
			playerObj.Position.X -= 2
			p.Dir = types.Left
			// playerObj.Shape.Move(-2, 0)
			moved = true
		}
	}
	if p.Input.ActionIsPressed(ActionMoveRight) {
		if collision := playerObj.Check(2, 0, "obstacle"); collision == nil {
			playerObj.Position.X += 2
			p.Dir = types.Right
			// playerObj.Shape.Move(2, 0)
			// fmt.Println(playerObj.Shape.Rotation())
			moved = true
		}
	}
	if p.Input.ActionIsPressed(ActionMoveUp) {
		if collision := playerObj.Check(0, -2, "obstacle"); collision == nil {
			playerObj.Position.Y -= 2
			// playerObj.Shape.Move(0, -2)
			moved = true
		}
	}
	if p.Input.ActionIsPressed(ActionMoveDown) {
		if collision := playerObj.Check(0, 10, "obstacle"); collision == nil {
			playerObj.Position.Y += 2
			// playerObj.Shape.Move(0, 2)
			moved = true
		}
	}
	if p.Input.ActionIsJustReleased(Fire) {
		p.Weapon.Fire(p.Src.Position, p.Dir, p.Name)
	}

	if moved {
		playerObj.Update()
	}

	p.Weapon.Update()

	if collision := playerObj.Check(0, 0, "bullet"); collision != nil {
		if collision.Objects[0].Data != string(p.Name) {
			playerObj.Position.Y = 70
			playerObj.Position.X = 60
			playerObj.Update()
		}
	}

}

func (p *Player) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Scale(1.8, 1.8)
	opts.GeoM.Translate(float64(p.Src.Position.X), float64(p.Src.Position.Y-10))
	screen.DrawImage(p.Image.SubImage(image.Rect(8, 5, 32, 32)).(*ebiten.Image), &opts)
	p.Weapon.Draw(screen, p.Src.Position)
	// debug code
	// vector.DrawFilledRect(screen, float32(p.Src.Position.X), float32(p.Src.Position.Y), float32(p.Src.Size.X), float32(p.Src.Size.Y), color.RGBA{0, 0, 255, 128}, true)
}
