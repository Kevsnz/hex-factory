package world

import (
	"hextopdown/renderer"
	"hextopdown/settings"
	"hextopdown/utils"
)

type WorldObject interface {
	GetObjectType() settings.ObjectType
	GetPos() utils.HexCoord
	DrawGroundLevel(r *renderer.GameRenderer)
	DrawOnGroundLevel(r *renderer.GameRenderer)
}
