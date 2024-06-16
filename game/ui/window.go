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
	HandleMouseAction(input.MouseButtonEvent) bool
	Draw(*renderer.GameRenderer, utils.ScreenCoord)
}

type Window struct {
	Pos      utils.ScreenCoord
	Size     utils.ScreenCoord
	Title    strings.StringID
	children []iControl
	Visible  bool
}

func (w *Window) AddChild(c iControl, ca ControlAlignment) {
	c.SetPos(ca.ConvertCoords(c.GetPos(), c.GetSize(), w.Size))
	w.children = append(w.children, c)
}

func (w *Window) Draw(r *renderer.GameRenderer) {
	if !w.Visible {
		return
	}

	r.DrawWindow(w.Pos, w.Size, w.Title)
	for _, child := range w.children {
		child.Draw(r, w.Pos)
	}
}

func (w *Window) within(mp utils.ScreenCoord) bool {
	return mp.X >= 0 && mp.X < w.Size.X && mp.Y >= 0 && mp.Y < w.Size.Y
}

func (w *Window) HandleMouseMovement(mp utils.ScreenCoord) {
	mp = mp.Sub(w.Pos)
	if !w.Visible || !w.within(mp) {
		return
	}

	for _, child := range w.children {
		child.HandleMouseMovement(mp)
	}
}

func (w *Window) HandleMouseAction(mbe input.MouseButtonEvent) bool {
	if !w.Visible {
		return false
	}

	mbe.Coord = mbe.Coord.Sub(w.Pos)
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

func (w *Window) Show() {
	w.Visible = true
}
func (w *Window) Hide() {
	w.Visible = false
}
