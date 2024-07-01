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
	slots  []*items.ItemStack
}

func NewChestBox(
	objType ss.ObjectType,
	pos utils.HexCoord,
	objParams *gd.ObjectParameters,
	params *gd.StorageParameters,
) *Storage {
	return &Storage{
		Object: Object{
			objType:   objType,
			pos:       pos,
			objParams: objParams,
		},
		params: params,
		slots:  make([]*items.ItemStack, params.Capacity),
	}
}

func (s *Storage) DrawGroundLevel(r *renderer.GameRenderer) {
	r.DrawObjectGround(s.pos.CenterToWorld(), s.objType, s.objParams.Shape, utils.DIR_LEFT)
}
func (s *Storage) DrawOnGroundLevel(r *renderer.GameRenderer) {}

func (s *Storage) TakeItemOut(pos utils.WorldCoord, allowedItems []ss.ItemType) (*items.ItemInWorld, bool) {
	for i, stack := range s.slots {
		if stack == nil {
			continue
		}
		if allowedItems != nil && !utils.ItemInList(stack.ItemType, allowedItems) {
			continue
		}

		item := items.NewItemInWorld2(stack.ItemType, s.pos.CenterToWorld())
		stack.Count--
		if stack.Count == 0 {
			s.slots[i] = nil
		}
		return &item, true
	}
	return nil, false
}

func (s *Storage) GetAcceptableItems() []ss.ItemType {
	for _, stack := range s.slots {
		if stack == nil {
			return nil
		}
	}

	info := []ss.ItemType{}
	for _, stack := range s.slots {
		if stack.Count < ss.StackMaxSizes[stack.ItemType] {
			info = append(info, stack.ItemType)
		}
	}
	return info
}

func (s *Storage) TakeItemIn(pos utils.WorldCoord, item items.ItemInWorld) (ok bool) {
	emptyIdx := -1
	for i, stack := range s.slots {
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
	s.slots[emptyIdx] = &stack
	return true
}

func (s *Storage) GetItemList() []utils.ItemInfo {
	var info []utils.ItemInfo
outer:
	for _, stack := range s.slots {
		if stack == nil {
			continue
		}

		for j, inf := range info {
			if inf.Type == stack.ItemType {
				info[j].Count += stack.Count
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

func (s *Storage) GetStorage() []*items.ItemStack {
	return s.slots
}
