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

var pos = WorldCoord{X: 0.0, Y: 0.0}
var zoom = float64(1.0) // multiply world coordinates to get screen coordinates

func GetViewZoom() float64 {
	return zoom
}

func ScreenToWorld(c ScreenCoord) WorldCoord {
	return WorldCoord{float64(c.X), float64(c.Y)}.Div(zoom).Add(pos)
}

func WorldToScreen(p WorldCoord) ScreenCoord {
	p = p.Sub(pos).Mul(zoom)
	return ScreenCoord{float32(p.X), float32(p.Y)}
}

func PctPosToScreen(pct ScreenCoord) ScreenCoord {
	return ScreenCoord{pct.X * ss.RES_X, pct.Y * ss.RES_Y}
}
func PctScaleToScreen(pct ScreenCoord) ScreenCoord {
	if ss.RES_X < ss.RES_Y {
		return ScreenCoord{pct.X * ss.RES_X, pct.Y * ss.RES_X}
	}
	return ScreenCoord{pct.X * ss.RES_Y, pct.Y * ss.RES_Y}
}
func ScreenToPctPos(c ScreenCoord) ScreenCoord {
	return ScreenCoord{c.X / ss.RES_X, c.Y / ss.RES_Y}
}

func GetZoomedHexWidth() float32 {
	return ss.HEX_WIDTH * float32(zoom)
}

func GetZoomedHexEdge() float32 {
	return ss.HEX_EDGE * float32(zoom)
}

func GetZoomedHexOffset() float32 {
	return ss.HEX_OFFSET * float32(zoom)
}

func GetZoomedDimension(dim float32) float32 {
	return dim * float32(zoom)
}

func SetView(target WorldCoord) {
	pos = target.Sub(cornerOffset.Div(float64(zoom)))
}

func ShiftView(target WorldCoord, dt uint64) {
	l := math.Min(1.0, VIEW_MOVE_LAMBDA*float64(dt*VIEW_MOVE_DT_MULT)/1000.0)

	// newpos = l * target + (1-l) * pos
	pos = target.Sub(cornerOffset.Div(float64(zoom))).Mul(l).Add(pos.Mul(1.0 - l))
}

func ZoomViewIn(dt uint64) {
	cc := pos.Add(cornerOffset.Div(float64(zoom)))

	newZoom := float64(zoom) * (1.0 + ZOOM_SPEED*math.Min(float64(dt), 1000.0)/1000.0)
	zoom = math.Min(newZoom, ZOOM_MAX)

	pos = cc.Sub(cornerOffset.Div(float64(zoom)))
}

func ZoomViewOut(dt uint64) {
	cc := pos.Add(cornerOffset.Div(float64(zoom)))

	newZoom := float64(zoom) / (1.0 + ZOOM_SPEED*math.Min(float64(dt), 1000.0)/1000.0)
	zoom = math.Max(newZoom, ZOOM_MIN)

	pos = cc.Sub(cornerOffset.Div(float64(zoom)))
}
