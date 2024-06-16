package ui

import (
	"hextopdown/input"
	"hextopdown/utils"
)

type Button struct {
	ControlBase
	down    bool
	onClick func()
}

func (b *Button) HandleMouseMovement(mp utils.ScreenCoord) {
	mp = mp.Sub(b.Pos)
	if !b.within(mp) {
		b.down = false
		return
	}
}

func (b *Button) HandleMouseAction(mbe input.MouseButtonEvent) bool {
	mbe.Coord = mbe.Coord.Sub(b.Pos)
	if !b.within(mbe.Coord) {
		return false
	}

	switch mbe.Type {
	case input.MOUSE_BUTTON_DOWN:
		b.down = true

	case input.MOUSE_BUTTON_UP:
		if !b.down {
			break
		}
		b.onClick()
		b.down = false
	}

	return true
}
