package objects

import (
	"hextopdown/game/items"
	"hextopdown/renderer"
	"hextopdown/settings"
	"hextopdown/settings/strings"
	"hextopdown/utils"
)

type ChestBox struct {
	Pos      utils.HexCoord
	Capacity int
	slots    []*items.ItemStack
}

func NewChestBox(pos utils.HexCoord, capacity int) *ChestBox {
	return &ChestBox{
		Pos:      pos,
		Capacity: capacity,
		slots:    make([]*items.ItemStack, capacity),
	}
}

func (cb *ChestBox) GetNameString() strings.StringID {
	switch cb.Capacity {
	case settings.CHESTBOX_CAPACITY_SMALL:
		return strings.STRING_OBJECT_CHESTBOX_SMALL
	case settings.CHESTBOX_CAPACITY_MEDIUM:
		return strings.STRING_OBJECT_CHESTBOX_MEDIUM
	case settings.CHESTBOX_CAPACITY_LARGE:
		return strings.STRING_OBJECT_CHESTBOX_LARGE
	}
	return strings.STRING_OBJECT_UNKNOWN
}

func (cb *ChestBox) GetPos() utils.HexCoord {
	return cb.Pos
}

func (cb *ChestBox) DrawGroundLevel(r *renderer.GameRenderer) {
	switch cb.Capacity {
	case settings.CHESTBOX_CAPACITY_SMALL:
		r.DrawStructureGround(cb.Pos, settings.STRUCTURE_TYPE_CHESHBOX_SMALL)
	case settings.CHESTBOX_CAPACITY_MEDIUM:
		r.DrawStructureGround(cb.Pos, settings.STRUCTURE_TYPE_CHESHBOX_MEDIUM)
	case settings.CHESTBOX_CAPACITY_LARGE:
		r.DrawStructureGround(cb.Pos, settings.STRUCTURE_TYPE_CHESHBOX_LARGE)
	}
}
func (cb *ChestBox) DrawOnGroundLevel(r *renderer.GameRenderer) {}

func (cb *ChestBox) TakeItemOut(pos utils.WorldCoord) (*items.ItemInWorld, bool) {
	for i, stack := range cb.slots {
		if stack == nil {
			continue
		}
		item := items.NewItemInWorld2(stack.ItemType, cb.Pos.CenterToWorld())
		stack.Count--
		if stack.Count == 0 {
			cb.slots[i] = nil
		}
		return &item, true
	}
	return nil, false
}

func (cb *ChestBox) TakeItemIn(pos utils.WorldCoord, item items.ItemInWorld) (ok bool) {
	emptyIdx := -1
	for i, stack := range cb.slots {
		if stack == nil {
			if emptyIdx == -1 {
				emptyIdx = i
			}
			continue
		}
		if stack.ItemType != item.ItemType {
			continue
		}
		if stack.AddOne() {
			return true
		}
	}
	if emptyIdx == -1 {
		return false
	}

	stack := items.NewSingleItemStack(item.ItemType)
	cb.slots[emptyIdx] = &stack
	return true
}

func (cb *ChestBox) GetItemList() []utils.ItemInfo {
	var info []utils.ItemInfo
outer:
	for _, stack := range cb.slots {
		if stack == nil {
			continue
		}

		for j, inf := range info {
			if inf.Type == stack.ItemType {
				info[j].Count++
				continue outer
			}
		}
		info = append(info, utils.ItemInfo{
			Type:  stack.ItemType,
			Count: stack.Count,
		})
	}
	return info
}
