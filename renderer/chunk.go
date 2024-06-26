package renderer

import (
	ss "hextopdown/settings"
	"hextopdown/utils"

	"github.com/veandco/go-sdl2/sdl"
)

const CHUNK_TEX_WIDTH = ss.CHUNK_SIZE*ss.HEX_WIDTH + (ss.CHUNK_SIZE-1)*ss.HEX_WIDTH/2
const CHUNK_TEX_HEIGHT = ss.CHUNK_SIZE*(ss.HEX_EDGE+ss.HEX_OFFSET) + ss.HEX_OFFSET

var chunkSize utils.ScreenCoord = utils.ScreenCoord{
	X: ss.CHUNK_SIZE*ss.HEX_WIDTH + (ss.CHUNK_SIZE-1)*ss.HEX_WIDTH/2,
	Y: ss.CHUNK_SIZE*(ss.HEX_EDGE+ss.HEX_OFFSET) + ss.HEX_OFFSET,
}

type ChunkRenderer struct {
	renderer *sdl.Renderer
	chunks   map[utils.ChunkCoord]*sdl.Texture
	resource [ss.RESOURCE_TYPE_COUNT]*sdl.Texture
}

func NewChunkRenderer(renderer *sdl.Renderer) *ChunkRenderer {
	return &ChunkRenderer{
		renderer: renderer,
		chunks:   make(map[utils.ChunkCoord]*sdl.Texture),
	}
}

func (r *ChunkRenderer) Destroy() {
	for _, tex := range r.chunks {
		if tex != nil {
			tex.Destroy()
		}
	}
	for _, tex := range r.resource {
		if tex != nil {
			tex.Destroy()
		}
	}
}

func (r *ChunkRenderer) DrawToScreen(chunk utils.ChunkCoord) {
	ch, ok := r.chunks[chunk]
	if !ok {
		return
	}

	c := chunk.TopLeftWorld().ToScreen()
	size := chunkSize.Mul(float32(utils.GetViewZoom()))

	if r.chunkOnScreen(c, size) {
		r.renderer.CopyF(ch, nil, &sdl.FRect{X: c.X, Y: c.Y, W: size.X, H: size.Y})
	}
}

func (r *ChunkRenderer) UpdateChunk(
	chunk utils.ChunkCoord,
	groundTypes [ss.CHUNK_SIZE * ss.CHUNK_SIZE]ss.GroundType,
	resourceTypes [ss.CHUNK_SIZE * ss.CHUNK_SIZE]ss.ResourceType,
	resourceAmounts [ss.CHUNK_SIZE * ss.CHUNK_SIZE]uint16,
) {
	ch, ok := r.chunks[chunk]
	if !ok {
		var err error
		ch, err = r.renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, CHUNK_TEX_WIDTH, CHUNK_TEX_HEIGHT)
		if err != nil {
			panic(err)
		}
		ch.SetBlendMode(sdl.BLENDMODE_BLEND)
		r.chunks[chunk] = ch
	}

	r.renderer.SetRenderTarget(ch)

	r.renderer.SetDrawColor(0, 0, 0, 0)
	r.renderer.Clear()

	for hy := int32(0); hy < ss.CHUNK_SIZE; hy++ {
		xo := float32(hy) * ss.HEX_WIDTH / 2
		y := float32(hy) * (ss.HEX_EDGE + ss.HEX_OFFSET) //+ ss.HEX_OFFSET

		for hx := int32(0); hx < ss.CHUNK_SIZE; hx++ {
			x := xo + float32(hx)*ss.HEX_WIDTH

			switch groundTypes[hy*ss.CHUNK_SIZE+hx] {
			case ss.GROUND_TYPE_GROUND:
				r.renderer.SetDrawColor(128, 96, 0, 127)
			case ss.GROUND_TYPE_WATER:
				r.renderer.SetDrawColor(127, 127, 255, 127)
			default:
				r.renderer.SetDrawColor(255, 0, 255, 127)
			}
			r.renderer.FillRectF(&sdl.FRect{
				X: x,
				Y: y + ss.HEX_OFFSET/2,
				W: ss.HEX_WIDTH,
				H: ss.HEX_OFFSET*3 + 1,
			})

			r.renderer.SetDrawColor(96, 96, 96, 255)

			if hy == 0 {
				r.renderer.SetDrawColor(140, 140, 140, 255)
			} else {
				r.renderer.SetDrawColor(96, 96, 96, 255)
			}
			r.renderer.DrawLineF(x, y+ss.HEX_OFFSET, x+ss.HEX_WIDTH/2, y)     // top left
			r.renderer.DrawLineF(x, y+ss.HEX_OFFSET+1, x+ss.HEX_WIDTH/2, y+1) // top left

			if hx == ss.CHUNK_SIZE-1 || hy == 0 {
				r.renderer.SetDrawColor(128, 128, 128, 255)
			} else {
				r.renderer.SetDrawColor(96, 96, 96, 255)
			}
			r.renderer.DrawLineF(x+ss.HEX_WIDTH/2, y, x+ss.HEX_WIDTH, y+ss.HEX_OFFSET)     // top right
			r.renderer.DrawLineF(x+ss.HEX_WIDTH/2, y+1, x+ss.HEX_WIDTH, y+ss.HEX_OFFSET+1) // top right

			if hx == 0 {
				r.renderer.SetDrawColor(128, 128, 128, 255)
			} else {
				r.renderer.SetDrawColor(96, 96, 96, 255)
			}
			r.renderer.DrawLineF(x, y+ss.HEX_OFFSET, x, y+ss.HEX_EDGE+ss.HEX_OFFSET)     // left
			r.renderer.DrawLineF(x+1, y+ss.HEX_OFFSET, x+1, y+ss.HEX_EDGE+ss.HEX_OFFSET) // left

			if resourceTypes[hy*ss.CHUNK_SIZE+hx] != ss.RESOURCE_TYPE_COUNT {
				tex := r.resource[resourceTypes[hy*ss.CHUNK_SIZE+hx]]
				if tex == nil {
					continue
				}
				r.renderer.CopyF(tex, nil, &sdl.FRect{
					X: x,
					Y: y,
					W: ss.HEX_WIDTH,
					H: ss.HEX_EDGE * 2,
				})
			}
		}
	}

	r.renderer.SetRenderTarget(nil)
}

func (r *ChunkRenderer) chunkOnScreen(c utils.ScreenCoord, size utils.ScreenCoord) bool {
	return isOnScreenBox(c, c.Add(size))
}

func (r *ChunkRenderer) GetVisibleChunkCoords() (utils.ChunkCoord, utils.ChunkCoord) {
	str := utils.ScreenCoord{X: 1, Y: 0}.PctPosToScreen().ToWorld().ToHex().GetChunkCoord()
	sbl := utils.ScreenCoord{X: 0, Y: 1}.PctPosToScreen().ToWorld().ToHex().GetChunkCoord()
	return utils.ChunkCoord{X: sbl.X, Y: str.Y}, utils.ChunkCoord{X: str.X, Y: sbl.Y}
}

func (r *ChunkRenderer) LoadTextures() {
	r.resource[ss.RESOURCE_TYPE_IRON] = loadCachedTexture("ground/res_iron_ore", r.renderer)
}
