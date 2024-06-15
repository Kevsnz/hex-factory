package ui

import (
	"hextopdown/input"
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
)

const (
	WINDOW_RECIPES = iota

	WINDOW_COUNT
)

type iWindow interface {
	HandleMouseMovement(utils.ScreenCoord)
	HandleMouseAction(input.MouseButtonEvent) bool
	Draw(*renderer.GameRenderer)
	Show()
	Hide()
}

type UI struct {
	show    bool
	windows [WINDOW_COUNT]iWindow
	scale   float32
}

func NewUI() *UI {
	w := NewWindowRecipes()
	w.Visible = false

	return &UI{
		windows: [WINDOW_COUNT]iWindow{w},
		scale:   1,
		show:    true,
	}
}

func (u *UI) Destroy() {}

func (u *UI) Show(b bool) {
	u.show = b
}
func (u *UI) ShowToggle() {
	u.show = !u.show
}
func (u *UI) WindowShow() {
	u.windows[WINDOW_RECIPES].Show()
}

func (u *UI) ShowRecipeWindow(recipes []ss.Recipe, onSelect func(ss.Recipe)) {
	u.windows[WINDOW_RECIPES].(*WindowRecipes).ShowSelector(recipes, onSelect)
}

func (u *UI) Draw(r *renderer.GameRenderer) {
	if !u.show {
		return
	}

	for _, w := range u.windows {
		w.Draw(r)
	}
}

func (u *UI) HandleMouseMovement(mp utils.ScreenCoord) {
	if !u.show {
		return
	}

	for _, w := range u.windows {
		w.HandleMouseMovement(mp)
	}
}

func (u *UI) HandleMouseAction(mbe input.MouseButtonEvent) bool {
	if !u.show {
		return false
	}

	for _, w := range u.windows {
		if w.HandleMouseAction(mbe) {
			return true
		}
	}
	return false
}
