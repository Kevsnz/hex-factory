package game

import (
	"hextopdown/game/char"
	gd "hextopdown/game/gamedata"
	"hextopdown/game/items"
	"hextopdown/game/objects"
	"hextopdown/game/ui"
	"hextopdown/game/world"
	"hextopdown/input"
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
	"math/rand"
)

type Game struct {
	masterRNG *rand.Rand
	player    char.Character
	ui        *ui.UI

	chunks   map[utils.ChunkCoord]*world.Chunk
	worldGen *world.WorldGen

	tick     uint64
	time     uint64 // time of last tick
	TickTime uint64

	mousePos         utils.ScreenCoord
	mousePosLast     utils.WorldCoord
	Running          bool
	paused           bool
	selectedObjType  ss.ObjectType
	selectedDir      utils.Dir
	preppedUnderConn [2]utils.HexCoord
	showPreppedUnder bool
}

func NewGame(seed int64) *Game {
	mrng := rand.New(rand.NewSource(seed))
	return &Game{
		masterRNG:       mrng,
		player:          char.NewCharacter(utils.WorldCoord{X: 450, Y: 450}),
		ui:              ui.NewUI(),
		Running:         true,
		paused:          false,
		worldGen:        world.NewWorldGen(int64(mrng.Uint64())),
		chunks:          make(map[utils.ChunkCoord]*world.Chunk),
		selectedObjType: ss.OBJECT_TYPE_COUNT,
	}
}

func (g *Game) Destroy() {
	g.ui.Destroy()
}

func (g *Game) GetPlayerPos() utils.WorldCoord {
	return g.player.GetPos().Pos
}

func (g *Game) SetTime(t uint64) {
	g.time = t
}

func (g *Game) ProcessInputFramed(ih *input.InputHandler) {
	for {
		actionEvent, ok := ih.KeyboardActionsFramed.Pop()
		if !ok {
			break
		}

		switch actionEvent.Type {
		case input.ACTION_TYPE_DOWN:
			switch actionEvent.Action {
			case input.ACTION_QUIT:
				g.Running = false
			case input.ACTION_PAUSE:
				g.paused = !g.paused
			case input.ACTION_TOGGLE_UI:
				g.ui.ShowToggle()
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
	if g.paused {
		ih.KeyboardActionsTicked.Clear()
		return
	}

	for {
		actionEvent, ok := ih.KeyboardActionsTicked.Pop()
		if !ok {
			break
		}

		if g.ui.HandleGameAction(actionEvent) {
			continue
		}

		hex := utils.HexCoordFromWorld(g.mousePos.ToWorld())
		switch actionEvent.Type {
		case input.ACTION_TYPE_DOWN:
			switch actionEvent.Action {
			case input.ACTION_OPEN_INVENTORY:
				g.ui.ShowInventoryWindow(g.player.GetInventory())
			case input.ACTION_PLACE_ITEM:
				if g.selectedObjType != ss.OBJECT_TYPE_COUNT {
					break
				}

				itemTaker, ok := g.GetItemInputAt(hex)
				if !ok {
					break
				}
				item := items.NewItemInWorld2(ss.ITEM_TYPE_IRON_ORE, g.mousePos.ToWorld())
				_ = itemTaker.TakeItemIn(g.mousePos.ToWorld(), item)
			case input.ACTION_ROTATE_CW:
				if g.selectedObjType != ss.OBJECT_TYPE_COUNT {
					g.selectedDir = g.selectedDir.NextCW()
					break
				}
				if t, ok := g.getWorldObject(hex); ok {
					if obj, ok := t.(DirectionalObject); ok {
						g.rotateObject(obj, true)
					}
				}
			case input.ACTION_ROTATE_CCW:
				if g.selectedObjType != ss.OBJECT_TYPE_COUNT {
					g.selectedDir = g.selectedDir.NextCCW()
					break
				}
				if t, ok := g.getWorldObject(hex); ok {
					if obj, ok := t.(DirectionalObject); ok {
						g.rotateObject(obj, false)
					}
				}
			case input.ACTION_SELECT_TOOL_1:
				g.selectObjType(ss.OBJECT_TYPE_BELT1)
			case input.ACTION_SELECT_TOOL_2:
				g.selectObjType(ss.OBJECT_TYPE_BELTSPLITTER1)
			case input.ACTION_SELECT_TOOL_3:
				g.selectObjType(ss.OBJECT_TYPE_BELTUNDER1)
			case input.ACTION_SELECT_TOOL_4:
				g.selectObjType(ss.OBJECT_TYPE_INSERTER1)
			case input.ACTION_SELECT_TOOL_5:
				g.selectObjType(ss.OBJECT_TYPE_CHESTBOX_SMALL)
			case input.ACTION_SELECT_TOOL_6:
				g.selectObjType(ss.OBJECT_TYPE_CHESTBOX_LARGE)
			case input.ACTION_SELECT_TOOL_7:
				g.selectObjType(ss.OBJECT_TYPE_MINER_STIRLING)
			case input.ACTION_SELECT_TOOL_8:
				g.selectObjType(ss.OBJECT_TYPE_FURNACE_STONE)
			case input.ACTION_SELECT_TOOL_9:
				g.selectObjType(ss.OBJECT_TYPE_ASSEMBLER_BASIC)
			case input.ACTION_CANCEL:
				if g.selectedObjType != ss.OBJECT_TYPE_COUNT {
					g.selectedObjType = ss.OBJECT_TYPE_COUNT
					break
				}
				g.Running = false
			}
		case input.ACTION_TYPE_UP:
		}
	}

	{
		dx, dy := int64(0), int64(0)
		if ih.GetActionState(input.ACTION_MOVE_LEFT) {
			dx -= 1
		}

		if ih.GetActionState(input.ACTION_MOVE_RIGHT) {
			dx += 1
		}

		if ih.GetActionState(input.ACTION_MOVE_UP) {
			dy -= 1
		}

		if ih.GetActionState(input.ACTION_MOVE_DOWN) {
			dy += 1
		}

		g.player.UpdateMovement(dx, dy)
	}
}

func (g *Game) processMouseActions(ih *input.InputHandler) {
	if g.paused {
		ih.MouseActions.Clear()
		return
	}

	for {
		mouseEvent, ok := ih.MouseActions.Pop()
		if !ok {
			break
		}

		if g.ui.HandleMouseAction(mouseEvent) {
			continue
		}

		hex := utils.HexCoordFromWorld(mouseEvent.Coord.ToWorld())
		switch mouseEvent.Type {
		case input.MOUSE_BUTTON_DOWN:
			switch mouseEvent.Button {
			case input.MOUSE_BUTTON_LEFT:
				if obj, ok := g.getWorldObject(hex); ok {
					g.interactWithWorldObject(obj)
				} else if g.selectedObjType != ss.OBJECT_TYPE_COUNT {
					g.useSelectedTool(hex)
				}
			case input.MOUSE_BUTTON_RIGHT:
				g.removeObjectAtHex(hex)
			}
		}
	}
}

func (g *Game) processMouseMovement(ih *input.InputHandler) {
	if g.paused {
		g.mousePos = ih.MousePos
		return
	}

	if ih.MousePos != g.mousePos {
		g.ui.HandleMouseMovement(ih.MousePos)
	}

	g.mousePos = ih.MousePos

	if ih.GetMouseButtonState(input.MOUSE_BUTTON_LEFT) {
		hex1 := utils.HexCoordFromWorld(g.mousePosLast)
		hex2 := utils.HexCoordFromWorld(g.mousePos.ToWorld())

		// TODO Switch to BELTLIKE
		if g.selectedObjType == ss.OBJECT_TYPE_BELT1 && hex1 != hex2 {
			g.placeConnectBelts(hex1, hex2, ss.OBJECT_TYPE_BELT1)
		}
	}

	g.mousePosLast = g.mousePos.ToWorld()

	if g.selectedObjType == ss.OBJECT_TYPE_COUNT || gd.ObjectParamsList[g.selectedObjType].BaseType != ss.STRUCTURE_BASETYPE_BELTLIKE {
		g.showPreppedUnder = false
		return
	}

	hex := utils.HexCoordFromWorld(g.mousePos.ToWorld())
	tier := gd.BeltlikeParamsList[g.selectedObjType].Tier
	bu := g.findUnderToJoin(hex, g.selectedDir, gd.BeltTierParamsList[tier].Reach)
	if bu == nil {
		g.showPreppedUnder = false
	} else {
		g.showPreppedUnder = true
		g.preppedUnderConn[0] = hex
		g.preppedUnderConn[1] = bu.GetPos()
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
	if g.paused {
		return
	}

	processedObjects := make(map[world.WorldObject]struct{})
	processedBgs := make(map[*objects.BeltGraphSegment]struct{})
	for _, ch := range g.chunks {
		for _, obj := range ch.GetWorldObjects() {
			if _, ok := processedObjects[obj]; ok {
				continue
			}

			if obj, ok := obj.(ItemMover); ok {
				obj.MoveItems(g.tick, processedBgs)
			}
			if obj, ok := obj.(Tickable); ok {
				obj.Update(g.tick, g)
			}

			processedObjects[obj] = struct{}{}
		}
	}
	g.tick++
}

func (g *Game) Draw(r *renderer.GameRenderer) {
	ctl, cbr := r.ChunkRenderer.GetVisibleChunkCoords()
	chunks := make([]*world.Chunk, 0, (cbr.X-ctl.X+1)*(cbr.Y-ctl.Y+1))

	for y := ctl.Y; y <= cbr.Y; y++ {
		for x := ctl.X; x <= cbr.X; x++ {
			coord := utils.ChunkCoord{X: x, Y: y}
			ch, ok := g.chunks[coord]
			if !ok {
				ch = world.NewChunk(coord, g.worldGen)
				g.chunks[coord] = ch
			}
			chunks = append(chunks, ch)
		}
	}

	for _, ch := range chunks {
		ch.DrawGround(r)
	}

	for _, ch := range chunks {
		ch.DrawObjectsGroundLevel(r)
	}

	if g.showPreppedUnder {
		r.DrawConnectionHexes(g.preppedUnderConn[0], g.preppedUnderConn[1])
	}

	for _, ch := range chunks {
		ch.DrawItems(r)
	}

	for _, ch := range chunks {
		ch.DrawObjectsOnGroundLevel(r)
	}

	// Draw Player
	r.DrawViewTarget(g.player.GetPos())

	// Draw UI
	g.ui.Draw(r)
	r.DrawArrowPct(0.9, 0.025, g.selectedDir)

	hex := utils.HexCoordFromWorld(g.mousePos.ToWorld())

	if obj, ok := g.getWorldObject(hex); ok {
		objType := obj.GetObjectType()
		var items []utils.ItemInfo
		if obj, ok := obj.(ItemHolder); ok {
			items = obj.GetItemList()
		}
		r.DrawObjectDetails(gd.ObjectParamsList[objType].Name, hex, items, ss.FONT_SIZE_PCT/3, 1-ss.FONT_SIZE_PCT*3.65)
	} else {
		r.DrawHexCoords(hex, ss.FONT_SIZE_PCT/3, 1-ss.FONT_SIZE_PCT*3.45)
		r.DrawHexCoords(utils.HexCoord(hex.GetChunkCoord()), ss.FONT_SIZE_PCT/3, 1-ss.FONT_SIZE_PCT*2.3)

		resType, amt, ok := g.getChunk(hex).GetResourceTypeAt(hex)
		if ok {
			r.DrawResourceAmount(resType, amt, ss.FONT_SIZE_PCT/3, 1-ss.FONT_SIZE_PCT*1.15)
		}
	}

	if g.selectedObjType != ss.OBJECT_TYPE_COUNT {
		r.DrawCurrentTool(gd.ObjectParamsList[g.selectedObjType].Name, 1-ss.FONT_SIZE_PCT/3, 0.1)
	}
}

func (g *Game) useSelectedTool(hex utils.HexCoord) {
	switch gd.ObjectParamsList[g.selectedObjType].BaseType {
	case ss.STRUCTURE_BASETYPE_BELTLIKE:
		_ = g.placeBeltlike(g.selectedObjType, hex, g.selectedDir)
		return
	case ss.STRUCTURE_BASETYPE_INSERTER:
		_ = g.placeInserter(hex, g.selectedDir, g.selectedObjType)
		return
	case ss.STRUCTURE_BASETYPE_STORAGE:
		_ = g.placeStorage(hex, g.selectedObjType)
		return
	case ss.STRUCTURE_BASETYPE_CONVERTER:
		_ = g.placeConverter(hex, g.selectedDir, g.selectedObjType)
		return
	case ss.STRUCTURE_BASETYPE_EXTRACTOR:
		_ = g.placeExtractor(hex, g.selectedDir, g.selectedObjType)
		return
	}
}

func (g *Game) interactWithWorldObject(wo world.WorldObject) {
	switch obj := wo.(type) {
	case *objects.Converter:
		if obj.RecipeChangeable() && !obj.HasRecipe() {
			g.ui.ShowRecipeWindow([]ss.Recipe{ss.RECIPE_IRON_GEAR}, func(r ss.Recipe) { obj.ChangeRecipe(r) })
			break
		}
		g.ui.ShowConverterWindow(gd.ObjectParamsList[obj.GetObjectType()].Name, g.player.GetInventory(), obj)
	case *objects.Storage:
		g.ui.ShowStorageWindow(gd.ObjectParamsList[obj.GetObjectType()].Name, g.player.GetInventory(), obj.GetStorage())
	}
}

func (g *Game) selectObjType(objType ss.ObjectType) {
	if g.selectedObjType == objType {
		g.selectedObjType = ss.OBJECT_TYPE_COUNT
	} else {
		g.selectedObjType = objType
	}
}

func (g *Game) placeBeltlike(objType ss.ObjectType, hex utils.HexCoord, dir utils.Dir) objects.BeltLike {
	if !g.canPlaceObject(hex, objType, dir) {
		return nil
	}

	tier := gd.BeltlikeParamsList[objType].Tier

	var newbelt objects.BeltLike
	switch gd.BeltlikeParamsList[objType].Type {
	case ss.BELTLIKE_TYPE_NORMAL:
		newbelt = objects.NewBelt(objType, hex, dir, gd.ObjectParamsList[objType], tier)
	case ss.BELTLIKE_TYPE_SPLITTER:
		newbelt = objects.NewBeltSplitter(objType, hex, dir, gd.ObjectParamsList[objType], tier)
	case ss.BELTLIKE_TYPE_UNDER:
		bu, sw := g.findUnderToJoin(hex, dir, gd.BeltTierParamsList[tier].Reach), false
		newbelt, sw = objects.NewBeltUnder(objType, hex, dir, gd.ObjectParamsList[objType], tier, bu)
		if sw {
			g.selectedDir = dir.Reverse()
		}
	}

	g.placeObject(newbelt.(world.WorldObject))
	return newbelt
}

func (g *Game) placeStorage(hex utils.HexCoord, objType ss.ObjectType) *objects.Storage {
	if !g.canPlaceObject(hex, objType, g.selectedDir) {
		return nil
	}
	chest := objects.NewChestBox(objType, hex, gd.ObjectParamsList[objType], gd.StorageParamsList[objType])
	g.placeObject(chest)
	return chest
}

func (g *Game) placeInserter(hex utils.HexCoord, dir utils.Dir, objType ss.ObjectType) *objects.Inserter {
	if !g.canPlaceObject(hex, objType, g.selectedDir) {
		return nil
	}
	inserter := objects.NewInserter(objType, hex, dir, gd.ObjectParamsList[objType], gd.InserterParamsList[objType])
	g.placeObject(inserter)
	return inserter
}

func (g *Game) placeConverter(hex utils.HexCoord, dir utils.Dir, objType ss.ObjectType) *objects.Converter {
	if !g.canPlaceObject(hex, objType, g.selectedDir) {
		return nil
	}
	converter := objects.NewConverter(objType, hex, dir, gd.ObjectParamsList[objType], gd.ConverterParamsList[objType])
	g.placeObject(converter)
	return converter
}

func (g *Game) placeExtractor(hex utils.HexCoord, dir utils.Dir, objType ss.ObjectType) *objects.Extractor {
	if !g.canPlaceObject(hex, objType, g.selectedDir) {
		return nil
	}
	converter := objects.NewExtractor(objType, hex, dir, gd.ObjectParamsList[objType], gd.ExtractorParamsList[objType])
	g.placeObject(converter)
	return converter
}

func (g *Game) canPlaceObject(hex utils.HexCoord, objType ss.ObjectType, dir utils.Dir) bool {
	objParams := gd.ObjectParamsList[objType]
	hexes := objParams.Shape.GetHexes(hex, dir)
	for _, h := range hexes {
		if _, ok := g.getWorldObject(h); ok {
			return false
		}
	}
	return true
}

func (g *Game) placeObject(obj world.WorldObject) {
	for _, h := range getObjectHexes(obj) {
		g.setWorldObject(h, obj)
	}
}

func (g *Game) removeObjectAtHex(hex utils.HexCoord) {
	obj, ok := g.getWorldObject(hex)
	if !ok {
		return
	}

	if belt, ok := obj.(objects.BeltLike); ok {
		belt.DisconnectAll()
	}

	for _, h := range getObjectHexes(obj) {
		g.clearHex(h)
	}
}

func (g *Game) placeConnectBelts(coord1, coord2 utils.HexCoord, objType ss.ObjectType) {
	var belt1, belt2 objects.BeltLike
	if obj, ok := g.getWorldObject(coord1); ok {
		if belt1, ok = obj.(objects.BeltLike); !ok {
			return
		}
	} else {
		return
	}

	if obj, ok := g.getWorldObject(coord2); ok {
		if belt2, ok = obj.(objects.BeltLike); !ok {
			return
		}
	} else {
		dir, err := coord1.DirTo(coord2)
		if err != nil {
			return
		}
		newBelt := g.placeBeltlike(objType, coord2, dir)
		//lint:ignore S1040 // it's a nil check!!!
		if _, ok := newBelt.(objects.BeltLike); !ok {
			return
		}
		belt2 = newBelt
	}

	if belt1.CanConnectTo(belt2) {
		belt1.ConnectTo(belt2)
	}
}

func (g *Game) findUnderToJoin(hex utils.HexCoord, dir utils.Dir, reach int32) *objects.BeltUnder {
	curHex := hex
	for i := int32(1); i < reach; i++ {
		curHex = curHex.Next(dir)
		obj, ok := g.getWorldObject(curHex)
		if !ok {
			continue
		}
		bu, ok := obj.(*objects.BeltUnder)
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

//************  HexGridWorldInteractor  **********************************************

func (g *Game) GetItemInputAt(hex utils.HexCoord) (obj objects.ItemInput, ok bool) {
	if obj, ok := g.getWorldObject(hex); ok {
		if ii, ok := obj.(objects.ItemInput); ok {
			return ii, true
		}
	}
	return nil, false
}

func (g *Game) GetItemOutputAt(hex utils.HexCoord) (obj objects.ItemOutput, ok bool) {
	if obj, ok := g.getWorldObject(hex); ok {
		if io, ok := obj.(objects.ItemOutput); ok {
			return io, true
		}
	}
	return nil, false
}

func (g *Game) GetResourceAt(hex utils.HexCoord) (resType ss.ResourceType, ok bool) {
	ch := g.getChunk(hex)
	t, _, ok := ch.GetResourceTypeAt(hex)
	return t, ok
}

func (g *Game) ExtractResourceAt(hex utils.HexCoord) ss.ItemType {
	ch := g.getChunk(hex)
	return ch.ExtractResourceAt(hex)
}

//***********************************************

func (g *Game) getChunk(hex utils.HexCoord) *world.Chunk {
	ch, ok := g.chunks[hex.GetChunkCoord()]
	if !ok {
		ch = world.NewChunk(hex.GetChunkCoord(), g.worldGen)
		g.chunks[hex.GetChunkCoord()] = ch
	}
	return ch
}

func (g *Game) getWorldObject(hex utils.HexCoord) (world.WorldObject, bool) {
	ch := g.getChunk(hex)
	obj, ok := ch.GetWorldObject(hex)
	return obj, ok
}
func (g *Game) setWorldObject(hex utils.HexCoord, obj world.WorldObject) {
	g.getChunk(hex).SetWorldObject(hex, obj)
}
func (g *Game) clearHex(hex utils.HexCoord) {
	g.getChunk(hex).RemoveWorldObject(hex)
}

func getObjectHexes(obj world.WorldObject) []utils.HexCoord {
	objParams := gd.ObjectParamsList[obj.GetObjectType()]
	dir := utils.DIR_LEFT
	if do, ok := obj.(DirectionalObject); ok {
		dir = do.GetDir()
	}
	return objParams.Shape.GetHexes(obj.GetPos(), dir)
}
