package items

import (
	"hextopdown/settings"
)

type ItemStack struct {
	ItemType settings.ItemType
	Count    int
}

func NewItemStack(itemType settings.ItemType, count int) ItemStack {
	return ItemStack{
		ItemType: itemType,
		Count:    count,
	}
}

func NewSingleItemStack(itemType settings.ItemType) ItemStack {
	return NewItemStack(itemType, 1)
}

func (i *ItemStack) AddOne() bool {
	if i.Count == settings.StackMaxSizes[i.ItemType] {
		return false
	}
	i.Count++
	return true
}

func (i *ItemStack) TakeOne() bool {
	if i.Count == 0 {
		return false
	}
	i.Count--
	return true
}
