package utils

import (
	"fmt"
	ss "hextopdown/settings"
	"math"
)

var nexts = [DIR_COUNT][2]int32{
	DIR_LEFT:       {-1, 0},
	DIR_UP_LEFT:    {0, -1},
	DIR_UP_RIGHT:   {1, -1},
	DIR_RIGHT:      {1, 0},
	DIR_DOWN_RIGHT: {0, 1},
	DIR_DOWN_LEFT:  {-1, 1},
}

type HexCoord struct {
	X, Y int32
}

func HexCoordFromWorld(pos WorldCoord) HexCoord {
	hy := int32(math.Floor(pos.Y / (ss.HEX_EDGE + ss.HEX_OFFSET)))
	pos.X -= float64(hy) * (ss.HEX_WIDTH / 2)
	hx := int32(math.Floor((pos.X - 1) / (ss.HEX_WIDTH)))

	lx := pos.X - 1 - float64(hx)*ss.HEX_WIDTH
	ly := pos.Y - float64(hy)*(ss.HEX_EDGE+ss.HEX_OFFSET)
	if ly > ss.HEX_EDGE {
		c := lx - ss.HEX_WIDTH/2
		ey := math.Abs(c) / (ss.HEX_WIDTH / 2) * ss.HEX_OFFSET
		if ly-0.5 > ss.HEX_EDGE+ss.HEX_OFFSET-ey {
			hx += int32(c) >> 31 // convert c to -1 or 0 depending on its sign
			hy += 1
		}
	}
	return HexCoord{X: hx, Y: hy}
}

func (hc HexCoord) LeftTopToWorld() WorldCoord {
	return WorldCoord{
		float64(hc.X)*ss.HEX_WIDTH + float64(hc.Y)*(ss.HEX_WIDTH/2),
		float64(hc.Y) * (ss.HEX_EDGE + ss.HEX_OFFSET),
	}
}

func (hc HexCoord) CenterToWorld() WorldCoord {
	return WorldCoord{
		float64(hc.X)*ss.HEX_WIDTH + float64(hc.Y)*(ss.HEX_WIDTH/2) + ss.HEX_WIDTH/2,
		float64(hc.Y)*(ss.HEX_EDGE+ss.HEX_OFFSET) + ss.HEX_EDGE/2,
	}
}

func (hc HexCoord) Next(dir Dir) HexCoord {
	return HexCoord{X: hc.X + nexts[dir][0], Y: hc.Y + nexts[dir][1]}
}
func (hc HexCoord) Shift(dir Dir, steps int) HexCoord {
	return HexCoord{X: hc.X + nexts[dir][0]*int32(steps), Y: hc.Y + nexts[dir][1]*int32(steps)}
}
func (hc HexCoord) DirTo(other HexCoord) (Dir, error) {
	dx := other.X - hc.X
	dy := other.Y - hc.Y
	if dx == 0 && dy == 0 {
		return 0, fmt.Errorf("same: %v, %v", hc, other)
	}
	if dx < -1 || dx > 1 || dy < -1 || dy > 1 {
		return 0, fmt.Errorf("not adjacent: %v, %v", hc, other)
	}

	dir := adjacency_dirs[dy+1][dx+1]
	if dir == DIR_COUNT {
		return 0, fmt.Errorf("not adjacent: %v, %v", hc, other)
	}
	return dir, nil
}

func (hc HexCoord) IsStraightTo(other HexCoord) bool {
	return hc.X == other.X || hc.Y == other.Y || hc.X+hc.Y == other.X+other.Y
}

func (hc HexCoord) DistanceTo(other HexCoord) int32 {
	dx := other.X - hc.X
	dy := other.Y - hc.Y

	if Sign(dx) == Sign(dy) {
		return Abs(dx + dy)
	}
	return max(Abs(dx), Abs(dy))
}

func (hc HexCoord) Add(other HexCoord) HexCoord {
	return HexCoord{X: hc.X + other.X, Y: hc.Y + other.Y}
}

func (hc HexCoord) GetChunkCoord() ChunkCoord {
	return ChunkCoord{X: hc.X / ss.CHUNK_SIZE, Y: hc.Y / ss.CHUNK_SIZE}
}

func (hc HexCoord) CoordsWithinChunk() HexCoord {
	return HexCoord{X: hc.X % ss.CHUNK_SIZE, Y: hc.Y % ss.CHUNK_SIZE}
}
