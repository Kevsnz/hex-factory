package ui

import (
	"hextopdown/settings/strings"
	"hextopdown/utils"
)

type ConfigurableWindow interface {
	AddChild(c iControl, ca ControlAlignment)
	Hide()
}

func WithCloseBox(w ConfigurableWindow) {
	closebox := NewButtonText(
		utils.ScreenCoord{X: CLOSE_BOX_PADDING_PCT, Y: CLOSE_BOX_PADDING_PCT}.PctScaleToScreen().Sub(wndTitleHeight),
		utils.ScreenCoord{X: CLOSE_BOX_SIZE_PCT, Y: CLOSE_BOX_SIZE_PCT}.PctScaleToScreen(),
		strings.STRING_X,
		w.Hide,
	)
	w.AddChild(closebox, CONTROL_ALIGN_TOPRIGHT)
}
