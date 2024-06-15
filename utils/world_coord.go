package utils

import (
	"math"
)

var sincos = [DIR_COUNT][2]float64{
	DIR_LEFT:       {math.Sin(math.Pi), math.Cos(math.Pi)},
	DIR_UP_LEFT:    {math.Sin(math.Pi * 2 / 3), math.Cos(math.Pi * 2 / 3)},
	DIR_UP_RIGHT:   {math.Sin(math.Pi / 3), math.Cos(math.Pi / 3)},
	DIR_RIGHT:      {math.Sin(0), math.Cos(0)},
	DIR_DOWN_RIGHT: {math.Sin(-math.Pi / 3), math.Cos(-math.Pi / 3)},
	DIR_DOWN_LEFT:  {math.Sin(-math.Pi * 2 / 3), math.Cos(-math.Pi * 2 / 3)},
}

type WorldCoord struct {
	X, Y float64
}

func (w WorldCoord) Add(w2 WorldCoord) WorldCoord {
	return WorldCoord{w.X + w2.X, w.Y + w2.Y}
}
func (w WorldCoord) Sub(w2 WorldCoord) WorldCoord {
	return WorldCoord{w.X - w2.X, w.Y - w2.Y}
}

func (w WorldCoord) Mul(d float64) WorldCoord {
	return WorldCoord{w.X * d, w.Y * d}
}
func (w WorldCoord) Div(d float64) WorldCoord {
	return WorldCoord{w.X / d, w.Y / d}
}

func (w WorldCoord) LengthSq() float64 {
	return w.X*w.X + w.Y*w.Y
}

// Rotates coordinate system so that dir is pointing to x+
func (w WorldCoord) RotateToPlusX(dir Dir) WorldCoord {
	sin := sincos[dir][0]
	cos := sincos[dir][1]
	return WorldCoord{w.X*cos - w.Y*sin, w.X*sin + w.Y*cos}
}

func (w WorldCoord) UnrotateFromPlusX(dir Dir) WorldCoord {
	sin := sincos[dir][0]
	cos := sincos[dir][1]
	return WorldCoord{w.X*cos + w.Y*sin, w.Y*cos - w.X*sin}
}

func (w WorldCoord) ShiftDir(dir Dir, d float64) WorldCoord {
	return w.Add(WorldCoord{d, 0}.UnrotateFromPlusX(dir))
}

func (w WorldCoord) Shift(x, y float64) WorldCoord {
	return w.Add(WorldCoord{x, y})
}

func (w WorldCoord) DistanceSqTo(w2 WorldCoord) float64 {
	return (w.X-w2.X)*(w.X-w2.X) + (w.Y-w2.Y)*(w.Y-w2.Y)
}

func (w WorldCoord) ToScreen() ScreenCoord {
	return WorldToScreen(w)
}

func (w WorldCoord) ToHex() HexCoord {
	return HexCoordFromWorld(w)
}
