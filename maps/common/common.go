package common

type MapLayerJSON struct {
	Data   []int  `json:"data"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Name   string `json:"name"`
}

type TiledMapJSON struct {
	Layers []MapLayerJSON `json:"layers"`
}
