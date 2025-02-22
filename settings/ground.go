package settings

type GroundType uint8

const (
	GROUND_TYPE_GROUND GroundType = iota
	GROUND_TYPE_WATER  GroundType = iota

	GROUND_TYPE_COUNT GroundType = iota
)

type ResourceType uint8

const (
	RESOURCE_TYPE_IRON    ResourceType = iota
	RESOURCE_TYPE_COPPER  ResourceType = iota
	RESOURCE_TYPE_COAL    ResourceType = iota
	RESOURCE_TYPE_STONE   ResourceType = iota
	RESOURCE_TYPE_URANIUM ResourceType = iota
	RESOURCE_TYPE_OIL     ResourceType = iota

	RESOURCE_TYPE_COUNT ResourceType = iota
)
