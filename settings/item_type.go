package settings

import "hextopdown/settings/strings"

type ItemType uint32

const (
	ITEM_TYPE_IRON_PLATE ItemType = iota
	ITEM_TYPE_IRON_GEAR  ItemType = iota
	ITEM_TYPE_COUNT               = iota
)

var StackMaxSizes = [ITEM_TYPE_COUNT]int{
	ITEM_TYPE_IRON_PLATE: 64,
	ITEM_TYPE_IRON_GEAR:  64,
}

var TypeStrings = [ITEM_TYPE_COUNT]strings.StringID{
	ITEM_TYPE_IRON_PLATE: strings.STRING_ITEM_IRON_PLATE,
}
