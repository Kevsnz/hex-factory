package objects

import (
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/settings/strings"
	"hextopdown/utils"
)

var dirHexes = [utils.DIR_COUNT][4]utils.HexCoord{
	utils.DIR_LEFT:       {utils.HexCoord{X: 0, Y: 0}, utils.HexCoord{X: 1, Y: -1}, utils.HexCoord{X: 1, Y: 0}, utils.HexCoord{X: 0, Y: 1}},
	utils.DIR_RIGHT:      {utils.HexCoord{X: 0, Y: 0}, utils.HexCoord{X: 1, Y: -1}, utils.HexCoord{X: 1, Y: 0}, utils.HexCoord{X: 0, Y: 1}},
	utils.DIR_UP_LEFT:    {utils.HexCoord{X: 0, Y: 0}, utils.HexCoord{X: 1, Y: 0}, utils.HexCoord{X: 0, Y: 1}, utils.HexCoord{X: -1, Y: 1}},
	utils.DIR_DOWN_RIGHT: {utils.HexCoord{X: 0, Y: 0}, utils.HexCoord{X: 1, Y: 0}, utils.HexCoord{X: 0, Y: 1}, utils.HexCoord{X: -1, Y: 1}},
	utils.DIR_UP_RIGHT:   {utils.HexCoord{X: 0, Y: 0}, utils.HexCoord{X: 0, Y: -1}, utils.HexCoord{X: 1, Y: -1}, utils.HexCoord{X: 1, Y: 0}},
	utils.DIR_DOWN_LEFT:  {utils.HexCoord{X: 0, Y: 0}, utils.HexCoord{X: 0, Y: -1}, utils.HexCoord{X: 1, Y: -1}, utils.HexCoord{X: 1, Y: 0}},
}

type Furnace struct {
	pos utils.HexCoord
	dir utils.Dir
}

func NewFurnace(pos utils.HexCoord, dir utils.Dir) *Furnace {
	return &Furnace{
		pos: pos,
		dir: utils.DIR_RIGHT,
	}
}

func (f *Furnace) GetNameString() strings.StringID {
	return strings.STRING_OBJECT_FURNACE
}

func (f *Furnace) GetPos() utils.HexCoord {
	return f.pos
}

func (f *Furnace) DrawGroundLevel(r *renderer.GameRenderer) {
	r.DrawStructureGround2(f.pos.CenterToWorld(), ss.OBJECT_TYPE_FURNACE_STONE, ss.SHAPE_DIAMOND, f.dir)
}

func (f *Furnace) DrawOnGroundLevel(r *renderer.GameRenderer) {}

func (f *Furnace) GetHexes() []utils.HexCoord {
	hexes := make([]utils.HexCoord, 0, 4)
	for _, o := range dirHexes[f.dir] {
		hexes = append(hexes, f.pos.Add(o))
	}
	return hexes
}
