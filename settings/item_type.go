package settings

type ItemType uint32

const (
	ITEM_TYPE_IRON_PLATE ItemType = iota
	ITEM_TYPE_COUNT               = iota
)

var StackMaxSizes = [ITEM_TYPE_COUNT]int{
	ITEM_TYPE_IRON_PLATE: 4,
}
