package world

import (
	ss "hextopdown/settings"
	"hextopdown/utils"
	"math"
	"math/rand"
)

const RING1_R_MAX = 3 * ss.CHUNK_SIZE * ss.HEX_WIDTH
const RING1_R_MIN = 1 * ss.CHUNK_SIZE * ss.HEX_WIDTH
const ORE_PATCH_R_MAX = 15 * ss.HEX_WIDTH
const ORE_PATCH_R_MIN = 13 * ss.HEX_WIDTH
const ORE_AMOUNT_VARIANCE = 0.125

type orePatch struct {
	Center   utils.WorldCoord
	Radius   float32
	Richness float32
	OreType  ss.ResourceType
}

type WorldGen struct {
	seed       int64
	rng        *rand.Rand
	orePatches []*orePatch
}

func NewWorldGen(seed int64) *WorldGen {
	rng := rand.New(rand.NewSource(seed))

	orePatches := []*orePatch{
		generateOrePatch(rng, ss.RESOURCE_TYPE_IRON),
		generateOrePatch(rng, ss.RESOURCE_TYPE_IRON),
	}

	return &WorldGen{
		seed:       seed,
		rng:        rng,
		orePatches: orePatches,
	}
}

func (wg *WorldGen) getChunkRng(pos utils.ChunkCoord) *rand.Rand {
	s1 := rand.New(rand.NewSource(wg.seed)).Uint64() * uint64(pos.X)
	s2 := rand.New(rand.NewSource(int64(s1))).Uint64() * uint64(pos.Y)
	return rand.New(rand.NewSource(int64(s2)))
}

func (wg *WorldGen) GenerateChunk(chunk *Chunk, pos utils.ChunkCoord) {
	rng := wg.getChunkRng(pos)

	for hy := int32(0); hy < ss.CHUNK_SIZE; hy++ {
		for hx := int32(0); hx < ss.CHUNK_SIZE; hx++ {
			i := hy*ss.CHUNK_SIZE + hx
			chunk.groundTypes[i] = ss.GroundType(rand.Uint32()) % ss.GROUND_TYPE_COUNT

			hex := chunk.pos.ToHexCoord().Add(utils.HexCoord{X: hx, Y: hy})
			wc := hex.CenterToWorld()

			rng1 := rng.Float64()
			orePatch := wg.getClosestOrePatch(wc)
			if orePatch == nil {
				chunk.groundResTypes[i] = ss.RESOURCE_TYPE_COUNT
				chunk.groundResAmounts[i] = 0
				continue
			}
			dist := math.Sqrt(orePatch.Center.DistanceSqTo(wc))
			if dist > float64(orePatch.Radius) {
				chunk.groundResTypes[i] = ss.RESOURCE_TYPE_COUNT
				chunk.groundResAmounts[i] = 0
				continue
			}

			dist_mult := min(1-(dist/float64(orePatch.Radius))+0.25, 1.0)
			amount := float64(orePatch.Richness) * dist_mult
			amount = amount * ((1 - ORE_AMOUNT_VARIANCE/2) + rng1*ORE_AMOUNT_VARIANCE)
			amount_u16 := max(1, uint16(math.Round(amount)))

			chunk.groundResTypes[i] = orePatch.OreType
			chunk.groundResAmounts[i] = amount_u16
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
	d := rng.Float64()*(RING1_R_MAX-RING1_R_MIN) + RING1_R_MIN
	x := math.Cos(alpha) * d
	y := math.Sin(alpha) * d
	cc := utils.WorldCoord{X: x, Y: y}

	return &orePatch{
		Center:   cc,
		Radius:   ORE_PATCH_R_MIN + rng.Float32()*(ORE_PATCH_R_MAX-ORE_PATCH_R_MIN),
		Richness: 750 + rng.Float32()*500,
		OreType:  oreType,
	}
}
