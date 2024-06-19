package renderer

import (
	ss "hextopdown/settings"
	"hextopdown/utils"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	RES_X = ss.RES_X
	RES_Y = ss.RES_Y
)

const LOR = ss.LANE_OFFSET_RATIO
const TEXTURE_SIZE_HEX = 400
const TEXTURE_ICON_SIZE = 64

var radiusOffsets = [utils.DIR_COUNT]utils.ScreenCoord{
	utils.DIR_LEFT:       {X: -ss.HEX_WIDTH / 2.0, Y: 0},
	utils.DIR_UP_LEFT:    {X: -ss.HEX_WIDTH / 4.0, Y: -(ss.HEX_EDGE/2.0 + ss.HEX_OFFSET/2.0)},
	utils.DIR_UP_RIGHT:   {X: ss.HEX_WIDTH / 4.0, Y: -(ss.HEX_EDGE/2.0 + ss.HEX_OFFSET/2.0)},
	utils.DIR_RIGHT:      {X: ss.HEX_WIDTH / 2.0, Y: 0},
	utils.DIR_DOWN_RIGHT: {X: ss.HEX_WIDTH / 4.0, Y: (ss.HEX_EDGE/2.0 + ss.HEX_OFFSET/2.0)},
	utils.DIR_DOWN_LEFT:  {X: -ss.HEX_WIDTH / 4.0, Y: (ss.HEX_EDGE/2.0 + ss.HEX_OFFSET/2.0)},
}

var lanesOffsetsLeft = [utils.DIR_COUNT]utils.ScreenCoord{
	utils.DIR_LEFT:       {X: 0.0, Y: ss.HEX_EDGE * LOR},
	utils.DIR_UP_LEFT:    {X: -ss.HEX_WIDTH / 2.0 * LOR, Y: ss.HEX_OFFSET * LOR},
	utils.DIR_UP_RIGHT:   {X: -ss.HEX_WIDTH / 2.0 * LOR, Y: -ss.HEX_OFFSET * LOR},
	utils.DIR_RIGHT:      {X: 0.0, Y: -ss.HEX_EDGE * LOR},
	utils.DIR_DOWN_RIGHT: {X: ss.HEX_WIDTH / 2.0 * LOR, Y: -ss.HEX_OFFSET * LOR},
	utils.DIR_DOWN_LEFT:  {X: ss.HEX_WIDTH / 2.0 * LOR, Y: ss.HEX_OFFSET * LOR},
}

type beltTypeFlip struct {
	type1 ss.BeltType
	flip  sdl.RendererFlip
}

var arrowDirMapping = map[utils.Dir]struct {
	idx  int
	flip sdl.RendererFlip
}{
	utils.DIR_LEFT:       {idx: 0, flip: sdl.FLIP_HORIZONTAL},
	utils.DIR_UP_LEFT:    {idx: 1, flip: sdl.FLIP_VERTICAL | sdl.FLIP_HORIZONTAL},
	utils.DIR_UP_RIGHT:   {idx: 1, flip: sdl.FLIP_VERTICAL},
	utils.DIR_RIGHT:      {idx: 0, flip: sdl.FLIP_NONE},
	utils.DIR_DOWN_RIGHT: {idx: 1, flip: sdl.FLIP_NONE},
	utils.DIR_DOWN_LEFT:  {idx: 1, flip: sdl.FLIP_HORIZONTAL},
}

var beltOnFlipMapping = [ss.BELT_ON_COUNT]beltTypeFlip{
	// Underground
	ss.BELT_ON_UNDER_IN_RIGHT:     {ss.BELT_ON_UNDER_IN_RIGHT, sdl.FLIP_NONE},
	ss.BELT_ON_UNDER_IN_LEFT:      {ss.BELT_ON_UNDER_IN_RIGHT, sdl.FLIP_HORIZONTAL},
	ss.BELT_ON_UNDER_IN_DOWNRIGHT: {ss.BELT_ON_UNDER_IN_DOWNRIGHT, sdl.FLIP_NONE},
	ss.BELT_ON_UNDER_IN_DOWNLEFT:  {ss.BELT_ON_UNDER_IN_DOWNRIGHT, sdl.FLIP_HORIZONTAL},
	ss.BELT_ON_UNDER_IN_UPLEFT:    {ss.BELT_ON_UNDER_IN_UPLEFT, sdl.FLIP_NONE},
	ss.BELT_ON_UNDER_IN_UPRIGHT:   {ss.BELT_ON_UNDER_IN_UPLEFT, sdl.FLIP_HORIZONTAL},

	ss.BELT_ON_UNDER_OUT_RIGHT:     {ss.BELT_ON_UNDER_OUT_RIGHT, sdl.FLIP_NONE},
	ss.BELT_ON_UNDER_OUT_LEFT:      {ss.BELT_ON_UNDER_OUT_RIGHT, sdl.FLIP_HORIZONTAL},
	ss.BELT_ON_UNDER_OUT_DOWNRIGHT: {ss.BELT_ON_UNDER_OUT_DOWNRIGHT, sdl.FLIP_NONE},
	ss.BELT_ON_UNDER_OUT_DOWNLEFT:  {ss.BELT_ON_UNDER_OUT_DOWNRIGHT, sdl.FLIP_HORIZONTAL},
	ss.BELT_ON_UNDER_OUT_UPLEFT:    {ss.BELT_ON_UNDER_OUT_UPLEFT, sdl.FLIP_NONE},
	ss.BELT_ON_UNDER_OUT_UPRIGHT:   {ss.BELT_ON_UNDER_OUT_UPLEFT, sdl.FLIP_HORIZONTAL},

	// Splitters
	ss.BELT_ON_SPLITTER_UPLEFTRIGHT_DOWNLEFTRIGHT: {ss.BELT_ON_SPLITTER_UPLEFTRIGHT_DOWNLEFTRIGHT, sdl.FLIP_NONE},
	ss.BELT_ON_SPLITTER_DOWNLEFTRIGHT_UPLEFTRIGHT: {ss.BELT_ON_SPLITTER_DOWNLEFTRIGHT_UPLEFTRIGHT, sdl.FLIP_NONE},
	ss.BELT_ON_SPLITTER_LEFTUPLEFT_RIGHTDOWNRIGHT: {ss.BELT_ON_SPLITTER_LEFTUPLEFT_RIGHTDOWNRIGHT, sdl.FLIP_NONE},
	ss.BELT_ON_SPLITTER_RIGHTUPRIGHT_LEFTDOWNLEFT: {ss.BELT_ON_SPLITTER_LEFTUPLEFT_RIGHTDOWNRIGHT, sdl.FLIP_HORIZONTAL},
	ss.BELT_ON_SPLITTER_RIGHTDOWNRIGHT_LEFTUPLEFT: {ss.BELT_ON_SPLITTER_RIGHTDOWNRIGHT_LEFTUPLEFT, sdl.FLIP_NONE},
	ss.BELT_ON_SPLITTER_LEFTDOWNLEFT_RIGHTUPRIGHT: {ss.BELT_ON_SPLITTER_RIGHTDOWNRIGHT_LEFTUPLEFT, sdl.FLIP_HORIZONTAL},
}

type ShapeParam struct {
	Size   utils.ScreenCoord
	Offset utils.ScreenCoord
}

func GetShapeParam(shape utils.Shape, dir utils.Dir) ShapeParam {
	switch shape {
	case utils.SHAPE_DIAMOND:
		if dir == utils.DIR_LEFT || dir == utils.DIR_RIGHT {
			return ShapeParam{
				Size:   utils.ScreenCoord{X: ss.HEX_WIDTH * 2, Y: ss.HEX_WIDTH * 3},
				Offset: utils.ScreenCoord{X: ss.HEX_WIDTH / 2, Y: ss.HEX_WIDTH*3 - ss.HEX_EDGE*5/2},
			}
		}
		if dir == utils.DIR_UP_LEFT || dir == utils.DIR_DOWN_RIGHT {
			return ShapeParam{
				Size:   utils.ScreenCoord{X: ss.HEX_WIDTH * 5 / 2, Y: ss.HEX_WIDTH * 5 / 2},
				Offset: utils.ScreenCoord{X: ss.HEX_WIDTH, Y: ss.HEX_WIDTH*5/2 - ss.HEX_EDGE*5/2},
			}
		}
		return ShapeParam{
			Size:   utils.ScreenCoord{X: ss.HEX_WIDTH * 5 / 2, Y: ss.HEX_WIDTH * 5 / 2},
			Offset: utils.ScreenCoord{X: ss.HEX_WIDTH, Y: ss.HEX_WIDTH*5/2 - ss.HEX_EDGE},
		}
	case utils.SHAPE_BIGHEX:
		return ShapeParam{
			Size:   utils.ScreenCoord{X: ss.HEX_WIDTH * 3, Y: ss.HEX_WIDTH * 3},
			Offset: utils.ScreenCoord{X: ss.HEX_WIDTH * 3 / 2, Y: ss.HEX_WIDTH*3 - ss.HEX_EDGE*5/2},
		}
	default: // single
		return ShapeParam{
			Size:   utils.ScreenCoord{X: ss.HEX_HEIGHT, Y: ss.HEX_HEIGHT},
			Offset: utils.ScreenCoord{X: ss.HEX_HEIGHT / 2, Y: ss.HEX_HEIGHT / 2},
		}
	}
}

var iconItemList = map[ss.ItemType]int{
	ss.ITEM_TYPE_IRON_GEAR:  0,
	ss.ITEM_TYPE_IRON_ORE:   1,
	ss.ITEM_TYPE_IRON_PLATE: 2,
}
