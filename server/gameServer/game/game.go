package game

import (
	"crypto/rand"
	"encoding/hex"
	"log"

	pb "github.com/hritesh04/shooter-game/proto"
	players "github.com/hritesh04/shooter-game/server/gameServer/player"
)

type Game struct {
	ID      string
	Players []players.Player
	Started bool
}

func NewGame(id string) *Game {
	return &Game{ID: id}
}

func (g *Game) EmitMove(name string, action pb.Action, direction pb.Direction, player []*pb.Player) {
	data := &pb.Data{
		Type:   action,
		Data:   direction,
		Name:   name,
		Player: player,
	}
	for _, p := range g.Players {
		if p.Name == name {
			continue
		}
		p.Conn.Send(data)
	}
}

func (g *Game) AddPlayer() *players.Player {
	var player players.Player
	if len(g.Players) == 0 {
		player = players.NewPlayer(generateSecureID(), 60, 70)
	} else {
		player = players.NewPlayer(generateSecureID(), 1172, 608)
		for _, p := range g.Players {
			log.Printf("Player %s data sent to %s\n", player.Name, p.Name)
			p.Conn.Send(&pb.Data{Type: pb.Action_Info, Player: []*pb.Player{{Name: player.Name, X: float32(player.X), Y: float32(player.Y)}}})
		}
	}
	g.Players = append(g.Players, player)
	log.Printf("Player %s joined room %s\n", player.Name, g.ID)
	return &player
}

func (g *Game) GetPlayer(name string) *players.Player {
	for i := range g.Players {
		if g.Players[i].Name == name {
			return &g.Players[i]
		}
	}
	return nil
}

func generateSecureID() string {
	b := make([]byte, 3) // 3 bytes = 6 hex characters
	rand.Read(b)
	return hex.EncodeToString(b)
}
