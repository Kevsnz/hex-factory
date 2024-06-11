package settings

type ObjectBaseType uint32

const (
	STRUCTURE_BASETYPE_BELTLIKE  ObjectBaseType = iota
	STRUCTURE_BASETYPE_INSERTER  ObjectBaseType = iota
	STRUCTURE_BASETYPE_STORAGE   ObjectBaseType = iota
	STRUCTURE_BASETYPE_CONVERTER ObjectBaseType = iota

	STRUCTURE_BASETYPE_COUNT ObjectBaseType = iota
)

type ObjectType uint32 // types of all the objects like "wooden chest" or "fast inserter"

const (
	OBJECT_TYPE_CHESTBOX_SMALL  ObjectType = iota
	OBJECT_TYPE_CHESTBOX_MEDIUM ObjectType = iota
	OBJECT_TYPE_CHESTBOX_LARGE  ObjectType = iota
	OBJECT_TYPE_INSERTER1       ObjectType = iota
	OBJECT_TYPE_FURNACE_STONE   ObjectType = iota

	OBJECT_TYPE_COUNT ObjectType = iota
)
