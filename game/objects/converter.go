package objects

import (
	gd "hextopdown/game/gamedata"
	"hextopdown/game/items"
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
)

const CONVERTER_OVERSTOCK_FACTOR = 2

type Converter struct {
	Object
	dir                utils.Dir
	params             *gd.ConverterParameters
	recipe             *gd.Recipe
	inputCounts        []int
	outputCounts       []int
	conversionProgress uint32
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
	r.DrawObjectOnGround(c.pos.CenterToWorld(), c.objType, c.objParams.Shape, c.dir)

	if c.recipe != nil {
		p := c.pos.CenterToWorld().Add(c.objParams.Shape.GetCenterOffset(c.dir))
		r.DrawDecal(p, 1.25, renderer.DECAL_BLACK_SPOT_FUZZY)
		r.DrawItemIconWorld(p, 0.8, c.recipe.Products[0].Type)
		r.DrawProgressBar(p, 1.25, c.conversionProgress, c.recipe.BuildPoints)
	}
}

func (c *Converter) RecipeChangeable() bool {
	return !c.params.AutoRecipe
}

func (c *Converter) ChangeRecipe(recipe ss.Recipe) {
	if c.params.AutoRecipe {
		panic("cannot change recipe for auto recipe converter")
	}

	if c.recipe == nil {
		c.recipe = &gd.RecipeList[recipe]
		c.inputCounts = make([]int, len(c.recipe.Ingredients))
		c.outputCounts = make([]int, len(c.recipe.Products))
		c.conversionProgress = 0
		return
	}

	newRecipe := &gd.RecipeList[recipe]
	newInputs := make([]int, len(c.recipe.Ingredients))
	for i, cnt := range c.inputCounts {
		if cnt == 0 {
			continue
		}

		for j, ing := range newRecipe.Ingredients {
			if ing.Type == c.recipe.Ingredients[i].Type {
				newInputs[j] += cnt
				c.inputCounts[i] = 0
				break
			}
		}

		// TODO Drop left over items (c.inputCounts > 0)
	}
	// TODO Drop left over items (c.outputCounts > 0)

	c.recipe = newRecipe
	c.inputCounts = newInputs
	c.outputCounts = make([]int, len(c.recipe.Products))
	c.conversionProgress = 0
}

func (c *Converter) GetAcceptableItems() []ss.ItemType {
	if c.params.AutoRecipe {
		items := []ss.ItemType{}

		if c.recipe != nil {
			for i, ing := range c.recipe.Ingredients {
				if c.inputCounts[i] < ing.Count*CONVERTER_OVERSTOCK_FACTOR {
					items = append(items, ing.Type)
				}
			}
			return items
		}

		list := gd.GetAvailableRecipes(c.objType)
		for _, id := range list {
			for _, ingr := range gd.RecipeList[id].Ingredients {
				items = append(items, ingr.Type)
			}
		}
		return items
	}

	if c.recipe == nil {
		return []ss.ItemType{}
	}

	items := []ss.ItemType{}
	for i, ing := range c.recipe.Ingredients {
		if c.inputCounts[i] < ing.Count*CONVERTER_OVERSTOCK_FACTOR {
			items = append(items, ing.Type)
		}
	}
	return items
}

func (c *Converter) TakeItemIn(pos utils.WorldCoord, item items.ItemInWorld) (ok bool) {
	if c.recipe == nil {
		if !c.params.AutoRecipe || !c.setFittingRecipe(item.ItemType) {
			return false
		}
	}

	for i, ing := range c.recipe.Ingredients {
		if item.ItemType != ing.Type {
			continue
		}
		c.inputCounts[i] += 1 // Change when item becomes item stack!!!!!!!!!!!!!!!!
		return true
	}

	return false
}

func (c *Converter) GetItemList() []utils.ItemInfo {
	if c.recipe == nil {
		return []utils.ItemInfo{}
	}

	items := []utils.ItemInfo{}
	for i, ingr := range c.recipe.Ingredients {
		if c.inputCounts[i] == 0 {
			continue
		}
		items = append(items, utils.ItemInfo{Type: ingr.Type, Count: c.inputCounts[i]})
	}
	for i, prod := range c.recipe.Products {
		if c.outputCounts[i] == 0 {
			continue
		}
		items = append(items, utils.ItemInfo{Type: prod.Type, Count: c.outputCounts[i]})
	}

	return items
}

func (c *Converter) Update(ticks uint64, world HexGridWorldInteractor) {
	if c.recipe == nil {
		return
	}

	isMaxxed := false
	for i, prod := range c.recipe.Products {
		if c.outputCounts[i]+prod.Count > ss.StackMaxSizes[prod.Type] {
			isMaxxed = true
			break
		}
	}

	if isMaxxed {
		return
	}

	enough := true
	for i, ingr := range c.recipe.Ingredients {
		if c.inputCounts[i] < ingr.Count {
			enough = false
			break
		}
	}

	if !enough {
		c.conversionProgress = 0
		return
	}

	c.conversionProgress += c.params.BuildPower

	if c.conversionProgress >= c.recipe.BuildPoints {
		c.conversionProgress -= c.recipe.BuildPoints

		for i, ingr := range c.recipe.Ingredients {
			c.inputCounts[i] -= ingr.Count
		}
		for i, prod := range c.recipe.Products {
			c.outputCounts[i] += prod.Count
		}
	}
}

func (c *Converter) TakeItemOut(pos utils.WorldCoord, allowedItems []ss.ItemType) (item *items.ItemInWorld, ok bool) {
	if c.recipe == nil {
		return nil, false
	}

	for i, prod := range c.recipe.Products {
		if c.outputCounts[i] == 0 {
			continue
		}
		item := items.NewItemInWorld2(prod.Type, pos)
		c.outputCounts[i] -= 1

		if c.params.AutoRecipe {
			c.checkResetRecipe()
		}
		return &item, true
	}
	return nil, false
}

func (c *Converter) setFittingRecipe(itemType ss.ItemType) bool {
	for _, id := range gd.GetAvailableRecipes(c.objType) {
		for _, ing := range gd.RecipeList[id].Ingredients {
			if ing.Type == itemType {
				c.recipe = &gd.RecipeList[id]
				c.inputCounts = make([]int, len(c.recipe.Ingredients))
				c.outputCounts = make([]int, len(c.recipe.Products))
				c.conversionProgress = 0
				return true
			}
		}
	}

	return false
}

func (c *Converter) checkResetRecipe() {
	if c.recipe == nil {
		return
	}
	if c.conversionProgress != 0 {
		return
	}

	for _, c := range c.inputCounts {
		if c > 0 {
			return
		}
	}

	for _, c := range c.outputCounts {
		if c > 0 {
			return
		}
	}

	c.recipe = nil
	c.conversionProgress = 0
	c.inputCounts = nil
	c.outputCounts = nil
}
