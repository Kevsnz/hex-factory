package gamedata

import (
	ss "hextopdown/settings"
	"hextopdown/utils"
)

type ExtractorParameters struct {
	Speed             uint16 // ticks per resource unit
	MineableResources []ss.ResourceType
}

var ExtractorParamsList = map[ss.ObjectType]*ExtractorParameters{
	ss.OBJECT_TYPE_MINER_STIRLING: {Speed: 30, MineableResources: []ss.ResourceType{ss.RESOURCE_TYPE_IRON}},
}

var ExtractionResourceItems = map[ss.ResourceType]ss.ItemType{
	ss.RESOURCE_TYPE_IRON: ss.ITEM_TYPE_IRON_ORE,
}

var ExtractorShapePushPositions = [utils.SHAPE_COUNT][utils.DIR_COUNT]utils.HexCoord{
	utils.SHAPE_SINGLE: {
		utils.DIR_LEFT:       {X: -1, Y: 0},
		utils.DIR_RIGHT:      {X: 1, Y: 0},
		utils.DIR_UP_LEFT:    {X: 0, Y: -1},
		utils.DIR_DOWN_RIGHT: {X: 0, Y: 1},
		utils.DIR_UP_RIGHT:   {X: 1, Y: -1},
		utils.DIR_DOWN_LEFT:  {X: -1, Y: 1},
	},
	utils.SHAPE_DIAMOND: {
		utils.DIR_LEFT:       {X: -1, Y: 0},
		utils.DIR_RIGHT:      {X: 2, Y: 0},
		utils.DIR_UP_LEFT:    {X: 0, Y: -1},
		utils.DIR_DOWN_RIGHT: {X: 0, Y: 2},
		utils.DIR_UP_RIGHT:   {X: 2, Y: -2},
		utils.DIR_DOWN_LEFT:  {X: -1, Y: 1},
	},
	utils.SHAPE_BIGHEX: {
		utils.DIR_LEFT:       {X: -2, Y: 0},
		utils.DIR_RIGHT:      {X: 2, Y: 0},
		utils.DIR_UP_LEFT:    {X: 0, Y: -2},
		utils.DIR_DOWN_RIGHT: {X: 0, Y: 2},
		utils.DIR_UP_RIGHT:   {X: 2, Y: -2},
		utils.DIR_DOWN_LEFT:  {X: -2, Y: 2},
	},
}

func GetExtractionHexes(pos utils.HexCoord, shape utils.Shape, dir utils.Dir) []utils.HexCoord {
	switch shape {
	case utils.SHAPE_SINGLE:
		return []utils.HexCoord{pos}
	case utils.SHAPE_DIAMOND:
		return shape.GetHexes(pos, dir)
	case utils.SHAPE_BIGHEX:
		hexes := shape.GetHexes(pos, dir)
		for _, d := range utils.AllDirs {
			hexes = append(hexes, pos.Shift(d, 2))
			hexes = append(hexes, pos.Shift(d, 2).Shift(d.NextCW().NextCW(), 1))
		}
		return hexes
	default:
		panic("unknown shape")
	}
}
