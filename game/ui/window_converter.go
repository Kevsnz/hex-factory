package ui

import (
	"hextopdown/game/items"
	"hextopdown/settings/strings"
	"hextopdown/utils"
)

type WindowConverter struct {
	Window
	inventoryPanel *GroupBox
	converterPanel *GroupBox
	inventory      items.Storage
	converter      items.Storage
}

func NewWindowConverter() *WindowConverter {
	size := utils.ScreenCoord{X: wndStorageWidth, Y: wndTitleHeight.Y}
	pos := utils.ScreenCoord{X: 0.5, Y: 0.5}.PctPosToScreen().Sub(size.Div(2))
	wnd := &WindowConverter{
		Window: Window{
			pos:     pos,
			size:    size,
			title:   strings.STRING_STORAGE,
			visible: false,
			dialog:  true,
		},
	}
	WithCloseBox(wnd)
	wnd.inventoryPanel = NewGroupBox(
		utils.ScreenCoord{X: 0, Y: 0},
		utils.ScreenCoord{X: wndStorageInvWidth, Y: groupBoxPadding.Y * 2},
		strings.STRING_INVENTORY,
	)
	wnd.converterPanel = NewGroupBox(
		utils.ScreenCoord{X: 0, Y: 0},
		utils.ScreenCoord{X: wndStorageStorWidth, Y: groupBoxPadding.Y * 2},
		strings.STRING_RECIPE,
	)
	wnd.AddChild(wnd.inventoryPanel, CONTROL_ALIGN_TOPLEFT)
	wnd.AddChild(wnd.converterPanel, CONTROL_ALIGN_TOPRIGHT)
	return wnd
}

func (w *WindowConverter) ShowConverter(
	objName strings.StringID, inventory items.Storage, inputSlots items.Storage, outputSlots items.Storage,
) {
	w.title = objName
	w.refillSlots(inventory, inputSlots, outputSlots)
	w.visible = true
}

func (w *WindowConverter) refillSlots(inventory items.Storage, inputSlots items.Storage, outputSlots items.Storage) {
	w.inventory = inventory
	w.converter = inputSlots

	w.inventoryPanel.Clear()
	for i, slot := range inventory {
		pos := utils.ScreenCoord{
			X: float32(i%SLOTS_IN_LINE) * (itemSlotSize.X + itemSlotGap),
			Y: float32(i/SLOTS_IN_LINE) * (itemSlotSize.Y + itemSlotGap),
		}
		is := NewItemSlot(pos, itemSlotSize, slot, nil) // TODO Callback!
		w.inventoryPanel.AddChild(is, CONTROL_ALIGN_TOPLEFT)
	}

	w.converterPanel.Clear()
	for i, slot := range inputSlots {
		pos := utils.ScreenCoord{
			X: float32(i%SLOTS_IN_LINE) * (itemSlotSize.X + itemSlotGap),
			Y: float32(i/SLOTS_IN_LINE) * (itemSlotSize.Y + itemSlotGap),
		}
		is := NewItemSlot(pos, itemSlotSize, slot, nil) // TODO Callback!
		w.converterPanel.AddChild(is, CONTROL_ALIGN_TOPLEFT)
	}
	for i, slot := range outputSlots {
		pos := utils.ScreenCoord{
			X: float32(i%SLOTS_IN_LINE) * (itemSlotSize.X + itemSlotGap),
			Y: float32(i/SLOTS_IN_LINE+2) * (itemSlotSize.Y + itemSlotGap),
		}
		is := NewItemSlot(pos, itemSlotSize, slot, nil) // TODO Callback!
		w.converterPanel.AddChild(is, CONTROL_ALIGN_TOPLEFT)
	}

	linesIn := (len(inventory) + SLOTS_IN_LINE - 1) / SLOTS_IN_LINE
	linesOut := (len(outputSlots)+SLOTS_IN_LINE-1)/SLOTS_IN_LINE + 2
	lines := max(linesIn, linesOut)

	slotsHeight := float32(lines) * (itemSlotSize.Y + itemSlotGap)
	w.inventoryPanel.Size.Y = slotsHeight + groupBoxPadding.Y*2
	w.converterPanel.Size.Y = slotsHeight + groupBoxPadding.Y*2
	w.size.Y = wndTitleHeight.Y + w.inventoryPanel.Size.Y

	w.pos = utils.ScreenCoord{X: 0.5, Y: 0.5}.PctPosToScreen().Sub(w.size.Div(2))
}
