package items

import (
	"hextopdown/renderer"
	"hextopdown/settings"
	"hextopdown/utils"
)

type ItemOnBelt struct {
	Item      ItemInWorld
	Offset    float64
	MovedTick uint64
}

func NewItemOnBelt2(item ItemInWorld, offset float64) *ItemOnBelt {
	return &ItemOnBelt{
		Item:   item,
		Offset: offset,
	}
}

func NewItemOnBelt(itemType settings.ItemType, pos *utils.WorldCoordInterpolated, offset float64) *ItemOnBelt {
	newPos := utils.NewWorldCoordInterpolated()
	if pos != nil {
		newPos = *pos
	}
	return &ItemOnBelt{
		Item:   NewItemInWorld(itemType, newPos),
		Offset: offset,
	}
}

func (iob *ItemOnBelt) Draw(r *renderer.GameRenderer) {
	iob.Item.Draw(r)
}
