package ui

import (
	"hextopdown/renderer"
	"hextopdown/utils"
)

type Window struct {
	Pos      utils.ScreenCoord
	Size     utils.ScreenCoord
	Children []*Button
}

func NewWindow(pos, size utils.ScreenCoord) *Window {
	return &Window{
		Pos:  pos,
		Size: size,
	}
}

func (w *Window) AddButton(b *Button) {
	w.Children = append(w.Children, b)
}

func (w *Window) Draw(r *renderer.GameRenderer) {
	r.DrawWindow(w.Pos, w.Size)
	for _, child := range w.Children {
		child.Draw(r, w.Pos)
	}
}
