package manager

import (
	"log"

	game "github.com/hritesh04/shooter-game/server/gameServer/game"
)

type GameManager struct {
	games map[string]*game.Game
}

func NewGameManager() *GameManager {
	game := &GameManager{
		games: make(map[string]*game.Game),
	}
	return game
}

func (g *GameManager) AddRoom(roomID string) {
	_, ok := g.games[roomID]
	if ok {
		log.Printf("dungeon already present %s", roomID)
		return
	}
	game := game.NewGame(roomID)
	g.games[roomID] = game
	log.Printf("created room : %s\n", roomID)
}

func (g *GameManager) GetRoom(roomID string) *game.Game {
	game, ok := g.games[roomID]
	if !ok {
		log.Printf("dungeon not present %s", roomID)
		return nil
	}
	return game
}
