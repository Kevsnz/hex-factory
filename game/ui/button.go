package ui

import (
	"hextopdown/renderer"
	"hextopdown/utils"
)

type Button struct {
	Pos  utils.ScreenCoord
	Size utils.ScreenCoord
}

func NewButton(pos, size utils.ScreenCoord) *Button {
	return &Button{
		Pos:  pos,
		Size: size,
	}
}

func (b *Button) Draw(r *renderer.GameRenderer, parentPos utils.ScreenCoord) {
	r.DrawButton(b.Pos.Add(parentPos), b.Size)
}
