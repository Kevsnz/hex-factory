package utils

import "math"

func Abs(x int32) int32 {
	return x * (x>>31 | 1)
}

func Sign(x int32) int32 {
	return x>>31 | 1
}

func AbsF32(x float32) float32 {
	return math.Float32frombits(math.Float32bits(x) &^ (1 << 31))
}

func GetFrameInterpolatedValue(x1, x2 float64, tickMs, timeMs, step_dt uint64) float64 {
	dx := (x2 - x1) * float64(timeMs-tickMs) / float64(step_dt)
	return x1 + dx
}
