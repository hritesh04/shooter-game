package maps

import (
	"encoding/json"
	"os"
)

const (
	DefaultMap = iota
)

type MapLayerJSON struct {
	Data   []int  `json:"data"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Name   string `json:"name"`
}

type TiledMapJSON struct {
	Layers []MapLayerJSON `json:"layers"`
}

var maps = map[int]string{
	DefaultMap: "maps/default_map.json",
}

func NewMap(mapType int) (*TiledMapJSON, error) {
	mapData, err := os.ReadFile(maps[mapType])
	if err != nil {
		return nil, err
	}

	var mapJSON TiledMapJSON
	if err := json.Unmarshal(mapData, &mapJSON); err != nil {
		return nil, err
	}
	return &mapJSON, nil
}
