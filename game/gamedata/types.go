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
