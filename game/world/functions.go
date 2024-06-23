package world

import (
	ss "hextopdown/settings"
)

func generateChunkGround(chunk *Chunk) {
	for i := 0; i < ss.CHUNK_SIZE*ss.CHUNK_SIZE; i++ {
		// TODOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
		chunk.ground[i] = &GroundHex{
			groundType:     ss.GROUND_TYPE_GROUND,
			resorceType:    ss.RESOURCE_TYPE_COUNT,
			resourceAmount: 0,
		}
	}
}
