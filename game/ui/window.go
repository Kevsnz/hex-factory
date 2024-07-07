package ui

import (
	"hextopdown/input"
	"hextopdown/renderer"
	"hextopdown/settings/strings"
	"hextopdown/utils"
)

type iControl interface {
	GetPos() utils.ScreenCoord
	SetPos(utils.ScreenCoord)
	GetSize() utils.ScreenCoord
	HandleMouseMovement(utils.ScreenCoord)
	HandleMouseAction(input.MouseButtonEvent) (handled bool)
	Draw(*renderer.GameRenderer, utils.ScreenCoord)
}

type Window struct {
	pos      utils.ScreenCoord
	size     utils.ScreenCoord
	title    strings.StringID
	children []iControl
	visible  bool
	dialog   bool
}

func (w *Window) AddChild(c iControl, ca ControlAlignment) {
	c.SetPos(ca.ConvertCoords(c.GetPos(), c.GetSize(), w.size).Add(wndTitleHeight))
	w.children = append(w.children, c)
}

func (w *Window) Draw(r *renderer.GameRenderer) {
	if !w.visible {
		return
	}

	r.DrawWindow(w.pos, w.size, w.title)
	for _, child := range w.children {
		child.Draw(r, w.pos)
	}
}

func (w *Window) within(mp utils.ScreenCoord) bool {
	return mp.X >= 0 && mp.X < w.size.X && mp.Y >= 0 && mp.Y < w.size.Y
}

func (w *Window) HandleMouseMovement(mp utils.ScreenCoord) {
	mp = mp.Sub(w.pos)
	if !w.visible || !w.within(mp) {
		return
	}

	for _, child := range w.children {
		child.HandleMouseMovement(mp)
	}
}

func (w *Window) HandleMouseAction(mbe input.MouseButtonEvent) bool {
	if !w.visible {
		return false
	}

	mbe.Coord = mbe.Coord.Sub(w.pos)
	if !w.within(mbe.Coord) {
		return false
	}

	for _, child := range w.children {
		if child.HandleMouseAction(mbe) {
			return true
		}
	}
	return true
}

func (w *Window) HandleGameAction(ae input.ActionEvent) bool {
	switch ae.Action {
	case input.ACTION_CANCEL:
		if w.dialog {
			w.Hide()
			return true
		}
	}
	return false
}

func (w *Window) Show() {
	w.visible = true
}
func (w *Window) Hide() {
	w.visible = false
}
func (w *Window) IsVisible() bool { return w.visible }
