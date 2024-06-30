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
	item  *items.ItemStack
}

func NewItemSlot(pos utils.ScreenCoord, size utils.ScreenCoord) *ItemSlot {
	return &ItemSlot{
		ControlBase: ControlBase{
			Pos:  pos,
			Size: size,
		},
	}
}

func (i *ItemSlot) HandleMouseMovement(mp utils.ScreenCoord) {
	mp = mp.Sub(i.Pos)
	i.hover = i.within(mp)
}

func (i *ItemSlot) HandleMouseAction(mbe input.MouseButtonEvent) bool {
	return false
}

func (i *ItemSlot) SetItem(item *items.ItemStack) {
	i.item = item
}

func (i *ItemSlot) Draw(r *renderer.GameRenderer, parentPos utils.ScreenCoord) {
	if i.item == nil {
		r.DrawItemSlot(i.Pos.Add(parentPos), i.Size, i.hover)
	} else {
		r.DrawItemSlotWithItem(i.Pos.Add(parentPos), i.Size, i.hover, i.item.ItemType, i.item.Count)
	}
}
