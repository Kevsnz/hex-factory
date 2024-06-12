package objects

import (
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
)

type Converter struct {
	Object
	dir utils.Dir
}

func NewConverter(objType ss.ObjectType, pos utils.HexCoord, dir utils.Dir) *Converter {
	return &Converter{
		Object: Object{
			objType: objType,
			pos:     pos,
		},
		dir: dir,
	}
}

func (c *Converter) GetDir() utils.Dir {
	return c.dir
}

func (c *Converter) Rotate(_ bool) {}

func (c *Converter) DrawGroundLevel(r *renderer.GameRenderer) {
	r.DrawObjectGround(c.pos.CenterToWorld(), c.objType, utils.SHAPE_DIAMOND, c.dir)
}

func (c *Converter) DrawOnGroundLevel(r *renderer.GameRenderer) {}
