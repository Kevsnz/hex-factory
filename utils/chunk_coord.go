package utils

import ss "hextopdown/settings"

type ChunkCoord struct {
	X, Y int32
}

func (cc ChunkCoord) ToHexCoord() HexCoord {
	return HexCoord{X: cc.X * ss.CHUNK_SIZE, Y: cc.Y * ss.CHUNK_SIZE}
}
