package ui

import (
	"hextopdown/renderer"
	"hextopdown/settings/strings"
	"hextopdown/utils"
)

type ButtonText struct {
	Button
	text strings.StringID
}

func NewButtonText(pos, size utils.ScreenCoord, text strings.StringID, onClick func()) *ButtonText {
	return &ButtonText{
		Button: Button{
			ControlBase: ControlBase{
				Pos:  pos,
				Size: size,
			},
			onClick: onClick,
		},
		text: text,
	}
}

func (b *ButtonText) Draw(r *renderer.GameRenderer, parentPos utils.ScreenCoord) {
	r.DrawButtonText(b.Pos.Add(parentPos), b.Size, b.text, b.down)
}
