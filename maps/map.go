package maps

import (
	def "github.com/hritesh04/shooter-game/maps/default"
	"github.com/hritesh04/shooter-game/types"
)

const (
	DefaultMap = iota
	NewDefMap
)

var maps = map[int]func(types.Game) types.IMap{
	NewDefMap: def.NewDefaultMap,
}

func NewMap(mapType int, game types.Game) types.IMap {
	return maps[mapType](game)
}
