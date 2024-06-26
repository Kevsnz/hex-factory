package world

import (
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
)

type GroundHex struct {
	resorceType    ss.ResourceType
	resourceAmount uint32
}

type Chunk struct {
	pos         utils.ChunkCoord
	ground      [ss.CHUNK_SIZE * ss.CHUNK_SIZE]*GroundHex
	groundTypes [ss.CHUNK_SIZE * ss.CHUNK_SIZE]ss.GroundType

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
	if hex.GetChunkCoord() != c.pos {
		panic("invalid coordinates")
	}

	obj, ok := c.objects[hex]
	return obj, ok
}
func (c *Chunk) SetWorldObject(hex utils.HexCoord, obj WorldObject) {
	if hex.GetChunkCoord() != c.pos {
		panic("invalid coordinates")
	}

	c.objects[hex] = obj
	c.dirty = true
}
func (c *Chunk) RemoveWorldObject(hex utils.HexCoord) {
	if hex.GetChunkCoord() != c.pos {
		panic("invalid coordinates")
	}

	delete(c.objects, hex)
	c.dirty = true
}

func (c *Chunk) DrawGround(r *renderer.GameRenderer) {
	if c.dirty {
		r.ChunkRenderer.UpdateChunk(c.pos, c.groundTypes)
		c.dirty = false
	}

	r.ChunkRenderer.DrawToScreen(c.pos)
}

func (c *Chunk) DrawObjectsGroundLevel(r *renderer.GameRenderer) {
	for hex, obj := range c.objects {
		if hex == obj.GetPos() {
			obj.DrawGroundLevel(r)
		}
	}
}

func (c *Chunk) DrawObjectsOnGroundLevel(r *renderer.GameRenderer) {
	for hex, obj := range c.objects {
		if hex == obj.GetPos() {
			obj.DrawOnGroundLevel(r)
		}
	}
}

func (c *Chunk) DrawItems(r *renderer.GameRenderer) {
	for _, obj := range c.objects {
		if !r.IsHexOnScreen(obj.GetPos()) {
			continue
		}
		switch drawer := obj.(type) {
		case ItemDrawer:
			drawer.DrawItems(r)
		}
	}
}

func (c *Chunk) GetWorldObjects() map[utils.HexCoord]WorldObject {
	return c.objects
}
