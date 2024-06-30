package ui

import (
	"hextopdown/game/items"
	ss "hextopdown/settings"
	"hextopdown/settings/strings"
	"hextopdown/utils"
)

const SLOTS_IN_LINE = 8

var itemSlotSize = utils.ScreenCoord{X: ss.FONT_SIZE_PCT * 2, Y: ss.FONT_SIZE_PCT * 2}.PctScaleToScreen()
var itemSlotGap = max(1, itemSlotSize.X*0.05)

type WindowInventory struct {
	Window
	slotCount int
}

func NewWindowInventory() *WindowInventory {
	wnd := &WindowInventory{
		Window: Window{
			Pos:     utils.ScreenCoord{X: 0.5, Y: 0.5}.PctPosToScreen(),
			Size:    utils.ScreenCoord{X: 0.5, Y: 0.5}.PctScaleToScreen(),
			Title:   strings.STRING_INVENTORY,
			Visible: false,
		},
	}
	WithCloseBox(wnd)
	return wnd
}

func (w *WindowInventory) ShowInventory(inventory []*items.ItemStack) {
	if w.slotCount != len(inventory) {
		w.refillSlots(inventory)
	}
	w.Visible = true
}

func (w *WindowInventory) refillSlots(inventory []*items.ItemStack) {
	newChildren := make([]iControl, 0, len(inventory)+2)
	for _, c := range w.children {
		if _, ok := c.(*ItemSlot); !ok {
			newChildren = append(newChildren, c)
		}
	}

	w.children = newChildren

	w.slotCount = len(inventory)
	for i, item := range inventory {
		row := i / SLOTS_IN_LINE
		col := i % SLOTS_IN_LINE
		pos := utils.ScreenCoord{
			X: float32(col) * (itemSlotSize.X + itemSlotGap),
			Y: float32(row)*(itemSlotSize.Y+itemSlotGap) + windowTitleHeight,
		}
		is := NewItemSlot(pos, itemSlotSize)
		w.AddChild(is, CONTROL_ALIGN_TOPLEFT)
		is.SetItem(item)
	}
}
