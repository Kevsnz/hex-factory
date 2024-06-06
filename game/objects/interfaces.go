package objects

import (
	"hextopdown/game/items"
	"hextopdown/renderer"
	"hextopdown/utils"
)

type BeltLike interface {
	GetPos() utils.HexCoord
	GetInConn(dir utils.Dir) *BeltConnection
	DrawItems(r *renderer.GameRenderer)
	Rotate(cw bool)
	CanConnectTo(b2 BeltLike) bool
	CanConnectIn(dir utils.Dir) bool
	ConnectTo(b2 BeltLike)
	ClearIn(dir utils.Dir)
	ClearOut(dir utils.Dir)
	CreateIn(inDir utils.Dir, b2 BeltLike)
	DisconnectAll()
}

type HexGridWorldInteractor interface {
	GetItemInputAt(hex utils.HexCoord) (obj ItemInput, ok bool)
	GetItemOutputAt(hex utils.HexCoord) (obj ItemOutput, ok bool)
}

type ItemOutput interface {
	TakeItemOut(pos utils.WorldCoord) (item *items.ItemInWorld, ok bool)
}
type ItemInput interface {
	TakeItemIn(pos utils.WorldCoord, item items.ItemInWorld) (ok bool)
}
