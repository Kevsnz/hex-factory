package ui

import (
	ss "hextopdown/settings"
	"hextopdown/utils"
)

type ControlAlignment int

const (
	CONTROL_ALIGN_TOPLEFT     ControlAlignment = iota
	CONTROL_ALIGN_TOPRIGHT    ControlAlignment = iota
	CONTROL_ALIGN_BOTTOMLEFT  ControlAlignment = iota
	CONTROL_ALIGN_BOTTOMRIGHT ControlAlignment = iota
)

func (ca ControlAlignment) ConvertCoords(pos, size, space utils.ScreenCoord) utils.ScreenCoord {
	switch ca {
	case CONTROL_ALIGN_TOPLEFT:
		return utils.ScreenCoord{X: pos.X, Y: pos.Y}
	case CONTROL_ALIGN_TOPRIGHT:
		return utils.ScreenCoord{X: space.X - pos.X - size.X, Y: pos.Y}
	case CONTROL_ALIGN_BOTTOMLEFT:
		return utils.ScreenCoord{X: pos.X, Y: space.Y - pos.Y - size.Y}
	case CONTROL_ALIGN_BOTTOMRIGHT:
		return utils.ScreenCoord{X: space.X - pos.X - size.X, Y: space.Y - pos.Y - size.Y}
	default:
		panic("unknown control alignment")
	}
}

const CLOSE_BOX_SIZE_PCT = ss.FONT_SIZE_PCT * 2 / 3
const CLOSE_BOX_PADDING_PCT = ss.FONT_SIZE_PCT / 8

const SLOTS_IN_LINE = 8

var itemSlotSize = utils.ScreenCoord{X: ss.FONT_SIZE_PCT * 2, Y: ss.FONT_SIZE_PCT * 2}.PctScaleToScreen()
var itemSlotGap = max(1, itemSlotSize.X*0.05)

var wndTitleHeight = utils.ScreenCoord{X: 0, Y: ss.FONT_SIZE_PCT * 1.2}.PctScaleToScreen()
var groupBoxPadding = utils.ScreenCoord{X: ss.FONT_SIZE_PCT, Y: ss.FONT_SIZE_PCT}.PctScaleToScreen()

var wndInventoryWidth = itemSlotSize.X*SLOTS_IN_LINE + itemSlotGap*(SLOTS_IN_LINE-1)
var wndStorageWidth = (itemSlotSize.X*SLOTS_IN_LINE+itemSlotGap*(SLOTS_IN_LINE-1))*2 + groupBoxPadding.X*4
var wndStorageInvWidth = (itemSlotSize.X*SLOTS_IN_LINE + itemSlotGap*(SLOTS_IN_LINE-1)) + groupBoxPadding.X*2
var wndStorageStorWidth = (itemSlotSize.X*SLOTS_IN_LINE + itemSlotGap*(SLOTS_IN_LINE-1)) + groupBoxPadding.X*2

const WINDOW_STORAGE_WIDTH_PCT = 0.75
const WINDOW_STORAGE_HEIGHT_PCT = 0.5
