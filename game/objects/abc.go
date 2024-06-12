package objects

import (
	gd "hextopdown/game/gamedata"
	ss "hextopdown/settings"
	"hextopdown/utils"
)

type Object struct {
	objType   ss.ObjectType
	pos       utils.HexCoord
	objParams *gd.ObjectParameters
}

func (o *Object) GetObjectType() ss.ObjectType {
	return o.objType
}

func (o *Object) GetPos() utils.HexCoord {
	return o.pos
}

type ObjectBeltlike struct {
	Object
	tier       ss.BeltTier
	tierParams *gd.BeltTierParameters
}

func (obl *ObjectBeltlike) GetTier() ss.BeltTier {
	return obl.tier
}
