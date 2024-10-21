package weapon

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hritesh04/shooter-game/types"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/resolv"
)

const (
	Pistol = iota
)

const (
	Fire = iota
	Reload
)

var (
	imagePath = map[int]string{
		Pistol: "assets/pistol/pistol.png",
	}
	bulletPath = map[int]string{
		Pistol: "assets/pistol/bullet.png",
	}
	keymap = input.Keymap{
		Fire:   {input.KeyEnter, input.KeySpace},
		Reload: {input.KeyR, input.KeyP},
	}
)

type Weapon struct {
	Input     *input.Handler
	Image     *ebiten.Image
	Bullet    *ebiten.Image
	Src       *resolv.Object
	Space     *resolv.Space
	obstacles []*resolv.Object
}

func NewWeapon(space *resolv.Space, weaponType int) *Weapon {
	pistolImage, _, err := ebitenutil.NewImageFromFile(imagePath[weaponType])
	if err != nil {
		log.Fatal(err)
	}
	bulletImage, _, err := ebitenutil.NewImageFromFile(bulletPath[weaponType])
	if err != nil {
		log.Fatal(err)
	}
	return &Weapon{
		Image:  pistolImage,
		Bullet: bulletImage,
		Space:  space,
	}
}

func (w *Weapon) Init() {
	inputSys := input.System{}
	inputSys.Init(input.SystemConfig{DevicesEnabled: input.AnyDevice})
	w.Input = inputSys.NewHandler(0, keymap)
}

func (w *Weapon) Fire(location resolv.Vector, direction types.Direction, by int) {
	var bullet *resolv.Object
	if direction == types.Left {
		bullet = resolv.NewObject(location.X-45, location.Y, 1, 1, "bullet", string(direction))
		bullet.Data = by
	} else if direction == types.Right {
		bullet = resolv.NewObject(location.X+45, location.Y, 1, 1, "bullet", string(direction))
		bullet.Data = by
	}
	w.Space.Add(bullet)
	w.obstacles = append(w.obstacles, bullet)
	// go w.Update(bullet)
}

// func (w *Weapon) Update(bullet *resolv.Object) {
func (w *Weapon) Update() {

	// for {
	for _, bullet := range w.obstacles {
		if collision := bullet.Check(0, 0, "obstacle"); collision != nil {
			w.Space.Remove(bullet)
			w.RemoveBullet(bullet)
			return
		}
		if collision := bullet.Check(0, 0, "player"); collision != nil {
			w.Space.Remove(bullet)
			w.RemoveBullet(bullet)
			return
		}

		if bullet.HasTags(string(types.Left)) {
			bullet.Position.X -= 3
			bullet.Update()
		} else if bullet.HasTags(string(types.Right)) {
			bullet.Position.X += 3
			bullet.Update()
		}
		bullet.Position.Y += 0.2
		bullet.Update()
		// sleep for 1 frame
		// time.Sleep(time.Second / 120)
		// }
	}
}
func (w *Weapon) RemoveBullet(bullet *resolv.Object) {
	for i := len(w.obstacles) - 1; i >= 0; i-- {
		if w.obstacles[i] == bullet {
			w.obstacles = append(w.obstacles[:i], w.obstacles[i+1:]...)
			break
		}
	}
}

func (w *Weapon) Draw(screen *ebiten.Image, location resolv.Vector) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Scale(0.7, 0.7)
	opts.GeoM.Translate(location.X, location.Y)
	screen.DrawImage(w.Image.SubImage(image.Rect(0, 0, 64, 64)).(*ebiten.Image), &opts)

	for _, bullet := range w.obstacles {
		opts := ebiten.DrawImageOptions{}
		opts.GeoM.Scale(0.7, 0.7)
		opts.GeoM.Translate(bullet.Position.X+5, bullet.Position.Y+2)
		screen.DrawImage(w.Bullet.SubImage(image.Rect(0, 0, 64, 64)).(*ebiten.Image), &opts)
	}
	// debug code
	// vector.DrawFilledRect(screen, float32(bullet.Position.X), float32(bullet.Position.Y), 16, 16, color.RGBA{0, 0, 255, 128}, true)
}
