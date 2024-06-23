package world

import (
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
)

type GroundHex struct {
	groundType     ss.GroundType
	resorceType    ss.ResourceType
	resourceAmount uint32
}

type Chunk struct {
	pos    utils.ChunkCoord
	ground [ss.CHUNK_SIZE * ss.CHUNK_SIZE]*GroundHex

	objects map[utils.HexCoord]WorldObject
	dirty   bool
}

func NewChunk(pos utils.ChunkCoord) *Chunk {
	c := &Chunk{
		pos:     pos,
		objects: make(map[utils.HexCoord]WorldObject),
		dirty:   true,
	}
	generateChunkGround(c)

	return c
}

func (c *Chunk) GetWorldObject(hex utils.HexCoord) (WorldObject, bool) {
	obj, ok := c.objects[hex]
	return obj, ok
}
func (c *Chunk) SetWorldObject(hex utils.HexCoord, obj WorldObject) {
	c.objects[hex] = obj
	c.dirty = true
}
func (c *Chunk) RemoveWorldObject(hex utils.HexCoord) {
	delete(c.objects, hex)
	c.dirty = true
}

func (c *Chunk) Draw(r *renderer.GameRenderer) {
	if c.dirty {
		r.ChunkRenderer.UpdateChunk(c.pos)
		c.dirty = false
	}

	r.ChunkRenderer.DrawToScreen(c.pos)
}
