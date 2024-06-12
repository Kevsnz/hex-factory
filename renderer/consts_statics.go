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

var radiusOffsets = [utils.DIR_COUNT][2]float32{
	utils.DIR_LEFT:       {-ss.HEX_WIDTH / 2.0, 0},
	utils.DIR_UP_LEFT:    {-ss.HEX_WIDTH / 4.0, -(ss.HEX_EDGE/2.0 + ss.HEX_OFFSET/2.0)},
	utils.DIR_UP_RIGHT:   {ss.HEX_WIDTH / 4.0, -(ss.HEX_EDGE/2.0 + ss.HEX_OFFSET/2.0)},
	utils.DIR_RIGHT:      {ss.HEX_WIDTH / 2.0, 0},
	utils.DIR_DOWN_RIGHT: {ss.HEX_WIDTH / 4.0, (ss.HEX_EDGE/2.0 + ss.HEX_OFFSET/2.0)},
	utils.DIR_DOWN_LEFT:  {-ss.HEX_WIDTH / 4.0, (ss.HEX_EDGE/2.0 + ss.HEX_OFFSET/2.0)},
}

var lanesOffsetsLeft = [utils.DIR_COUNT][2]float32{
	utils.DIR_LEFT:       {0.0, ss.HEX_EDGE * LOR},
	utils.DIR_UP_LEFT:    {-ss.HEX_WIDTH / 2.0 * LOR, ss.HEX_OFFSET * LOR},
	utils.DIR_UP_RIGHT:   {-ss.HEX_WIDTH / 2.0 * LOR, -ss.HEX_OFFSET * LOR},
	utils.DIR_RIGHT:      {0.0, -ss.HEX_EDGE * LOR},
	utils.DIR_DOWN_RIGHT: {ss.HEX_WIDTH / 2.0 * LOR, -ss.HEX_OFFSET * LOR},
	utils.DIR_DOWN_LEFT:  {ss.HEX_WIDTH / 2.0 * LOR, ss.HEX_OFFSET * LOR},
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
	Width, Height    float64
	OffsetX, OffsetY float64
}

func GetShapeParam(shape utils.Shape, dir utils.Dir) ShapeParam {
	switch shape {
	case utils.SHAPE_DIAMOND:
		if dir == utils.DIR_LEFT || dir == utils.DIR_RIGHT {
			return ShapeParam{
				Width:   ss.HEX_WIDTH * 2,
				Height:  ss.HEX_WIDTH * 3,
				OffsetX: ss.HEX_WIDTH / 2,
				OffsetY: ss.HEX_WIDTH*3 - ss.HEX_EDGE*5/2,
			}
		}
		if dir == utils.DIR_UP_LEFT || dir == utils.DIR_DOWN_RIGHT {
			return ShapeParam{
				Width:   ss.HEX_WIDTH * 5 / 2,
				Height:  ss.HEX_WIDTH * 5 / 2,
				OffsetX: ss.HEX_WIDTH,
				OffsetY: ss.HEX_WIDTH*5/2 - ss.HEX_EDGE*5/2,
			}
		}
		return ShapeParam{
			Width:   ss.HEX_WIDTH * 5 / 2,
			Height:  ss.HEX_WIDTH * 5 / 2,
			OffsetX: ss.HEX_WIDTH,
			OffsetY: ss.HEX_WIDTH*5/2 - ss.HEX_EDGE,
		}
	case utils.SHAPE_BIGHEX:
		return ShapeParam{
			Width:   ss.HEX_WIDTH * 3,
			Height:  ss.HEX_WIDTH * 3,
			OffsetX: ss.HEX_WIDTH * 3 / 2,
			OffsetY: ss.HEX_WIDTH*3 - ss.HEX_EDGE*5/2,
		}
	default: // single
		return ShapeParam{
			Width:   ss.HEX_HEIGHT,
			Height:  ss.HEX_HEIGHT,
			OffsetX: ss.HEX_HEIGHT / 2,
			OffsetY: ss.HEX_HEIGHT / 2,
		}
	}
}
