package types

import "github.com/hajimehoshi/ebiten/v2"

type Game interface {
	Update() error
	Draw(*ebiten.Image)
	GetSize() (float64, float64)
}

type IMap interface {
	Init()
	Update() error
	Draw(*ebiten.Image)
}

type Direction string

const (
	Left  Direction = "left"
	Right Direction = "right"
	// Up    Direction = "up"
	// Down  Direction = "down"
)
