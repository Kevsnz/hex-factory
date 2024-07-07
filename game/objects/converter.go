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
	dir     utils.Dir
	params  *gd.ConverterParameters
	recipe  *gd.Recipe
	inputs  []*items.StorageSlot
	outputs []*items.StorageSlot
	// inputCounts        []int
	// outputCounts       []int
	conversionProgress uint32
}

func NewConverter(
	objType ss.ObjectType,
	pos utils.HexCoord,
	dir utils.Dir,
	objParams *gd.ObjectParameters,
	params *gd.ConverterParameters,
) *Converter {
	maxInputs := 0
	maxOutputs := 0
	for _, id := range gd.GetAvailableRecipes(objType) {
		recipe := gd.RecipeList[id]
		maxInputs = max(maxInputs, len(recipe.Ingredients))
		maxOutputs = max(maxOutputs, len(recipe.Products))
	}

	return &Converter{
		Object: Object{
			objType:   objType,
			pos:       pos,
			objParams: objParams,
		},
		dir:     dir,
		params:  params,
		inputs:  items.NewStorage(maxInputs, false),
		outputs: items.NewStorage(maxOutputs, false),
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

func (c *Converter) HasRecipe() bool {
	return c.recipe != nil
}

func (c *Converter) ChangeRecipe(recipe ss.Recipe) {
	if c.params.AutoRecipe {
		panic("cannot change recipe for auto recipe converter")
	}

	if c.recipe == nil {
		c.recipe = &gd.RecipeList[recipe]
		for i, slot := range c.inputs {
			slot.Active = i < len(c.recipe.Ingredients)
		}
		for i, slot := range c.outputs {
			slot.Active = i < len(c.recipe.Products)
		}
		c.conversionProgress = 0
		return
	}

	c.recipe = &gd.RecipeList[recipe]
	for i, slot := range c.inputs {
		// TODO Drop left over items (slot.Item.Count > 0)
		slot.Item = nil
		slot.Active = i < len(c.recipe.Ingredients)
	}
	for i, slot := range c.outputs {
		// TODO Drop left over items (slot.Item.Count > 0)
		slot.Item = nil
		slot.Active = i < len(c.recipe.Products)
	}
	c.conversionProgress = 0
}

func (c *Converter) GetAcceptableItems() []ss.ItemType {
	if c.recipe != nil {
		items := []ss.ItemType{}
		for i, ing := range c.recipe.Ingredients {
			if c.inputs[i].Item == nil || c.inputs[i].Item.Count < ing.Count*CONVERTER_OVERSTOCK_FACTOR {
				items = append(items, ing.Type)
			}
		}
		return items
	}

	if !c.params.AutoRecipe {
		return []ss.ItemType{}
	}

	items := []ss.ItemType{}
	list := gd.GetAvailableRecipes(c.objType)
	for _, id := range list {
		for _, ingr := range gd.RecipeList[id].Ingredients {
			items = append(items, ingr.Type)
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

		if c.inputs[i].Item != nil && c.inputs[i].Item.ItemType != item.ItemType && c.inputs[i].Item.Count > 0 {
			panic("converter inputs are messed up")
		}
		if c.inputs[i].Item == nil {
			c.inputs[i].Item = &items.ItemStack{ItemType: item.ItemType, Count: 0}
		}
		c.inputs[i].Item.ItemType = item.ItemType
		c.inputs[i].Item.Count += 1
		return true
	}

	return false
}

func (c *Converter) GetItemList() []utils.ItemInfo {
	if c.recipe == nil {
		return []utils.ItemInfo{}
	}

	items := []utils.ItemInfo{}
	for _, slot := range c.inputs {
		if slot.Item != nil {
			items = append(items, utils.ItemInfo{Type: slot.Item.ItemType, Count: slot.Item.Count})
		}
	}
	for _, slot := range c.outputs {
		if slot.Item != nil {
			items = append(items, utils.ItemInfo{Type: slot.Item.ItemType, Count: slot.Item.Count})
		}
	}

	return items
}

func (c *Converter) Update(ticks uint64, world HexGridWorldInteractor) {
	if c.recipe == nil {
		return
	}

	isMaxxed := false
	for i, prod := range c.recipe.Products {
		if c.outputs[i].Item != nil && c.outputs[i].Item.ItemType != prod.Type && c.outputs[i].Item.Count > 0 {
			panic("converter outputs are messed up")
		}
		if c.outputs[i].Item != nil && c.outputs[i].Item.Count+prod.Count > ss.StackMaxSizes[prod.Type] {
			isMaxxed = true
			break
		}
	}

	if isMaxxed {
		return
	}

	enough := true
	for i, ingr := range c.recipe.Ingredients {
		if c.inputs[i].Item == nil || c.inputs[i].Item.Count < ingr.Count {
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
			if c.inputs[i].Item == nil || c.inputs[i].Item.ItemType != ingr.Type {
				panic("converter inputs are messed up")
			}
			c.inputs[i].Item.Count -= ingr.Count
			if c.inputs[i].Item.Count == 0 {
				c.inputs[i].Item = nil
			}
		}

		for i, prod := range c.recipe.Products {
			if c.outputs[i].Item != nil && c.outputs[i].Item.ItemType != prod.Type && c.outputs[i].Item.Count > 0 {
				panic("converter outputs are messed up")
			}
			if c.outputs[i].Item == nil {
				newItem := items.NewItemStack(prod.Type, 0)
				c.outputs[i].Item = &newItem
			}

			c.outputs[i].Item.Count += prod.Count
		}
	}
}

func (c *Converter) TakeItemOut(pos utils.WorldCoord, allowedItems []ss.ItemType) (item *items.ItemInWorld, ok bool) {
	if c.recipe == nil {
		return nil, false
	}

	for i, prod := range c.recipe.Products {
		if c.outputs[i].Item == nil {
			continue
		}
		if c.outputs[i].Item.ItemType != prod.Type && c.outputs[i].Item.Count > 0 {
			panic("converter outputs are messed up")
		}
		if !utils.ItemInList(prod.Type, allowedItems) {
			continue
		}

		item := items.NewItemInWorld2(prod.Type, pos)
		c.outputs[i].Item.Count -= 1
		if c.outputs[i].Item.Count == 0 {
			c.outputs[i].Item = nil
		}

		if c.params.AutoRecipe {
			c.checkResetRecipe()
		}
		return &item, true
	}
	return nil, false
}

func (c *Converter) setFittingRecipe(itemType ss.ItemType) bool {
	for _, id := range gd.GetAvailableRecipes(c.objType) {
		for _, ingr := range gd.RecipeList[id].Ingredients {
			if ingr.Type == itemType {
				c.recipe = &gd.RecipeList[id]
				for i, slot := range c.inputs {
					slot.Item = nil
					slot.Active = i < len(c.recipe.Ingredients)
				}
				for i, slot := range c.outputs {
					slot.Item = nil
					slot.Active = i < len(c.recipe.Products)
				}
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

	for _, slot := range c.inputs {
		if slot.Item != nil && slot.Item.Count > 0 {
			return
		}
	}

	for _, slot := range c.outputs {
		if slot.Item != nil && slot.Item.Count > 0 {
			return
		}
	}

	c.recipe = nil
	c.conversionProgress = 0
	for _, slot := range c.inputs {
		slot.Item = nil
		slot.Active = false
	}
	for _, slot := range c.outputs {
		slot.Item = nil
		slot.Active = false
	}
}

func (c *Converter) GetInputSlots() []*items.StorageSlot {
	return c.inputs
}

func (c *Converter) GetOutputSlots() []*items.StorageSlot {
	return c.outputs
}
