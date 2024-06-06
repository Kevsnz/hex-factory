package items

import (
	"hextopdown/renderer"
	"hextopdown/settings"
	"hextopdown/utils"
)

type ItemInWorld struct {
	itemType settings.ItemType
	Pos      utils.WorldCoordInterpolated
}

func NewItemInWorld(itemType settings.ItemType, pos utils.WorldCoordInterpolated) ItemInWorld {
	return ItemInWorld{
		itemType: itemType,
		Pos:      pos,
	}
}

func NewItemInWorld2(itemType settings.ItemType, pos utils.WorldCoord) ItemInWorld {
	return ItemInWorld{
		itemType: itemType,
		Pos:      utils.NewWorldCoordInterpolated(),
	}
}

func (i *ItemInWorld) Draw(r *renderer.GameRenderer) {
	r.DrawItem(i.Pos, i.itemType)
}
