package player

import (
	pb "github.com/hritesh04/shooter-game/proto"
)

type Player struct {
	Name string
	Conn pb.MovementEmitter_SendMoveServer
	X    float64
	Y    float64
}

func NewPlayer(name string, x, y int64) Player {
	return Player{
		Name: name,
		X:    float64(x),
		Y:    float64(y),
	}
}

func (p *Player) AddStream(stream pb.MovementEmitter_SendMoveServer) {
	p.Conn = stream
}

func (p *Player) UpdateLoc(player *pb.Player) {
	p.X = float64(player.GetX())
	p.Y = float64(player.GetY())
}
