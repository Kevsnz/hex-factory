package strings

type StringID int

const (
	STRING_SPACE      StringID = iota
	STRING_COMMA      StringID = iota
	STRING_COMMASPACE StringID = iota
	STRING_PERIOD     StringID = iota
	STRING_DASH       StringID = iota

	STRING_FPS           StringID = iota
	STRING_TPS           StringID = iota
	STRING_PLAYER_COORDS StringID = iota

	STRING_OBJECT_UNKNOWN         StringID = iota
	STRING_OBJECT_INSERTER        StringID = iota
	STRING_OBJECT_CHESTBOX_SMALL  StringID = iota
	STRING_OBJECT_CHESTBOX_MEDIUM StringID = iota
	STRING_OBJECT_CHESTBOX_LARGE  StringID = iota
	STRING_OBJECT_BELT            StringID = iota
	STRING_OBJECT_BELT_UNDER      StringID = iota
	STRING_OBJECT_BELT_SPLITTER   StringID = iota

	STRING_ITEM_IRON_PLATE StringID = iota

	STRING_COUNT StringID = iota
)

var Strings = [STRING_COUNT]string{
	STRING_SPACE:      " ",
	STRING_COMMA:      ",",
	STRING_COMMASPACE: ", ",
	STRING_PERIOD:     ".",
	STRING_DASH:       "-",

	STRING_FPS:           "FPS: ",
	STRING_TPS:           "TPS: ",
	STRING_PLAYER_COORDS: "Player coords: ",

	STRING_OBJECT_UNKNOWN:         "Unknown object",
	STRING_OBJECT_INSERTER:        "Inserter",
	STRING_OBJECT_CHESTBOX_SMALL:  "Small Chest",
	STRING_OBJECT_CHESTBOX_MEDIUM: "Medium Chest",
	STRING_OBJECT_CHESTBOX_LARGE:  "Large Chest",
	STRING_OBJECT_BELT:            "Transport Belt",
	STRING_OBJECT_BELT_UNDER:      "Underground Belt",
	STRING_OBJECT_BELT_SPLITTER:   "Belt Splitter",

	STRING_ITEM_IRON_PLATE: "Iron Plate",
}
