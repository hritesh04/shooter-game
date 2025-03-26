package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type KeyboardInput struct {
	ID            int
	ParentScene   int
	width, height float64
	cooldown      int
	runes         []rune
	text          string
	placeholder   string
	roomID        string
	counter       int
	editable      bool
	doneFunc      func() func(string) error
	doneFlag      int
	updateFlag    bool
}

func NewKeyBoardInput(ID, ParentScene int, placeholder string, w, h float64, doneFunc func() func(string) error, doneFlag int) *KeyboardInput {
	return &KeyboardInput{
		ID:          ID,
		ParentScene: ParentScene,
		width:       w,
		height:      h,
		cooldown:    0,
		text:        placeholder,
		placeholder: placeholder,
		counter:     0,
		editable:    true,
		doneFunc:    doneFunc,
		doneFlag:    doneFlag,
		updateFlag:  true,
	}
}

func (k *KeyboardInput) Init() {

}

func (k *KeyboardInput) Update() int {
	if k.editable {
		k.runes = ebiten.AppendInputChars(k.runes[:0])
		k.roomID += string(k.runes)
	}
	if repeatingKeyPressed(ebiten.KeyBackspace) {
		if len(k.roomID) >= 1 {
			k.roomID = k.roomID[:len(k.roomID)-1]
		}
	}
	k.counter++
	if k.updateFlag {
		if len(k.roomID) == 6 {
			k.text = "joining room "
			k.editable = false
			joinRoom := k.doneFunc()
			k.updateFlag = false
			if err := joinRoom(k.roomID); err != nil {
				k.text = err.Error()
				k.cooldown = 60
				k.editable = false
				return k.ID
			} else {
				k.text = k.placeholder
				k.editable = true
				k.roomID = ""
				return k.doneFlag
			}
		} else {
			return k.ID
		}
	} else {
		k.cooldown--
		if k.cooldown <= 0 {
			k.text = k.placeholder
			k.roomID = ""
			k.editable = true
			k.updateFlag = true
			k.counter = 0
			return k.ParentScene
		} else {
			return k.ID
		}

	}
}

func (k *KeyboardInput) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(k.width-10), float32(k.height-10), float32(150), float32(50), color.Black, true)
	t := k.text
	w := k.roomID
	if k.editable {
		if k.counter%60 < 30 {
			w += "_"
		}
	} else {
		dots := k.counter % 6
		for range dots {
			t += "."
		}
	}
	ebitenutil.DebugPrintAt(screen, t, int(k.width), int(k.height))
	ebitenutil.DebugPrintAt(screen, w, int(k.width), int(k.height+20))
}

func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}
