package objects

import (
	"hextopdown/game/items"
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
)

var splitterTypeMapping = map[utils.Dir]ss.BeltType{
	utils.DIR_LEFT:       ss.BELT_TYPE_SPLITTER_LEFTUPLEFT_RIGHTDOWNRIGHT,
	utils.DIR_UP_LEFT:    ss.BELT_TYPE_SPLITTER_UPLEFTRIGHT_DOWNLEFTRIGHT,
	utils.DIR_UP_RIGHT:   ss.BELT_TYPE_SPLITTER_RIGHTUPRIGHT_LEFTDOWNLEFT,
	utils.DIR_RIGHT:      ss.BELT_TYPE_SPLITTER_RIGHTDOWNRIGHT_LEFTUPLEFT,
	utils.DIR_DOWN_RIGHT: ss.BELT_TYPE_SPLITTER_DOWNLEFTRIGHT_UPLEFTRIGHT,
	utils.DIR_DOWN_LEFT:  ss.BELT_TYPE_SPLITTER_LEFTDOWNLEFT_RIGHTUPRIGHT,
}
var splitterOnMapping = map[utils.Dir]ss.BeltType{
	utils.DIR_LEFT:       ss.BELT_ON_SPLITTER_LEFTUPLEFT_RIGHTDOWNRIGHT,
	utils.DIR_UP_LEFT:    ss.BELT_ON_SPLITTER_UPLEFTRIGHT_DOWNLEFTRIGHT,
	utils.DIR_UP_RIGHT:   ss.BELT_ON_SPLITTER_RIGHTUPRIGHT_LEFTDOWNLEFT,
	utils.DIR_RIGHT:      ss.BELT_ON_SPLITTER_RIGHTDOWNRIGHT_LEFTUPLEFT,
	utils.DIR_DOWN_RIGHT: ss.BELT_ON_SPLITTER_DOWNLEFTRIGHT_UPLEFTRIGHT,
	utils.DIR_DOWN_LEFT:  ss.BELT_ON_SPLITTER_LEFTDOWNLEFT_RIGHTUPRIGHT,
}

type BeltSplitter struct {
	Pos          utils.HexCoord
	inConns      [2]*BeltConnection
	outConn      [2]*BeltConnection
	speed        float64
	outPrioLeft  bool
	outPrioRight bool
	beltType     ss.BeltType
	onType       ss.BeltType
}

func NewBeltSplitter(pos utils.HexCoord, dir utils.Dir, speed float64) *BeltSplitter {
	beltType, ok := splitterTypeMapping[dir.Reverse()]
	if !ok {
		panic("invalid belt type")
	}
	onType, ok := splitterOnMapping[dir.Reverse()]
	if !ok {
		panic("invalid belt on-ground type")
	}

	inl := newInConn(pos, dir.Reverse(), speed)
	inr := newInConn(pos, dir.Reverse().NextCW(), speed)

	outl := NewBeltConnection(pos, dir, speed, false)
	outr := NewBeltConnection(pos, dir.NextCW(), speed, false)

	return &BeltSplitter{
		Pos:      pos,
		inConns:  [2]*BeltConnection{inl, inr},
		outConn:  [2]*BeltConnection{outl, outr},
		speed:    speed,
		beltType: beltType,
		onType:   onType,
	}
}

func newInConn(pos utils.HexCoord, dir utils.Dir, speed float64) *BeltConnection {
	conn := NewBeltConnectionWithDist(pos, dir, speed, true, 0.5+ss.ITEM_DW/2)
	conn.LaneLeft.End = 0.5 + ss.ITEM_DW/2
	conn.LaneRight.End = 0.5 + ss.ITEM_DW/2
	return conn
}

func (b *BeltSplitter) GetObjectType() ss.ObjectType {
	return ss.OBJECT_TYPE_BELTSPLITTER1
}

func (b *BeltSplitter) GetPos() utils.HexCoord {
	return b.Pos
}

func (b *BeltSplitter) GetDir() utils.Dir {
	return b.inConns[0].Dir
}

func (b *BeltSplitter) updateType() {
	beltType, ok := splitterTypeMapping[b.inConns[0].Dir]
	if !ok {
		panic("invalid belt type")
	}
	b.beltType = beltType

	onType, ok := splitterOnMapping[b.inConns[0].Dir]
	if !ok {
		panic("invalid belt on-ground type")
	}
	b.onType = onType
}

func (b *BeltSplitter) DrawGroundLevel(r *renderer.GameRenderer) {
	if !r.IsHexOnScreen(b.Pos) {
		return
	}
	r.DrawAnimatedBelt(b.Pos, b.beltType, b.speed*ss.TPS)

	// b.outConn[0].Draw(b.Pos, r)
	// b.outConn[1].Draw(b.Pos, r)

	// b.inConns[0].Draw(b.Pos, r)
	// b.inConns[1].Draw(b.Pos, r)
}

func (b *BeltSplitter) DrawOnGroundLevel(r *renderer.GameRenderer) {
	r.DrawBeltOnGround(b.Pos, b.onType)
}

func (b *BeltSplitter) Update(ticks uint64, world HexGridWorldInteractor) {
	b.jumpItems(b.inConns[0])
	b.jumpItems(b.inConns[1])
}

func (b *BeltSplitter) jumpItems(inConn *BeltConnection) {
	b.jumpLaneItems(&inConn.LaneLeft, &b.outConn[0].LaneLeft, &b.outConn[1].LaneLeft, &b.outPrioLeft)
	b.jumpLaneItems(&inConn.LaneRight, &b.outConn[0].LaneRight, &b.outConn[1].LaneRight, &b.outPrioRight)
}

func (b *BeltSplitter) jumpLaneItems(inLane *BeltGraphSegment, outLane1, outLane2 *BeltGraphSegment, outPrio *bool) {
	items := &inLane.Items
	iob, ok := items.PeekLast()
	if !ok {
		return
	}

	nextOffset := iob.Offset + b.speed
	if nextOffset < inLane.End-ss.ITEM_DW/2 {
		return
	}
	nextOffset -= inLane.End - ss.ITEM_DW/2

	if !*outPrio {
		outLane1, outLane2 = outLane2, outLane1
	}

	if nextOffset < outLane1.GetFirstOffset()-ss.ITEM_DW {
		outLane1.TakeItemOnBelt(iob, nextOffset)
		items.PopLast()
		*outPrio = !*outPrio
		return
	}

	if nextOffset < outLane2.GetFirstOffset()-ss.ITEM_DW {
		outLane2.TakeItemOnBelt(iob, nextOffset)
		items.PopLast()
	}
}

func (b *BeltSplitter) DrawItems(r *renderer.GameRenderer) {
	b.outConn[0].DrawItem(b.Pos, r)
	b.outConn[1].DrawItem(b.Pos, r)

	b.inConns[0].DrawItem(b.Pos, r)
	b.inConns[1].DrawItem(b.Pos, r)
}

func (b *BeltSplitter) GetInConn(dir utils.Dir) *BeltConnection {
	if dir == b.inConns[0].Dir {
		return b.inConns[0]
	}
	if dir == b.inConns[1].Dir {
		return b.inConns[1]
	}

	return nil
}

func (b *BeltSplitter) ClearIn(dir utils.Dir) {
	if dir == b.inConns[0].Dir {
		b.inConns[0].Belt = nil
	} else if dir == b.inConns[1].Dir {
		b.inConns[1].Belt = nil
	} else {
		panic("belt connectivity inconsistency")
	}
}

func (b *BeltSplitter) ClearOut(dir utils.Dir) {
	if dir == b.outConn[0].Dir {
		b.outConn[0].DisconnectNext()
	} else if dir == b.outConn[1].Dir {
		b.outConn[1].DisconnectNext()
	} else {
		panic("belt connectivity inconsistency")
	}
}

func (b *BeltSplitter) CreateIn(dir utils.Dir, b2 BeltLike) {
	if dir == b.inConns[0].Dir {
		b.inConns[0].Belt = b2
	} else if dir == b.inConns[1].Dir {
		b.inConns[1].Belt = b2
	} else {
		panic("invalid incoming direction")
	}
}

func (b *BeltSplitter) DisconnectAll() {
	if b.outConn[0].Belt != nil {
		b.outConn[0].Belt.ClearIn(b.outConn[0].Dir.Reverse())
		b.ClearOut(b.outConn[0].Dir)
	}
	if b.outConn[1].Belt != nil {
		b.outConn[1].Belt.ClearIn(b.outConn[1].Dir.Reverse())
		b.ClearOut(b.outConn[1].Dir)
	}

	if b.inConns[0].Belt != nil {
		b.inConns[0].Belt.ClearOut(b.inConns[0].Dir.Reverse())
		b.ClearIn(b.inConns[0].Dir)
	}
	if b.inConns[1].Belt != nil {
		b.inConns[1].Belt.ClearOut(b.inConns[1].Dir.Reverse())
		b.ClearIn(b.inConns[1].Dir)
	}
}

func (b *BeltSplitter) CanConnectIn(dir utils.Dir) bool {
	return dir == b.inConns[0].Dir || dir == b.inConns[1].Dir
}

func (b *BeltSplitter) CanConnectTo(b2 BeltLike) bool {
	dir, err := b.Pos.DirTo(b2.GetPos())
	if err != nil {
		return false
	}

	if dir != b.outConn[0].Dir && dir != b.outConn[1].Dir {
		return false
	}
	if !b2.CanConnectIn(dir.Reverse()) {
		return false
	}

	if dir == b.outConn[0].Dir && b.outConn[0].Belt != nil {
		if b.outConn[0].Belt != b2 {
			panic("belt connectivity inconsistency")
		}
		return false
	}
	if dir == b.outConn[1].Dir && b.outConn[1].Belt != nil {
		if b.outConn[1].Belt != b2 {
			panic("belt connectivity inconsistency")
		}
		return false
	}

	return true
}

func (b *BeltSplitter) Rotate(cw bool) {
	dir := b.outConn[0].Dir
	if cw {
		dir = dir.NextCW()
	} else {
		dir = dir.NextCCW()
	}

	if cw {
		if b.outConn[0].Belt != nil {
			b.outConn[0].Belt.ClearIn(b.outConn[0].Dir.Reverse())
			b.ClearOut(b.outConn[0].Dir)
		}
		b.outConn[0] = b.outConn[1]
		b.outConn[1] = NewBeltConnection(b.Pos, dir.NextCW(), b.speed, false)

		if b.inConns[0].Belt != nil {
			b.inConns[0].Belt.ClearOut(b.inConns[0].Dir.Reverse())
			b.ClearIn(b.inConns[0].Dir)
		}
		b.inConns[0] = b.inConns[1]
		b.inConns[1] = newInConn(b.Pos, dir.Reverse().NextCW(), b.speed) // NewBeltConnection(b.Pos, dir.Reverse().NextCW(), b.speed, true)
	} else {
		if b.outConn[1].Belt != nil {
			b.outConn[1].Belt.ClearIn(b.outConn[1].Dir.Reverse())
			b.ClearOut(b.outConn[1].Dir)
		}
		b.outConn[1] = b.outConn[0]
		b.outConn[0] = NewBeltConnection(b.Pos, dir, b.speed, false)

		if b.inConns[1].Belt != nil {
			b.inConns[1].Belt.ClearOut(b.inConns[1].Dir.Reverse())
			b.ClearIn(b.inConns[1].Dir)
		}
		b.inConns[1] = b.inConns[0]
		b.inConns[0] = newInConn(b.Pos, dir.Reverse(), b.speed) //NewBeltConnection(b.Pos, dir.Reverse(), b.speed, true)
	}

	b.updateType()
}

func (b *BeltSplitter) ConnectTo(b2 BeltLike) {
	dir, err := b.Pos.DirTo(b2.GetPos())
	if err != nil {
		panic(err)
	}

	var bc *BeltConnection
	if dir == b.outConn[0].Dir {
		bc = b.outConn[0]
	} else {
		bc = b.outConn[1]
	}
	bc.Belt = b2
	b2.CreateIn(dir.Reverse(), b)
	bc.ReconnectTo(b2.GetInConn(dir.Reverse()))
}

func (b *BeltSplitter) MoveItems(ticks uint64, processed map[*BeltGraphSegment]struct{}) {
	b.outConn[0].MoveItems(ticks, processed)
	b.outConn[1].MoveItems(ticks, processed)
}

func (b *BeltSplitter) TakeItemOut(pos utils.WorldCoord) (*items.ItemInWorld, bool) {
	var closestItem *items.ItemOnBelt
	closestDistSq := 99999999999.0
	var closestConn *BeltConnection

	for i := 0; i < 2; i++ {
		if iob, distSq := b.outConn[i].FindClosestItem(pos); iob != nil {
			if distSq < closestDistSq {
				closestItem = iob
				closestDistSq = distSq
				closestConn = b.outConn[i]
			}
		}
	}

	for i := 0; i < 2; i++ {
		if b.inConns[i] == nil {
			continue
		}
		if iob, distSq := b.inConns[i].FindClosestItem(pos); iob != nil {
			if distSq < closestDistSq {
				closestItem = iob
				closestDistSq = distSq
				closestConn = b.inConns[i]
			}
		}
	}

	if closestItem != nil {
		if !closestConn.TakeItemOut(closestItem) {
			panic("items are messed up")
		}
		return &closestItem.Item, true
	}
	return nil, false
}
