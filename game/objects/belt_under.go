package objects

import (
	gd "hextopdown/game/gamedata"
	"hextopdown/game/items"
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
)

var underBeltTypeMapping = [utils.DIR_COUNT]ss.BeltType{
	utils.DIR_LEFT:       ss.BELT_TYPE_IN_LEFT,
	utils.DIR_UP_LEFT:    ss.BELT_TYPE_IN_UPLEFT,
	utils.DIR_UP_RIGHT:   ss.BELT_TYPE_IN_UPRIGHT,
	utils.DIR_RIGHT:      ss.BELT_TYPE_IN_RIGHT,
	utils.DIR_DOWN_RIGHT: ss.BELT_TYPE_IN_DOWNRIGHT,
	utils.DIR_DOWN_LEFT:  ss.BELT_TYPE_IN_DOWNLEFT,
}

type onBeltUnderTypeMappingKey struct {
	dir     utils.Dir
	isEntry bool
}

var underBeltOnTypeMapping = map[onBeltUnderTypeMappingKey]ss.BeltType{
	{dir: utils.DIR_LEFT, isEntry: true}:       ss.BELT_ON_UNDER_IN_LEFT,
	{dir: utils.DIR_UP_LEFT, isEntry: true}:    ss.BELT_ON_UNDER_IN_UPLEFT,
	{dir: utils.DIR_UP_RIGHT, isEntry: true}:   ss.BELT_ON_UNDER_IN_UPRIGHT,
	{dir: utils.DIR_RIGHT, isEntry: true}:      ss.BELT_ON_UNDER_IN_RIGHT,
	{dir: utils.DIR_DOWN_RIGHT, isEntry: true}: ss.BELT_ON_UNDER_IN_DOWNRIGHT,
	{dir: utils.DIR_DOWN_LEFT, isEntry: true}:  ss.BELT_ON_UNDER_IN_DOWNLEFT,

	{dir: utils.DIR_LEFT, isEntry: false}:       ss.BELT_ON_UNDER_OUT_LEFT,
	{dir: utils.DIR_UP_LEFT, isEntry: false}:    ss.BELT_ON_UNDER_OUT_UPLEFT,
	{dir: utils.DIR_UP_RIGHT, isEntry: false}:   ss.BELT_ON_UNDER_OUT_UPRIGHT,
	{dir: utils.DIR_RIGHT, isEntry: false}:      ss.BELT_ON_UNDER_OUT_RIGHT,
	{dir: utils.DIR_DOWN_RIGHT, isEntry: false}: ss.BELT_ON_UNDER_OUT_DOWNRIGHT,
	{dir: utils.DIR_DOWN_LEFT, isEntry: false}:  ss.BELT_ON_UNDER_OUT_DOWNLEFT,
}

type BeltUnder struct {
	ObjectBeltlike
	speed      float64
	beltType   ss.BeltType
	onType     ss.BeltType
	inConn     *BeltConnection
	outConn    *BeltConnection
	JoinedBelt *BeltUnder
	IsEntry    bool
	Reach      int32
}

func NewBeltUnder(
	objType ss.ObjectType,
	pos utils.HexCoord,
	dir utils.Dir,
	objParams *gd.ObjectParameters,
	tier ss.BeltTier,
	isEntry bool,
) *BeltUnder {
	speed := 1 / float64(gd.BeltTierParamsList[tier].Speed)

	var inConn, outConn *BeltConnection
	if isEntry {
		inConn = NewBeltConnection(pos, dir.Reverse(), speed, true)
	} else {
		outConn = NewBeltConnection(pos, dir, speed, false)
	}

	newBelt := &BeltUnder{
		ObjectBeltlike: ObjectBeltlike{
			Object: Object{
				objType:   objType,
				pos:       pos,
				objParams: objParams,
			},
			tier:       tier,
			tierParams: &gd.BeltTierParamsList[tier],
		},
		speed:   speed,
		IsEntry: isEntry,
		inConn:  inConn,
		outConn: outConn,
		Reach:   gd.BeltTierParamsList[tier].Reach,
	}

	newBelt.setBeltType()
	return newBelt
}

func (b *BeltUnder) setBeltType() {
	if b.IsEntry {
		b.beltType = underBeltTypeMapping[b.inConn.Dir.Reverse()]
		b.onType = underBeltOnTypeMapping[onBeltUnderTypeMappingKey{b.inConn.Dir.Reverse(), b.IsEntry}]
	} else {
		b.beltType = beltTypeMapping[typeMappingKey{b.outConn.Dir, [3]bool{false, false, false}}]
		b.onType = underBeltOnTypeMapping[onBeltUnderTypeMappingKey{b.outConn.Dir, b.IsEntry}]
	}
}

func (b *BeltUnder) GetDir() utils.Dir {
	if b.IsEntry {
		return b.inConn.Dir.Reverse()
	}
	return b.outConn.Dir
}

func (b *BeltUnder) GetInConn(dir utils.Dir) *BeltConnection {
	if !b.IsEntry {
		panic("not an entry belt")
	}
	return b.inConn
}

func (b *BeltUnder) Update(ticks uint64, world HexGridWorldInteractor) {}

func (b *BeltUnder) DrawGroundLevel(r *renderer.GameRenderer) {
	r.DrawAnimatedBelt(b.pos, b.beltType, b.speed*ss.TPS)

	// if b.IsEntry {
	// 	b.inConn.Draw(b.pos, r)
	// } else {
	// 	b.outConn.Draw(b.pos, r)
	// }
}

func (b *BeltUnder) DrawOnGroundLevel(r *renderer.GameRenderer) {
	r.DrawBeltOnGround(b.pos, b.onType)
}

func (b *BeltUnder) DrawItems(r *renderer.GameRenderer) {
	if b.IsEntry {
		b.inConn.DrawItem(b.pos, r)
		return
	}

	// if b.inConn != nil {
	// 	b.inConn.DrawItem(b.Pos, r)
	// }
	b.outConn.DrawItem(b.pos, r)
}

func (b *BeltUnder) Reverse() {
	if !b.IsEntry {
		if b.JoinedBelt != nil {
			b.JoinedBelt.Reverse()
			return
		}
		if b.inConn != nil {
			panic("belt connectivity inconsistency")
		}
		b.outConn.DisconnectNext()
		b.outConn.Reverse(b.pos)
		b.inConn, b.outConn = b.outConn, nil
		return
	}

	oldInConn := b.inConn
	midConn := b.outConn

	if oldInConn.Belt != nil {
		oldInConn.Belt.ClearOut(oldInConn.Dir.Reverse())
		b.ClearIn(oldInConn.Dir)
	}
	oldInConn.DisconnectNext()
	oldInConn.Reverse(b.pos)
	b.outConn = b.inConn
	b.IsEntry = false
	b.setBeltType()

	b2 := b.JoinedBelt
	if b2 != nil {
		oldb2OutConn := b2.outConn

		dist := b.pos.DistanceTo(b2.GetPos())
		midConn.DisconnectNext()
		midConn.UpdateDir(b2.pos, midConn.Dir.Reverse())
		midConn.ReconnectToWithDist(oldInConn, dist)
		midConn.SwapReverseItems()

		if oldb2OutConn.Belt != nil {
			oldb2OutConn.Belt.ClearIn(oldb2OutConn.Dir.Reverse())
			b2.ClearOut(oldb2OutConn.Dir)
		}
		oldb2OutConn.DisconnectNext()
		oldb2OutConn.Reverse(b2.pos)
		oldb2OutConn.ReconnectTo(midConn)

		b.inConn = midConn
		b2.outConn = midConn
		b2.inConn = oldb2OutConn
		b2.IsEntry = true
		b2.setBeltType()
	}
}

func (b *BeltUnder) Rotate(cw bool) {
	if b.JoinedBelt != nil {
		panic("rotation of joined underground belt is forbidden")
	}

	if b.IsEntry {
		var newDir utils.Dir
		if cw {
			newDir = b.inConn.Dir.NextCW()
		} else {
			newDir = b.inConn.Dir.NextCCW()
		}

		if b.inConn.Belt != nil {
			b.inConn.Belt.ClearOut(b.inConn.Dir.Reverse())
			b.ClearIn(b.inConn.Dir)
		}
		b.inConn.UpdateDir(b.pos, newDir)
		b.setBeltType()
		return
	}

	var newDir utils.Dir
	if cw {
		newDir = b.outConn.Dir.NextCW()
	} else {
		newDir = b.outConn.Dir.NextCCW()
	}

	if b.outConn.Belt != nil {
		b.outConn.Belt.ClearIn(b.outConn.Dir.Reverse())
		b.ClearOut(b.outConn.Dir)
	}
	b.outConn.UpdateDir(b.pos, newDir)
	b.setBeltType()
}

func (b *BeltUnder) CanConnectTo(b2 BeltLike) bool {
	if b.IsEntry {
		return false
	}

	dir, err := b.pos.DirTo(b2.GetPos())
	if err != nil {
		return false
	}

	if dir != b.outConn.Dir || !b2.CanConnectIn(dir.Reverse()) || b.outConn.Belt == b2 {
		return false
	}

	if b.outConn.Belt != nil {
		panic("belt connectivity inconsistency")
	}
	return true
}

func (b *BeltUnder) CanConnectIn(dir utils.Dir) bool {
	return b.IsEntry && dir == b.inConn.Dir
}

func (b *BeltUnder) ConnectTo(b2 BeltLike) {
	if b.IsEntry {
		panic("it's an entry belt")
	}

	dir, err := b.pos.DirTo(b2.GetPos())
	if err != nil {
		panic(err)
	}

	b.outConn.Belt = b2
	b2.CreateIn(dir.Reverse(), b)
	b.outConn.ReconnectTo(b2.GetInConn(dir.Reverse()))
}

func (b *BeltUnder) ClearIn(dir utils.Dir) {
	if !b.IsEntry {
		panic("not an entry belt")
	}
	b.inConn.Belt = nil
}

func (b *BeltUnder) ClearOut(dir utils.Dir) {
	if b.IsEntry {
		panic("not an exit belt")
	}
	b.outConn.DisconnectNext()
}

func (b *BeltUnder) CreateIn(inDir utils.Dir, b2 BeltLike) {
	if !b.IsEntry {
		panic("not an entry belt")
	}
	b.inConn.Belt = b2
}

func (b *BeltUnder) DisconnectAll() {
	if b.IsEntry {
		if b.inConn.Belt != nil {
			b.inConn.Belt.ClearOut(b.inConn.Dir.Reverse())
			b.ClearIn(b.inConn.Dir)
		}

		b.DisjoinUnder()
		return
	}

	if b.outConn.Belt != nil {
		b.outConn.Belt.ClearIn(b.outConn.Dir.Reverse())
		b.ClearOut(b.outConn.Dir)
	}

	if b.JoinedBelt != nil {
		b.JoinedBelt.DisjoinUnder()
	}
}

func (b *BeltUnder) CanJoinUnder(hex utils.HexCoord, dir utils.Dir) bool {
	if !b.pos.IsStraightTo(hex) {
		return false
	}

	if b.IsEntry && dir != b.inConn.Dir.Reverse() {
		return false
	}
	if !b.IsEntry && dir != b.outConn.Dir.Reverse() {
		return false
	}

	d := b.pos.DistanceTo(hex)
	if d+1 > b.Reach {
		return false
	}

	if b.IsEntry && b.outConn != nil {
		d2 := b.pos.DistanceTo(b.JoinedBelt.GetPos())
		if d >= d2 {
			return false
		}
	}
	if !b.IsEntry && b.inConn != nil {
		d2 := b.pos.DistanceTo(b.JoinedBelt.GetPos())
		if d >= d2 {
			return false
		}
	}
	return true
}

func (b *BeltUnder) JoinUnder(b2 *BeltUnder) {
	if !b.IsEntry {
		panic("not an entry belt")
	}

	dist := b.pos.DistanceTo(b2.pos)
	b.outConn = NewBeltConnectionWithDist(b.pos, b.inConn.Dir.Reverse(), b.speed, false, float64(dist))
	b.outConn.Belt = b2
	b.outConn.ReconnectToWithDist(b2.outConn, dist)
	b.inConn.ReconnectTo(b.outConn)
	b.JoinedBelt = b2
	b2.inConn = b.outConn
	b2.JoinedBelt = b
}

func (b *BeltUnder) DisjoinUnder() {
	if !b.IsEntry {
		panic("not an entry belt")
	}

	b2 := b.JoinedBelt
	if b2 == nil {
		return
	}

	b.inConn.DisconnectNext()
	b.outConn.DisconnectNext()
	b.outConn = nil
	b.JoinedBelt = nil

	b2.inConn = nil
	b2.JoinedBelt = nil
}

func (b *BeltUnder) MoveItems(ticks uint64, processed map[*BeltGraphSegment]struct{}) {
	if b.IsEntry {
		b.inConn.MoveItems(ticks, processed)
	} else {
		b.outConn.MoveItems(ticks, processed)
	}
}

func (b *BeltUnder) TakeItemOut(pos utils.WorldCoord) (*items.ItemInWorld, bool) {
	var closestItem *items.ItemOnBelt
	var closestConn *BeltConnection

	if b.IsEntry {
		if iob, _ := b.inConn.FindClosestItem(pos); iob != nil {
			closestItem = iob
			closestConn = b.inConn
		}
	} else {
		if iob, _ := b.outConn.FindClosestItem(pos); iob != nil {
			closestItem = iob
			closestConn = b.outConn
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
