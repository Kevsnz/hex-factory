package objects

import (
	gd "hextopdown/game/gamedata"
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
)

type Converter struct {
	Object
	dir    utils.Dir
	params *gd.ConverterParameters
	recipe ss.Recipe
}

func NewConverter(
	objType ss.ObjectType,
	pos utils.HexCoord,
	dir utils.Dir,
	objParams *gd.ObjectParameters,
	params *gd.ConverterParameters,
) *Converter {
	return &Converter{
		Object: Object{
			objType:   objType,
			pos:       pos,
			objParams: objParams,
		},
		dir:    dir,
		params: params,
		recipe: ss.RECIPE_COUNT,
	}
}

func (c *Converter) GetDir() utils.Dir {
	return c.dir
}

func (c *Converter) Rotate(_ bool) {}

func (c *Converter) DrawGroundLevel(r *renderer.GameRenderer) {
	r.DrawObjectGround(c.pos.CenterToWorld(), c.objType, c.objParams.Shape, c.dir)
}

func (c *Converter) DrawOnGroundLevel(r *renderer.GameRenderer) {
	if c.recipe != ss.RECIPE_COUNT {
		iconSize := utils.ScreenCoord{X: ss.FONT_SIZE_PCT * 2, Y: ss.FONT_SIZE_PCT * 2}.PctScaleToScreen()
		r.DrawItemIcon(
			c.pos.CenterToWorld().ToScreen().Sub(iconSize.Div(2)),
			iconSize.X,
			gd.RecipeList[c.recipe].Products[0].Type,
		)
	}
}

func (c *Converter) GetRecipe() ss.Recipe {
	return c.recipe
}

func (c *Converter) RecipeChangeable() bool {
	return !c.params.AutoRecipe
}

func (c *Converter) ChangeRecipe(recipe ss.Recipe) {
	if c.params.AutoRecipe {
		panic("cannot change recipe for auto recipe converter")
	}
	c.recipe = recipe
	// TODO Changing recipe logic!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
}
