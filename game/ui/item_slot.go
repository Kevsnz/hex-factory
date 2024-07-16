package ui

import (
	"hextopdown/game/items"
	"hextopdown/input"
	"hextopdown/renderer"
	"hextopdown/settings"
	"hextopdown/utils"
)

type ItemSlot struct {
	ControlBase
	hover        bool
	Slot         *items.StorageSlot
	onTakeOutAll func(*items.StorageSlot)
}

func NewItemSlot(pos utils.ScreenCoord, size utils.ScreenCoord, slot *items.StorageSlot, onTakeOutAll func(*items.StorageSlot)) *ItemSlot {
	return &ItemSlot{
		ControlBase: ControlBase{
			Pos:  pos,
			Size: size,
		},
		Slot:         slot,
		onTakeOutAll: onTakeOutAll,
	}
}

func (i *ItemSlot) HandleMouseMovement(mp utils.ScreenCoord) {
	mp = mp.Sub(i.Pos)
	i.hover = i.within(mp)
}

func (i *ItemSlot) HandleMouseAction(mbe input.MouseButtonEvent) bool {
	mbe.Coord = mbe.Coord.Sub(i.Pos)
	if !i.within(mbe.Coord) {
		return false
	}

	if mbe.Type != input.MOUSE_BUTTON_DOWN {
		return true
	}

	switch mbe.Button {
	case input.MOUSE_BUTTON_LEFT:
		if i.onTakeOutAll != nil {
			i.onTakeOutAll(i.Slot)
		}
	}

	return true
}

func (i *ItemSlot) Draw(r *renderer.GameRenderer, parentPos utils.ScreenCoord) {
	if i.Slot == nil {
		panic("slot is not initialized")
	}
	if i.Slot.Item == nil {
		if i.Slot.FixedItemType != settings.ITEM_TYPE_COUNT {
			r.DrawItemSlotWithItemShadow(i.Pos.Add(parentPos), i.Size, i.hover, i.Slot.FixedItemType)
		} else {
			r.DrawItemSlot(i.Pos.Add(parentPos), i.Size, i.hover)
		}
	} else {
		r.DrawItemSlotWithItem(i.Pos.Add(parentPos), i.Size, i.hover, i.Slot.Item.ItemType, i.Slot.Item.Count)
	}
}
