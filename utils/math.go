package utils

import (
	"hextopdown/settings"
	"math"
)

func Abs(x int32) int32 {
	return x * (x>>31 | 1)
}

func Sign(x int32) int32 {
	return x>>31 | 1
}
func SignF(x float64) float64 {
	return float64(1 - 2*(math.Float64bits(x)>>63))
}
func Float64WithSign(x float64, sign float64) float64 {
	return math.Float64frombits(math.Float64bits(x)&^(1<<63) | math.Float64bits(sign)&(1<<63))
}

func AbsF32(x float32) float32 {
	return math.Float32frombits(math.Float32bits(x) &^ (1 << 31))
}

func GetFrameInterpolatedValue(x1, x2 float64, tickMs, timeMs, step_dt uint64) float64 {
	dx := (x2 - x1) * float64(timeMs-tickMs) / float64(step_dt)
	return x1 + dx
}

func ItemInList(itemType settings.ItemType, itemList []settings.ItemType) bool {
	for _, i := range itemList {
		if i == itemType {
			return true
		}
	}
	return false
}
