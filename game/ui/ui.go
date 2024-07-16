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
	WINDOW_CONVERTER

	WINDOW_COUNT
)

type iWindow interface {
	HandleMouseMovement(utils.ScreenCoord)
	HandleMouseAction(input.MouseButtonEvent) bool
	HandleGameAction(ae input.ActionEvent) bool
	Draw(*renderer.GameRenderer)
	Show()
	Hide()
	IsVisible() bool
}

type UI struct {
	show                bool
	windows             [WINDOW_COUNT]iWindow
	scale               float32
	currentDialogWindow iWindow
}

func NewUI() *UI {
	wr := NewWindowRecipes()
	wi := NewWindowInventory()
	ws := NewWindowStorage()
	wc := NewWindowConverter()

	return &UI{
		windows: [WINDOW_COUNT]iWindow{
			WINDOW_RECIPES:   wr,
			WINDOW_INVENTORY: wi,
			WINDOW_STORAGE:   ws,
			WINDOW_CONVERTER: wc,
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

func (u *UI) ShowInventoryWindow(inventory items.Storage) {
	//lint:ignore S1040 // it's a nil check!!!
	if w, ok := u.currentDialogWindow.(iWindow); ok {
		if w.IsVisible() {
			return
		}
	}
	u.windows[WINDOW_INVENTORY].(*WindowInventory).ShowInventory(inventory)
	u.currentDialogWindow = u.windows[WINDOW_INVENTORY]
}

func (u *UI) ShowStorageWindow(
	objName strings.StringID,
	inventory items.Storage,
	storage items.Storage,
) {
	//lint:ignore S1040 // it's a nil check!!!
	if w, ok := u.currentDialogWindow.(iWindow); ok {
		if w.IsVisible() {
			return
		}
	}
	u.windows[WINDOW_STORAGE].(*WindowStorage).ShowStorage(objName, inventory, storage)
	u.currentDialogWindow = u.windows[WINDOW_STORAGE]
}

func (u *UI) ShowConverterWindow(
	objName strings.StringID,
	inventory items.Storage,
	converter InteractableStorageConverter,
) {
	//lint:ignore S1040 // it's a nil check!!!
	if w, ok := u.currentDialogWindow.(iWindow); ok {
		if w.IsVisible() {
			return
		}
	}
	u.windows[WINDOW_CONVERTER].(*WindowConverter).ShowConverter(objName, inventory, converter)
	u.currentDialogWindow = u.windows[WINDOW_CONVERTER]
}

func (u *UI) ShowRecipeWindow(recipes []ss.Recipe, onSelect func(ss.Recipe)) {
	//lint:ignore S1040 // it's a nil check!!!
	if w, ok := u.currentDialogWindow.(iWindow); ok {
		if w.IsVisible() {
			return
		}
	}
	u.windows[WINDOW_RECIPES].(*WindowRecipes).ShowSelector(recipes, onSelect)
	u.currentDialogWindow = u.windows[WINDOW_RECIPES]
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

func (u *UI) HandleGameAction(ae input.ActionEvent) bool {
	//lint:ignore S1040 // it's a nil check!!!
	if w, ok := u.currentDialogWindow.(iWindow); ok {
		if !w.IsVisible() {
			return false
		}

		return w.HandleGameAction(ae)
	}
	return false
}
