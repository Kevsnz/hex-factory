package ui

import "hextopdown/utils"

type ControlBase struct {
	Pos  utils.ScreenCoord
	Size utils.ScreenCoord
}

func (cb *ControlBase) GetPos() utils.ScreenCoord  { return cb.Pos }
func (cb *ControlBase) SetPos(p utils.ScreenCoord) { cb.Pos = p }
func (cb *ControlBase) GetSize() utils.ScreenCoord { return cb.Size }

func (b *ControlBase) within(mp utils.ScreenCoord) bool {
	return mp.X >= 0 && mp.X < b.Size.X && mp.Y >= 0 && mp.Y < b.Size.Y
}
