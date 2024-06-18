package gamedata

import (
	ss "hextopdown/settings"
	"hextopdown/utils"
)

type Recipe struct {
	Ingredients []utils.ItemInfo
	Products    []utils.ItemInfo
	BuildPoints uint32
	Energy      uint32
	Converters  []ss.ObjectType
}

var RecipeList = [ss.RECIPE_COUNT]Recipe{
	ss.RECIPE_IRON_GEAR: {
		Ingredients: []utils.ItemInfo{{Type: ss.ITEM_TYPE_IRON_PLATE, Count: 1}},
		Products:    []utils.ItemInfo{{Type: ss.ITEM_TYPE_IRON_GEAR, Count: 1}},
		BuildPoints: 200,
		Energy:      0,
		Converters:  []ss.ObjectType{ss.OBJECT_TYPE_ASSEMBLER_BASIC},
	},
}
