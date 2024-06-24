package utils

import ss "hextopdown/settings"

type ChunkCoord struct {
	X, Y int32
}

func (cc ChunkCoord) ToHexCoord() HexCoord {
	return HexCoord{X: cc.X * ss.CHUNK_SIZE, Y: cc.Y * ss.CHUNK_SIZE}
}

func (cc ChunkCoord) TopLeftWorld() WorldCoord {
	return WorldCoord{
		float64(cc.X*ss.CHUNK_SIZE)*ss.HEX_WIDTH + float64(cc.Y*ss.CHUNK_SIZE)*(ss.HEX_WIDTH/2),
		float64(cc.Y*ss.CHUNK_SIZE)*(ss.HEX_EDGE+ss.HEX_OFFSET) - ss.HEX_OFFSET,
	}
}
