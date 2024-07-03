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
	inventory      items.Storage
	storage        items.Storage
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

func (w *WindowStorage) ShowStorage(objName strings.StringID, inventory []*items.StorageSlot, storage []*items.StorageSlot) {
	w.Title = objName
	w.refillSlots(inventory, storage)
	w.Visible = true
}

func (w *WindowStorage) refillSlots(inventory items.Storage, storage items.Storage) {
	w.inventory = inventory
	w.storage = storage

	w.inventoryPanel.Clear()
	for i, slot := range inventory {
		pos := utils.ScreenCoord{
			X: float32(i%SLOTS_IN_LINE) * (itemSlotSize.X + itemSlotGap),
			Y: float32(i/SLOTS_IN_LINE) * (itemSlotSize.Y + itemSlotGap),
		}
		is := NewItemSlot(pos, itemSlotSize, slot, w.moveStackToStorage)
		w.inventoryPanel.AddChild(is, CONTROL_ALIGN_TOPLEFT)
	}

	w.storagePanel.Clear()
	for i, slot := range storage {
		pos := utils.ScreenCoord{
			X: float32(i%SLOTS_IN_LINE) * (itemSlotSize.X + itemSlotGap),
			Y: float32(i/SLOTS_IN_LINE) * (itemSlotSize.Y + itemSlotGap),
		}
		is := NewItemSlot(pos, itemSlotSize, slot, w.moveStackToInventory)
		w.storagePanel.AddChild(is, CONTROL_ALIGN_TOPLEFT)
	}

	lines := (max(len(inventory), len(storage)) + SLOTS_IN_LINE - 1) / SLOTS_IN_LINE
	slotsHeight := float32(lines) * (itemSlotSize.Y + itemSlotGap)
	w.inventoryPanel.Size.Y = slotsHeight + groupBoxPadding.Y*2
	w.storagePanel.Size.Y = slotsHeight + groupBoxPadding.Y*2
	w.Size.Y = wndTitleHeight.Y + w.inventoryPanel.Size.Y

	w.Pos = utils.ScreenCoord{X: 0.5, Y: 0.5}.PctPosToScreen().Sub(w.Size.Div(2))
}

func (w *WindowStorage) moveStackToStorage(slot *items.StorageSlot) {
	if slot.Item == nil {
		return
	}

	w.storage.TakeItemStackAnywhere(slot.Item)
	if slot.Item.Count == 0 {
		slot.Item = nil
	}
}

func (w *WindowStorage) moveStackToInventory(slot *items.StorageSlot) {
	if slot.Item == nil {
		return
	}

	w.inventory.TakeItemStackAnywhere(slot.Item)
	if slot.Item.Count == 0 {
		slot.Item = nil
	}
}
