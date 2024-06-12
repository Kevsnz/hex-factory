package objects

import (
	"hextopdown/game/items"
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
)

type Storage struct {
	Pos      utils.HexCoord
	Capacity int
	slots    []*items.ItemStack
}

func NewChestBox(pos utils.HexCoord, capacity int) *Storage {
	return &Storage{
		Pos:      pos,
		Capacity: capacity,
		slots:    make([]*items.ItemStack, capacity),
	}
}

func (s *Storage) GetObjectType() ss.ObjectType {
	switch s.Capacity {
	case ss.CHESTBOX_CAPACITY_SMALL:
		return ss.OBJECT_TYPE_CHESTBOX_SMALL
	case ss.CHESTBOX_CAPACITY_MEDIUM:
		return ss.OBJECT_TYPE_CHESTBOX_MEDIUM
	case ss.CHESTBOX_CAPACITY_LARGE:
		return ss.OBJECT_TYPE_CHESTBOX_LARGE
	}
	panic("invalid chestbox capacity")
}

func (s *Storage) GetPos() utils.HexCoord {
	return s.Pos
}

func (s *Storage) DrawGroundLevel(r *renderer.GameRenderer) {
	switch s.Capacity {
	case ss.CHESTBOX_CAPACITY_SMALL:
		r.DrawObjectGround(s.Pos.CenterToWorld(), ss.OBJECT_TYPE_CHESTBOX_SMALL, utils.SHAPE_SINGLE, utils.DIR_LEFT)
	case ss.CHESTBOX_CAPACITY_MEDIUM:
		r.DrawObjectGround(s.Pos.CenterToWorld(), ss.OBJECT_TYPE_CHESTBOX_MEDIUM, utils.SHAPE_SINGLE, utils.DIR_LEFT)
	case ss.CHESTBOX_CAPACITY_LARGE:
		r.DrawObjectGround(s.Pos.CenterToWorld(), ss.OBJECT_TYPE_CHESTBOX_LARGE, utils.SHAPE_SINGLE, utils.DIR_LEFT)
	}
}
func (s *Storage) DrawOnGroundLevel(r *renderer.GameRenderer) {}

func (s *Storage) TakeItemOut(pos utils.WorldCoord) (*items.ItemInWorld, bool) {
	for i, stack := range s.slots {
		if stack == nil {
			continue
		}
		item := items.NewItemInWorld2(stack.ItemType, s.Pos.CenterToWorld())
		stack.Count--
		if stack.Count == 0 {
			s.slots[i] = nil
		}
		return &item, true
	}
	return nil, false
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
