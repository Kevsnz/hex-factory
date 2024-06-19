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
	ss.RECIPE_IRON_PLATE: {
		Ingredients: []utils.ItemInfo{{Type: ss.ITEM_TYPE_IRON_ORE, Count: 1}},
		Products:    []utils.ItemInfo{{Type: ss.ITEM_TYPE_IRON_PLATE, Count: 1}},
		BuildPoints: 200,
		Energy:      0,
		Converters:  []ss.ObjectType{ss.OBJECT_TYPE_FURNACE_STONE},
	},
}

func GetAvailableRecipes(objType ss.ObjectType) []ss.Recipe {
	list := []ss.Recipe{}
	for id, recipe := range RecipeList {
		for _, converter := range recipe.Converters {
			if converter == objType {
				list = append(list, ss.Recipe(id))
				break
			}
		}
	}
	return list
}
