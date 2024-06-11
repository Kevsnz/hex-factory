package game

import (
	"hextopdown/game/items"
	"hextopdown/game/objects"
	"hextopdown/input"
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
)

const MOVE_SPEED = 8

type Game struct {
	Pos utils.WorldCoordInterpolated

	worldObjects map[utils.HexCoord]WorldObject

	tick     uint64
	time     uint64 // time of last tick
	TickTime uint64

	mousePos         utils.WorldCoord
	Running          bool
	Paused           bool
	selectedDir      utils.Dir
	preppedUnderConn [2]utils.HexCoord
	showPreppedUnder bool
}

func NewGame() *Game {
	return &Game{
		Running:      true,
		Paused:       false,
		worldObjects: make(map[utils.HexCoord]WorldObject),
	}
}

func (g *Game) Destroy() {}

func (g *Game) SetTime(t uint64) {
	g.time = t
}

func (g *Game) ProcessInputFramed(ih *input.InputHandler, r *renderer.GameRenderer) {
	for {
		actionEvent, ok := ih.KeyboardActionsFramed.Pop()
		if !ok {
			break
		}

		switch actionEvent.Type {
		case input.ACTION_TYPE_DOWN:
			if actionEvent.Action == input.ACTION_QUIT {
				g.Running = false
			}
			if actionEvent.Action == input.ACTION_PAUSE {
				g.Paused = !g.Paused
			}
		}
	}
}

func (g *Game) processInputTicked(ih *input.InputHandler) {
	g.processGameActions(ih)
	g.processMouseActions(ih)
	g.processMouseMovement(ih)
}

func (g *Game) processGameActions(ih *input.InputHandler) {
	if g.Paused {
		ih.KeyboardActionsTicked.Clear()
		return
	}

	for {
		actionEvent, ok := ih.KeyboardActionsTicked.Pop()
		if !ok {
			break
		}

		hex := utils.HexCoordFromWorld(g.mousePos)
		switch actionEvent.Type {
		case input.ACTION_TYPE_DOWN:
			switch actionEvent.Action {
			case input.ACTION_PLACE_ITEM:
				itemTaker, ok := g.GetItemInputAt(hex)
				if !ok {
					break
				}
				item := items.NewItemInWorld2(ss.ITEM_TYPE_IRON_PLATE, g.mousePos)
				_ = itemTaker.TakeItemIn(g.mousePos, item)
			case input.ACTION_ROTATE_CW:
				if t, ok := g.worldObjects[hex]; ok {
					if obj, ok := t.(DirectionalObject); ok {
						g.rotateObject(obj, true)
						break
					}
				}
				g.selectedDir = g.selectedDir.NextCW()
			case input.ACTION_ROTATE_CCW:
				if t, ok := g.worldObjects[hex]; ok {
					if obj, ok := t.(DirectionalObject); ok {
						g.rotateObject(obj, false)
						break
					}
				}
				g.selectedDir = g.selectedDir.NextCCW()
			case input.ACTION_PLOP_SPLITTER:
				if !g.canPlaceObject(hex, ss.OBJECT_TYPE_BELTSPLITTER1, g.selectedDir) {
					break
				}
				bs := objects.NewBeltSplitter(hex, g.selectedDir, ss.BELT_SPEED_TICK)
				g.placeObject(bs)
			case input.ACTION_PLOP_UNDERGROUND:
				if !g.canPlaceObject(hex, ss.OBJECT_TYPE_BELTUNDER1, g.selectedDir) {
					break
				}
				bu := g.findUnderToJoin(hex, g.selectedDir, ss.BELT_UNDER_REACH)
				if bu == nil {
					newBelt := objects.NewBeltUnder(hex, g.selectedDir, ss.BELT_SPEED_TICK, true, ss.BELT_UNDER_REACH)
					g.worldObjects[hex] = newBelt
					g.selectedDir = g.selectedDir.Reverse()
					break
				}
				var newBelt *objects.BeltUnder
				if bu.IsEntry {
					newBelt = objects.NewBeltUnder(hex, g.selectedDir.Reverse(), ss.BELT_SPEED_TICK, false, ss.BELT_UNDER_REACH)
				} else {
					newBelt = objects.NewBeltUnder(hex, g.selectedDir, ss.BELT_SPEED_TICK, true, ss.BELT_UNDER_REACH)
				}
				g.placeObject(newBelt)
				if bu.IsEntry {
					bu.JoinUnder(newBelt)
				} else {
					newBelt.JoinUnder(bu)
				}
			case input.ACTION_PLOP_INSERTER:
				if !g.canPlaceObject(hex, ss.OBJECT_TYPE_INSERTER1, g.selectedDir) {
					break
				}
				ins := objects.NewInserter(hex, g.selectedDir, ss.INSERTER_SPEED_TICK)
				g.placeObject(ins)
			case input.ACTION_PLOP_CHESTBOX_SMALL:
				if !g.canPlaceObject(hex, ss.OBJECT_TYPE_CHESTBOX_SMALL, g.selectedDir) {
					break
				}
				chest := objects.NewChestBox(hex, ss.CHESTBOX_CAPACITY_SMALL)
				g.placeObject(chest)
			case input.ACTION_PLOP_CHESTBOX_MEDIUM:
				if !g.canPlaceObject(hex, ss.OBJECT_TYPE_CHESTBOX_MEDIUM, g.selectedDir) {
					break
				}
				chest := objects.NewChestBox(hex, ss.CHESTBOX_CAPACITY_MEDIUM)
				g.placeObject(chest)
			case input.ACTION_PLOP_CHESTBOX_LARGE:
				if !g.canPlaceObject(hex, ss.OBJECT_TYPE_CHESTBOX_LARGE, g.selectedDir) {
					break
				}
				chest := objects.NewChestBox(hex, ss.CHESTBOX_CAPACITY_LARGE)
				g.placeObject(chest)
			case input.ACTION_PLOP_FURNACE:
				if !g.canPlaceObject(hex, ss.OBJECT_TYPE_FURNACE_STONE, g.selectedDir) {
					break
				}
				fur := objects.NewFurnace(hex, g.selectedDir)
				g.placeObject(fur)
			}
		case input.ACTION_TYPE_UP:
		}
	}

	sx, sy := 0.0, 0.0
	if ih.GetActionState(input.ACTION_MOVE_LEFT) {
		sx -= MOVE_SPEED
	}

	if ih.GetActionState(input.ACTION_MOVE_RIGHT) {
		sx += MOVE_SPEED
	}

	if ih.GetActionState(input.ACTION_MOVE_UP) {
		sy -= MOVE_SPEED
	}

	if ih.GetActionState(input.ACTION_MOVE_DOWN) {
		sy += MOVE_SPEED
	}

	g.Pos.UpdatePosition(g.Pos.Pos.Shift(sx, sy), false)
}

func (g *Game) processMouseActions(ih *input.InputHandler) {
	if g.Paused {
		ih.MouseActions.Clear()
		return
	}

	for {
		mouseEvent, ok := ih.MouseActions.Pop()
		if !ok {
			break
		}

		switch mouseEvent.Type {
		case input.MOUSE_BUTTON_DOWN:
			if mouseEvent.Button == input.MOUSE_BUTTON_LEFT {
				hex := utils.HexCoordFromWorld(mouseEvent.Coord)
				_ = g.placeBelt(hex, g.selectedDir)
			}
			if mouseEvent.Button == input.MOUSE_BUTTON_RIGHT {
				hex := utils.HexCoordFromWorld(mouseEvent.Coord)
				g.removeHex(hex)
			}
		}
	}
}

func (g *Game) processMouseMovement(ih *input.InputHandler) {
	if g.Paused {
		g.mousePos = ih.MousePos
		return
	}

	lastMousePos := g.mousePos
	g.mousePos = ih.MousePos

	if ih.GetMouseButtonState(input.MOUSE_BUTTON_LEFT) {
		hex1 := utils.HexCoordFromWorld(lastMousePos)
		hex2 := utils.HexCoordFromWorld(g.mousePos)
		if hex1 != hex2 {
			g.placeConnectBelts(hex1, hex2)
		}
	}

	hex := utils.HexCoordFromWorld(g.mousePos)
	bu := g.findUnderToJoin(hex, g.selectedDir, ss.BELT_UNDER_REACH)
	if bu == nil {
		g.showPreppedUnder = false
	} else {
		g.showPreppedUnder = true
		g.preppedUnderConn[0] = hex
		g.preppedUnderConn[1] = bu.Pos
	}
}

func (g *Game) Update(t uint64, ih *input.InputHandler) {
	lastTime := g.time

	if g.time+ss.TICK_DT < t {
		ss.G_currentTickTimeMs = t
		g.doTick()
		g.processInputTicked(ih)

		g.time += ss.TICK_DT
		if g.time+ss.TICK_DT < t {
			g.time = t - ss.TICK_DT
		}
		g.TickTime = g.time - lastTime
	}
}

func (g *Game) doTick() {
	if g.Paused {
		return
	}

	processedObjects := make(map[WorldObject]struct{})
	processedBgs := make(map[*objects.BeltGraphSegment]struct{})
	for _, tickable := range g.worldObjects {
		if _, ok := processedObjects[tickable]; ok {
			continue
		}

		if obj, ok := tickable.(ItemMover); ok {
			obj.MoveItems(g.tick, processedBgs)
		}
		if obj, ok := tickable.(Tickable); ok {
			obj.Update(g.tick, g)
		}

		processedObjects[tickable] = struct{}{}
	}
	g.tick++
}

func (g *Game) Draw(r *renderer.GameRenderer) {
	r.DrawScreen()

	for _, belt := range g.worldObjects {
		belt.DrawGroundLevel(r)
	}

	if g.showPreppedUnder {
		r.DrawConnectionHexes(g.preppedUnderConn[0], g.preppedUnderConn[1])
	}

	for _, obj := range g.worldObjects {
		if !r.IsHexOnScreen(obj.GetPos()) {
			continue
		}
		switch drawer := obj.(type) {
		case ItemDrawer:
			drawer.DrawItems(r)
		}
	}

	for _, belt := range g.worldObjects {
		belt.DrawOnGroundLevel(r)
	}

	r.DrawViewTarget(g.Pos)
	r.DrawArrow(0.9, 0.025, g.selectedDir)
	r.DrawPlayerCoords(g.Pos.Pos, 0.01, 0.03)

	hex := utils.HexCoordFromWorld(g.mousePos)

	if obj, ok := g.worldObjects[hex]; ok {
		objType := obj.GetObjectType()
		var items []utils.ItemInfo
		if obj, ok := obj.(ItemHolder); ok {
			items = obj.GetItemList()
		}
		r.DrawObjectDetails(objectParamsList[objType].Name, hex, items, 0.01, 0.90)
	} else {
		r.DrawHexCoords(hex, 0.01, 0.96)
	}
}

func (g *Game) placeBelt(hex utils.HexCoord, dir utils.Dir) *objects.Belt {
	if _, ok := g.worldObjects[hex]; ok {
		return nil
	}
	belt := objects.NewBelt(hex, dir, ss.BELT_SPEED_TICK)
	g.worldObjects[hex] = belt
	return belt
}

func (g *Game) canPlaceObject(hex utils.HexCoord, objType ss.ObjectType, dir utils.Dir) bool {
	objParams := objectParamsList[objType]
	hexes := objParams.Shape.GetHexes(hex, dir)
	for _, h := range hexes {
		if _, ok := g.worldObjects[h]; ok {
			return false
		}
	}
	return true
}

func (g *Game) placeObject(obj WorldObject) {
	for _, h := range getObjectHexes(obj) {
		g.worldObjects[h] = obj
	}
}

func (g *Game) removeHex(hex utils.HexCoord) {
	obj, ok := g.worldObjects[hex]
	if !ok {
		return
	}

	if belt, ok := obj.(objects.BeltLike); ok {
		belt.DisconnectAll()
	}

	for _, h := range getObjectHexes(obj) {
		delete(g.worldObjects, h)
	}
}

func (g *Game) placeConnectBelts(coord1, coord2 utils.HexCoord) {
	var belt1, belt2 objects.BeltLike
	if t, ok := g.worldObjects[coord1]; ok {
		if belt1, ok = t.(objects.BeltLike); !ok {
			return
		}
	} else {
		return
	}

	if t, ok := g.worldObjects[coord2]; ok {
		if belt2, ok = t.(objects.BeltLike); !ok {
			return
		}
	} else {
		dir, err := coord1.DirTo(coord2)
		if err != nil {
			return
		}
		newBelt := g.placeBelt(coord2, dir)
		if newBelt == nil {
			return
		}
		belt2 = newBelt
	}

	if belt1.CanConnectTo(belt2) {
		belt1.ConnectTo(belt2)
	}
}

func (g *Game) findUnderToJoin(hex utils.HexCoord, dir utils.Dir, reach int) *objects.BeltUnder {
	curHex := hex
	for i := 0; i < reach; i++ {
		curHex = curHex.Next(dir)
		tickable, ok := g.worldObjects[curHex]
		if !ok {
			continue
		}
		bu, ok := tickable.(*objects.BeltUnder)
		if !ok {
			continue
		}
		if bu.CanJoinUnder(hex, dir.Reverse()) {
			return bu
		}
	}
	return nil
}

func (g *Game) rotateObject(obj DirectionalObject, cw bool) {
	if bu, ok := obj.(*objects.BeltUnder); ok {
		if bu.JoinedBelt != nil {
			bu.Reverse()
			return
		}

		bu.Rotate(cw)

		bu2 := g.findUnderToJoin(bu.GetPos(), bu.GetDir(), bu.Reach)
		if bu2 == nil || bu.IsEntry == bu2.IsEntry { // TODO Also check belt tier!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			return
		}
		if bu2.IsEntry {
			bu2.JoinUnder(bu)
		} else {
			bu.JoinUnder(bu2)
		}
		return

	}

	obj.Rotate(cw)
}

func (g *Game) GetItemInputAt(hex utils.HexCoord) (obj objects.ItemInput, ok bool) {
	if obj, ok := g.worldObjects[hex]; ok {
		if ii, ok := obj.(objects.ItemInput); ok {
			return ii, true
		}
	}
	return nil, false
}

func (g *Game) GetItemOutputAt(hex utils.HexCoord) (obj objects.ItemOutput, ok bool) {
	if obj, ok := g.worldObjects[hex]; ok {
		if io, ok := obj.(objects.ItemOutput); ok {
			return io, true
		}
	}
	return nil, false
}

func getObjectHexes(obj WorldObject) []utils.HexCoord {
	objParams := objectParamsList[obj.GetObjectType()]
	dir := utils.DIR_LEFT
	if do, ok := obj.(DirectionalObject); ok {
		dir = do.GetDir()
	}
	return objParams.Shape.GetHexes(obj.GetPos(), dir)
}
