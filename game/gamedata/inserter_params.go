package gamedata

import (
	ss "hextopdown/settings"
)

var InserterParamsList = map[ss.ObjectType]*InserterParameters{
	ss.OBJECT_TYPE_INSERTER1: {
		SwingSpeed: 20,
		Reach:      1,
		StackSize:  1,
		Filtering:  false,
	},
}
