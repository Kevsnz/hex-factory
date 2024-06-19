package utils

import (
	"hextopdown/settings"
)

var hexesDiamond = [DIR_COUNT][4]HexCoord{
	DIR_LEFT:       {{X: 0, Y: 0}, {X: 1, Y: -1}, {X: 1, Y: 0}, {X: 0, Y: 1}},
	DIR_RIGHT:      {{X: 0, Y: 0}, {X: 1, Y: -1}, {X: 1, Y: 0}, {X: 0, Y: 1}},
	DIR_UP_LEFT:    {{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 1}},
	DIR_DOWN_RIGHT: {{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 1}},
	DIR_UP_RIGHT:   {{X: 0, Y: 0}, {X: 0, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 0}},
	DIR_DOWN_LEFT:  {{X: 0, Y: 0}, {X: 0, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 0}},
}
var hexesBighex = [7]HexCoord{
	{0, 0},
	{1, 0}, {-1, 0},
	{0, -1}, {0, 1},
	{1, -1}, {-1, 1},
}

type Shape int

func (s Shape) GetHexes(pos HexCoord, dir Dir) []HexCoord {
	switch s {
	case SHAPE_SINGLE:
		return []HexCoord{pos}
	case SHAPE_DIAMOND:
		hexes := make([]HexCoord, 0, 4)
		for _, o := range hexesDiamond[dir] {
			hexes = append(hexes, pos.Add(o))
		}
		return hexes
	case SHAPE_BIGHEX:
		hexes := make([]HexCoord, 0, 7)
		for _, o := range hexesBighex {
			hexes = append(hexes, pos.Add(o))
		}
		return hexes
	}
	panic("invalid shape")
}

func (s Shape) GetCenterOffset(dir Dir) WorldCoord {
	if s != SHAPE_DIAMOND {
		return WorldCoord{X: 0, Y: 0}
	}

	if dir == DIR_LEFT || dir == DIR_RIGHT {
		return WorldCoord{X: settings.HEX_WIDTH / 2, Y: 0}
	}

	if dir == DIR_UP_LEFT || dir == DIR_DOWN_RIGHT {
		return WorldCoord{X: settings.HEX_WIDTH / 4, Y: settings.HEX_EDGE * 3 / 4}
	}

	return WorldCoord{X: settings.HEX_WIDTH / 4, Y: -settings.HEX_EDGE * 3 / 4}
}

const (
	SHAPE_SINGLE  Shape = iota // single hex
	SHAPE_DIAMOND Shape = iota // diamond (4 hexes)
	SHAPE_BIGHEX  Shape = iota // big hex

	SHAPE_COUNT Shape = iota
)
