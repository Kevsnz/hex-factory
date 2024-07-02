package ui

import (
	"hextopdown/game/items"
	"hextopdown/input"
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/settings/strings"
	"hextopdown/utils"
)

const (
	WINDOW_RECIPES = iota
	WINDOW_INVENTORY
	WINDOW_STORAGE

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
	wr := NewWindowRecipes()
	wi := NewWindowInventory()
	ws := NewWindowStorage()

	return &UI{
		windows: [WINDOW_COUNT]iWindow{
			WINDOW_RECIPES:   wr,
			WINDOW_INVENTORY: wi,
			WINDOW_STORAGE:   ws,
		},
		scale: 1,
		show:  true,
	}
}

func (u *UI) Destroy() {}

func (u *UI) Show(b bool) {
	u.show = b
}
func (u *UI) ShowToggle() {
	u.show = !u.show
}

func (u *UI) ShowInventoryWindow(inventory []*items.StorageSlot) {
	u.windows[WINDOW_INVENTORY].(*WindowInventory).ShowInventory(inventory)
}

func (u *UI) ShowStorageWindow(objName strings.StringID, inventory []*items.StorageSlot, storage []*items.StorageSlot) {
	u.windows[WINDOW_STORAGE].(*WindowStorage).ShowStorage(objName, inventory, storage)
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
