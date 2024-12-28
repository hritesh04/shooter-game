package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// repeatingKeyPressed return true when key is pressed considering the repeat state.

type KeyboardInput struct {
	ID            int
	ParentScene   int
	width, height float64
	cooldown      int
	runes         []rune
	text          string
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
	// Add runes that are input by the user by AppendInputChars.
	// Note that AppendInputChars result changes every frame, so you need to call this
	// every frame.
	if k.editable {
		k.runes = ebiten.AppendInputChars(k.runes[:0])
		k.roomID += string(k.runes)
	}

	// Adjust the string to be at most 10 lines.
	// ss := strings.Split(g.text, "\n")
	// if len(ss) > 10 {
	// 	g.text = strings.Join(ss[len(ss)-10:], "\n")
	// }

	// If the enter key is pressed, add a line break.
	// if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter) {
	// 	k.text += "\n"
	// }

	// If the backspace key is pressed, remove one character.
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
			// ctx := context.Background()
			if err := joinRoom(k.roomID); err != nil {
				k.text = err.Error()
				k.cooldown = 60
				k.editable = false
				return k.ID
			} else {
				k.text = "Enter the dungeon ID\n"
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
			k.text = "Enter the dungeon ID\n"
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
	// Blink the cursor.
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
		// dots := strings.Repeat(".", k.counter%4)
		// t += dots

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

// func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	return screenWidth, screenHeight
// }

// func main() {
// 	g := &Game{
// 		text:    "Type on the keyboard:\n",
// 		counter: 0,
// 	}

// 	ebiten.SetWindowSize(screenWidth, screenHeight)
// 	ebiten.SetWindowTitle("TypeWriter (Ebitengine Demo)")
// 	if err := ebiten.RunGame(g); err != nil {
// 		log.Fatal(err)
// 	}
// }
