package objects

import (
	"hextopdown/game/items"
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
	"math"
)

type BeltConnection struct {
	Belt      BeltLike
	Dir       utils.Dir
	LaneLeft  BeltGraphSegment
	LaneRight BeltGraphSegment
	isIn      bool
}

func NewBeltConnection(hex utils.HexCoord, dir utils.Dir, speed float64, isIn bool) *BeltConnection {
	return NewBeltConnectionWithDist(hex, dir, speed, isIn, 0.5)
}

func NewBeltConnectionWithDist(hex utils.HexCoord, dir utils.Dir, speed float64, isIn bool, dist float64) *BeltConnection {
	lz, lho := getZeroOffsetPos(hex, dir, isIn, true)
	rz, rho := getZeroOffsetPos(hex, dir, isIn, false)
	capacity := int(math.Ceil(dist / ss.ITEM_DW))

	return &BeltConnection{
		Dir:       dir,
		LaneLeft:  NewBeltGraphSegment(lz, lho, speed, capacity),
		LaneRight: NewBeltGraphSegment(rz, rho, speed, capacity),
		isIn:      isIn,
	}
}

func (bc *BeltConnection) UpdateDir(hex utils.HexCoord, dir utils.Dir) {
	bc.Dir = dir
	bc.updateLanePos(hex)
}

func (bc *BeltConnection) UpdateShifts(isKink, isLeft bool) {
	if !isKink {
		bc.LaneLeft.Shift = utils.WorldCoord{X: 0, Y: 0}
		bc.LaneRight.Shift = utils.WorldCoord{X: 0, Y: 0}
		bc.LaneLeft.UpdateAllItemsPos()
		bc.LaneRight.UpdateAllItemsPos()
		return
	}

	var rotInside, rotOutside utils.WorldCoord
	if bc.isIn {
		rotInside = utils.WorldCoord{X: ss.LANE_OFFSET_WORLD * 0.8660254, Y: ss.LANE_OFFSET_WORLD * 0.5}
		rotOutside = utils.WorldCoord{X: 0, Y: ss.LANE_OFFSET_WORLD}
	} else {
		rotInside = utils.WorldCoord{X: ss.LANE_OFFSET_WORLD * 0.8660254, Y: -ss.LANE_OFFSET_WORLD * 0.5}
		rotOutside = utils.WorldCoord{X: 0, Y: -ss.LANE_OFFSET_WORLD}
	}
	if !isLeft {
		rotInside.Y = -rotInside.Y
		rotOutside.Y = -rotOutside.Y
	}

	if isLeft {
		bc.LaneLeft.Shift = rotInside.UnrotateFromPlusX(bc.Dir)
		bc.LaneRight.Shift = rotOutside.UnrotateFromPlusX(bc.Dir)
	} else {
		bc.LaneLeft.Shift = rotOutside.UnrotateFromPlusX(bc.Dir)
		bc.LaneRight.Shift = rotInside.UnrotateFromPlusX(bc.Dir)
	}

	if bc.isIn {
		bc.LaneLeft.ShiftOffset = 0
		bc.LaneRight.ShiftOffset = 0
	} else {
		bc.LaneLeft.ShiftOffset = 0.5
		bc.LaneRight.ShiftOffset = 0.5
	}
	bc.LaneLeft.UpdateAllItemsPos()
	bc.LaneRight.UpdateAllItemsPos()
}

func (bc *BeltConnection) ReconnectTo(bc2 *BeltConnection) {
	bc.LaneLeft.Reconnect(&bc2.LaneLeft, 0, 0.5, 0)
	bc.LaneRight.Reconnect(&bc2.LaneRight, 0, 0.5, 0)
}

func (bc *BeltConnection) ReconnectToWithDist(bc2 *BeltConnection, hexDist int32) {
	bc.LaneLeft.Reconnect(&bc2.LaneLeft, 0, float64(hexDist), 0)
	bc.LaneRight.Reconnect(&bc2.LaneRight, 0, float64(hexDist), 0)
}

func (bc *BeltConnection) DisconnectNext() {
	bc.LaneLeft.DisconnectNext()
	bc.LaneRight.DisconnectNext()
	bc.Belt = nil
}

func (bc *BeltConnection) Reverse(hex utils.HexCoord) {
	bc.isIn = !bc.isIn
	bc.LaneLeft, bc.LaneRight = bc.LaneRight, bc.LaneLeft

	bc.LaneLeft.ReverseItemOffsets()
	bc.LaneRight.ReverseItemOffsets()

	bc.updateLanePos(hex)
}

func (bc *BeltConnection) GetClosestSegment(localPos utils.WorldCoord) (distSq float64, isLeft bool, offset float64) {
	rotPos := localPos.RotateToPlusX(bc.Dir)

	if !bc.isIn {
		d2left, oLeft := distSqToX(rotPos.X, rotPos.Y+ss.LANE_OFFSET_WORLD, bc.LaneLeft.Start*ss.HEX_WIDTH, bc.LaneLeft.End*ss.HEX_WIDTH)
		d2right, oRight := distSqToX(rotPos.X, rotPos.Y-ss.LANE_OFFSET_WORLD, bc.LaneRight.Start*ss.HEX_WIDTH, bc.LaneRight.End*ss.HEX_WIDTH)
		if d2left < d2right {
			return d2left, true, oLeft/ss.HEX_WIDTH + bc.LaneLeft.Start
		} else {
			return d2right, false, oRight/ss.HEX_WIDTH + bc.LaneRight.Start
		}
	}

	d2left, oLeft := distSqToX(ss.HEX_WIDTH/2-rotPos.X, rotPos.Y-ss.LANE_OFFSET_WORLD, bc.LaneLeft.Start*ss.HEX_WIDTH, bc.LaneLeft.End*ss.HEX_WIDTH)
	d2right, oRight := distSqToX(ss.HEX_WIDTH/2-rotPos.X, rotPos.Y+ss.LANE_OFFSET_WORLD, bc.LaneRight.Start*ss.HEX_WIDTH, bc.LaneRight.End*ss.HEX_WIDTH)
	if d2left < d2right {
		return d2left, true, oLeft/ss.HEX_WIDTH + bc.LaneLeft.Start
	} else {
		return d2right, false, oRight/ss.HEX_WIDTH + bc.LaneRight.Start
	}
}

func (bc *BeltConnection) CanPlaceItem(offset float64, isLeft bool) bool {
	if isLeft {
		return bc.LaneLeft.CanPlaceItem(offset)
	}
	return bc.LaneRight.CanPlaceItem(offset)
}

func (bc *BeltConnection) PlaceItem(item items.ItemInWorld, offset float64, isLeft bool) {
	if isLeft {
		bc.LaneLeft.PlaceItem(item, offset)
	} else {
		bc.LaneRight.PlaceItem(item, offset)
	}
}

func (bc *BeltConnection) Draw(hex utils.HexCoord, r *renderer.GameRenderer) {
	if !bc.isIn {
		r.DrawBeltConnectionsOutgoing(hex, bc.Dir)
	} else {
		r.DrawBeltConnectionIncoming(hex, bc.Dir, true, bc.LaneLeft.Start, bc.LaneLeft.End)
		r.DrawBeltConnectionIncoming(hex, bc.Dir, false, bc.LaneRight.Start, bc.LaneRight.End)
	}
}

func (bc *BeltConnection) DrawItem(hex utils.HexCoord, r *renderer.GameRenderer) {
	bc.LaneLeft.DrawItems(r)
	bc.LaneRight.DrawItems(r)
}

func (bc *BeltConnection) MoveItems(ticks uint64, processed map[*BeltGraphSegment]struct{}) {
	if _, ok := processed[&bc.LaneLeft]; !ok {
		graphStart := findBeltSegmentGraphStart(&bc.LaneLeft)
		moveItemsWholeGraph(graphStart, ticks, processed)
	}

	if _, ok := processed[&bc.LaneRight]; !ok {
		graphStart := findBeltSegmentGraphStart(&bc.LaneRight)
		moveItemsWholeGraph(graphStart, ticks, processed)
	}
}

func (bc *BeltConnection) SwapReverseItems() {
	iobs1 := bc.LaneLeft.PopAllItems()
	iobs2 := bc.LaneRight.PopAllItems()

	bc.LaneLeft.PushItemsReversed(iobs2)
	bc.LaneRight.PushItemsReversed(iobs1)
}

func (bc *BeltConnection) updateLanePos(hex utils.HexCoord) {
	lz, lho := getZeroOffsetPos(hex, bc.Dir, bc.isIn, true)
	rz, rho := getZeroOffsetPos(hex, bc.Dir, bc.isIn, false)

	bc.LaneLeft.PosZero = lz
	bc.LaneLeft.PosHalfOffset = lho
	bc.LaneRight.PosZero = rz
	bc.LaneRight.PosHalfOffset = rho

	bc.LaneLeft.UpdateAllItemsPos()
	bc.LaneRight.UpdateAllItemsPos()
}

func distSqToX(x, y float64, x1, x2 float64) (float64, float64) {
	if x1 > x2 {
		panic("invalid x1 and x2")
	}

	if x < x1 {
		dx := x1 - x
		return dx*dx + y*y, x1
	}
	if x > x2 {
		dx := x - x2
		return dx*dx + y*y, x2
	}
	return y * y, x - x1
}

func getZeroOffsetPos(hex utils.HexCoord, dir utils.Dir, isIn, isLeft bool) (utils.WorldCoord, utils.WorldCoord) {
	var rxz, rxh, ry float64

	if isIn {
		rxz = 0.5 * ss.HEX_WIDTH
		rxh = 0
		ry = ss.LANE_OFFSET_WORLD
	} else {
		rxz = 0
		rxh = 0.5 * ss.HEX_WIDTH
		ry = -ss.LANE_OFFSET_WORLD
	}
	if !isLeft {
		ry = -ry
	}

	lzl := utils.WorldCoord{X: rxz, Y: ry}.UnrotateFromPlusX(dir)
	lhl := utils.WorldCoord{X: rxh, Y: ry}.UnrotateFromPlusX(dir)
	center := hex.CenterToWorld()

	return center.Add(lzl), lhl.Sub(lzl)
}

func findBeltSegmentGraphStart(bgs *BeltGraphSegment) *BeltGraphSegment {
	current := bgs
	visited := make(map[*BeltGraphSegment]struct{})
	for {
		visited[current] = struct{}{}
		nextSegment := current.NextSegment
		if nextSegment == nil {
			break
		}
		if _, ok := visited[nextSegment]; ok {
			break
		}
		current = nextSegment
	}
	return current
}

func moveItemsWholeGraph(graphStart *BeltGraphSegment, ticks uint64, processed map[*BeltGraphSegment]struct{}) {
	toProcess := []*BeltGraphSegment{graphStart}
	for len(toProcess) > 0 {
		bgs := toProcess[len(toProcess)-1]
		toProcess = toProcess[:len(toProcess)-1]
		if _, ok := processed[bgs]; ok {
			continue
		}
		bgs.MoveItems(ticks)

		processed[bgs] = struct{}{}
		if bgs.PrevSegment1 != nil {
			toProcess = append(toProcess, bgs.PrevSegment1)
		}
		if bgs.PrevSegment2 != nil {
			toProcess = append(toProcess, bgs.PrevSegment2)
		}
	}
}

func (bc *BeltConnection) FindClosestItem(pos utils.WorldCoord) (*items.ItemOnBelt, float64) {
	var closestItem *items.ItemOnBelt
	minDistSq := 99999999999.0
	if item, distSq := bc.LaneLeft.FindClosestItem(pos); item != nil && distSq < minDistSq {
		minDistSq = distSq
		closestItem = item
	}
	if item, distSq := bc.LaneRight.FindClosestItem(pos); item != nil && distSq < minDistSq {
		minDistSq = distSq
		closestItem = item
	}
	return closestItem, minDistSq
}

func (bc *BeltConnection) TakeItemOut(item *items.ItemOnBelt) bool {
	if bc.LaneLeft.TakeItemOut(item) {
		return true
	}
	return bc.LaneRight.TakeItemOut(item)
}
