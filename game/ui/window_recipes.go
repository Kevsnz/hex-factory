package ui

import (
	ss "hextopdown/settings"
	"hextopdown/settings/strings"
	"hextopdown/utils"
)

const WND_RECIPES_W_PCT = 0.4
const WND_RECIPES_H_PCT = 0.6
const WND_RECIPES_X_PCT = 0.5 - WND_RECIPES_W_PCT/2
const WND_RECIPES_Y_PCT = 0.5 - WND_RECIPES_H_PCT/2

type WindowRecipes struct {
	Window
	recipeList []ss.Recipe
	onSelect   func(ss.Recipe)
}

func NewWindowRecipes() *WindowRecipes {
	wnd := &WindowRecipes{
		Window: Window{
			Pos:   utils.ScreenCoord{X: WND_RECIPES_X_PCT, Y: WND_RECIPES_Y_PCT}.PctPosToScreen(),
			Size:  utils.ScreenCoord{X: WND_RECIPES_W_PCT, Y: WND_RECIPES_H_PCT}.PctScaleToScreen(),
			Title: strings.STRING_RECIPE,
		},
	}
	WithCloseBox(wnd)

	return wnd
}

func (w *WindowRecipes) ShowSelector(recipeList []ss.Recipe, onSelect func(ss.Recipe)) {
	w.recipeList = make([]ss.Recipe, 0, len(recipeList))
	w.recipeList = append(w.recipeList, recipeList...)
	w.onSelect = onSelect
	w.Visible = true
}
