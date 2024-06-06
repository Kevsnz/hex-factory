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

func (i *ItemStack) Split() (ItemStack, bool) {
	if i.Count < 2 {
		return ItemStack{}, false
	}

	outCount := i.Count / 2
	i.Count -= outCount
	return NewItemStack(i.ItemType, outCount), true
}

func (i *ItemStack) TakeOne(i2 *ItemStack) (*ItemStack, bool) {
	if i.Count < 2 {
		return nil, false
	}
	i.Count--

	if i2 == nil {
		i22 := NewSingleItemStack(i.ItemType)
		return &i22, true
	}

	i2.Count++
	return i2, true
}
