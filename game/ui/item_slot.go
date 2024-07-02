package ui

import (
	"hextopdown/game/items"
	"hextopdown/input"
	"hextopdown/renderer"
	"hextopdown/utils"
)

type ItemSlot struct {
	ControlBase
	hover bool
	slot  *items.StorageSlot
}

func NewItemSlot(pos utils.ScreenCoord, size utils.ScreenCoord, slot *items.StorageSlot) *ItemSlot {
	return &ItemSlot{
		ControlBase: ControlBase{
			Pos:  pos,
			Size: size,
		},
		slot: slot,
	}
}

func (i *ItemSlot) HandleMouseMovement(mp utils.ScreenCoord) {
	mp = mp.Sub(i.Pos)
	i.hover = i.within(mp)
}

func (i *ItemSlot) HandleMouseAction(mbe input.MouseButtonEvent) bool {
	return false
}

func (i *ItemSlot) Draw(r *renderer.GameRenderer, parentPos utils.ScreenCoord) {
	if i.slot == nil {
		panic("slot is not initialized")
	}
	if i.slot.Item == nil {
		r.DrawItemSlot(i.Pos.Add(parentPos), i.Size, i.hover)
	} else {
		r.DrawItemSlotWithItem(i.Pos.Add(parentPos), i.Size, i.hover, i.slot.Item.ItemType, i.slot.Item.Count)
	}
}
