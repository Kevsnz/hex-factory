package gamedata

import (
	"hextopdown/settings"
	"hextopdown/settings/strings"
	"hextopdown/utils"
)

type ObjectParameters struct {
	Name     strings.StringID
	BaseType settings.ObjectBaseType
	Shape    utils.Shape
}

type BeltLikeParameters struct {
	Type settings.BeltLikeType
	Tier settings.BeltTier
}

type BeltTierParameters struct {
	Speed uint32 // ticks per width
	Reach int32
}

type StorageParameters struct {
	Capacity int
}

type InserterParameters struct {
	SwingSpeed uint32 // ticks per 180deg
	Reach      int32
	StackSize  int
	Filtering  bool
}

type ConverterParameters struct {
	BuildPower uint32 // BP per tick
	AutoRecipe bool
}
