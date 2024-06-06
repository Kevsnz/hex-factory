package utils

import (
	ss "hextopdown/settings"
	"math"
)

const (
	ZOOM_MAX   = 2.0
	ZOOM_MIN   = 0.25
	ZOOM_SPEED = 0.5
	ZOOM_MULT  = ZOOM_SPEED + 1.0
)

const VIEW_MOVE_LAMBDA = 1.0
const VIEW_MOVE_DT_MULT = 2.0

var cornerOffset = WorldCoord{X: ss.RES_X / 2.0, Y: ss.RES_Y / 2.0}

type Viewport struct {
	Pos  WorldCoord
	Zoom float64 // multiply world coordinated to get screen coordinates
}

func NewViewport(view WorldCoord, zoom float64) *Viewport {
	return &Viewport{
		Pos:  view.Sub(cornerOffset),
		Zoom: zoom,
	}
}

func (v *Viewport) ScreenToWorld(x, y float32) WorldCoord {
	return WorldCoord{float64(x), float64(y)}.Div(v.Zoom).Add(v.Pos)
}

func (v *Viewport) WorldToScreen(pos WorldCoord) (float32, float32) {
	pos = pos.Sub(v.Pos).Mul(v.Zoom)
	return float32(pos.X), float32(pos.Y)
}

func (v *Viewport) GetHexWidth() float32 {
	return ss.HEX_WIDTH * float32(v.Zoom)
}

func (v *Viewport) GetHexEdge() float32 {
	return ss.HEX_EDGE * float32(v.Zoom)
}

func (v *Viewport) GetHexOffset() float32 {
	return ss.HEX_OFFSET * float32(v.Zoom)
}

func (v *Viewport) GetZoomedDimension(dim float32) float32 {
	return dim * float32(v.Zoom)
}

func (v *Viewport) ShiftViewport(target WorldCoord, dt uint64) {
	l := math.Min(1.0, VIEW_MOVE_LAMBDA*float64(dt*VIEW_MOVE_DT_MULT)/1000.0)

	// newpos = l * target + (1-l) * pos
	v.Pos = target.Sub(cornerOffset.Div(float64(v.Zoom))).Mul(l).Add(v.Pos.Mul(1.0 - l))
}

func (v *Viewport) ZoomIn(dt uint64) {
	cc := v.Pos.Add(cornerOffset.Div(float64(v.Zoom)))

	newZoom := float64(v.Zoom) * (1.0 + ZOOM_SPEED*math.Min(float64(dt), 1000.0)/1000.0)
	v.Zoom = math.Min(newZoom, ZOOM_MAX)

	v.Pos = cc.Sub(cornerOffset.Div(float64(v.Zoom)))
}

func (v *Viewport) ZoomOut(dt uint64) {
	cc := v.Pos.Add(cornerOffset.Div(float64(v.Zoom)))

	newZoom := float64(v.Zoom) / (1.0 + ZOOM_SPEED*math.Min(float64(dt), 1000.0)/1000.0)
	v.Zoom = math.Max(newZoom, ZOOM_MIN)

	v.Pos = cc.Sub(cornerOffset.Div(float64(v.Zoom)))
}
