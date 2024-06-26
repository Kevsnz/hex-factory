package world

import (
	ss "hextopdown/settings"
	"math/rand"
)

func generateChunkGround(chunk *Chunk) {
	for i := 0; i < ss.CHUNK_SIZE*ss.CHUNK_SIZE; i++ {
		// TODOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO
		chunk.groundTypes[i] = ss.GroundType(rand.Uint32()) % ss.GROUND_TYPE_COUNT
		chunk.ground[i] = &GroundHex{
			resorceType:    ss.RESOURCE_TYPE_COUNT,
			resourceAmount: 0,
		}
	}
}
