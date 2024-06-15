package settings

const (
	RES_X = 1280
	RES_Y = 960
)

const TPS = 20
const FPS = 60

const TICK_DT = 1000 / TPS

const (
	HEX_EDGE   = 64
	HEX_WIDTH  = ((HEX_EDGE*1732/1000 + 1) >> 1) << 1 // between parallel sides = Edge * sqrt(3)
	HEX_OFFSET = HEX_EDGE / 2
	HEX_HEIGHT = HEX_EDGE + 2*HEX_OFFSET
)

const ITEM_R = (HEX_WIDTH + 0.01) / 12.0 // Item Radius (in world coordinates)
const ITEM_D = ITEM_R * 2.0              // Min distance between items (in world coordinates)
const ITEM_DW = ITEM_D / HEX_WIDTH       // Min distance between items (in hex widths)

const GAP = 0.05

const LANE_OFFSET_RATIO = (1.0-2.0*GAP)/6.0 + ITEM_R/(3.0*HEX_EDGE) // Lane offset as part of hex edge (y / E)
const LANE_OFFSET_WORLD = LANE_OFFSET_RATIO * HEX_EDGE

const JOIN1 = 0.5 + LANE_OFFSET_RATIO   // first join offset (closer to entry)
const JOIN2 = 0.5 - LANE_OFFSET_RATIO/3 // second join offset (closer to exit)

const INSERTER_ARM_LENGTH = HEX_WIDTH

const HEX_DRAW_R = HEX_EDGE + 1
const BELT_DRAW_R = HEX_EDGE
const ITEM_DRAW_R = 1.75 * ITEM_R // item icon size

const DASH_LEN = 15.0 // line dash length in pixels

const ANIM_BELT_STEPS = 120 // belt animation steps per hex
const ANIM_BELT_FRAMES = 40 // belt animation frames

const TEXTURE_CACHE_DIR = "resources/cache/"
const TEXTURE_DIR = "resources/"
const TEXTURE_CACHE_EXT = ".tex"
const TEXTURE_SOURCE_EXT = ".png"

const CHESTBOX_CAPACITY_SMALL = 8
const CHESTBOX_CAPACITY_MEDIUM = 16
const CHESTBOX_CAPACITY_LARGE = 32

const CHAR_MAX_SPEED = float64(HEX_WIDTH) * 3 / TPS
const CHAR_ACCEL = float64(HEX_WIDTH) / 2 / TPS
const CHAR_DECCEL = CHAR_ACCEL * 2

const FONT_SIZE_PCT = 0.04 // 25 lines
