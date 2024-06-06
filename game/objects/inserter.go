package objects

import (
	"hextopdown/game/items"
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
	"math"
)

var typeMapping = [utils.DIR_COUNT]ss.StructureType{
	utils.DIR_LEFT:       ss.STRUCTURE_TYPE_INSERTER_LEFT,
	utils.DIR_UP_LEFT:    ss.STRUCTURE_TYPE_INSERTER_UPLEFT,
	utils.DIR_UP_RIGHT:   ss.STRUCTURE_TYPE_INSERTER_UPRIGHT,
	utils.DIR_RIGHT:      ss.STRUCTURE_TYPE_INSERTER_RIGHT,
	utils.DIR_DOWN_RIGHT: ss.STRUCTURE_TYPE_INSERTER_DOWNRIGHT,
	utils.DIR_DOWN_LEFT:  ss.STRUCTURE_TYPE_INSERTER_DOWNLEFT,
}

type Inserter struct {
	pos          utils.HexCoord
	dir          utils.Dir
	speed        float64
	inserterType ss.StructureType
	armAngle     float64
	itemOnHand   *items.ItemInWorld
}

func NewInserter(pos utils.HexCoord, dir utils.Dir, speed float64) *Inserter {
	return &Inserter{
		pos:          pos,
		dir:          dir,
		speed:        speed,
		inserterType: typeMapping[dir],
		armAngle:     math.Pi / 2,
	}
}

func (i *Inserter) GetPos() utils.HexCoord {
	return i.pos
}

func (i *Inserter) Update(ticks uint64, world HexGridWorldInteractor) {
	if i.itemOnHand == nil {
		i.armAngle = i.armAngle - i.speed
		if i.armAngle <= 0 {
			i.armAngle = 0
			otherPos := i.pos.Shift(i.dir.Reverse(), 1)
			obj, ok := world.GetItemOutputAt(otherPos)
			if !ok {
				return
			}
			item, ok := obj.TakeItemOut(otherPos.CenterToWorld().Shift(i.dir.Reverse(), ss.LANE_OFFSET_WORLD))
			if !ok {
				return
			}
			i.itemOnHand = item
			i.updateHeldItemPosition()
		}
		return
	}

	i.armAngle = i.armAngle + i.speed
	if i.armAngle >= math.Pi {
		i.armAngle = math.Pi
		otherPos := i.pos.Shift(i.dir, 1)
		obj, ok := world.GetItemInputAt(otherPos)
		if !ok {
			return
		}
		ok = obj.TakeItemIn(otherPos.CenterToWorld().Shift(i.dir, ss.LANE_OFFSET_WORLD), *i.itemOnHand)
		if ok {
			i.itemOnHand = nil
			return
		}
	}

	i.updateHeldItemPosition()
}

func (i *Inserter) DrawGroundLevel(r *renderer.GameRenderer) {
	r.DrawStructure(i.pos, i.inserterType, i.dir)
}

func (i *Inserter) DrawOnGroundLevel(r *renderer.GameRenderer) {
	p1 := i.pos.CenterToWorld()
	p2 := p1.Shift(i.dir, -ss.INSERTER_ARM_LENGTH*math.Cos(i.armAngle))
	p2.Y -= ss.INSERTER_ARM_LENGTH / 2 * math.Sin(i.armAngle)
	p1.Y -= ss.HEX_EDGE / 7

	r.DrawWorldLine(p1, p2)
}

func (i *Inserter) DrawItems(r *renderer.GameRenderer) {
	if i.itemOnHand == nil {
		return
	}
	i.itemOnHand.Draw(r)
}

func (i *Inserter) Rotate(cw bool) {
	i.dir = i.dir.Next(cw)
	i.inserterType = typeMapping[i.dir]
}

func (i *Inserter) updateHeldItemPosition() {
	if i.itemOnHand == nil {
		return
	}
	itemPos := i.pos.CenterToWorld().Shift(i.dir, -ss.INSERTER_ARM_LENGTH*math.Cos(i.armAngle))
	itemPos.Y -= ss.INSERTER_ARM_LENGTH / 2 * math.Sin(i.armAngle)
	i.itemOnHand.Pos.UpdatePosition(itemPos, false)
}
