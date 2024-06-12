package gamedata

import (
	ss "hextopdown/settings"
)

var ConverterParamsList = map[ss.ObjectType]*ConverterParameters{
	ss.OBJECT_TYPE_FURNACE_STONE: {
		BuildPower: 10,
		AutoRecipe: true,
	},
	ss.OBJECT_TYPE_ASSEMBLER_BASIC: {
		BuildPower: 8,
		AutoRecipe: false,
	},
}
