package game

import (
	"hextopdown/game/objects"
	"hextopdown/renderer"
	"hextopdown/utils"
)

type DirectionalObject interface {
	GetDir() utils.Dir
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
