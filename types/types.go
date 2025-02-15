package types

import (
	"context"
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
	pb "github.com/hritesh04/shooter-game/proto"
	"google.golang.org/grpc"
)

type Game interface {
	Update() error
	Draw(*ebiten.Image)
	Layout(int, int) (int, int)
	GetSize() (float64, float64)
	GetDevice() Device
	GetFS() embed.FS
	GetClient() *pb.MovementEmitterClient //
	SetServerInfo(string, string)
	TogglePopUp(bool)
}

type IMap interface {
	Init()
	Update() error
	Draw(*ebiten.Image)
	JoinRoom(string, string) error
	ListenCommand(string, string)
}

type IConnection interface {
	JoinRoom(string) (*pb.Room, error)
	GetEventConn() pb.MovementEmitter_SendMoveClient
	SendMove(grpc.BidiStreamingClient[pb.Data, pb.Data], *pb.Data) error
}

type GrpcFunc func(context.Context, *pb.Room) (*pb.Player, error)

type IScreen interface {
	Init()
	Update() int
	Draw(*ebiten.Image)
}

const (
	Map = iota
	Onboarding
	Winner
	JoinDungeon
	CreateDungeon
	InputBox
)

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
