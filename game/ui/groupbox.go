package ui

import (
	"hextopdown/input"
	"hextopdown/renderer"
	"hextopdown/settings/strings"
	"hextopdown/utils"
)

type GroupBox struct {
	ControlBase
	text     strings.StringID
	children []iControl
}

func NewGroupBox(pos utils.ScreenCoord, size utils.ScreenCoord, text strings.StringID) *GroupBox {
	return &GroupBox{
		ControlBase: ControlBase{
			Pos:  pos,
			Size: size,
		},
		text: text,
	}
}

func (gb *GroupBox) Draw(r *renderer.GameRenderer, parentPos utils.ScreenCoord) {
	pos := gb.Pos.Add(parentPos)
	r.DrawGroupBox(pos, gb.Size, gb.text, groupBoxPadding)
	for _, child := range gb.children {
		child.Draw(r, pos)
	}
}

func (gb *GroupBox) AddChild(c iControl, ca ControlAlignment) {
	c.SetPos(ca.ConvertCoords(c.GetPos().Add(groupBoxPadding), c.GetSize(), gb.Size))
	gb.children = append(gb.children, c)
}

func (gb *GroupBox) Clear() {
	gb.children = make([]iControl, 0)
}

func (gb *GroupBox) HandleMouseMovement(mp utils.ScreenCoord) {
	mp = mp.Sub(gb.Pos)
	if !gb.within(mp) {
		return
	}

	for _, child := range gb.children {
		child.HandleMouseMovement(mp)
	}
}

func (gb *GroupBox) HandleMouseAction(mbe input.MouseButtonEvent) bool {
	mbe.Coord = mbe.Coord.Sub(gb.Pos)
	if !gb.within(mbe.Coord) {
		return false
	}

	for _, child := range gb.children {
		if child.HandleMouseAction(mbe) {
			return true
		}
	}
	return true
}
