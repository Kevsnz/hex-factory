package utils

type ScreenCoord struct {
	X, Y float32
}

func (s ScreenCoord) ToWorld() WorldCoord {
	return ScreenToWorld(s)
}
func (s ScreenCoord) ScreenToPctPos() ScreenCoord {
	return ScreenToPctPos(s)
}
func (s ScreenCoord) PctPosToScreen() ScreenCoord {
	return PctPosToScreen(s)
}
func (s ScreenCoord) PctScaleToScreen() ScreenCoord {
	return PctScaleToScreen(s)
}

func (s ScreenCoord) Add(w2 ScreenCoord) ScreenCoord {
	return ScreenCoord{s.X + w2.X, s.Y + w2.Y}
}
func (s ScreenCoord) Sub(w2 ScreenCoord) ScreenCoord {
	return ScreenCoord{s.X - w2.X, s.Y - w2.Y}
}

func (s ScreenCoord) Mul(d float32) ScreenCoord {
	return ScreenCoord{s.X * d, s.Y * d}
}
func (s ScreenCoord) Div(d float32) ScreenCoord {
	return ScreenCoord{s.X / d, s.Y / d}
}

func (s ScreenCoord) LengthSq() float32 {
	return s.X*s.X + s.Y*s.Y
}

func (s ScreenCoord) Inverse() ScreenCoord {
	return ScreenCoord{-s.X, -s.Y}
}

func (s ScreenCoord) SwapXY() ScreenCoord {
	return ScreenCoord{X: s.Y, Y: s.X}
}
