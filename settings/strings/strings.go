package strings

type StringID int

const (
	// Symbols
	STRING_X          StringID = iota
	STRING_SPACE      StringID = iota
	STRING_COMMA      StringID = iota
	STRING_COMMASPACE StringID = iota
	STRING_PERIOD     StringID = iota
	STRING_DASH       StringID = iota

	// UI Strings
	STRING_FPS           StringID = iota
	STRING_TPS           StringID = iota
	STRING_PLAYER_COORDS StringID = iota
	STRING_TOOL          StringID = iota
	STRING_INVENTORY     StringID = iota
	STRING_RECIPE        StringID = iota
	STRING_NOTEXTURE     StringID = iota

	// Object Name Strings
	STRING_OBJECT_UNKNOWN          StringID = iota
	STRING_OBJECT_INSERTER         StringID = iota
	STRING_OBJECT_CHESTBOX_SMALL   StringID = iota
	STRING_OBJECT_CHESTBOX_MEDIUM  StringID = iota
	STRING_OBJECT_CHESTBOX_LARGE   StringID = iota
	STRING_OBJECT_BELT             StringID = iota
	STRING_OBJECT_BELT_UNDER       StringID = iota
	STRING_OBJECT_BELT_SPLITTER    StringID = iota
	STRING_OBJECT_FURNACE_STONE    StringID = iota
	STRING_OBJECT_ASSSEMBLER_BASIC StringID = iota
	STRING_OBJECT_MINER_STIRLING   StringID = iota

	// Item Name Strings
	STRING_ITEM_IRON_ORE   StringID = iota
	STRING_ITEM_IRON_PLATE StringID = iota
	STRING_ITEM_IRON_GEAR  StringID = iota

	STRING_COUNT StringID = iota
)

var Strings = [STRING_COUNT]string{
	STRING_X:          "X",
	STRING_SPACE:      " ",
	STRING_COMMA:      ",",
	STRING_COMMASPACE: ", ",
	STRING_PERIOD:     ".",
	STRING_DASH:       "-",

	STRING_FPS:           "FPS: ",
	STRING_TPS:           "TPS: ",
	STRING_PLAYER_COORDS: "Player coords: ",
	STRING_TOOL:          "Current tool: ",
	STRING_INVENTORY:     "Inventory",
	STRING_RECIPE:        "Recipe",
	STRING_NOTEXTURE:     "No Texture",

	STRING_OBJECT_UNKNOWN:          "Unknown object",
	STRING_OBJECT_INSERTER:         "Inserter",
	STRING_OBJECT_CHESTBOX_SMALL:   "Small Chest",
	STRING_OBJECT_CHESTBOX_MEDIUM:  "Medium Chest",
	STRING_OBJECT_CHESTBOX_LARGE:   "Large Chest",
	STRING_OBJECT_BELT:             "Transport Belt",
	STRING_OBJECT_BELT_UNDER:       "Underground Belt",
	STRING_OBJECT_BELT_SPLITTER:    "Belt Splitter",
	STRING_OBJECT_FURNACE_STONE:    "Stone Furnace",
	STRING_OBJECT_ASSSEMBLER_BASIC: "Basic Assembling Machine",
	STRING_OBJECT_MINER_STIRLING:   "Stirling Miner",

	STRING_ITEM_IRON_ORE:   "Iron Ore",
	STRING_ITEM_IRON_PLATE: "Iron Plate",
	STRING_ITEM_IRON_GEAR:  "Iron Gear",
}
