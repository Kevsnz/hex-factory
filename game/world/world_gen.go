package world

import (
	ss "hextopdown/settings"
	"hextopdown/utils"
	"math"
	"math/rand"
)

const RING1_R = 3 * ss.CHUNK_SIZE * ss.HEX_WIDTH
const RING1_R_MIN = 1 * ss.CHUNK_SIZE * ss.HEX_WIDTH
const ORE_PATCH_R = 4 * ss.HEX_WIDTH

type orePatch struct {
	Center   utils.WorldCoord
	Radius   float32
	Richness float32
	OreType  ss.ResourceType
}

type WorldGen struct {
	rng        *rand.Rand
	orePatches []*orePatch
}

func NewWorldGen() *WorldGen {
	rng := rand.New(rand.NewSource(99))
	orePatches := []*orePatch{
		generateOrePatch(rng, ss.RESOURCE_TYPE_IRON),
		generateOrePatch(rng, ss.RESOURCE_TYPE_IRON),
	}
	return &WorldGen{
		rng:        rng,
		orePatches: orePatches,
	}
}

func (wg *WorldGen) GenerateChunkGround(chunk *Chunk, pos utils.ChunkCoord) {
	for hy := int32(0); hy < ss.CHUNK_SIZE; hy++ {
		for hx := int32(0); hx < ss.CHUNK_SIZE; hx++ {
			i := hy*ss.CHUNK_SIZE + hx
			chunk.groundTypes[i] = ss.GroundType(rand.Uint32()) % ss.GROUND_TYPE_COUNT

			hex := chunk.pos.ToHexCoord().Add(utils.HexCoord{X: hx, Y: hy})
			wc := hex.CenterToWorld()

			orePatch := wg.getClosestOrePatch(wc)
			if orePatch == nil || orePatch.Center.DistanceSqTo(wc) > float64(orePatch.Radius*orePatch.Radius) {
				chunk.groundResTypes[i] = ss.RESOURCE_TYPE_COUNT
				chunk.groundResAmounts[i] = 0
				continue
			}

			chunk.groundResTypes[i] = orePatch.OreType
			chunk.groundResAmounts[i] = uint16(orePatch.Richness)
		}
	}
}

func (wg *WorldGen) getClosestOrePatch(pos utils.WorldCoord) *orePatch {
	var closestPatch *orePatch
	closestDist := math.MaxFloat32

	for _, patch := range wg.orePatches {
		dist := pos.DistanceSqTo(patch.Center)
		if dist < closestDist {
			closestDist = dist
			closestPatch = patch
		}
	}
	return closestPatch
}

func generateOrePatch(rng *rand.Rand, oreType ss.ResourceType) *orePatch {
	alpha := rng.Float64() * math.Pi * 2
	d := rng.Float64()*(RING1_R-RING1_R_MIN) + RING1_R_MIN
	x := math.Cos(alpha) * d
	y := math.Sin(alpha) * d
	cc := utils.WorldCoord{X: x, Y: y}

	return &orePatch{
		Center:   cc,
		Radius:   ORE_PATCH_R,
		Richness: 750 + rng.Float32()*500,
		OreType:  oreType,
	}
}
