package ui

import (
	"hextopdown/game/items"
	"hextopdown/settings/strings"
	"hextopdown/utils"
)

type WindowInventory struct {
	Window
	storage items.Storage
}

func NewWindowInventory() *WindowInventory {
	size := utils.ScreenCoord{X: wndInventoryWidth, Y: wndTitleHeight.Y}
	pos := utils.ScreenCoord{X: 0.5, Y: 0.5}.PctPosToScreen().Sub(size.Div(2))
	wnd := &WindowInventory{
		Window: Window{
			pos:     pos,
			size:    size,
			title:   strings.STRING_INVENTORY,
			visible: false,
			dialog:  true,
		},
	}
	WithCloseBox(wnd)
	return wnd
}

func (w *WindowInventory) ShowInventory(inventory []*items.StorageSlot) {
	w.refillSlots(inventory)
	w.visible = true
}

func (w *WindowInventory) refillSlots(inventory items.Storage) {
	newChildren := make([]iControl, 0, len(inventory)+2)
	for _, c := range w.children {
		if _, ok := c.(*ItemSlot); !ok {
			newChildren = append(newChildren, c)
		}
	}

	w.children = newChildren

	w.storage = inventory
	for i, slot := range inventory {
		pos := utils.ScreenCoord{
			X: float32(i%SLOTS_IN_LINE) * (itemSlotSize.X + itemSlotGap),
			Y: float32(i/SLOTS_IN_LINE) * (itemSlotSize.Y + itemSlotGap),
		}
		is := NewItemSlot(pos, itemSlotSize, slot, nil)
		w.AddChild(is, CONTROL_ALIGN_TOPLEFT)
	}

	w.size.Y = (itemSlotSize.X+itemSlotGap)*float32(len(inventory))/SLOTS_IN_LINE + wndTitleHeight.Y
	if len(inventory)%SLOTS_IN_LINE != 0 {
		w.size.Y += itemSlotSize.X + itemSlotGap
	}
	w.pos.Y = utils.ScreenCoord{X: 0.5, Y: 0.5}.PctPosToScreen().Y - w.size.Y/2
}
