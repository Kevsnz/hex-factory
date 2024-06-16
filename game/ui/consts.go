package ui

import (
	"hextopdown/settings"
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

const CLOSE_BOX_SIZE_PCT = settings.FONT_SIZE_PCT * 2 / 3
const CLOSE_BOX_PADDING_PCT = settings.FONT_SIZE_PCT / 8
