package settings

type StructureType uint32

const (
	STRUCTURE_TYPE_INSERTER_LEFT      StructureType = iota
	STRUCTURE_TYPE_INSERTER_RIGHT     StructureType = iota
	STRUCTURE_TYPE_INSERTER_UPLEFT    StructureType = iota
	STRUCTURE_TYPE_INSERTER_UPRIGHT   StructureType = iota
	STRUCTURE_TYPE_INSERTER_DOWNLEFT  StructureType = iota
	STRUCTURE_TYPE_INSERTER_DOWNRIGHT StructureType = iota

	STRUCTURE_TYPE_CHESHBOX_SMALL  StructureType = iota
	STRUCTURE_TYPE_CHESHBOX_MEDIUM StructureType = iota
	STRUCTURE_TYPE_CHESHBOX_LARGE  StructureType = iota

	STRUCTURE_TYPE_FURNACE_RIGHT     StructureType = iota
	STRUCTURE_TYPE_FURNACE_UPRIGHT   StructureType = iota
	STRUCTURE_TYPE_FURNACE_DOWNRIGHT StructureType = iota

	STRUCTURE_TYPE_COUNT StructureType = iota
)
