package world

import (
	ss "hextopdown/settings"
	"math/rand"
)

func generateChunkGround(chunk *Chunk) {
	for i := 0; i < ss.CHUNK_SIZE*ss.CHUNK_SIZE; i++ {
		// TODOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
		chunk.groundTypes[i] = ss.GroundType(rand.Uint32()) % ss.GROUND_TYPE_COUNT
		if rand.Uint32()%5 == 0 {
			chunk.groundResTypes[i] = ss.RESOURCE_TYPE_IRON
			chunk.groundResAmounts[i] = 100
		} else {
			chunk.groundResTypes[i] = ss.RESOURCE_TYPE_COUNT
			chunk.groundResAmounts[i] = 0
		}
	}
}
