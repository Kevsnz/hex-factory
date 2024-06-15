package ui

import (
	"hextopdown/renderer"
	"hextopdown/utils"
)

type UI struct {
	show    bool
	windows []*Window
}

func NewUI() *UI {
	b := NewButton(utils.ScreenCoord{X: 20, Y: 20}, utils.ScreenCoord{X: 20, Y: 10})
	w := NewWindow(utils.ScreenCoord{X: 90, Y: 90}, utils.ScreenCoord{X: 50, Y: 50})
	w.AddButton(b)
	return &UI{
		windows: []*Window{w},
	}
}

func (u *UI) Destroy() {}

func (u *UI) Show(b bool) {
	u.show = b
}
func (u *UI) ShowToggle() {
	u.show = !u.show
}

func (u *UI) Draw(r *renderer.GameRenderer) {
	if !u.show {
		return
	}

	for _, w := range u.windows {
		w.Draw(r)
	}
}
