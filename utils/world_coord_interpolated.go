package utils

import (
	ss "hextopdown/settings"
	"math"
)

type WorldCoordInterpolated struct {
	lastPos WorldCoord
	lastMs  uint64
	Pos     WorldCoord
}

func NewWorldCoordInterpolated() WorldCoordInterpolated {
	return WorldCoordInterpolated{
		lastPos: WorldCoord{math.NaN(), math.NaN()},
	}
}
func NewWorldCoordInterpolated2(pos WorldCoord) WorldCoordInterpolated {
	return WorldCoordInterpolated{
		Pos:     pos,
		lastPos: pos,
	}
}

func (w *WorldCoordInterpolated) GetInterpolatedPos(t, stepDt uint64) WorldCoord {
	if math.IsNaN(w.lastPos.X) {
		panic("coordinates not initialized")
	}
	dt := t - w.lastMs
	if dt > stepDt {
		return w.Pos
	}
	dd := w.Pos.Sub(w.lastPos)
	dd = dd.Mul(float64(dt) / float64(stepDt))
	return w.lastPos.Add(dd)
}

func (w *WorldCoordInterpolated) UpdatePosition(pos WorldCoord, reset bool) {
	if math.IsNaN(w.lastPos.X) || reset {
		w.lastPos = pos
	} else {
		w.lastPos = w.Pos
	}
	w.Pos = pos
	w.lastMs = ss.G_currentTickTimeMs
}
