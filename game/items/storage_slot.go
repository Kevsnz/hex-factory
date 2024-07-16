package items

import "hextopdown/settings"

type StorageSlot struct {
	Item          *ItemStack
	Active        bool
	FixedItemType settings.ItemType
}

type Storage []*StorageSlot

func NewStorage(length int, active bool) Storage {
	storage := make(Storage, length)
	for i := 0; i < length; i++ {
		storage[i] = &StorageSlot{Item: nil, Active: active, FixedItemType: settings.ITEM_TYPE_COUNT}
	}
	return storage
}

func (s Storage) TakeItemStackAnywhere(i *ItemStack, maxCount int) int {
	count := maxCount
	for _, slot := range s {
		if !slot.Active || slot.Item == nil || slot.Item.ItemType != i.ItemType {
			continue
		}
		count = slot.Item.TakeWithRemainder(count)
		if count == 0 {
			return maxCount
		}
	}

	for _, slot := range s {
		if !slot.Active || slot.Item != nil {
			continue
		}
		if slot.FixedItemType != settings.ITEM_TYPE_COUNT && slot.FixedItemType != i.ItemType {
			continue
		}
		newItemStack := NewItemStack(i.ItemType, count)
		slot.Item = &newItemStack
		return maxCount
	}

	return maxCount - count
}
