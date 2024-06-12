package gamedata

import ss "hextopdown/settings"

var StorageParamsList = map[ss.ObjectType]*StorageParameters{
	ss.OBJECT_TYPE_CHESTBOX_SMALL:  {Capacity: 8},
	ss.OBJECT_TYPE_CHESTBOX_MEDIUM: {Capacity: 16},
	ss.OBJECT_TYPE_CHESTBOX_LARGE:  {Capacity: 24},
}
