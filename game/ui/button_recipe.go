package ui

import (
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
)

const BTN_RECIPE_SIZE = ss.FONT_SIZE_PCT * 2

type ButtonRecipe struct {
	Button
	item ss.ItemType
}

func NewButtonRecipe(pos utils.ScreenCoord, item ss.ItemType, onClick func()) *ButtonRecipe {
	return &ButtonRecipe{
		Button: Button{
			ControlBase: ControlBase{
				Pos:  pos,
				Size: utils.ScreenCoord{X: BTN_RECIPE_SIZE, Y: BTN_RECIPE_SIZE}.PctScaleToScreen(),
			},
			onClick: onClick,
		},
		item: item,
	}
}

func (b *ButtonRecipe) Draw(r *renderer.GameRenderer, parentPos utils.ScreenCoord) {
	r.DrawButtonIcon(b.Pos.Add(parentPos), b.Size, b.item, b.down)
}
