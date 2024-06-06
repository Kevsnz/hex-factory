package objects

import (
	"fmt"
	"hextopdown/game/items"
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
)

type typeMappingKey struct {
	outDir utils.Dir
	ins    [3]bool
}

var beltTypeMapping = map[typeMappingKey]ss.BeltType{
	// Output only
	{utils.DIR_LEFT, [3]bool{false, false, false}}:       ss.BELT_TYPE_LEFT,
	{utils.DIR_RIGHT, [3]bool{false, false, false}}:      ss.BELT_TYPE_RIGHT,
	{utils.DIR_UP_LEFT, [3]bool{false, false, false}}:    ss.BELT_TYPE_UPLEFT,
	{utils.DIR_UP_RIGHT, [3]bool{false, false, false}}:   ss.BELT_TYPE_UPRIGHT,
	{utils.DIR_DOWN_LEFT, [3]bool{false, false, false}}:  ss.BELT_TYPE_DOWNLEFT,
	{utils.DIR_DOWN_RIGHT, [3]bool{false, false, false}}: ss.BELT_TYPE_DOWNRIGHT,

	// Straight through
	{utils.DIR_LEFT, [3]bool{true, false, false}}:       ss.BELT_TYPE_RIGHT_LEFT,
	{utils.DIR_RIGHT, [3]bool{true, false, false}}:      ss.BELT_TYPE_LEFT_RIGHT,
	{utils.DIR_UP_LEFT, [3]bool{true, false, false}}:    ss.BELT_TYPE_DOWNRIGHT_UPLEFT,
	{utils.DIR_UP_RIGHT, [3]bool{true, false, false}}:   ss.BELT_TYPE_DOWNLEFT_UPRIGHT,
	{utils.DIR_DOWN_LEFT, [3]bool{true, false, false}}:  ss.BELT_TYPE_UPRIGHT_DOWNLEFT,
	{utils.DIR_DOWN_RIGHT, [3]bool{true, false, false}}: ss.BELT_TYPE_UPLEFT_DOWNRIGHT,

	// Bend
	{utils.DIR_DOWN_RIGHT, [3]bool{false, false, true}}: ss.BELT_TYPE_LEFT_DOWNRIGHT,
	{utils.DIR_UP_RIGHT, [3]bool{false, true, false}}:   ss.BELT_TYPE_LEFT_UPRIGHT,
	{utils.DIR_DOWN_LEFT, [3]bool{false, true, false}}:  ss.BELT_TYPE_RIGHT_DOWNLEFT,
	{utils.DIR_UP_LEFT, [3]bool{false, false, true}}:    ss.BELT_TYPE_RIGHT_UPLEFT,

	{utils.DIR_LEFT, [3]bool{false, true, false}}:  ss.BELT_TYPE_DOWNRIGHT_LEFT,
	{utils.DIR_LEFT, [3]bool{false, false, true}}:  ss.BELT_TYPE_UPRIGHT_LEFT,
	{utils.DIR_RIGHT, [3]bool{false, false, true}}: ss.BELT_TYPE_DOWNLEFT_RIGHT,
	{utils.DIR_RIGHT, [3]bool{false, true, false}}: ss.BELT_TYPE_UPLEFT_RIGHT,

	{utils.DIR_DOWN_LEFT, [3]bool{false, false, true}}:  ss.BELT_TYPE_UPLEFT_DOWNLEFT,
	{utils.DIR_UP_LEFT, [3]bool{false, true, false}}:    ss.BELT_TYPE_DOWNLEFT_UPLEFT,
	{utils.DIR_DOWN_RIGHT, [3]bool{false, true, false}}: ss.BELT_TYPE_UPRIGHT_DOWNRIGHT,
	{utils.DIR_UP_RIGHT, [3]bool{false, false, true}}:   ss.BELT_TYPE_DOWNRIGHT_UPRIGHT,

	// Straight join
	{utils.DIR_DOWN_RIGHT, [3]bool{true, false, true}}: ss.BELT_TYPE_UPLEFT_DOWNRIGHT_LEFT,
	{utils.DIR_UP_RIGHT, [3]bool{true, true, false}}:   ss.BELT_TYPE_DOWNLEFT_UPRIGHT_LEFT,
	{utils.DIR_UP_LEFT, [3]bool{true, false, true}}:    ss.BELT_TYPE_DOWNRIGHT_UPLEFT_RIGHT,
	{utils.DIR_DOWN_LEFT, [3]bool{true, true, false}}:  ss.BELT_TYPE_UPRIGHT_DOWNLEFT_RIGHT,

	{utils.DIR_DOWN_RIGHT, [3]bool{true, true, false}}: ss.BELT_TYPE_UPLEFT_DOWNRIGHT_UPRIGHT,
	{utils.DIR_UP_RIGHT, [3]bool{true, false, true}}:   ss.BELT_TYPE_DOWNLEFT_UPRIGHT_DOWNRIGHT,
	{utils.DIR_DOWN_LEFT, [3]bool{true, false, true}}:  ss.BELT_TYPE_UPRIGHT_DOWNLEFT_UPLEFT,
	{utils.DIR_UP_LEFT, [3]bool{true, true, false}}:    ss.BELT_TYPE_DOWNRIGHT_UPLEFT_DOWNLEFT,

	{utils.DIR_RIGHT, [3]bool{true, true, false}}: ss.BELT_TYPE_LEFT_RIGHT_UPLEFT,
	{utils.DIR_RIGHT, [3]bool{true, false, true}}: ss.BELT_TYPE_LEFT_RIGHT_DOWNLEFT,
	{utils.DIR_LEFT, [3]bool{true, false, true}}:  ss.BELT_TYPE_RIGHT_LEFT_UPRIGHT,
	{utils.DIR_LEFT, [3]bool{true, true, false}}:  ss.BELT_TYPE_RIGHT_LEFT_DOWNRIGHT,

	// Two join into one
	{utils.DIR_RIGHT, [3]bool{false, true, true}}: ss.BELT_TYPE_UPLEFT_DOWNLEFT_RIGHT,
	{utils.DIR_LEFT, [3]bool{false, true, true}}:  ss.BELT_TYPE_DOWNRIGHT_UPRIGHT_LEFT,

	{utils.DIR_DOWN_RIGHT, [3]bool{false, true, true}}: ss.BELT_TYPE_UPRIGHT_LEFT_DOWNRIGHT,
	{utils.DIR_UP_RIGHT, [3]bool{false, true, true}}:   ss.BELT_TYPE_LEFT_DOWNRIGHT_UPRIGHT,
	{utils.DIR_DOWN_LEFT, [3]bool{false, true, true}}:  ss.BELT_TYPE_RIGHT_UPLEFT_DOWNLEFT,
	{utils.DIR_UP_LEFT, [3]bool{false, true, true}}:    ss.BELT_TYPE_DOWNLEFT_RIGHT_UPLEFT,

	// Three join into one
	{utils.DIR_RIGHT, [3]bool{true, true, true}}: ss.BELT_TYPE_ALL_RIGHT,
	{utils.DIR_LEFT, [3]bool{true, true, true}}:  ss.BELT_TYPE_ALL_LEFT,

	{utils.DIR_DOWN_RIGHT, [3]bool{true, true, true}}: ss.BELT_TYPE_ALL_DOWNRIGHT,
	{utils.DIR_UP_RIGHT, [3]bool{true, true, true}}:   ss.BELT_TYPE_ALL_UPRIGHT,
	{utils.DIR_DOWN_LEFT, [3]bool{true, true, true}}:  ss.BELT_TYPE_ALL_DOWNLEFT,
	{utils.DIR_UP_LEFT, [3]bool{true, true, true}}:    ss.BELT_TYPE_ALL_UPLEFT,
}

type Belt struct {
	Pos      utils.HexCoord
	inConns  [3]*BeltConnection
	outConn  *BeltConnection
	speed    float64
	beltType ss.BeltType
}

func NewBelt(pos utils.HexCoord, dir utils.Dir, speed float64) *Belt {
	beltType, ok := beltTypeMapping[typeMappingKey{dir, [3]bool{false, false, false}}]
	if !ok {
		panic("invalid belt type")
	}
	return &Belt{
		Pos:      pos,
		outConn:  NewBeltConnection(pos, dir, speed, false),
		speed:    speed,
		beltType: beltType,
	}
}

func (b *Belt) GetPos() utils.HexCoord {
	return b.Pos
}

func (b *Belt) updateType() {
	beltType, ok := beltTypeMapping[typeMappingKey{b.outConn.Dir, [3]bool{
		b.inConns[0] != nil && b.inConns[0].Belt != nil, b.inConns[1] != nil, b.inConns[2] != nil}}]
	if !ok {
		panic("invalid belt type")
	}
	b.beltType = beltType
}

func (b *Belt) DrawGroundLevel(r *renderer.GameRenderer) {
	if !r.IsHexOnScreen(b.Pos) {
		return
	}
	r.DrawAnimatedBelt(b.Pos, b.beltType, b.speed*ss.TPS)

	// r.DrawHexCenter(b.Pos)

	// b.outConn.Draw(b.Pos, r)
	// for i := 0; i < 3; i++ {
	// 	if b.inConns[i] != nil {
	// 		b.inConns[i].Draw(b.Pos, r)
	// 	}
	// }
}

func (b *Belt) DrawOnGroundLevel(r *renderer.GameRenderer) {}

func (b *Belt) DrawItems(r *renderer.GameRenderer) {
	b.outConn.DrawItem(b.Pos, r)
	for i := 0; i < 3; i++ {
		if b.inConns[i] != nil {
			b.inConns[i].DrawItem(b.Pos, r)
		}
	}
}

func (b *Belt) Update(ticks uint64, world HexGridWorldInteractor) {}

func (b *Belt) MoveItems(ticks uint64, processed map[*BeltGraphSegment]struct{}) {
	b.outConn.MoveItems(ticks, processed)
}

func (b *Belt) Rotate(cw bool) {
	var newDir utils.Dir
	if cw {
		newDir = b.outConn.Dir.NextCW()
	} else {
		newDir = b.outConn.Dir.NextCCW()
	}

	b.changeOut(newDir, nil)
}

func (b *Belt) CanConnectTo(b2 BeltLike) bool {
	dir, err := b.Pos.DirTo(b2.GetPos())
	if err != nil {
		return false
	}

	if !b2.CanConnectIn(dir.Reverse()) {
		return false
	}

	return b.outConn.Belt != b2
}

func (b *Belt) CanConnectIn(dir utils.Dir) bool {
	return !b.outConn.Dir.IsAcute(dir)
}

func (b *Belt) ConnectTo(b2 BeltLike) {
	dir, err := b.Pos.DirTo(b2.GetPos())
	if err != nil {
		panic(fmt.Sprintf("unable to connect: %v", err))
	}

	if b.outConn.Belt == b2 {
		return
	}

	b2.CreateIn(dir.Reverse(), b)
	b.changeOut(dir, b2)
}

func (b *Belt) fixIncomingConnections() {
	bc2 := b.outConn
	isKink, isLeft := false, false

	bcCenter := b.inConns[0]
	bc := b.inConns[1]
	if bc != nil {
		if bc.Belt == nil {
			panic("impossible")
		}

		if b.inConns[0] == nil && b.inConns[2] == nil {
			bc.LaneLeft.Reconnect(&bc2.LaneLeft, 0, 0.5, 0)
			bc.LaneRight.Reconnect(&bc2.LaneRight, 0, 0.5, 0)
			bc.UpdateShifts(true, true)
			isKink, isLeft = true, true
		} else {
			bc.LaneLeft.Reconnect(&bc2.LaneLeft, 0, ss.JOIN2, 0.5-ss.JOIN2)
			bc.LaneRight.Reconnect(&bcCenter.LaneLeft, 0, 1-ss.JOIN1, 1-ss.JOIN1)
			bc.UpdateShifts(false, false)
		}
	}

	bc = b.inConns[2]
	if bc != nil {
		if bc.Belt == nil {
			panic("impossible")
		}
		if b.inConns[0] == nil && b.inConns[1] == nil {
			bc.LaneLeft.Reconnect(&bc2.LaneLeft, 0, 0.5, 0)
			bc.LaneRight.Reconnect(&bc2.LaneRight, 0, 0.5, 0)
			bc.UpdateShifts(true, false)
			isKink = true
		} else {
			bc.LaneLeft.Reconnect(&bcCenter.LaneRight, 0, 1-ss.JOIN1, 1-ss.JOIN1)
			bc.LaneRight.Reconnect(&bc2.LaneRight, 0, ss.JOIN2, 0.5-ss.JOIN2)
			bc.UpdateShifts(false, false)
		}
	}

	if bcCenter != nil {
		bcCenter.UpdateShifts(false, false)
		if bcCenter.Belt != nil {
			bcCenter.LaneLeft.Reconnect(&bc2.LaneLeft, 0, 0.5, 0)
			bcCenter.LaneRight.Reconnect(&bc2.LaneRight, 0, 0.5, 0)
		} else if b.inConns[1] != nil && b.inConns[2] != nil {
			bcCenter.LaneLeft.Reconnect(&bc2.LaneLeft, 1-ss.JOIN1, 0.5, 0)
			bcCenter.LaneRight.Reconnect(&bc2.LaneRight, 1-ss.JOIN1, 0.5, 0)
		} else {
			panic("invalid center incoming connection")
		}
	}

	b.outConn.UpdateShifts(isKink, isLeft)
}

func (b *Belt) DisconnectAll() {
	b.disconnectOut()

	for i := 0; i < 3; i++ {
		if b.inConns[i] != nil {
			b.disconnectIn(i)
		}
	}
}

func (b *Belt) disconnectOut() {
	b2 := b.outConn.Belt
	if b2 == nil {
		return
	}

	b.ClearOut(b.outConn.Dir)
	b2.ClearIn(b.outConn.Dir.Reverse())
}

func (b *Belt) disconnectIn(idx int) {
	bc := b.inConns[idx]
	if bc == nil {
		return
	}

	b2 := bc.Belt
	if b2 != nil {
		b2.ClearOut(b.inConns[idx].Dir.Reverse())
	}

	b.clearIn(idx)
}

func (b *Belt) ClearOut(dir utils.Dir) {
	if dir != b.outConn.Dir {
		panic("belt connectivity inconsistency")
	}
	b.outConn.DisconnectNext()
}

func (b *Belt) ClearIn(dir utils.Dir) {
	idx := b.getInBeltsIndex(dir)
	b.clearIn(idx)
	b.fixIncomingConnections()
	b.updateType()
}

func (b *Belt) clearIn(idx int) {
	bc := b.inConns[idx]
	if bc == nil {
		panic("belt connectivity inconsistency")
	}

	bc.Belt = nil

	if idx != 0 || b.inConns[1] == nil || b.inConns[2] == nil {
		bc.LaneLeft.DisconnectNext()
		bc.LaneLeft.DropAllItems()
		bc.LaneRight.DisconnectNext()
		bc.LaneRight.DropAllItems()
		b.inConns[idx] = nil
	}

	if idx != 0 {
		if b.inConns[0] != nil && b.inConns[0].Belt == nil {
			b.clearIn(0)
		}
	}
}

func (b *Belt) removeAcuteIns(outDir utils.Dir) {
	for i := 0; i < 3; i++ {
		if b.inConns[i] != nil && outDir.IsAcute(b.inConns[i].Dir) {
			if b.inConns[i].Belt != nil {
				b.disconnectIn(i)
			}
		}
	}
}

func (b *Belt) CreateIn(inDir utils.Dir, b2 BeltLike) {
	idx := b.getInBeltsIndex(inDir)

	if b.inConns[idx] != nil {
		if idx != 0 {
			panic("belt connectivity inconsistency")
		}
		b.inConns[idx].Belt = b2
		b.fixIncomingConnections()
		b.updateType()
		return
	}

	b.inConns[idx] = NewBeltConnection(b.Pos, inDir, b.speed, true)
	b.inConns[idx].Belt = b2

	if idx == 1 && b.inConns[2] != nil && b.inConns[0] == nil {
		b.inConns[0] = NewBeltConnection(b.Pos, b.outConn.Dir.Reverse(), b.speed, true)
	}
	if idx == 2 && b.inConns[1] != nil && b.inConns[0] == nil {
		b.inConns[0] = NewBeltConnection(b.Pos, b.outConn.Dir.Reverse(), b.speed, true)
	}

	b.fixIncomingConnections()
	b.updateType()
}

func (b *Belt) changeOut(dir utils.Dir, b2 BeltLike) {
	if b.outConn.Dir == dir && b.outConn.Belt == b2 {
		return
	}

	b.disconnectOut()
	b.removeAcuteIns(dir)

	dirs := [3]utils.Dir{b.outConn.Dir.Reverse(), b.outConn.Dir.Reverse().NextCW(), b.outConn.Dir.Reverse().NextCCW()}
	oldIns := [3]*BeltConnection{b.inConns[0], b.inConns[1], b.inConns[2]}

	b.inConns[0] = nil
	b.inConns[1] = nil
	b.inConns[2] = nil

	newDirRev := dir.Reverse()
	for i := 0; i < 3; i++ {
		if oldIns[i] != nil && dirs[i] == newDirRev {
			b.inConns[0] = oldIns[i]
			break
		}
	}

	for i := 0; i < 3; i++ {
		if oldIns[i] != nil && dirs[i] == newDirRev.NextCW() {
			b.inConns[1] = oldIns[i]
			break
		}
	}

	for i := 0; i < 3; i++ {
		if oldIns[i] != nil && dirs[i] == newDirRev.NextCCW() {
			b.inConns[2] = oldIns[i]
			break
		}
	}

	b.outConn.UpdateDir(b.Pos, dir)
	b.outConn.Belt = b2
	if b2 != nil {
		b.outConn.ReconnectTo(b2.GetInConn(b.outConn.Dir.Reverse()))
	}
	b.fixIncomingConnections()
	b.updateType()
}

func (b *Belt) GetInConn(dir utils.Dir) *BeltConnection {
	idx := b.getInBeltsIndex(dir)
	return b.inConns[idx]
}

// 0 - opposite, 1 - clockwise from opposite, 2 - counter-clockwise from opposite
func (b *Belt) getInBeltsIndex(inDir utils.Dir) int {
	outDirRev := b.outConn.Dir.Reverse()
	if inDir == outDirRev {
		return 0
	}
	if inDir == outDirRev.NextCW() {
		return 1
	}
	if inDir == outDirRev.NextCCW() {
		return 2
	}
	panic("invalid directions combination")
}

func (b *Belt) GetSegmentOffset(pos utils.WorldCoord) (dir utils.Dir, offset float64, isLeft bool, ok bool) {
	center := b.Pos.CenterToWorld()
	localPos := pos.Sub(center)

	minDistSq, minDistIsLeft, minDistOffset := b.outConn.GetClosestSegment(localPos)
	minDistDir := b.outConn.Dir

	for i := 0; i < 3; i++ {
		if b.inConns[i] == nil {
			continue
		}
		distSq, segment, offset := b.inConns[i].GetClosestSegment(localPos)
		if distSq < minDistSq {
			minDistSq, minDistIsLeft, minDistOffset = distSq, segment, offset
			minDistDir = b.inConns[i].Dir
		}
	}

	if minDistSq > 1.5*ss.ITEM_R*ss.ITEM_R {
		return utils.DIR_LEFT, 0, false, false
	}

	return minDistDir, minDistOffset, minDistIsLeft, true
}

func (b *Belt) CanPlaceItem(dir utils.Dir, isLeft bool, offset float64) bool {
	if dir == b.outConn.Dir {
		return b.outConn.CanPlaceItem(offset, isLeft)
	}
	idx := b.getInBeltsIndex(dir)
	return b.inConns[idx].CanPlaceItem(offset, isLeft)
}

func (b *Belt) TakeItemIn(pos utils.WorldCoord, item items.ItemInWorld) (ok bool) {
	dir, offset, isLeft, ok := b.GetSegmentOffset(pos)
	if !ok {
		return false
	}
	if !b.CanPlaceItem(dir, isLeft, offset) {
		return false
	}

	if dir == b.outConn.Dir {
		b.outConn.PlaceItem(item, offset, isLeft)
	} else {
		idx := b.getInBeltsIndex(dir)
		b.inConns[idx].PlaceItem(item, offset, isLeft)
	}
	return true
}

func (b *Belt) TakeItemOut(pos utils.WorldCoord) (*items.ItemInWorld, bool) {
	var closestItem *items.ItemOnBelt
	closestDistSq := 99999999999.0
	var closestConn *BeltConnection

	if iob, distSq := b.outConn.FindClosestItem(pos); iob != nil {
		closestItem = iob
		closestDistSq = distSq
		closestConn = b.outConn
	}

	for i := 0; i < 3; i++ {
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
