package game

import (
	"hextopdown/game/objects"
	"hextopdown/renderer"
	"hextopdown/settings/strings"
	"hextopdown/utils"
)

type WorldObject interface {
	GetNameString() strings.StringID
	GetPos() utils.HexCoord
	DrawGroundLevel(r *renderer.GameRenderer)
	DrawOnGroundLevel(r *renderer.GameRenderer)
}
type DirectionalObject interface {
	WorldObject
	// GetDir() utils.Dir
	Rotate(cw bool)
}

type Tickable interface {
	Update(ticks uint64, world objects.HexGridWorldInteractor)
}

type ItemMover interface {
	MoveItems(ticks uint64, processed map[*objects.BeltGraphSegment]struct{})
}
type ItemDrawer interface {
	DrawItems(r *renderer.GameRenderer)
}

type ItemHolder interface {
	GetItemList() []utils.ItemInfo
}
