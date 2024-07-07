package ui

import (
	"hextopdown/game/gamedata"
	ss "hextopdown/settings"
	"hextopdown/settings/strings"
	"hextopdown/utils"
)

const WND_RECIPES_W_PCT = ss.FONT_SIZE_PCT * 10
const WND_RECIPES_H_PCT = ss.FONT_SIZE_PCT * 15
const WND_RECIPES_X_PCT = 0.5 - WND_RECIPES_W_PCT/2
const WND_RECIPES_Y_PCT = 0.5 - WND_RECIPES_H_PCT/2
const WND_RECIPES_BUTTONS_OFFSET_X_PCT = ss.FONT_SIZE_PCT / 2
const WND_RECIPES_BUTTONS_OFFSET_Y_PCT = ss.FONT_SIZE_PCT / 2

type WindowRecipes struct {
	Window
	recipeList []ss.Recipe
	onSelect   func(ss.Recipe)
}

func NewWindowRecipes() *WindowRecipes {
	wnd := &WindowRecipes{
		Window: Window{
			pos:     utils.ScreenCoord{X: WND_RECIPES_X_PCT, Y: WND_RECIPES_Y_PCT}.PctPosToScreen(),
			size:    utils.ScreenCoord{X: WND_RECIPES_W_PCT, Y: WND_RECIPES_H_PCT}.PctScaleToScreen(),
			title:   strings.STRING_RECIPE,
			visible: false,
			dialog:  true,
		},
	}
	WithCloseBox(wnd)

	return wnd
}

func (w *WindowRecipes) removeRecipeButtons() {
	newChildren := make([]iControl, 0, 1)
	for _, c := range w.children {
		if _, ok := c.(*ButtonRecipe); !ok {
			newChildren = append(newChildren, c)
		}
	}
	w.children = newChildren
}

func (w *WindowRecipes) ShowSelector(recipeList []ss.Recipe, onSelect func(ss.Recipe)) {
	w.removeRecipeButtons()

	w.recipeList = make([]ss.Recipe, 0, len(recipeList))
	w.recipeList = append(w.recipeList, recipeList...)

	for i, r := range w.recipeList {
		b := NewButtonRecipe(
			utils.ScreenCoord{
				X: WND_RECIPES_BUTTONS_OFFSET_X_PCT,
				Y: WND_RECIPES_BUTTONS_OFFSET_Y_PCT + float32(i)*ss.FONT_SIZE_PCT,
			}.PctScaleToScreen(),
			gamedata.RecipeList[r].Products[0].Type,
			func() {
				w.SelectRecipe(r)
			})
		w.AddChild(b, CONTROL_ALIGN_TOPLEFT)
	}

	w.onSelect = onSelect
	w.visible = true
}

func (w *WindowRecipes) SelectRecipe(recipe ss.Recipe) {
	w.onSelect(recipe)
	w.Hide()
}
