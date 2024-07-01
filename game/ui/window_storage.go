package ui

import (
	"hextopdown/game/items"
	"hextopdown/settings/strings"
	"hextopdown/utils"
)

type WindowStorage struct {
	Window
	inventoryPanel *GroupBox
	storagePanel   *GroupBox
	// inventorySlots []*ItemSlot
	// storageSlots   []*ItemSlot
}

func NewWindowStorage() *WindowStorage {
	size := utils.ScreenCoord{X: wndStorageWidth, Y: wndTitleHeight.Y}
	pos := utils.ScreenCoord{X: 0.5, Y: 0.5}.PctPosToScreen().Sub(size.Div(2))
	wnd := &WindowStorage{
		Window: Window{
			Pos:     pos,
			Size:    size,
			Title:   strings.STRING_STORAGE,
			Visible: false,
		},
	}
	WithCloseBox(wnd)
	wnd.inventoryPanel = NewGroupBox(
		utils.ScreenCoord{X: 0, Y: 0},
		utils.ScreenCoord{X: wndStorageInvWidth, Y: groupBoxPadding.Y * 2},
		strings.STRING_INVENTORY,
	)
	wnd.storagePanel = NewGroupBox(
		utils.ScreenCoord{X: 0, Y: 0},
		utils.ScreenCoord{X: wndStorageStorWidth, Y: groupBoxPadding.Y * 2},
		strings.STRING_STORAGE,
	)
	wnd.AddChild(wnd.inventoryPanel, CONTROL_ALIGN_TOPLEFT)
	wnd.AddChild(wnd.storagePanel, CONTROL_ALIGN_TOPRIGHT)
	return wnd
}

func (w *WindowStorage) ShowStorage(objName strings.StringID, inventory []*items.ItemStack, storage []*items.ItemStack) {
	w.Title = objName
	w.refillSlots(inventory, storage)
	w.Visible = true
}

func (w *WindowStorage) refillSlots(inventory []*items.ItemStack, storage []*items.ItemStack) {
	w.inventoryPanel.Clear()
	for i, item := range inventory {
		pos := utils.ScreenCoord{
			X: float32(i%SLOTS_IN_LINE) * (itemSlotSize.X + itemSlotGap),
			Y: float32(i/SLOTS_IN_LINE) * (itemSlotSize.Y + itemSlotGap),
		}
		is := NewItemSlot(pos, itemSlotSize)
		is.SetItem(item)
		w.inventoryPanel.AddChild(is, CONTROL_ALIGN_TOPLEFT)
	}

	w.storagePanel.Clear()
	for i, item := range storage {
		pos := utils.ScreenCoord{
			X: float32(i%SLOTS_IN_LINE) * (itemSlotSize.X + itemSlotGap),
			Y: float32(i/SLOTS_IN_LINE) * (itemSlotSize.Y + itemSlotGap),
		}
		is := NewItemSlot(pos, itemSlotSize)
		is.SetItem(item)
		w.storagePanel.AddChild(is, CONTROL_ALIGN_TOPLEFT)
	}

	lines := (max(len(inventory), len(storage)) + SLOTS_IN_LINE - 1) / SLOTS_IN_LINE
	slotsHeight := float32(lines) * (itemSlotSize.Y + itemSlotGap)
	w.inventoryPanel.Size.Y = slotsHeight + groupBoxPadding.Y*2
	w.storagePanel.Size.Y = slotsHeight + groupBoxPadding.Y*2
	w.Size.Y = wndTitleHeight.Y + w.inventoryPanel.Size.Y

	w.Pos = utils.ScreenCoord{X: 0.5, Y: 0.5}.PctPosToScreen().Sub(w.Size.Div(2))
}
