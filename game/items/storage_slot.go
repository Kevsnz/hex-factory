package items

type StorageSlot struct {
	Item *ItemStack
}

type Storage []*StorageSlot

func NewStorage(length int) Storage {
	storage := make(Storage, length)
	for i := 0; i < length; i++ {
		storage[i] = &StorageSlot{Item: nil}
	}
	return storage
}

func (s Storage) TakeItemStackAnywhere(i *ItemStack) {
	for _, slot := range s {
		if slot.Item == nil || slot.Item.ItemType != i.ItemType {
			continue
		}
		remainder := slot.Item.TakeWithRemainder(i.Count)
		i.Count = remainder
		if i.Count == 0 {
			return
		}
	}

	for _, slot := range s {
		if slot.Item != nil {
			continue
		}
		newItemStack := NewItemStack(i.ItemType, i.Count)
		slot.Item = &newItemStack
		i.Count = 0
		return
	}
}
