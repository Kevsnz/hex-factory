package ui

import "hextopdown/game/items"

type InteractableStorageConverter interface {
	GetInputSlots() items.Storage
	GetOutputSlots() items.Storage
	PutIn(*items.ItemStack, int) int
	TakeOut(*items.StorageSlot, int)
}
