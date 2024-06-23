package renderer

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"path"

	ss "hextopdown/settings"
	"hextopdown/settings/strings"
	"hextopdown/utils"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type GameRenderer struct {
	renderer                  *sdl.Renderer
	ChunkRenderer             *ChunkRenderer
	stringManager             *StringManager
	beltTextures              [ss.BELT_TYPE_COUNT]*sdl.Texture
	beltAnimationTextures     [ss.BELT_TYPE_COUNT]*sdl.Texture
	beltOnGroundTextures      [ss.BELT_ON_COUNT]*sdl.Texture
	objectGroundTextures      [ss.OBJECT_TYPE_COUNT]*sdl.Texture
	objectOnGroundTextures    [ss.OBJECT_TYPE_COUNT]*sdl.Texture
	objectGroundDirTextures   [ss.OBJECT_TYPE_COUNT][utils.DIR_COUNT]*sdl.Texture
	objectOnGroundDirTextures [ss.OBJECT_TYPE_COUNT][utils.DIR_COUNT]*sdl.Texture
	itemTextures              [ss.ITEM_TYPE_COUNT]*sdl.Texture
	arrowTextures             [2]*sdl.Texture
	iconsItems                *sdl.Texture
	decalTextures             [DECAL_COUNT]*sdl.Texture
	timeMs                    uint64
}

func NewGameRenderer(window *sdl.Window) *GameRenderer {
	renderer, err := sdl.CreateRenderer(window, 0, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	sm := NewStringManager()
	sm.Prerender(renderer)

	return &GameRenderer{
		renderer:      renderer,
		ChunkRenderer: NewChunkRenderer(renderer),
		stringManager: sm,
	}
}

func (r *GameRenderer) Destroy() {
	for _, tex := range r.beltTextures {
		tex.Destroy()
	}
	for _, tex := range r.beltAnimationTextures {
		tex.Destroy()
	}
	for _, tex := range r.beltOnGroundTextures {
		tex.Destroy()
	}
	for _, tex := range r.itemTextures {
		tex.Destroy()
	}
	for _, tex := range r.arrowTextures {
		tex.Destroy()
	}
	for _, tex := range r.objectGroundTextures {
		tex.Destroy()
	}
	for _, tex := range r.objectOnGroundTextures {
		tex.Destroy()
	}
	for _, texs := range r.objectGroundDirTextures {
		for _, tex := range texs {
			tex.Destroy()
		}
	}
	for _, texs := range r.objectOnGroundDirTextures {
		for _, tex := range texs {
			tex.Destroy()
		}
	}
	for _, tex := range r.decalTextures {
		tex.Destroy()
	}
	r.iconsItems.Destroy()
	r.stringManager.Destroy()
	r.ChunkRenderer.Destroy()
	r.renderer.Destroy()
}

func (r *GameRenderer) StartNewFrame(timeMs uint64) {
	r.timeMs = timeMs
	r.renderer.SetDrawColor(44, 48, 48, 255)
	r.renderer.Clear()
}

func (r *GameRenderer) DrawScreen() {
	r.drawHexGrid()
}

func (r *GameRenderer) DrawViewTarget(pos utils.WorldCoordInterpolated) {
	drawPos := pos.GetInterpolatedPos(r.timeMs, ss.TICK_DT)
	c := drawPos.ToScreen()
	if !isOnScreenRadius(c, 10) {
		return
	}
	r.renderer.SetDrawColor(127, 127, 255, 255)
	r.renderer.FillRectF(fRectFromScreenOffset(c, -10, -10, 20, 20))
}

func (r *GameRenderer) drawHexGrid() {
	hex1 := utils.ScreenCoord{X: 0, Y: 0}.ToWorld().ToHex()
	hex2 := utils.ScreenCoord{X: RES_X, Y: 0}.ToWorld().ToHex()
	hex3 := utils.ScreenCoord{X: 0, Y: RES_Y}.ToWorld().ToHex()
	w := utils.GetZoomedHexWidth()
	e := utils.GetZoomedHexEdge()
	o := utils.GetZoomedHexOffset()

	r.renderer.SetDrawColor(96, 96, 96, 255)
	s := hex1.LeftTopToWorld().ToScreen()
	for hy := int32(0); hy <= hex3.Y-hex1.Y+1; hy++ {
		yo := float32(hy) * (e + o)
		xo := float32(hy&1) * w / 2
		for hx := int32(-1); hx <= hex2.X-hex1.X+1; hx++ {
			r.renderer.DrawLineF(s.X+float32(hx)*w+xo, s.Y+yo, s.X+xo+float32(hx)*w+w/2, s.Y+yo-o)   // top left
			r.renderer.DrawLineF(s.X+float32(hx)*w+xo+w/2, s.Y+yo-o, s.X+xo+float32(hx)*w+w, s.Y+yo) // top right
			r.renderer.DrawLineF(s.X+float32(hx)*w+xo, s.Y+yo, s.X+float32(hx)*w+xo, s.Y+yo+e)       // left
		}
	}
}

func (r *GameRenderer) DrawHexCenter(hex utils.HexCoord) {
	c := hexCenterToScreen(hex)
	if !isOnScreen(c) {
		return
	}
	r.renderer.SetDrawColor(0, 96, 0, 255)
	rect := fRectFromScreen(c, 1, 1)
	r.renderer.FillRectF(rect)
}

func (r *GameRenderer) DrawBeltConnectionIncoming(hex utils.HexCoord, dir utils.Dir, left bool, start, end float64) {
	c := hexCenterToScreen(hex)
	z := float32(utils.GetViewZoom())

	outOffset := radiusOffsets[dir].Mul(z)
	laneOffset := lanesOffsetsLeft[dir].Mul(z)

	if left {
		laneOffset = laneOffset.Inverse()
	}

	if left {
		r.renderer.SetDrawColor(192, 0, 0, 255)
	} else {
		r.renderer.SetDrawColor(0, 0, 192, 255)
	}

	s := c.Add(outOffset.Mul(1 - 2*float32(start))).Add(laneOffset)
	e := c.Add(outOffset.Mul(1 - 2*float32(end))).Add(laneOffset)

	if isOnScreenRadius(s, 2) {
		r.renderer.FillRectF(fRectFromScreenOffset(s, -2, -2, 4, 4))
	}
	if isOnScreenBox(s, e) {
		r.renderer.DrawLineF(s.X, s.Y, e.X, e.Y)
	}
}

func (r *GameRenderer) DrawBeltConnectionsOutgoing(hex utils.HexCoord, dir utils.Dir) {
	c := hexCenterToScreen(hex)
	z := float32(utils.GetViewZoom())

	outOffset := radiusOffsets[dir].Mul(z)
	laneOffset := lanesOffsetsLeft[dir].Mul(z)

	r.renderer.SetDrawColor(192, 0, 0, 255)
	// Left lane
	s := c.Add(laneOffset)
	e := s.Add(outOffset)

	if isOnScreenRadius(s, 2) {
		r.renderer.FillRectF(fRectFromScreenOffset(s, -2, -2, 4, 4))
	}
	if isOnScreenBox(s, e) {
		r.renderer.DrawLineF(s.X, s.Y, e.X, e.Y)
	}

	// Right lane
	laneOffset = laneOffset.Inverse()
	r.renderer.SetDrawColor(0, 0, 192, 255)

	s = c.Add(laneOffset)
	e = s.Add(outOffset)

	if isOnScreenRadius(s, 2) {
		r.renderer.FillRectF(fRectFromScreenOffset(s, -2, -2, 4, 4))
	}
	if isOnScreenBox(s, e) {
		r.renderer.DrawLineF(s.X, s.Y, e.X, e.Y)
	}
}

func (r *GameRenderer) DrawAnimatedBelt(hex utils.HexCoord, beltType ss.BeltType, speed float64) {
	tex := r.beltAnimationTextures[beltType]
	if tex == nil {
		panic("no animation texture for belt type")
	}

	c := hexCenterToScreen(hex)
	if !isOnScreenRadius(c, utils.GetZoomedDimension(ss.BELT_DRAW_R)) {
		return
	}
	_, tsegm := math.Modf(speed * float64(r.timeMs) / 1000)
	frame := int32(math.Floor(tsegm*ss.ANIM_BELT_STEPS)) % ss.ANIM_BELT_FRAMES

	e := utils.GetZoomedHexEdge()
	r.renderer.CopyF(
		tex,
		&sdl.Rect{X: 0, Y: frame * TEXTURE_SIZE_HEX, W: TEXTURE_SIZE_HEX, H: TEXTURE_SIZE_HEX},
		fRectFromScreenOffset(c, -e, -e, 2*e, 2*e),
	)
}

func (r *GameRenderer) DrawBeltOnGround(hex utils.HexCoord, beltType ss.BeltType) {
	if beltType == ss.BELT_ON_COUNT {
		return
	}

	c := hexCenterToScreen(hex)
	if !isOnScreenRadius(c, utils.GetZoomedDimension(ss.BELT_DRAW_R)) {
		return
	}
	typeFlip := beltOnFlipMapping[beltType]

	tex := r.beltOnGroundTextures[typeFlip.type1]
	if tex == nil {
		panic(fmt.Sprintf("no texture for belt type %d", typeFlip.type1))
	}
	e := utils.GetZoomedHexEdge()
	r.renderer.CopyExF(tex, nil, fRectFromScreenOffset(c, -e, -e, 2*e, 2*e), 0, nil, typeFlip.flip)
}

func (r *GameRenderer) DrawObjectGround(pos utils.WorldCoord, objectType ss.ObjectType, shape utils.Shape, dir utils.Dir) {
	c := pos.ToScreen()
	z := float32(utils.GetViewZoom())
	sp := GetShapeParam(shape, dir)

	c = c.Sub(sp.Offset.Mul(z))
	size := sp.Size.Mul(z)

	if !isOnScreenBox(c, c.Add(size)) {
		return
	}

	tex := r.objectGroundTextures[objectType]
	if tex == nil {
		tex = r.objectGroundDirTextures[objectType][dir]
		if tex == nil {
			tex = r.objectGroundDirTextures[objectType][dir.Reverse()] // mirrored shape
			if tex == nil {
				r.DrawString(strings.STRING_NOTEXTURE, c.Add(size.Div(2)), TEXT_ALIGN_CENTER)
				return
			}
		}
	}

	r.renderer.CopyF(tex, nil, fRectFromScreen(c, size.X, size.Y))
}

func (r *GameRenderer) DrawObjectOnGround(pos utils.WorldCoord, objectType ss.ObjectType, shape utils.Shape, dir utils.Dir) {
	c := pos.ToScreen()
	z := float32(utils.GetViewZoom())
	sp := GetShapeParam(shape, dir)

	c = c.Sub(sp.Offset.Mul(z))
	size := sp.Size.Mul(z)

	if !isOnScreenBox(c, c.Add(size)) {
		return
	}

	tex := r.objectOnGroundTextures[objectType]
	if tex == nil {
		tex = r.objectOnGroundDirTextures[objectType][dir]
		if tex == nil {
			tex = r.objectOnGroundDirTextures[objectType][dir.Reverse()] // mirrored shape
			if tex == nil {
				r.DrawString(strings.STRING_NOTEXTURE, c.Add(size.Div(2)), TEXT_ALIGN_CENTER)
				return
			}
		}
	}

	r.renderer.CopyF(tex, nil, fRectFromScreen(c, size.X, size.Y))
}

func (r *GameRenderer) DrawItem(pos utils.WorldCoordInterpolated, itemType ss.ItemType) {
	drawPos := pos.GetInterpolatedPos(r.timeMs, ss.TICK_DT)
	s := drawPos.ToScreen()
	idr := utils.GetZoomedDimension(ss.ITEM_DRAW_R)
	if !isOnScreenRadius(s, idr) {
		return
	}

	tex := r.itemTextures[itemType]
	if tex == nil {
		r.renderer.SetDrawColor(255, 0, 255, 255)
		r.renderer.FillRectF(fRectFromScreenOffset(s, -idr, -idr, 2*idr, 2*idr))
		return
	}
	r.renderer.CopyF(tex, nil, fRectFromScreenOffset(s, -idr, -idr, 2*idr, 2*idr))
}

func (r *GameRenderer) DrawArrow(pctX, pctY float32, dir utils.Dir) {
	p := utils.ScreenCoord{X: pctX * RES_X, Y: pctY * RES_Y}
	s := utils.ScreenCoord{X: ss.FONT_SIZE_PCT * 2, Y: ss.FONT_SIZE_PCT * 2}.PctScaleToScreen()
	idxFlip := arrowDirMapping[dir]

	r.renderer.CopyExF(r.arrowTextures[idxFlip.idx], nil, fRectFromScreen(p, s.X, s.Y), 0, nil, idxFlip.flip)
}

func (r *GameRenderer) DrawItemIconWorld(pos utils.WorldCoord, sizeHexes float32, itemType ss.ItemType) {
	i, ok := iconItemList[itemType]
	if !ok {
		panic("invalid item type")
	}

	c := pos.ToScreen()
	s := utils.GetZoomedHexWidth()
	if !isOnScreenRadius(c, s*sizeHexes) {
		return
	}
	ss := utils.ScreenCoord{X: s, Y: s}.Mul(sizeHexes)

	_, _, w, _, err := r.iconsItems.Query()
	if err != nil {
		panic("Failed to query icons items texture")
	}
	w /= TEXTURE_ICON_SIZE

	r.renderer.CopyF(r.iconsItems, &sdl.Rect{
		X: (int32(i) % w) * TEXTURE_ICON_SIZE,
		Y: (int32(i) / w) * TEXTURE_ICON_SIZE,
		W: TEXTURE_ICON_SIZE,
		H: TEXTURE_ICON_SIZE,
	}, fRectFromScreen(c.Sub(ss.Div(2)), ss.X, ss.Y))
}
func (r *GameRenderer) DrawItemIconScreen(pos utils.ScreenCoord, size float32, itemType ss.ItemType) {
	_, _, w, _, err := r.iconsItems.Query()
	if err != nil {
		panic("Failed to query icons items texture")
	}
	w /= TEXTURE_ICON_SIZE

	i, ok := iconItemList[itemType]
	if !ok {
		panic("invalid item type")
	}

	r.renderer.CopyF(
		r.iconsItems,
		&sdl.Rect{
			X: (int32(i) % w) * TEXTURE_ICON_SIZE,
			Y: (int32(i) / w) * TEXTURE_ICON_SIZE,
			W: TEXTURE_ICON_SIZE,
			H: TEXTURE_ICON_SIZE,
		},
		&sdl.FRect{X: pos.X, Y: pos.Y, W: size, H: size},
	)
}

func (r *GameRenderer) DrawDecal(pos utils.WorldCoord, sizeHexes float32, decal DecalId) {
	c := pos.ToScreen()

	s := utils.GetZoomedHexWidth()
	if !isOnScreenRadius(c, s*sizeHexes) {
		return
	}
	ss := utils.ScreenCoord{X: s, Y: s}.Mul(sizeHexes)

	tex := r.decalTextures[decal]
	if tex == nil {
		panic(fmt.Sprintf("no texture for decal type %d", decal))
	}

	r.renderer.CopyF(tex, nil, fRectFromScreen(c.Sub(ss.Div(2)), ss.X, ss.Y))
}

func (r *GameRenderer) DrawProgressBar(pos utils.WorldCoord, sizeHexes float32, part, max uint32) {
	c := pos.ToScreen()

	s := utils.GetZoomedHexWidth()
	if !isOnScreenRadius(c, s*sizeHexes) {
		return
	}
	ss := utils.ScreenCoord{X: s, Y: s / 8}.Mul(sizeHexes)
	c.Y += s / 2

	r.renderer.SetDrawColor(0, 0, 0, 255)
	r.renderer.FillRectF(fRectFromScreen(c.Sub(ss.Div(2)), ss.X, ss.Y))
	r.renderer.SetDrawColor(0, 255, 0, 255)
	r.renderer.FillRectF(fRectFromScreen(c.Sub(ss.Div(2)), ss.X*float32(part)/float32(max), ss.Y))
}

func (r *GameRenderer) Finish() {
	r.renderer.Present()
}

func (r *GameRenderer) IsHexOnScreen(coord utils.HexCoord) bool {
	s := hexCenterToScreen(coord)
	return isOnScreenRadius(s, utils.GetZoomedDimension(ss.HEX_DRAW_R))
}

func (r *GameRenderer) DrawConnectionHexes(hex1, hex2 utils.HexCoord) {
	if hex1 == hex2 {
		return
	}
	r.renderer.SetDrawColor(255, 255, 0, 255)

	p1 := hexCenterToScreen(hex1)
	p2 := hexCenterToScreen(hex2)

	r.drawDashedLine(p1, p2)
}

func (r *GameRenderer) drawDashedLine(p1, p2 utils.ScreenCoord) {
	dp := p2.Sub(p1)
	dSq := dp.LengthSq()
	if dSq < 4*ss.DASH_LEN*ss.DASH_LEN {
		r.renderer.DrawLineF(p1.X, p1.Y, p2.X, p2.Y)
		return
	}
	d := math.Sqrt(float64(dSq))

	swap := false
	if utils.AbsF32(dp.Y) > utils.AbsF32(dp.X) {
		swap = true
		dp = dp.SwapXY()
		p1 = p1.SwapXY()
		p2 = p2.SwapXY()
	}

	dashes := int(math.Round(d / (2 * ss.DASH_LEN)))
	dStep := dp.Div(float32(dashes*2 + 1))
	for i := 0; i < dashes; i++ {
		pd1 := p1.Add(dStep.Mul(float32(i * 2)))
		pd2 := p1.Add(dStep.Mul(float32(i*2 + 1)))

		if swap {
			r.renderer.DrawLineF(pd1.Y, pd1.X, pd2.Y, pd2.X)
		} else {
			r.renderer.DrawLineF(pd1.X, pd1.Y, pd2.X, pd2.Y)
		}
	}

	pd1 := p1.Add(dStep.Mul(float32(dashes * 2)))
	if swap {
		r.renderer.DrawLineF(pd1.Y, pd1.X, p2.Y, p2.X)
	} else {
		r.renderer.DrawLineF(pd1.X, pd1.Y, p2.X, p2.Y)
	}
}

func (r *GameRenderer) DrawWorldLine(p1, p2 utils.WorldCoord) {
	s1 := p1.ToScreen()
	s2 := p2.ToScreen()

	r.renderer.SetDrawColor(255, 0, 0, 255)
	r.renderer.DrawLineF(s1.X, s1.Y, s2.X, s2.Y)
}

func hexCenterToScreen(hex utils.HexCoord) utils.ScreenCoord {
	return hex.CenterToWorld().ToScreen()
}

func fRectFromScreen(c utils.ScreenCoord, w, h float32) *sdl.FRect {
	return &sdl.FRect{X: c.X, Y: c.Y, W: w, H: h}
}
func fRectFromScreenOffset(c utils.ScreenCoord, ox, oy, w, h float32) *sdl.FRect {
	return &sdl.FRect{X: c.X + ox, Y: c.Y + oy, W: w, H: h}
}

func isOnScreen(c utils.ScreenCoord) bool {
	return c.X >= 0 && c.X < RES_X && c.Y >= 0 && c.Y < RES_Y
}
func isOnScreenRadius(c utils.ScreenCoord, radius float32) bool {
	return c.X >= -radius && c.X < RES_X+radius && c.Y >= -radius && c.Y < RES_Y+radius
}
func isOnScreenBox(c1 utils.ScreenCoord, c2 utils.ScreenCoord) bool {
	return max(c1.X, c2.X) >= 0 && min(c1.X, c2.X) < RES_X && max(c1.Y, c2.Y) >= 0 && min(c1.Y, c2.Y) < RES_Y
}

func (r *GameRenderer) LoadTextures() {
	r.loadBeltTextures()
	r.loadOnGroundTextures()
	r.loadItemTextures()
	r.loadArrowTextures()
	r.loadStructureGroundTextures()
	r.loadStructureOnGroundTextures()
	r.loadDecalTextures()

	r.iconsItems = r.loadCachedTexture("icons/item_icons")
}

func (r *GameRenderer) loadArrowTextures() {
	r.arrowTextures[0] = r.loadCachedTexture("arrow-l-r")
	r.arrowTextures[1] = r.loadCachedTexture("arrow-tl-br")
}

func (r *GameRenderer) loadOnGroundTextures() {
	r.beltOnGroundTextures[ss.BELT_ON_UNDER_IN_RIGHT] = r.loadCachedTexture("on-beltunder-in-r")
	r.beltOnGroundTextures[ss.BELT_ON_UNDER_IN_DOWNRIGHT] = r.loadCachedTexture("on-beltunder-in-br")
	r.beltOnGroundTextures[ss.BELT_ON_UNDER_IN_UPLEFT] = r.loadCachedTexture("on-beltunder-in-tl")

	r.beltOnGroundTextures[ss.BELT_ON_UNDER_OUT_RIGHT] = r.loadCachedTexture("on-beltunder-out-r")
	r.beltOnGroundTextures[ss.BELT_ON_UNDER_OUT_DOWNRIGHT] = r.loadCachedTexture("on-beltunder-out-br")
	r.beltOnGroundTextures[ss.BELT_ON_UNDER_OUT_UPLEFT] = r.loadCachedTexture("on-beltunder-out-tl")

	r.beltOnGroundTextures[ss.BELT_ON_SPLITTER_UPLEFTRIGHT_DOWNLEFTRIGHT] = r.loadCachedTexture("on-beltsplitter-tlr-blr")
	r.beltOnGroundTextures[ss.BELT_ON_SPLITTER_LEFTUPLEFT_RIGHTDOWNRIGHT] = r.loadCachedTexture("on-beltsplitter-ltl-rbr")
	r.beltOnGroundTextures[ss.BELT_ON_SPLITTER_DOWNLEFTRIGHT_UPLEFTRIGHT] = r.loadCachedTexture("on-beltsplitter-blr-tlr")
	r.beltOnGroundTextures[ss.BELT_ON_SPLITTER_RIGHTDOWNRIGHT_LEFTUPLEFT] = r.loadCachedTexture("on-beltsplitter-rbr-ltl")
}

func (r *GameRenderer) loadBeltTextures() {
	// Straights
	r.beltAnimationTextures[ss.BELT_TYPE_LEFT_RIGHT] = r.loadCachedTexture("belts/L_R")
	r.beltAnimationTextures[ss.BELT_TYPE_RIGHT_LEFT] = r.loadCachedTexture("belts/R_L")
	r.beltAnimationTextures[ss.BELT_TYPE_UPLEFT_DOWNRIGHT] = r.loadCachedTexture("belts/TL_BR")
	r.beltAnimationTextures[ss.BELT_TYPE_UPRIGHT_DOWNLEFT] = r.loadCachedTexture("belts/TR_BL")
	r.beltAnimationTextures[ss.BELT_TYPE_DOWNRIGHT_UPLEFT] = r.loadCachedTexture("belts/BR_TL")
	r.beltAnimationTextures[ss.BELT_TYPE_DOWNLEFT_UPRIGHT] = r.loadCachedTexture("belts/BL_TR")

	r.beltAnimationTextures[ss.BELT_TYPE_IN_LEFT] = r.loadCachedTexture("belts/IN_L")
	r.beltAnimationTextures[ss.BELT_TYPE_IN_RIGHT] = r.loadCachedTexture("belts/IN_R")
	r.beltAnimationTextures[ss.BELT_TYPE_IN_UPLEFT] = r.loadCachedTexture("belts/IN_TL")
	r.beltAnimationTextures[ss.BELT_TYPE_IN_UPRIGHT] = r.loadCachedTexture("belts/IN_TR")
	r.beltAnimationTextures[ss.BELT_TYPE_IN_DOWNLEFT] = r.loadCachedTexture("belts/IN_BL")
	r.beltAnimationTextures[ss.BELT_TYPE_IN_DOWNRIGHT] = r.loadCachedTexture("belts/IN_BR")

	r.beltAnimationTextures[ss.BELT_TYPE_LEFT] = r.loadCachedTexture("belts/OUT_L")
	r.beltAnimationTextures[ss.BELT_TYPE_RIGHT] = r.loadCachedTexture("belts/OUT_R")
	r.beltAnimationTextures[ss.BELT_TYPE_UPLEFT] = r.loadCachedTexture("belts/OUT_TL")
	r.beltAnimationTextures[ss.BELT_TYPE_UPRIGHT] = r.loadCachedTexture("belts/OUT_TR")
	r.beltAnimationTextures[ss.BELT_TYPE_DOWNLEFT] = r.loadCachedTexture("belts/OUT_BL")
	r.beltAnimationTextures[ss.BELT_TYPE_DOWNRIGHT] = r.loadCachedTexture("belts/OUT_BR")

	// Bends
	r.beltAnimationTextures[ss.BELT_TYPE_LEFT_UPRIGHT] = r.loadCachedTexture("belts/L_TR")
	r.beltAnimationTextures[ss.BELT_TYPE_LEFT_DOWNRIGHT] = r.loadCachedTexture("belts/L_BR")
	r.beltAnimationTextures[ss.BELT_TYPE_RIGHT_DOWNLEFT] = r.loadCachedTexture("belts/R_BL")
	r.beltAnimationTextures[ss.BELT_TYPE_RIGHT_UPLEFT] = r.loadCachedTexture("belts/R_TL")
	r.beltAnimationTextures[ss.BELT_TYPE_UPLEFT_RIGHT] = r.loadCachedTexture("belts/TL_R")
	r.beltAnimationTextures[ss.BELT_TYPE_UPLEFT_DOWNLEFT] = r.loadCachedTexture("belts/TL_BL")
	r.beltAnimationTextures[ss.BELT_TYPE_UPRIGHT_LEFT] = r.loadCachedTexture("belts/TR_L")
	r.beltAnimationTextures[ss.BELT_TYPE_UPRIGHT_DOWNRIGHT] = r.loadCachedTexture("belts/TR_BR")
	r.beltAnimationTextures[ss.BELT_TYPE_DOWNLEFT_RIGHT] = r.loadCachedTexture("belts/BL_R")
	r.beltAnimationTextures[ss.BELT_TYPE_DOWNLEFT_UPLEFT] = r.loadCachedTexture("belts/BL_TL")
	r.beltAnimationTextures[ss.BELT_TYPE_DOWNRIGHT_LEFT] = r.loadCachedTexture("belts/BR_L")
	r.beltAnimationTextures[ss.BELT_TYPE_DOWNRIGHT_UPRIGHT] = r.loadCachedTexture("belts/BR_TR")

	// Side joins
	r.beltAnimationTextures[ss.BELT_TYPE_LEFT_RIGHT_UPLEFT] = r.loadCachedTexture("belts/L_TL_R")
	r.beltAnimationTextures[ss.BELT_TYPE_LEFT_RIGHT_DOWNLEFT] = r.loadCachedTexture("belts/L_BL_R")
	r.beltAnimationTextures[ss.BELT_TYPE_RIGHT_LEFT_UPRIGHT] = r.loadCachedTexture("belts/R_TR_L")
	r.beltAnimationTextures[ss.BELT_TYPE_RIGHT_LEFT_DOWNRIGHT] = r.loadCachedTexture("belts/R_BR_L")
	r.beltAnimationTextures[ss.BELT_TYPE_UPLEFT_DOWNRIGHT_LEFT] = r.loadCachedTexture("belts/TL_L_BR")
	r.beltAnimationTextures[ss.BELT_TYPE_UPLEFT_DOWNRIGHT_UPRIGHT] = r.loadCachedTexture("belts/TL_TR_BR")
	r.beltAnimationTextures[ss.BELT_TYPE_UPRIGHT_DOWNLEFT_RIGHT] = r.loadCachedTexture("belts/TR_R_BL")
	r.beltAnimationTextures[ss.BELT_TYPE_UPRIGHT_DOWNLEFT_UPLEFT] = r.loadCachedTexture("belts/TR_TL_BL")
	r.beltAnimationTextures[ss.BELT_TYPE_DOWNRIGHT_UPLEFT_DOWNLEFT] = r.loadCachedTexture("belts/BR_BL_TL")
	r.beltAnimationTextures[ss.BELT_TYPE_DOWNRIGHT_UPLEFT_RIGHT] = r.loadCachedTexture("belts/BR_R_TL")
	r.beltAnimationTextures[ss.BELT_TYPE_DOWNLEFT_UPRIGHT_DOWNRIGHT] = r.loadCachedTexture("belts/BL_BR_TR")
	r.beltAnimationTextures[ss.BELT_TYPE_DOWNLEFT_UPRIGHT_LEFT] = r.loadCachedTexture("belts/BL_L_TR")

	// Merges
	r.beltAnimationTextures[ss.BELT_TYPE_DOWNRIGHT_UPRIGHT_LEFT] = r.loadCachedTexture("belts/BR_TR_L")
	r.beltAnimationTextures[ss.BELT_TYPE_UPLEFT_DOWNLEFT_RIGHT] = r.loadCachedTexture("belts/TL_BL_R")
	r.beltAnimationTextures[ss.BELT_TYPE_UPRIGHT_LEFT_DOWNRIGHT] = r.loadCachedTexture("belts/TR_L_BR")
	r.beltAnimationTextures[ss.BELT_TYPE_LEFT_DOWNRIGHT_UPRIGHT] = r.loadCachedTexture("belts/L_BR_TR")
	r.beltAnimationTextures[ss.BELT_TYPE_DOWNLEFT_RIGHT_UPLEFT] = r.loadCachedTexture("belts/BL_R_TL")
	r.beltAnimationTextures[ss.BELT_TYPE_RIGHT_UPLEFT_DOWNLEFT] = r.loadCachedTexture("belts/R_TL_BL")

	// 3-2-1
	r.beltAnimationTextures[ss.BELT_TYPE_ALL_LEFT] = r.loadCachedTexture("belts/3_L")
	r.beltAnimationTextures[ss.BELT_TYPE_ALL_RIGHT] = r.loadCachedTexture("belts/3_R")
	r.beltAnimationTextures[ss.BELT_TYPE_ALL_UPLEFT] = r.loadCachedTexture("belts/3_TL")
	r.beltAnimationTextures[ss.BELT_TYPE_ALL_DOWNLEFT] = r.loadCachedTexture("belts/3_BL")
	r.beltAnimationTextures[ss.BELT_TYPE_ALL_UPRIGHT] = r.loadCachedTexture("belts/3_TR")
	r.beltAnimationTextures[ss.BELT_TYPE_ALL_DOWNRIGHT] = r.loadCachedTexture("belts/3_BR")

	// Splitters
	r.beltAnimationTextures[ss.BELT_TYPE_SPLITTER_LEFTUPLEFT_RIGHTDOWNRIGHT] = r.loadCachedTexture("belts/SP_LTL_RBR")
	r.beltAnimationTextures[ss.BELT_TYPE_SPLITTER_UPLEFTRIGHT_DOWNLEFTRIGHT] = r.loadCachedTexture("belts/SP_TLR_BRL")
	r.beltAnimationTextures[ss.BELT_TYPE_SPLITTER_RIGHTUPRIGHT_LEFTDOWNLEFT] = r.loadCachedTexture("belts/SP_TRR_BLL")
	r.beltAnimationTextures[ss.BELT_TYPE_SPLITTER_RIGHTDOWNRIGHT_LEFTUPLEFT] = r.loadCachedTexture("belts/SP_RBR_LTL")
	r.beltAnimationTextures[ss.BELT_TYPE_SPLITTER_DOWNLEFTRIGHT_UPLEFTRIGHT] = r.loadCachedTexture("belts/SP_BRL_TLR")
	r.beltAnimationTextures[ss.BELT_TYPE_SPLITTER_LEFTDOWNLEFT_RIGHTUPRIGHT] = r.loadCachedTexture("belts/SP_BLL_TRR")
}

func (r *GameRenderer) loadItemTextures() {
	r.itemTextures[ss.ITEM_TYPE_IRON_ORE] = r.loadCachedTexture("items/iron_ore")
	r.itemTextures[ss.ITEM_TYPE_IRON_PLATE] = r.loadCachedTexture("items/iron_plate")
	r.itemTextures[ss.ITEM_TYPE_IRON_GEAR] = r.loadCachedTexture("items/iron_gear")
}

func (r *GameRenderer) loadStructureGroundTextures() {
	r.objectGroundTextures[ss.OBJECT_TYPE_CHESTBOX_SMALL] = r.loadCachedTexture("chests/chest_small")
	r.objectGroundTextures[ss.OBJECT_TYPE_CHESTBOX_MEDIUM] = r.loadCachedTexture("chests/chest_medium")
	r.objectGroundTextures[ss.OBJECT_TYPE_CHESTBOX_LARGE] = r.loadCachedTexture("chests/chest_large")

	r.objectGroundTextures[ss.OBJECT_TYPE_ASSEMBLER_BASIC] = r.loadCachedTexture("shape_bighex")

	r.objectGroundDirTextures[ss.OBJECT_TYPE_INSERTER1] = [utils.DIR_COUNT]*sdl.Texture{
		utils.DIR_LEFT:       r.loadCachedTexture("inserter/base_l"),
		utils.DIR_RIGHT:      r.loadCachedTexture("inserter/base_r"),
		utils.DIR_UP_LEFT:    r.loadCachedTexture("inserter/base_tl"),
		utils.DIR_DOWN_LEFT:  r.loadCachedTexture("inserter/base_bl"),
		utils.DIR_UP_RIGHT:   r.loadCachedTexture("inserter/base_tr"),
		utils.DIR_DOWN_RIGHT: r.loadCachedTexture("inserter/base_br"),
	}

	r.objectGroundDirTextures[ss.OBJECT_TYPE_FURNACE_STONE] = [utils.DIR_COUNT]*sdl.Texture{
		utils.DIR_LEFT:     r.loadCachedTexture("shape_diamond_lr"),
		utils.DIR_UP_LEFT:  r.loadCachedTexture("shape_diamond_tl_br"),
		utils.DIR_UP_RIGHT: r.loadCachedTexture("shape_diamond_tr_bl"),
	}
}

func (r *GameRenderer) loadStructureOnGroundTextures() {
	r.objectOnGroundTextures[ss.OBJECT_TYPE_ASSEMBLER_BASIC] = r.loadCachedTexture("objects/onground/assembly_machine")

	r.objectOnGroundDirTextures[ss.OBJECT_TYPE_FURNACE_STONE] = [utils.DIR_COUNT]*sdl.Texture{
		utils.DIR_LEFT:     r.loadCachedTexture("objects/onground/stone_furnace_lr"),
		utils.DIR_UP_LEFT:  r.loadCachedTexture("objects/onground/stone_furnace_tl_br"),
		utils.DIR_UP_RIGHT: r.loadCachedTexture("objects/onground/stone_furnace_tr_bl"),
	}
}

func (r *GameRenderer) loadDecalTextures() {
	r.decalTextures[DECAL_BLACK_SPOT_FUZZY] = r.loadCachedTexture("decals/black_spot_fuzzy")
}

func (r *GameRenderer) loadCachedTexture(name string) *sdl.Texture {
	surface, err := loadWithCaching(name)
	if err != nil {
		panic(err)
	}
	defer surface.Free()

	tex, err := r.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}

	return tex
}

func loadWithCaching(name string) (*sdl.Surface, error) {
	cachedPath := path.Join(ss.TEXTURE_CACHE_DIR, name) + ss.TEXTURE_CACHE_EXT
	if _, err := os.Stat(cachedPath); err == nil {
		surface, err := loadFromCache(cachedPath)
		if err != nil {
			return nil, err
		}
		return surface, nil
	}

	srcPath := path.Join(ss.TEXTURE_DIR, name) + ss.TEXTURE_SOURCE_EXT

	surface, err := img.Load(srcPath)
	if err != nil {
		return nil, err
	}

	err = writeToCache(surface, cachedPath)
	if err != nil {
		return nil, err
	}
	return surface, nil
}

func loadFromCache(cachedPath string) (*sdl.Surface, error) {
	fp, err := os.Open(cachedPath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	var w, h int32
	err = binary.Read(fp, binary.LittleEndian, &w)
	if err != nil {
		return nil, err
	}
	err = binary.Read(fp, binary.LittleEndian, &h)
	if err != nil {
		return nil, err
	}
	var format uint32
	err = binary.Read(fp, binary.LittleEndian, &format)
	if err != nil {
		return nil, err
	}

	surface, err := sdl.CreateRGBSurfaceWithFormat(0, w, h, 32, format)
	if err != nil {
		return nil, err
	}

	err = surface.Lock()
	if err != nil {
		return nil, err
	}
	defer surface.Unlock()

	_, err = fp.Read(surface.Pixels())
	if err != nil {
		return nil, err
	}

	return surface, nil
}

func writeToCache(surface *sdl.Surface, cachedPath string) error {
	err := surface.Lock()
	if err != nil {
		return err
	}
	defer surface.Unlock()

	err = os.MkdirAll(path.Dir(cachedPath), 0755)
	if err != nil {
		return err
	}

	fp, err := os.Create(cachedPath)
	if err != nil {
		return err
	}
	defer fp.Close()

	binary.Write(fp, binary.LittleEndian, surface.W)
	binary.Write(fp, binary.LittleEndian, surface.H)
	binary.Write(fp, binary.LittleEndian, surface.Format.Format)
	fp.Write(surface.Pixels())
	return nil
}

func (r *GameRenderer) DrawString(stringID strings.StringID, pos utils.ScreenCoord, align TextAlignment) {
	cs := CompoundString{}
	cs.AddString(stringID, r.stringManager)
	r.stringManager.RenderCompoundString(r.renderer, &cs, int32(pos.X), int32(pos.Y), align)
}

func (r *GameRenderer) DrawWorldCoords(coord utils.WorldCoord, precision int, pctX, pctY float32) {
	x, y := int32(math.Round(float64(pctX*RES_X))), int32(math.Round(float64(pctY*RES_Y)))

	cs := CompoundString{}
	cs.AddFloat(coord.X, precision, r.stringManager)
	cs.AddString(strings.STRING_COMMASPACE, r.stringManager)
	cs.AddFloat(coord.Y, precision, r.stringManager)
	r.stringManager.RenderCompoundString(r.renderer, &cs, x, y, TEXT_ALIGN_LEFT)
}
func (r *GameRenderer) DrawHexCoords(hex utils.HexCoord, pctX, pctY float32) {
	x, y := int32(math.Round(float64(pctX*RES_X))), int32(math.Round(float64(pctY*RES_Y)))

	cs := CompoundString{}
	cs.AddInt(int(hex.X), 1, r.stringManager)
	cs.AddString(strings.STRING_COMMASPACE, r.stringManager)
	cs.AddInt(int(hex.Y), 1, r.stringManager)

	r.stringManager.RenderCompoundString(r.renderer, &cs, x, y, TEXT_ALIGN_LEFT)
}

func (r *GameRenderer) DrawFpsTps(fps, tps float64, pctX, pctY float32) {
	x, y := int32(math.Round(float64(pctX*RES_X))), int32(math.Round(float64(pctY*RES_Y)))

	cs := CompoundString{}
	cs.AddString(strings.STRING_FPS, r.stringManager)
	cs.AddFloat(fps, 1, r.stringManager)

	cs.AddString(strings.STRING_SPACE, r.stringManager)

	cs.AddString(strings.STRING_TPS, r.stringManager)
	cs.AddFloat(tps, 1, r.stringManager)

	r.stringManager.RenderCompoundString(r.renderer, &cs, x, y, TEXT_ALIGN_LEFT)
}

func (r *GameRenderer) DrawPlayerCoords(coord utils.WorldCoord, pctX, pctY float32) {
	x, y := int32(math.Round(float64(pctX*RES_X))), int32(math.Round(float64(pctY*RES_Y)))

	cs := CompoundString{}
	cs.AddString(strings.STRING_PLAYER_COORDS, r.stringManager)
	cs.AddFloat(coord.X, 2, r.stringManager)
	cs.AddString(strings.STRING_COMMASPACE, r.stringManager)
	cs.AddFloat(coord.Y, 2, r.stringManager)

	r.stringManager.RenderCompoundString(r.renderer, &cs, x, y, TEXT_ALIGN_LEFT)
}

func (r *GameRenderer) DrawCurrentTool(toolStr strings.StringID, pctX, pctY float32) {
	x, y := int32(math.Round(float64(pctX*RES_X))), int32(math.Round(float64(pctY*RES_Y)))

	cs := CompoundString{}
	cs.AddString(strings.STRING_TOOL, r.stringManager)
	cs.AddString(toolStr, r.stringManager)

	r.stringManager.RenderCompoundString(r.renderer, &cs, x, y, TEXT_ALIGN_RIGHT)
}

func (r *GameRenderer) DrawObjectDetails(
	name strings.StringID, hex utils.HexCoord, items []utils.ItemInfo, pctX, pctY float32,
) {
	x, y := int32(math.Round(float64(pctX*RES_X))), int32(math.Round(float64(pctY*RES_Y)))

	cs := CompoundString{}
	cs.AddString(name, r.stringManager)
	r.stringManager.RenderCompoundString(r.renderer, &cs, x, y, TEXT_ALIGN_LEFT)
	y += cs.H

	cs = CompoundString{}
	cs.AddInt(int(hex.X), 1, r.stringManager)
	cs.AddString(strings.STRING_COMMASPACE, r.stringManager)
	cs.AddInt(int(hex.Y), 1, r.stringManager)
	r.stringManager.RenderCompoundString(r.renderer, &cs, x, y, TEXT_ALIGN_LEFT)
	y += cs.H / 2

	for _, itemInfo := range items {
		r.DrawItemIconScreen(utils.ScreenCoord{X: float32(x), Y: float32(y)}, float32(fontHeight*1.5), itemInfo.Type)
		r.stringManager.RenderInt(r.renderer, itemInfo.Count, 1, x, y+int32(fontHeight/2))
		x += int32(fontHeight * 2)
	}
}

func (r *GameRenderer) DrawWindow(pos utils.ScreenCoord, size utils.ScreenCoord, title strings.StringID) {
	rect := fRectFromScreen(pos, size.X, size.Y)

	r.renderer.SetDrawColor(0, 0, 0, 255)
	r.renderer.FillRectF(rect)

	cs := CompoundString{}
	cs.AddString(title, r.stringManager)
	posX := math.Round(float64(pos.Add(size.Div(2)).X))
	r.stringManager.RenderCompoundString(r.renderer, &cs, int32(posX), int32(pos.Y)+cs.H/2, TEXT_ALIGN_CENTER)
}

func (r *GameRenderer) DrawButton(pos utils.ScreenCoord, size utils.ScreenCoord, hover bool) {
	rect := fRectFromScreen(pos, size.X, size.Y)

	r.renderer.SetDrawColor(32, 32, 32, 255)
	r.renderer.FillRectF(rect)
	if hover {
		r.renderer.SetDrawColor(32, 127, 32, 255)
		r.renderer.DrawRectF(rect)
	}
}
func (r *GameRenderer) DrawButtonText(pos utils.ScreenCoord, size utils.ScreenCoord, text strings.StringID, hover bool) {
	r.DrawButton(pos, size, hover)

	cs := CompoundString{}
	cs.AddString(text, r.stringManager)
	pos = pos.Add(size.Div(2))
	posX := math.Round(float64(pos.X))
	posY := math.Round(float64(pos.Y))
	r.stringManager.RenderCompoundString(r.renderer, &cs, int32(posX), int32(posY), TEXT_ALIGN_CENTER)
}
func (r *GameRenderer) DrawButtonIcon(pos utils.ScreenCoord, size utils.ScreenCoord, item ss.ItemType, hover bool) {
	r.DrawButton(pos, size, hover)
	r.DrawItemIconScreen(pos, size.Y, item)
}
