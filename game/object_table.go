package game

import (
	ss "hextopdown/settings"
	"hextopdown/settings/strings"
	"hextopdown/utils"
)

type ObjectParameters struct {
	Name       strings.StringID
	ObjectType ss.ObjectBaseType
	Shape      utils.Shape
}

var objectParamsList = [ss.OBJECT_TYPE_COUNT]ObjectParameters{
	ss.OBJECT_TYPE_BELT1: {
		Name:       strings.STRING_OBJECT_BELT,
		ObjectType: ss.STRUCTURE_BASETYPE_BELTLIKE,
		Shape:      utils.SHAPE_SINGLE,
	},
	ss.OBJECT_TYPE_BELTSPLITTER1: {
		Name:       strings.STRING_OBJECT_BELT_SPLITTER,
		ObjectType: ss.STRUCTURE_BASETYPE_BELTLIKE,
		Shape:      utils.SHAPE_SINGLE,
	},
	ss.OBJECT_TYPE_BELTUNDER1: {
		Name:       strings.STRING_OBJECT_BELT_UNDER,
		ObjectType: ss.STRUCTURE_BASETYPE_BELTLIKE,
		Shape:      utils.SHAPE_SINGLE,
	},
	ss.OBJECT_TYPE_CHESTBOX_SMALL: {
		Name:       strings.STRING_OBJECT_CHESTBOX_SMALL,
		ObjectType: ss.STRUCTURE_BASETYPE_STORAGE,
		Shape:      utils.SHAPE_SINGLE,
	},
	ss.OBJECT_TYPE_CHESTBOX_MEDIUM: {
		Name:       strings.STRING_OBJECT_CHESTBOX_MEDIUM,
		ObjectType: ss.STRUCTURE_BASETYPE_STORAGE,
		Shape:      utils.SHAPE_SINGLE,
	},
	ss.OBJECT_TYPE_CHESTBOX_LARGE: {
		Name:       strings.STRING_OBJECT_CHESTBOX_LARGE,
		ObjectType: ss.STRUCTURE_BASETYPE_STORAGE,
		Shape:      utils.SHAPE_SINGLE,
	},
	ss.OBJECT_TYPE_INSERTER1: {
		Name:       strings.STRING_OBJECT_INSERTER,
		ObjectType: ss.STRUCTURE_BASETYPE_INSERTER,
		Shape:      utils.SHAPE_SINGLE,
	},
	ss.OBJECT_TYPE_FURNACE_STONE: {
		Name:       strings.STRING_OBJECT_FURNACE,
		ObjectType: ss.STRUCTURE_BASETYPE_CONVERTER,
		Shape:      utils.SHAPE_DIAMOND,
	},
}
