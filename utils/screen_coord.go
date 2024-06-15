package utils

type ScreenCoord struct {
	X, Y float32
}

func (s ScreenCoord) ToScreenCoord() WorldCoord {
	return ScreenToWorld2(s)
}

func (w ScreenCoord) Add(w2 ScreenCoord) ScreenCoord {
	return ScreenCoord{w.X + w2.X, w.Y + w2.Y}
}
func (w ScreenCoord) Sub(w2 ScreenCoord) ScreenCoord {
	return ScreenCoord{w.X - w2.X, w.Y - w2.Y}
}

func (w ScreenCoord) Mul(d float32) ScreenCoord {
	return ScreenCoord{w.X * d, w.Y * d}
}
func (w ScreenCoord) Div(d float32) ScreenCoord {
	return ScreenCoord{w.X / d, w.Y / d}
}

func (w ScreenCoord) LengthSq() float32 {
	return w.X*w.X + w.Y*w.Y
}

func (w ScreenCoord) Inverse() ScreenCoord {
	return ScreenCoord{-w.X, -w.Y}
}

func (s ScreenCoord) SwapXY() ScreenCoord {
	return ScreenCoord{X: s.Y, Y: s.X}
}
