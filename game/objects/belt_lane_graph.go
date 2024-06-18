package objects

import (
	"hextopdown/game/items"
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
	"math"
)

type BeltGraphSegment struct {
	PosZero            utils.WorldCoord
	PosHalfOffset      utils.WorldCoord
	Items              utils.RingBuffer[*items.ItemOnBelt]
	Start              float64
	End                float64 // from 0.0 (from hexagon's side or center) (NOT from start!!!)
	speed              float64
	NextSegment        *BeltGraphSegment
	nextSegmentOffset  float64
	PrevSegment1       *BeltGraphSegment // connected to offset 0.0
	PrevSegment2       *BeltGraphSegment // joined to prevSegment2Offset
	prevSegment2Offset float64
	Shift              utils.WorldCoord
	ShiftOffset        float64
}

func NewBeltGraphSegment(localZeroPos, localHalfOffsetPos utils.WorldCoord, speed float64, capacity int) BeltGraphSegment {
	return BeltGraphSegment{
		Items:         utils.NewRingBuffer[*items.ItemOnBelt](capacity),
		PosZero:       localZeroPos,
		PosHalfOffset: localHalfOffsetPos,
		Start:         0,
		End:           0.5,
		speed:         speed,
	}
}

func (bgs *BeltGraphSegment) Connect(bgs2 *BeltGraphSegment, start, end, nextOffset float64) {
	bgs.Change(start, end)
	bgs.NextSegment = bgs2
	bgs.nextSegmentOffset = nextOffset
	if nextOffset == 0 {
		if bgs2.PrevSegment1 != nil && bgs2.PrevSegment1 != bgs {
			panic("belt connectivity inconsistency")
		}
		bgs2.PrevSegment1 = bgs
	} else {
		if bgs2.PrevSegment2 != nil && bgs2.PrevSegment2 != bgs {
			panic("belt connectivity inconsistency")
		}
		bgs2.PrevSegment2 = bgs
		bgs2.prevSegment2Offset = nextOffset
	}
}

func (bgs *BeltGraphSegment) Reconnect(bgs2 *BeltGraphSegment, start, end, nextOffset float64) {
	bgs.DisconnectNext()
	bgs.Connect(bgs2, start, end, nextOffset)
}

func (bgs *BeltGraphSegment) DropAllItems() {
	// TODO Drop items?????????????????????????
	bgs.Items.Clear()
}

func (bgs *BeltGraphSegment) Change(start, end float64) {
	bgs.Start = start
	bgs.End = end

	// TODO Drop items?????????????????????????
	for bgs.Items.Len() > 0 {
		iob, _ := bgs.Items.PeekFirst()
		if iob.Offset >= bgs.Start {
			break
		}
		bgs.Items.Pop()
	}
	for bgs.Items.Len() > 0 {
		iob, _ := bgs.Items.PeekLast()
		if iob.Offset <= bgs.End {
			break
		}
		bgs.Items.PopLast()
	}
}

func (bgs *BeltGraphSegment) DisconnectNext() {
	if bgs.NextSegment == nil {
		return
	}

	if bgs.nextSegmentOffset == 0 {
		if bgs.NextSegment.PrevSegment1 != bgs {
			panic("belt connectivity inconsistency")
		}
		bgs.NextSegment.PrevSegment1 = nil
	} else {
		if bgs.NextSegment.PrevSegment2 != bgs {
			panic("belt connectivity inconsistency")
		}
		bgs.NextSegment.PrevSegment2 = nil
	}
	bgs.NextSegment = nil
}

func (bgs *BeltGraphSegment) CanPlaceItem(offset float64) bool {
	for i := 0; i < bgs.Items.Len(); i++ {
		do := bgs.Items.Peek(i).Offset - offset
		if -ss.ITEM_DW < do && do < ss.ITEM_DW {
			return false
		}
	}

	if offset < bgs.End && bgs.End-offset < ss.ITEM_DW && bgs.NextSegment != nil {
		if !bgs.NextSegment.CanPlaceItem(bgs.nextSegmentOffset - (bgs.End - offset)) {
			return false
		}
	}

	if offset >= 0 && offset < ss.ITEM_DW && bgs.PrevSegment1 != nil {
		if !bgs.PrevSegment1.CanPlaceItem(bgs.PrevSegment1.End + offset) {
			return false
		}
	}

	if bgs.PrevSegment2 != nil {
		do := math.Abs(bgs.prevSegment2Offset - offset)
		if do < ss.ITEM_DW {
			if !bgs.PrevSegment2.CanPlaceItem(bgs.PrevSegment2.End + do) {
				return false
			}
		}
	}
	return true
}

func (bgs *BeltGraphSegment) insertItem(iob *items.ItemOnBelt) *items.ItemOnBelt {
	i := 0
	for ; i < bgs.Items.Len(); i++ {
		if bgs.Items.Peek(i).Offset < iob.Offset {
			continue

		}
		break
	}

	err := bgs.Items.Insert(i, iob)
	if err != nil {
		panic(err)
	}
	return iob
}

func (bgs *BeltGraphSegment) PlaceItem(item items.ItemInWorld, offset float64) {
	iob := items.NewItemOnBelt2(item, offset)
	bgs.TakeItemOnBelt(iob, offset)
}

func (bgs *BeltGraphSegment) TakeItemOnBelt(iob *items.ItemOnBelt, offset float64) {
	iob.Offset = offset
	bgs.insertItem(iob)
	bgs.SetItemPos(iob, false)
}

func (bgs *BeltGraphSegment) UpdateAllItemsPos() {
	for i := 0; i < bgs.Items.Len(); i++ {
		iob := bgs.Items.Peek(i)
		bgs.SetItemPos(iob, true)
	}
}

func (bgs *BeltGraphSegment) SetItemPos(iob *items.ItemOnBelt, reset bool) {
	mult := math.Abs(iob.Offset - bgs.ShiftOffset)
	newPos := bgs.PosZero.Add(bgs.PosHalfOffset.Mul(2 * iob.Offset)).Add(bgs.Shift.Mul(4 * mult * mult))
	iob.Item.Pos.UpdatePosition(newPos, reset)
}

func (bgs *BeltGraphSegment) DrawItems(r *renderer.GameRenderer) {
	for i := 0; i < bgs.Items.Len(); i++ {
		item := bgs.Items.Peek(i)
		item.Draw(r)
	}
}

func (bgs *BeltGraphSegment) GetFirstOffset() float64 {
	iob, ok := bgs.Items.PeekFirst()
	if !ok {
		return ss.ITEM_DW * 2
	}
	return iob.Offset
}
func (bgs *BeltGraphSegment) GetLastOffset() float64 {
	iob, ok := bgs.Items.PeekLast()
	if !ok {
		return -ss.ITEM_DW * 2
	}
	return iob.Offset
}

func (bgs *BeltGraphSegment) GetClosestOffset(beginning bool) (float64, bool) {
	if beginning {
		sideOffset := bgs.End + ss.ITEM_DW
		if bgs.PrevSegment2 != nil {
			sideOffset = bgs.PrevSegment2.End - bgs.PrevSegment2.GetLastOffset() + bgs.prevSegment2Offset
		}
		if iob, ok := bgs.Items.PeekFirst(); ok {
			return min(iob.Offset, sideOffset), false
		}
		return sideOffset, false
	}

	closestOffset := bgs.End - bgs.prevSegment2Offset + ss.ITEM_DW
	blocked := false
	if bgs.PrevSegment1 != nil {
		prevOffset := bgs.PrevSegment1.GetLastOffset() - bgs.PrevSegment1.End - bgs.prevSegment2Offset
		if -prevOffset < ss.ITEM_DW*2 {
			blocked = true
		}
	}
	for i := 0; i < bgs.Items.Len(); i++ {
		offset := bgs.Items.Peek(i).Offset - bgs.prevSegment2Offset
		if offset < 0 && -offset < ss.ITEM_DW*2 {
			blocked = true
			continue
		}
		closestOffset = min(offset, closestOffset)
	}
	return closestOffset, blocked
}

func (bgs *BeltGraphSegment) MoveItems(ticks uint64) {
	maxOffset := bgs.End + ss.ITEM_DW/2
	blocked := false
	if bgs.NextSegment != nil {
		var offset float64
		offset, blocked = bgs.NextSegment.GetClosestOffset(bgs.nextSegmentOffset == 0.0)
		maxOffset = bgs.End + offset
	}

	for i := bgs.Items.Len() - 1; i >= 0; i-- {
		iob := bgs.Items.Peek(i)
		if iob.MovedTick == ticks {
			continue
		}
		iob.MovedTick = ticks

		newOffset := min(iob.Offset+bgs.speed, maxOffset-ss.ITEM_DW)
		if blocked {
			if iob.Offset < bgs.End-ss.ITEM_DW {
				newOffset = min(newOffset, bgs.End-ss.ITEM_DW-0.0001)
			}
		}

		if newOffset > bgs.End {
			if bgs.NextSegment != nil {
				bgs.NextSegment.TakeItemOnBelt(iob, newOffset-bgs.End+bgs.nextSegmentOffset)
				if i, ok := bgs.Items.PopLast(); !ok || i != iob {
					panic("item order is messed up")
				}
				maxOffset = min(newOffset, maxOffset)
				continue
			}
			newOffset = bgs.End
		}

		iob.Offset = max(iob.Offset, newOffset)
		bgs.SetItemPos(iob, false)
		maxOffset = min(newOffset, maxOffset)
	}
}

func (bgs *BeltGraphSegment) ReverseItemOffsets() {
	iobs := make([]*items.ItemOnBelt, 0, bgs.Items.Len())
	for bgs.Items.Len() > 0 {
		iob, _ := bgs.Items.Pop()
		iobs = append(iobs, iob)
	}

	for len(iobs) > 0 {
		iob := iobs[len(iobs)-1]
		iob.Offset = bgs.End - iob.Offset
		_ = bgs.Items.Push(iob)
		iobs = iobs[:len(iobs)-1]
	}
}

func (bgs *BeltGraphSegment) PopAllItems() []*items.ItemOnBelt {
	iobs := make([]*items.ItemOnBelt, 0, bgs.Items.Len())
	for bgs.Items.Len() > 0 {
		iob, _ := bgs.Items.Pop()
		iobs = append(iobs, iob)
	}
	return iobs
}

func (bgs *BeltGraphSegment) PushItemsReversed(iobs []*items.ItemOnBelt) {
	for len(iobs) > 0 {
		iob := iobs[len(iobs)-1]
		iob.Offset = bgs.End - iob.Offset
		_ = bgs.Items.Push(iob)
		bgs.SetItemPos(iob, true)
		iobs = iobs[:len(iobs)-1]
	}
}

func (bgs *BeltGraphSegment) FindClosestItem(pos utils.WorldCoord, allowedItems []ss.ItemType) (*items.ItemOnBelt, float64) {
	if bgs.Items.Len() == 0 {
		return nil, 0
	}

	var closestIdx int = -1
	closestDistSq := ss.ITEM_D * ss.ITEM_D * 2.0
	for i := 0; i < bgs.Items.Len(); i++ {
		iob := bgs.Items.Peek(i)
		if allowedItems != nil && !utils.ItemInList(iob.Item.ItemType, allowedItems) {
			continue
		}
		if closestIdx == -1 || iob.Item.Pos.Pos.DistanceSqTo(pos) < closestDistSq {
			closestDistSq = iob.Item.Pos.Pos.DistanceSqTo(pos)
			closestIdx = i
		}
	}

	if closestIdx == -1 {
		return nil, 0
	}
	return bgs.Items.Peek(closestIdx), closestDistSq
}

func (bgs *BeltGraphSegment) TakeItemOut(item *items.ItemOnBelt) bool {
	for i := 0; i < bgs.Items.Len(); i++ {
		if bgs.Items.Peek(i) == item {
			_, ok := bgs.Items.PopAt(i)
			if !ok {
				panic("item not found")
			}
			return true
		}
	}
	return false
}
