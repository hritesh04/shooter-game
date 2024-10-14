package player

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/resolv"
)

const (
	ActionMoveUp input.Action = iota
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight
)

var Keymap = input.Keymap{
	ActionMoveUp:    {input.KeyGamepadUp, input.KeyUp, input.KeyW},
	ActionMoveDown:  {input.KeyGamepadBack, input.KeyDown, input.KeyS},
	ActionMoveLeft:  {input.KeyGamepadLeft, input.KeyLeft, input.KeyA},
	ActionMoveRight: {input.KeyGamepadRight, input.KeyRight, input.KeyD},
}

type Player struct {
	Input  *input.Handler
	Keymap input.Keymap
	Image  *ebiten.Image
	Src    *resolv.Object
	// Gun
}

func NewPlayer(space *resolv.Space, playerImage *ebiten.Image) *Player {
	// player := resolv.NewObject(25*1.8, 50*1.8, 30, 48, "player")
	player := resolv.NewObject(25, 50, 16, 16, "player")
	space.Add(player)
	return &Player{
		Src:   player,
		Image: playerImage,
	}
}

// func (p *Player) Update() {
// 	if p.Input.ActionIsPressed(ActionMoveLeft) {
// 		p.Pos.X -= 2
// 	}
// 	if p.Input.ActionIsPressed(ActionMoveRight) {
// 		p.Pos.X += 2
// 	}
// 	if p.Input.ActionIsPressed(ActionMoveUp) {
// 		p.Pos.Y -= 2
// 	}
// 	if p.Input.ActionIsPressed(ActionMoveDown) {
// 		p.Pos.Y += 2
// 	}
// }

func (p *Player) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Scale(1.8, 1.8)
	opts.GeoM.Translate(float64(p.Src.Position.X), float64(p.Src.Position.Y))
	screen.DrawImage(p.Image.SubImage(image.Rect(8, 5, 32, 32)).(*ebiten.Image), &opts)
}
