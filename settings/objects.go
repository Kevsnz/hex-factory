package settings

type ObjectBaseType uint32

const (
	STRUCTURE_BASETYPE_BELTLIKE  ObjectBaseType = iota
	STRUCTURE_BASETYPE_INSERTER  ObjectBaseType = iota
	STRUCTURE_BASETYPE_STORAGE   ObjectBaseType = iota
	STRUCTURE_BASETYPE_CONVERTER ObjectBaseType = iota
	STRUCTURE_BASETYPE_EXTRACTOR ObjectBaseType = iota

	STRUCTURE_BASETYPE_COUNT ObjectBaseType = iota
)

type ObjectType uint32 // concrete types of all the objects like "wooden chest" or "fast inserter"

const (
	OBJECT_TYPE_BELT1         ObjectType = iota
	OBJECT_TYPE_BELTUNDER1    ObjectType = iota
	OBJECT_TYPE_BELTSPLITTER1 ObjectType = iota

	OBJECT_TYPE_CHESTBOX_SMALL  ObjectType = iota
	OBJECT_TYPE_CHESTBOX_MEDIUM ObjectType = iota
	OBJECT_TYPE_CHESTBOX_LARGE  ObjectType = iota

	OBJECT_TYPE_INSERTER1 ObjectType = iota

	OBJECT_TYPE_FURNACE_STONE ObjectType = iota

	OBJECT_TYPE_ASSEMBLER_BASIC ObjectType = iota

	OBJECT_TYPE_MINER_STIRLING ObjectType = iota

	OBJECT_TYPE_COUNT ObjectType = iota
)

type BeltLikeType uint32

const (
	BELTLIKE_TYPE_NORMAL   BeltLikeType = iota
	BELTLIKE_TYPE_UNDER    BeltLikeType = iota
	BELTLIKE_TYPE_SPLITTER BeltLikeType = iota

	BELTLIKE_TYPE_COUNT BeltLikeType = iota
)

type BeltTier uint32

const (
	BELT_TIER_NORMAL  BeltTier = iota
	BELT_TIER_FAST    BeltTier = iota
	BELT_TIER_EXPRESS BeltTier = iota

	BELT_TIER_COUNT BeltTier = iota
)
