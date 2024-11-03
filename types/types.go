package types

import (
	"context"
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
	pb "github.com/hritesh04/shooter-game/stubs"
)

type Game interface {
	Update() error
	Draw(*ebiten.Image)
	Layout(int, int) (int, int)
	GetSize() (float64, float64)
	GetDevice() Device
	GetFS() embed.FS
	GetClient() *pb.MovementEmitterClient
	SetServerInfo(string, string)
}

type IMap interface {
	Init()
	Update() error
	Draw(*ebiten.Image)
	ListenCommand(string, string)
}

type GrpcFunc func(context.Context, *pb.Room) (*pb.Player, error)

type Screen int

const (
	Onboarding = iota
	Winner
)

const (
	JoinDungeon = iota
	CreateDungeon
)

type IScreen interface {
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

type Device int

const (
	Desktop = iota
	Web
	Mobile
)
