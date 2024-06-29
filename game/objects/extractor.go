package objects

import (
	gd "hextopdown/game/gamedata"
	"hextopdown/game/items"
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
)

type Extractor struct {
	Object
	dir               utils.Dir
	params            *gd.ExtractorParameters
	extractionHex     int
	extractionProgres uint16
	extractedItem     *ss.ItemType
}

func NewExtractor(
	objType ss.ObjectType,
	pos utils.HexCoord,
	dir utils.Dir,
	objParams *gd.ObjectParameters,
	params *gd.ExtractorParameters,
) *Extractor {
	return &Extractor{
		Object: Object{
			objType:   objType,
			pos:       pos,
			objParams: objParams,
		},
		dir:    dir,
		params: params,
	}
}

func (e *Extractor) GetDir() utils.Dir {
	return e.dir
}

func (e *Extractor) Rotate(_ bool) {}

func (e *Extractor) DrawGroundLevel(r *renderer.GameRenderer) {
	r.DrawObjectGround(e.pos.CenterToWorld(), e.objType, e.objParams.Shape, e.dir)
}

func (e *Extractor) DrawOnGroundLevel(r *renderer.GameRenderer) {
	r.DrawObjectOnGround(e.pos.CenterToWorld(), e.objType, e.objParams.Shape, e.dir)

	pushPos := e.pos.Add(gd.ExtractorShapePushPositions[e.objParams.Shape][e.dir])
	r.DrawArrowWorld(pushPos.CenterToWorld().ShiftDir(e.dir.Reverse(), ss.HEX_WIDTH/2), e.dir, ss.HEX_WIDTH/4)

	p := e.pos.CenterToWorld().Add(e.objParams.Shape.GetCenterOffset(e.dir))
	r.DrawProgressBar(p, 1.25, uint32(e.params.Speed-e.extractionProgres-1), uint32(e.params.Speed-1))
}

func (e *Extractor) Update(ticks uint64, world HexGridWorldInteractor) {
	if e.extractedItem != nil {
		pushPos := e.pos.Add(gd.ExtractorShapePushPositions[e.objParams.Shape][e.dir])
		obj, ok := world.GetItemInputAt(pushPos)
		if ok {
			itemPos := pushPos.CenterToWorld().ShiftDir(e.dir.Reverse(), ss.LANE_OFFSET_WORLD)
			ok = obj.TakeItemIn(
				itemPos,
				items.NewItemInWorld2(*e.extractedItem, itemPos),
			)
			if ok {
				e.extractedItem = nil
			}
		}
	}

	if e.extractionProgres == 0 {
		hexes := gd.GetExtractionHexes(e.pos, e.objParams.Shape, e.dir)
		for i := 0; i < len(hexes); i++ {
			hex := hexes[(i+e.extractionHex)%len(hexes)]
			_, ok := world.GetResourceAt(hex)
			if !ok {
				continue
			}

			e.extractionHex = (i + e.extractionHex) % len(hexes)
			e.extractionProgres = e.params.Speed
			break
		}
	}

	if e.extractionProgres == 0 {
		return
	}
	if e.extractionProgres == 1 && e.extractedItem != nil {
		return
	}

	e.extractionProgres--
	if e.extractionProgres > 0 {
		return
	}

	hexes := gd.GetExtractionHexes(e.pos, e.objParams.Shape, e.dir)
	if _, ok := world.GetResourceAt(hexes[e.extractionHex]); ok {
		item := world.ExtractResourceAt(hexes[e.extractionHex])
		e.extractedItem = &item
	}
}

func (e *Extractor) GetItemList() []utils.ItemInfo {
	if e.extractedItem == nil {
		return []utils.ItemInfo{}
	}
	return []utils.ItemInfo{
		{Type: *e.extractedItem, Count: 1},
	}
}
