package objects

import (
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
)

type Converter struct {
	objType ss.ObjectType
	pos     utils.HexCoord
	dir     utils.Dir
}

func NewConverter(objType ss.ObjectType, pos utils.HexCoord, dir utils.Dir) *Converter {
	return &Converter{
		objType: objType,
		pos:     pos,
		dir:     dir,
	}
}

func (c *Converter) GetObjectType() ss.ObjectType {
	return c.objType
}

func (c *Converter) GetPos() utils.HexCoord {
	return c.pos
}

func (c *Converter) GetDir() utils.Dir {
	return c.dir
}

func (c *Converter) Rotate(_ bool) {}

func (c *Converter) DrawGroundLevel(r *renderer.GameRenderer) {
	r.DrawObjectGround(c.pos.CenterToWorld(), ss.OBJECT_TYPE_FURNACE_STONE, utils.SHAPE_DIAMOND, c.dir)
}

func (c *Converter) DrawOnGroundLevel(r *renderer.GameRenderer) {}
