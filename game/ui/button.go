package ui

import (
	"hextopdown/input"
	"hextopdown/renderer"
	"hextopdown/settings/strings"
	"hextopdown/utils"
)

type Button struct {
	ControlBase
	text    strings.StringID
	down    bool
	onClick func()
}

func NewButton(pos, size utils.ScreenCoord, text strings.StringID, onClick func()) *Button {
	return &Button{
		ControlBase: ControlBase{
			Pos:  pos,
			Size: size,
		},
		text:    text,
		onClick: onClick,
	}
}

func (b *Button) within(mp utils.ScreenCoord) bool {
	return mp.X >= 0 && mp.X < b.Size.X && mp.Y >= 0 && mp.Y < b.Size.Y
}

func (b *Button) Draw(r *renderer.GameRenderer, parentPos utils.ScreenCoord) {
	r.DrawButton(b.Pos.Add(parentPos), b.Size, b.down)
	r.DrawString(b.text, b.Pos.Add(parentPos).Add(b.Size.Div(2)), renderer.TEXT_ALIGN_CENTER)
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
