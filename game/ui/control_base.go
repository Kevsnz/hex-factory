package ui

import "hextopdown/utils"

type ControlBase struct {
	Pos  utils.ScreenCoord
	Size utils.ScreenCoord
}

func (cb *ControlBase) GetPos() utils.ScreenCoord  { return cb.Pos }
func (cb *ControlBase) SetPos(p utils.ScreenCoord) { cb.Pos = p }
func (cb *ControlBase) GetSize() utils.ScreenCoord { return cb.Size }
