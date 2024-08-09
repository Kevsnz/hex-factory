package world

import (
	gd "hextopdown/game/gamedata"
	"hextopdown/renderer"
	ss "hextopdown/settings"
	"hextopdown/utils"
)

type Chunk struct {
	pos              utils.ChunkCoord
	groundResTypes   [ss.CHUNK_SIZE * ss.CHUNK_SIZE]ss.ResourceType
	groundResAmounts [ss.CHUNK_SIZE * ss.CHUNK_SIZE]uint16
	groundTypes      [ss.CHUNK_SIZE * ss.CHUNK_SIZE]ss.GroundType

	objects map[utils.HexCoord]WorldObject
	dirty   bool
}

func NewChunk(pos utils.ChunkCoord, wg *WorldGen) *Chunk {
	c := &Chunk{
		pos:     pos,
		objects: make(map[utils.HexCoord]WorldObject),
		dirty:   true,
	}
	wg.GenerateChunk(c, pos)

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
		r.ChunkRenderer.UpdateChunk(c.pos, c.groundTypes, c.groundResTypes, c.groundResAmounts)
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

func (c *Chunk) GetResourceTypeAt(hex utils.HexCoord) (ss.ResourceType, uint16, bool) {
	pos := hex.CoordsWithinChunk()
	res := c.groundResTypes[pos.X+pos.Y*ss.CHUNK_SIZE]
	amt := c.groundResAmounts[pos.X+pos.Y*ss.CHUNK_SIZE]
	return res, amt, res != ss.RESOURCE_TYPE_COUNT
}

func (c *Chunk) ExtractResourceAt(hex utils.HexCoord) ss.ItemType {
	pos := hex.CoordsWithinChunk()
	res := c.groundResTypes[pos.X+pos.Y*ss.CHUNK_SIZE]
	if res == ss.RESOURCE_TYPE_COUNT {
		panic("resource not found")
	}

	c.groundResAmounts[pos.X+pos.Y*ss.CHUNK_SIZE]--
	if c.groundResAmounts[pos.X+pos.Y*ss.CHUNK_SIZE] == 0 {
		c.groundResTypes[pos.X+pos.Y*ss.CHUNK_SIZE] = ss.RESOURCE_TYPE_COUNT
		c.dirty = true
	}

	return gd.ExtractionResourceItems[res]
}
