package char

import (
	"hextopdown/settings"
	"hextopdown/utils"
	"math"
)

type Character struct {
	pos      utils.WorldCoordInterpolated
	velocity utils.WorldCoord
}

func NewCharacter(pos utils.WorldCoord) Character {
	return Character{
		pos:      utils.NewWorldCoordInterpolated2(pos),
		velocity: utils.WorldCoord{},
	}
}

func (c *Character) GetPos() utils.WorldCoordInterpolated {
	return c.pos
}

func (c *Character) UpdateMovement(dx, dy int64) {
	// TODO Fix Duke Nukem 3D's diagonal move exploit without using sqrts
	dvx, dvy := 0.0, 0.0

	if dx == 0 || math.Signbit(float64(dx)) != math.Signbit(c.velocity.X) {
		dvx = utils.Float64WithSign(settings.CHAR_DECCEL, -c.velocity.X)
	} else {
		dvx = float64(dx) * settings.CHAR_ACCEL
	}
	if dy == 0 || math.Signbit(float64(dy)) != math.Signbit(c.velocity.Y) {
		dvy = utils.Float64WithSign(settings.CHAR_DECCEL, -c.velocity.Y)
	} else {
		dvy = float64(dy) * settings.CHAR_ACCEL
	}

	if dx == 0 && math.Abs(dvx) > math.Abs(c.velocity.X) {
		c.velocity.X = 0
	} else {
		c.velocity.X = max(-settings.CHAR_MAX_SPEED, min(c.velocity.X+dvx, settings.CHAR_MAX_SPEED))
	}
	if dy == 0 && math.Abs(dvy) > math.Abs(c.velocity.Y) {
		c.velocity.Y = 0
	} else {
		c.velocity.Y = max(-settings.CHAR_MAX_SPEED, min(c.velocity.Y+dvy, settings.CHAR_MAX_SPEED))
	}

	c.pos.UpdatePosition(c.pos.Pos.Add(c.velocity), false)
}
