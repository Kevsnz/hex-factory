package objects

import (
	gd "hextopdown/game/gamedata"
	"hextopdown/game/items"
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
)

type Storage struct {
	Object
	params *gd.StorageParameters
	slots  []*items.StorageSlot
}

func NewChestBox(
	objType ss.ObjectType,
	pos utils.HexCoord,
	objParams *gd.ObjectParameters,
	params *gd.StorageParameters,
) *Storage {
	slots := make([]*items.StorageSlot, params.Capacity)
	for i := range slots {
		slots[i] = &items.StorageSlot{}
	}
	return &Storage{
		Object: Object{
			objType:   objType,
			pos:       pos,
			objParams: objParams,
		},
		params: params,
		slots:  slots,
	}
}

func (s *Storage) DrawGroundLevel(r *renderer.GameRenderer) {
	r.DrawObjectGround(s.pos.CenterToWorld(), s.objType, s.objParams.Shape, utils.DIR_LEFT)
}
func (s *Storage) DrawOnGroundLevel(r *renderer.GameRenderer) {}

func (s *Storage) TakeItemOut(pos utils.WorldCoord, allowedItems []ss.ItemType) (*items.ItemInWorld, bool) {
	for _, slot := range s.slots {
		item := slot.Item
		if item == nil {
			continue
		}
		if allowedItems != nil && !utils.ItemInList(item.ItemType, allowedItems) {
			continue
		}

		if !item.TakeOne() {
			panic("item counts are messed up") // TODO Change to continue????????????
		}
		newItem := items.NewItemInWorld2(item.ItemType, s.pos.CenterToWorld())
		if item.Count == 0 {
			slot.Item = nil
		}
		return &newItem, true
	}
	return nil, false
}

func (s *Storage) GetAcceptableItems() []ss.ItemType {
	for _, slot := range s.slots {
		if slot.Item == nil {
			return nil
		}
	}

	info := []ss.ItemType{}
	for _, slot := range s.slots {
		if slot.Item.Count < ss.StackMaxSizes[slot.Item.ItemType] {
			info = append(info, slot.Item.ItemType)
		}
	}
	return info
}

func (s *Storage) TakeItemIn(pos utils.WorldCoord, item items.ItemInWorld) (ok bool) {
	var emptySlot *items.StorageSlot = nil
	for _, slot := range s.slots {
		if slot.Item == nil {
			if emptySlot == nil {
				emptySlot = slot
			}
			continue
		}
		if slot.Item.ItemType != item.ItemType {
			continue
		}
		if slot.Item.AddOne() {
			return true
		}
	}
	if emptySlot == nil {
		return false
	}

	stack := items.NewSingleItemStack(item.ItemType)
	emptySlot.Item = &stack
	return true
}

func (s *Storage) GetItemList() []utils.ItemInfo {
	var info []utils.ItemInfo
outer:
	for _, slot := range s.slots {
		item := slot.Item
		if item == nil {
			continue
		}

		for j, inf := range info {
			if inf.Type == item.ItemType {
				info[j].Count += item.Count
				continue outer
			}
		}
		info = append(info, utils.ItemInfo{
			Type:  item.ItemType,
			Count: item.Count,
		})
	}
	return info
}

func (s *Storage) GetStorage() []*items.StorageSlot {
	return s.slots
}
