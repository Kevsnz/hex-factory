package objects

import (
	gd "hextopdown/game/gamedata"
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
)

type Converter struct {
	Object
	dir    utils.Dir
	params *gd.ConverterParameters
}

func NewConverter(
	objType ss.ObjectType,
	pos utils.HexCoord,
	dir utils.Dir,
	objParams *gd.ObjectParameters,
	params *gd.ConverterParameters,
) *Converter {
	return &Converter{
		Object: Object{
			objType:   objType,
			pos:       pos,
			objParams: objParams,
		},
		dir:    dir,
		params: params,
	}
}

func (c *Converter) GetDir() utils.Dir {
	return c.dir
}

func (c *Converter) Rotate(_ bool) {}

func (c *Converter) DrawGroundLevel(r *renderer.GameRenderer) {
	r.DrawObjectGround(c.pos.CenterToWorld(), c.objType, c.objParams.Shape, c.dir)
}

func (c *Converter) DrawOnGroundLevel(r *renderer.GameRenderer) {}
