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
}

func (r *ChunkRenderer) DrawToScreen(chunk utils.ChunkCoord) {
	ch, ok := r.chunks[chunk]
	if !ok {
		return
	}

	c := chunk.TopLeftWorld().ToScreen()
	size := chunkSize.Mul(float32(utils.GetViewZoom()))

	if r.chunkOnScreen(c, size) {
		r.renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
		r.renderer.CopyF(ch, nil, &sdl.FRect{X: c.X, Y: c.Y, W: size.X, H: size.Y})
	}
}

func (r *ChunkRenderer) UpdateChunk(chunk utils.ChunkCoord) {
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

	r.renderer.SetDrawColor(96, 96, 96, 255)

	for hy := int32(0); hy < ss.CHUNK_SIZE; hy++ {
		xo := float32(hy) * ss.HEX_WIDTH / 2
		yo := float32(hy)*(ss.HEX_EDGE+ss.HEX_OFFSET) + ss.HEX_OFFSET

		for hx := int32(0); hx < ss.CHUNK_SIZE; hx++ {
			x := xo + float32(hx)*ss.HEX_WIDTH
			r.renderer.DrawLineF(x, yo, x+ss.HEX_WIDTH/2, yo-ss.HEX_OFFSET)              // top left
			r.renderer.DrawLineF(x+ss.HEX_WIDTH/2, yo-ss.HEX_OFFSET, x+ss.HEX_WIDTH, yo) // top right
			r.renderer.DrawLineF(x, yo, x, yo+ss.HEX_EDGE)                               // left
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
