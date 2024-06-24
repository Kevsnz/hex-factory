package objects

import (
	gd "hextopdown/game/gamedata"
	"hextopdown/game/items"
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
	"math"
)

type Inserter struct {
	Object
	dir        utils.Dir
	params     *gd.InserterParameters
	armPos     uint32
	itemOnHand *items.ItemInWorld
}

func NewInserter(
	objType ss.ObjectType,
	pos utils.HexCoord,
	dir utils.Dir,
	objParams *gd.ObjectParameters,
	params *gd.InserterParameters,
) *Inserter {
	return &Inserter{
		Object: Object{
			objType:   objType,
			pos:       pos,
			objParams: objParams,
		},
		dir:    dir,
		params: params,
		armPos: params.SwingSpeed / 2,
	}
}

func (i *Inserter) GetDir() utils.Dir {
	return i.dir
}

func (i *Inserter) Update(ticks uint64, world HexGridWorldInteractor) {
	if i.itemOnHand == nil {
		if i.armPos > 0 {
			i.armPos--
			return
		}

		i.armPos = 0
		otherPos := i.pos.Shift(i.dir.Reverse(), int(i.params.Reach))
		obj, ok := world.GetItemOutputAt(otherPos)
		if !ok {
			return
		}
		item, ok := obj.TakeItemOut(
			otherPos.CenterToWorld().ShiftDir(i.dir.Reverse(), ss.LANE_OFFSET_WORLD),
			i.getAllowedItems(world),
		)
		if !ok {
			return
		}
		i.itemOnHand = item
		i.updateHeldItemPosition()
		return
	}

	if i.armPos < i.params.SwingSpeed {
		i.armPos++
		i.updateHeldItemPosition()
		return
	}

	i.armPos = i.params.SwingSpeed
	otherPos := i.pos.Shift(i.dir, int(i.params.Reach))
	obj, ok := world.GetItemInputAt(otherPos)
	if !ok {
		return
	}
	ok = obj.TakeItemIn(otherPos.CenterToWorld().ShiftDir(i.dir, ss.LANE_OFFSET_WORLD), *i.itemOnHand)
	if ok {
		i.itemOnHand = nil
		return
	}

	if !isItemAllowed(i.itemOnHand.ItemType, i.getAllowedItems(world)) {
		// TODO Drop held items!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		i.itemOnHand = nil
	}
}

func (i *Inserter) DrawGroundLevel(r *renderer.GameRenderer) {
	r.DrawObjectGround(i.pos.CenterToWorld(), i.objType, i.objParams.Shape, i.dir)
}

func (i *Inserter) DrawOnGroundLevel(r *renderer.GameRenderer) {
	if i.itemOnHand != nil {
		i.itemOnHand.Draw(r)
	}

	angle := math.Pi * float64(i.armPos) / float64(i.params.SwingSpeed)
	armLen := ss.INSERTER_ARM_LENGTH * float64(i.params.Reach)

	p1 := i.pos.CenterToWorld()
	p2 := p1.ShiftDir(i.dir, -armLen*math.Cos(angle))
	p2.Y -= armLen / 2 * math.Sin(angle)
	p1.Y -= ss.HEX_EDGE / 7

	r.DrawWorldLine(p1, p2)
}

func (i *Inserter) Rotate(cw bool) {
	i.dir = i.dir.Next(cw)
}

func (i *Inserter) updateHeldItemPosition() {
	if i.itemOnHand == nil {
		return
	}
	angle := math.Pi * float64(i.armPos) / float64(i.params.SwingSpeed)
	armLen := ss.INSERTER_ARM_LENGTH * float64(i.params.Reach)

	itemPos := i.pos.CenterToWorld().ShiftDir(i.dir, -armLen*math.Cos(angle))
	itemPos.Y -= armLen / 2 * math.Sin(angle)
	i.itemOnHand.Pos.UpdatePosition(itemPos, false)
}

func (i *Inserter) getAllowedItems(world HexGridWorldInteractor) []ss.ItemType {
	otherPos := i.pos.Shift(i.dir, int(i.params.Reach))
	obj, ok := world.GetItemInputAt(otherPos)
	if !ok {
		return nil
	}

	return obj.GetAcceptableItems()
}

func isItemAllowed(iType ss.ItemType, allowedItems []ss.ItemType) bool {
	if allowedItems == nil {
		return true
	}
	for _, item := range allowedItems {
		if iType == item {
			return true
		}
	}
	return false
}
