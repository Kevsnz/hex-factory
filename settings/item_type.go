package settings

import "hextopdown/settings/strings"

type ItemType uint32

const (
	ITEM_TYPE_IRON_ORE   ItemType = iota
	ITEM_TYPE_IRON_PLATE ItemType = iota
	ITEM_TYPE_IRON_GEAR  ItemType = iota
	ITEM_TYPE_COUNT               = iota
)

var StackMaxSizes = [ITEM_TYPE_COUNT]int{
	ITEM_TYPE_IRON_ORE:   24,
	ITEM_TYPE_IRON_PLATE: 24,
	ITEM_TYPE_IRON_GEAR:  24,
}

var TypeStrings = [ITEM_TYPE_COUNT]strings.StringID{
	ITEM_TYPE_IRON_ORE:   strings.STRING_ITEM_IRON_ORE,
	ITEM_TYPE_IRON_PLATE: strings.STRING_ITEM_IRON_PLATE,
	ITEM_TYPE_IRON_GEAR:  strings.STRING_ITEM_IRON_GEAR,
}
