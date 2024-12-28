package weapon

import (
	"embed"
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
		Pistol: "pistol/pistol.png",
	}
	bulletPath = map[int]string{
		Pistol: "pistol/bullet.png",
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

func NewWeapon(space *resolv.Space, weaponType int, device types.Device, assets embed.FS) *Weapon {
	var pistolImage *ebiten.Image
	var bulletImage *ebiten.Image
	if device == types.Desktop {
		gun, err := assets.Open("assets/" + imagePath[weaponType])
		if err != nil {
			log.Fatal(err)
		}
		pistolImage, _, err = ebitenutil.NewImageFromReader(gun)
		if err != nil {
			log.Fatal(err)
		}
		bullet, err := assets.Open("assets/" + bulletPath[weaponType])
		if err != nil {
			log.Fatal(err)
		}
		bulletImage, _, err = ebitenutil.NewImageFromReader(bullet)
		if err != nil {
			log.Fatal(err)
		}
	} else if device == types.Web {
		gun, err := assets.Open("assets/" + imagePath[weaponType])
		if err != nil {
			log.Fatal(err)
		}
		pistolImage, _, err = ebitenutil.NewImageFromReader(gun)
		if err != nil {
			log.Fatal(err)
		}
		bullet, err := assets.Open("assets/" + bulletPath[weaponType])
		if err != nil {
			log.Fatal(err)
		}
		bulletImage, _, err = ebitenutil.NewImageFromReader(bullet)
		if err != nil {
			log.Fatal(err)
		}
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

func (w *Weapon) Fire(location resolv.Vector, direction types.Direction, by string) {
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

func (w *Weapon) Draw(screen *ebiten.Image, location resolv.Vector, dir types.Direction) {
	opts := ebiten.DrawImageOptions{}
	if dir == types.Left {
		opts.GeoM.Scale(-0.7, 0.7)
		opts.GeoM.Translate(location.X+20, location.Y)
	} else {
		opts.GeoM.Scale(0.7, 0.7)
		opts.GeoM.Translate(location.X, location.Y)
	}
	screen.DrawImage(w.Image.SubImage(image.Rect(0, 0, 64, 64)).(*ebiten.Image), &opts)

	for _, bullet := range w.obstacles {
		opts := ebiten.DrawImageOptions{}
		opts.GeoM.Scale(0.7, 0.7)
		opts.GeoM.Translate(bullet.Position.X+5, bullet.Position.Y+2)
		screen.DrawImage(w.Bullet.SubImage(image.Rect(0, 0, 64, 64)).(*ebiten.Image), &opts)
	}
	opts.GeoM.Reset()
	// debug code
	// vector.DrawFilledRect(screen, float32(bullet.Position.X), float32(bullet.Position.Y), 16, 16, color.RGBA{0, 0, 255, 128}, true)
}
