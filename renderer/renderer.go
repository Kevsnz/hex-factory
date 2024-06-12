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
	"github.com/veandco/go-sdl2/ttf"
)

type GameRenderer struct {
	renderer              *sdl.Renderer
	stringManager         *StringManager
	font                  *ttf.Font
	Viewport              *utils.Viewport
	beltTextures          [ss.BELT_TYPE_COUNT]*sdl.Texture
	beltAnimationTextures [ss.BELT_TYPE_COUNT]*sdl.Texture
	beltOnGroundTextures  [ss.BELT_ON_COUNT]*sdl.Texture
	objectTextures        [ss.OBJECT_TYPE_COUNT]*sdl.Texture
	objectDirTextures     [ss.OBJECT_TYPE_COUNT][utils.DIR_COUNT]*sdl.Texture
	itemTextures          [ss.ITEM_TYPE_COUNT]*sdl.Texture
	arrowTextures         [2]*sdl.Texture
	timeMs                uint64
}

func NewGameRenderer(window *sdl.Window, view utils.WorldCoord) *GameRenderer {
	renderer, err := sdl.CreateRenderer(window, 0, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	font, err := ttf.OpenFont(path.Join("resources", "Roboto-Regular.ttf"), 20)
	if err != nil {
		panic(err)
	}
	sm := NewStringManager()
	sm.Prerender(renderer)

	return &GameRenderer{
		renderer:      renderer,
		stringManager: sm,
		font:          font,
		Viewport:      utils.NewViewport(view, 1.0),
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
	for _, tex := range r.objectTextures {
		tex.Destroy()
	}
	for _, texs := range r.objectDirTextures {
		for _, tex := range texs {
			tex.Destroy()
		}
	}
	r.font.Close()
	r.stringManager.Destroy()
	r.renderer.Destroy()
}

func (r *GameRenderer) MoveTheView(pos utils.WorldCoord, dt uint64) {
	r.Viewport.ShiftViewport(pos, dt)
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
	cx, cy := r.Viewport.WorldToScreen(drawPos)
	if !isOnScreenRadius(cx, cy, 10) {
		return
	}
	r.renderer.SetDrawColor(127, 127, 255, 255)
	r.renderer.FillRectF(&sdl.FRect{X: cx - 10, Y: cy - 10, W: 20, H: 20})
}

func (r *GameRenderer) drawHexGrid() {
	hex1 := utils.HexCoordFromWorld(r.Viewport.ScreenToWorld(0, 0))
	hex2 := utils.HexCoordFromWorld(r.Viewport.ScreenToWorld(RES_X, 0))
	hex3 := utils.HexCoordFromWorld(r.Viewport.ScreenToWorld(0, RES_Y))
	w := r.Viewport.GetHexWidth()
	e := r.Viewport.GetHexEdge()
	o := r.Viewport.GetHexOffset()

	r.renderer.SetDrawColor(96, 96, 96, 255)
	sx1, sy1 := r.Viewport.WorldToScreen(hex1.LeftTopToWorld())
	for hy := int32(0); hy <= hex3.Y-hex1.Y+1; hy++ {
		yo := float32(hy) * (e + o)
		xo := float32(hy&1) * w / 2
		for hx := int32(-1); hx <= hex2.X-hex1.X+1; hx++ {
			r.renderer.DrawLineF(sx1+float32(hx)*w+xo, sy1+yo, sx1+xo+float32(hx)*w+w/2, sy1+yo-o)   // top left
			r.renderer.DrawLineF(sx1+float32(hx)*w+xo+w/2, sy1+yo-o, sx1+xo+float32(hx)*w+w, sy1+yo) // top right
			r.renderer.DrawLineF(sx1+float32(hx)*w+xo, sy1+yo, sx1+float32(hx)*w+xo, sy1+yo+e)       // left
		}
	}
}

func (r *GameRenderer) DrawHexCenter(hex utils.HexCoord) {
	cx, cy := hexCenterToScreen(hex, r.Viewport)
	if !isOnScreen(cx, cy) {
		return
	}
	r.renderer.SetDrawColor(0, 96, 0, 255)
	rect := &sdl.FRect{X: float32(cx - 0), Y: float32(cy - 0), W: 1, H: 1}
	r.renderer.FillRectF(rect)
}

func (r *GameRenderer) DrawBeltConnectionIncoming(hex utils.HexCoord, dir utils.Dir, left bool, start, end float64) {
	cx, cy := hexCenterToScreen(hex, r.Viewport)

	outOffset := radiusOffsets[dir]
	laneOffset := lanesOffsetsLeft[dir]

	if left {
		laneOffset[0] = -laneOffset[0]
		laneOffset[1] = -laneOffset[1]
	}

	if left {
		r.renderer.SetDrawColor(192, 0, 0, 255)
	} else {
		r.renderer.SetDrawColor(0, 0, 192, 255)
	}

	sx := cx + outOffset[0]*(1-2*float32(start)) + laneOffset[0]
	sy := cy + outOffset[1]*(1-2*float32(start)) + laneOffset[1]

	ex := cx + outOffset[0]*(1-2*float32(end)) + laneOffset[0]
	ey := cy + outOffset[1]*(1-2*float32(end)) + laneOffset[1]

	if isOnScreenRadius(sx, sy, 2) {
		r.renderer.FillRectF(&sdl.FRect{X: sx - 2, Y: sy - 2, W: 4, H: 4})
	}
	if isOnScreenBox(sx, sy, ex, ey) {
		r.renderer.DrawLineF(sx, sy, ex, ey)
	}
}

func (r *GameRenderer) DrawBeltConnectionsOutgoing(hex utils.HexCoord, dir utils.Dir) {
	cx, cy := hexCenterToScreen(hex, r.Viewport)

	outOffset := radiusOffsets[dir]
	laneOffset := lanesOffsetsLeft[dir]

	// Left lane
	r.renderer.SetDrawColor(192, 0, 0, 255)
	sx := cx + laneOffset[0]
	sy := cy + laneOffset[1]

	ex := sx + outOffset[0]
	ey := sy + outOffset[1]

	if isOnScreenRadius(sx, sy, 2) {
		r.renderer.FillRectF(&sdl.FRect{X: sx - 2, Y: sy - 2, W: 4, H: 4})
	}
	if isOnScreenBox(sx, sy, ex, ey) {
		r.renderer.DrawLineF(sx, sy, ex, ey)
	}

	// Right lane
	laneOffset[0] = -laneOffset[0]
	laneOffset[1] = -laneOffset[1]
	r.renderer.SetDrawColor(0, 0, 192, 255)

	sx = cx + laneOffset[0]
	sy = cy + laneOffset[1]

	ex = sx + outOffset[0]
	ey = sy + outOffset[1]

	if isOnScreenRadius(sx, sy, 2) {
		r.renderer.FillRectF(&sdl.FRect{X: sx - 2, Y: sy - 2, W: 4, H: 4})
	}
	if isOnScreenBox(sx, sy, ex, ey) {
		r.renderer.DrawLineF(sx, sy, ex, ey)
	}
}

func (r *GameRenderer) DrawAnimatedBelt(hex utils.HexCoord, beltType ss.BeltType, speed float64) {
	tex := r.beltAnimationTextures[beltType]
	if tex == nil {
		panic("no animation texture for belt type")
	}

	cx, cy := hexCenterToScreen(hex, r.Viewport)
	if !isOnScreenRadius(cx, cy, r.Viewport.GetZoomedDimension(ss.BELT_DRAW_R)) {
		return
	}
	_, tsegm := math.Modf(speed * float64(r.timeMs) / 1000)
	frame := int32(math.Floor(tsegm*ss.ANIM_BELT_STEPS)) % ss.ANIM_BELT_FRAMES

	e := r.Viewport.GetHexEdge()
	r.renderer.CopyF(
		tex,
		&sdl.Rect{X: 0, Y: frame * TEXTURE_SIZE_HEX, W: TEXTURE_SIZE_HEX, H: TEXTURE_SIZE_HEX},
		&sdl.FRect{X: cx - e, Y: cy - e, W: 2 * e, H: 2 * e},
	)
}

func (r *GameRenderer) DrawBeltOnGround(hex utils.HexCoord, beltType ss.BeltType) {
	if beltType == ss.BELT_ON_COUNT {
		return
	}

	cx, cy := hexCenterToScreen(hex, r.Viewport)
	if !isOnScreenRadius(cx, cy, r.Viewport.GetZoomedDimension(ss.BELT_DRAW_R)) {
		return
	}
	typeFlip := beltOnFlipMapping[beltType]

	tex := r.beltOnGroundTextures[typeFlip.type1]
	if tex == nil {
		panic(fmt.Sprintf("no texture for belt type %d", typeFlip.type1))
	}
	e := r.Viewport.GetHexEdge()
	r.renderer.CopyExF(tex, nil, &sdl.FRect{
		X: cx - e,
		Y: cy - e,
		W: 2 * e,
		H: 2 * e,
	}, 0, nil, typeFlip.flip)
}

func (r *GameRenderer) DrawObjectGround(pos utils.WorldCoord, objectType ss.ObjectType, shape utils.Shape, dir utils.Dir) {
	x, y := r.Viewport.WorldToScreen(pos)
	z := r.Viewport.Zoom
	sp := GetShapeParam(shape, dir)
	x -= float32(sp.OffsetX * z)
	y -= float32(sp.OffsetY * z)
	if !isOnScreenBox(x, y, x+float32(sp.Width*z), y+float32(sp.Height*z)) {
		return
	}

	tex := r.objectTextures[objectType]
	if tex == nil {
		tex = r.objectDirTextures[objectType][dir]
		if tex == nil {
			tex = r.objectDirTextures[objectType][dir.Reverse()] // mirrored shape
			if tex == nil {
				panic(fmt.Sprintf("no texture for object type %d, dir %d", objectType, dir))
			}
		}
	}

	r.renderer.CopyF(tex, nil, &sdl.FRect{X: x, Y: y, W: float32(sp.Width * z), H: float32(sp.Height * z)})
}

func (r *GameRenderer) DrawItem(pos utils.WorldCoordInterpolated, itemType ss.ItemType) {
	drawPos := pos.GetInterpolatedPos(r.timeMs, ss.TICK_DT)
	sx, sy := r.Viewport.WorldToScreen(drawPos)
	idr := r.Viewport.GetZoomedDimension(ss.ITEM_DRAW_R)
	if !isOnScreenRadius(sx, sy, idr) {
		return
	}

	tex := r.itemTextures[itemType]
	if tex == nil {
		r.renderer.SetDrawColor(255, 0, 255, 255)
		r.renderer.FillRectF(&sdl.FRect{X: sx - idr, Y: sy - idr, W: 2 * idr, H: 2 * idr})
		return
	}
	r.renderer.CopyF(tex, nil, &sdl.FRect{X: sx - idr, Y: sy - idr, W: 2 * idr, H: 2 * idr})
}

func (r *GameRenderer) DrawArrow(pctX, pctY float32, dir utils.Dir) {
	x := pctX * RES_X
	y := pctY * RES_Y
	idxFlip := arrowDirMapping[dir]

	r.renderer.CopyExF(r.arrowTextures[idxFlip.idx], nil, &sdl.FRect{
		X: x,
		Y: y,
		W: 64,
		H: 64,
	}, 0, nil, idxFlip.flip)
}

func (r *GameRenderer) Finish() {
	r.renderer.Present()
}

func (r *GameRenderer) IsHexOnScreen(coord utils.HexCoord) bool {
	sx, sy := hexCenterToScreen(coord, r.Viewport)
	return isOnScreenRadius(sx, sy, r.Viewport.GetZoomedDimension(ss.HEX_DRAW_R))
}

func (r *GameRenderer) DrawConnectionHexes(hex1, hex2 utils.HexCoord) {
	if hex1 == hex2 {
		return
	}
	r.renderer.SetDrawColor(255, 255, 0, 255)

	x1, y1 := hexCenterToScreen(hex1, r.Viewport)
	x2, y2 := hexCenterToScreen(hex2, r.Viewport)

	r.drawDashedLine(x1, y1, x2, y2)
}

func (r *GameRenderer) drawDashedLine(x1, y1, x2, y2 float32) {
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	dSq := dx*dx + dy*dy
	if dSq < 4*ss.DASH_LEN*ss.DASH_LEN {
		r.renderer.DrawLineF(x1, y1, x2, y2)
		return
	}
	d := math.Sqrt(dSq)

	swap := false
	if math.Abs(dy) > math.Abs(dx) {
		swap = true
		dx, dy = dy, dx
		x1, x2, y1, y2 = y1, y2, x1, x2
	}

	dashes := int(math.Round(d / (2 * ss.DASH_LEN)))
	dxStep := dx / float64(dashes*2+1)
	dyStep := dy / float64(dashes*2+1)

	for i := 0; i < dashes; i++ {
		xd1 := x1 + float32(dxStep*float64(i*2))
		yd1 := y1 + float32(dyStep*float64(i*2))
		xd2 := x1 + float32(dxStep*float64(i*2+1))
		yd2 := y1 + float32(dyStep*float64(i*2+1))
		if swap {
			r.renderer.DrawLineF(yd1, xd1, yd2, xd2)
		} else {
			r.renderer.DrawLineF(xd1, yd1, xd2, yd2)
		}
	}

	xd1 := x1 + float32(dxStep*float64(dashes*2))
	yd1 := y1 + float32(dyStep*float64(dashes*2))
	if swap {
		r.renderer.DrawLineF(yd1, xd1, y2, x2)
	} else {
		r.renderer.DrawLineF(xd1, yd1, x2, y2)
	}
}

func (r *GameRenderer) DrawWorldLine(p1, p2 utils.WorldCoord) {
	x1, y1 := r.Viewport.WorldToScreen(p1)
	x2, y2 := r.Viewport.WorldToScreen(p2)

	r.renderer.SetDrawColor(255, 0, 0, 255)
	r.renderer.DrawLineF(x1, y1, x2, y2)
}

func hexCenterToScreen(hex utils.HexCoord, viewport *utils.Viewport) (float32, float32) {
	return viewport.WorldToScreen(hex.CenterToWorld())
}

func isOnScreen(x, y float32) bool {
	return x >= 0 && x < RES_X && y >= 0 && y < RES_Y
}

func isOnScreenRadius(x, y, radius float32) bool {
	return x >= -radius && x < RES_X+radius && y >= -radius && y < RES_Y+radius
}

func isOnScreenBox(x1, y1, x2, y2 float32) bool {
	return max(x1, x2) >= 0 && min(x1, x2) < RES_X && max(y1, y2) >= 0 && min(y1, y2) < RES_Y
}

func (r *GameRenderer) LoadTextures() {
	r.LoadBeltTextures()
	r.LoadOnGroundTextures()
	r.LoadItemTextures()
	r.LoadArrowTextures()
	r.LoadStructureGroundTextures()
}

func (r *GameRenderer) LoadArrowTextures() {
	r.arrowTextures[0] = r.loadCachedTexture("arrow-l-r")
	r.arrowTextures[1] = r.loadCachedTexture("arrow-tl-br")
}

func (r *GameRenderer) LoadOnGroundTextures() {
	r.LoadBeltOnGroundTexture(ss.BELT_ON_UNDER_IN_RIGHT, "on-beltunder-in-r")
	r.LoadBeltOnGroundTexture(ss.BELT_ON_UNDER_IN_DOWNRIGHT, "on-beltunder-in-br")
	r.LoadBeltOnGroundTexture(ss.BELT_ON_UNDER_IN_UPLEFT, "on-beltunder-in-tl")

	r.LoadBeltOnGroundTexture(ss.BELT_ON_UNDER_OUT_RIGHT, "on-beltunder-out-r")
	r.LoadBeltOnGroundTexture(ss.BELT_ON_UNDER_OUT_DOWNRIGHT, "on-beltunder-out-br")
	r.LoadBeltOnGroundTexture(ss.BELT_ON_UNDER_OUT_UPLEFT, "on-beltunder-out-tl")

	r.LoadBeltOnGroundTexture(ss.BELT_ON_SPLITTER_UPLEFTRIGHT_DOWNLEFTRIGHT, "on-beltsplitter-tlr-blr")
	r.LoadBeltOnGroundTexture(ss.BELT_ON_SPLITTER_LEFTUPLEFT_RIGHTDOWNRIGHT, "on-beltsplitter-ltl-rbr")
	r.LoadBeltOnGroundTexture(ss.BELT_ON_SPLITTER_DOWNLEFTRIGHT_UPLEFTRIGHT, "on-beltsplitter-blr-tlr")
	r.LoadBeltOnGroundTexture(ss.BELT_ON_SPLITTER_RIGHTDOWNRIGHT_LEFTUPLEFT, "on-beltsplitter-rbr-ltl")
}

func (r *GameRenderer) LoadBeltTextures() {
	// Straights
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_LEFT_RIGHT, "belts/L_R")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_RIGHT_LEFT, "belts/R_L")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_UPLEFT_DOWNRIGHT, "belts/TL_BR")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_UPRIGHT_DOWNLEFT, "belts/TR_BL")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_DOWNRIGHT_UPLEFT, "belts/BR_TL")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_DOWNLEFT_UPRIGHT, "belts/BL_TR")

	r.LoadBeltAnimationTexture(ss.BELT_TYPE_IN_LEFT, "belts/IN_L")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_IN_RIGHT, "belts/IN_R")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_IN_UPLEFT, "belts/IN_TL")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_IN_UPRIGHT, "belts/IN_TR")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_IN_DOWNLEFT, "belts/IN_BL")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_IN_DOWNRIGHT, "belts/IN_BR")

	r.LoadBeltAnimationTexture(ss.BELT_TYPE_LEFT, "belts/OUT_L")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_RIGHT, "belts/OUT_R")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_UPLEFT, "belts/OUT_TL")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_UPRIGHT, "belts/OUT_TR")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_DOWNLEFT, "belts/OUT_BL")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_DOWNRIGHT, "belts/OUT_BR")

	// Bends
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_LEFT_UPRIGHT, "belts/L_TR")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_LEFT_DOWNRIGHT, "belts/L_BR")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_RIGHT_DOWNLEFT, "belts/R_BL")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_RIGHT_UPLEFT, "belts/R_TL")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_UPLEFT_RIGHT, "belts/TL_R")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_UPLEFT_DOWNLEFT, "belts/TL_BL")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_UPRIGHT_LEFT, "belts/TR_L")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_UPRIGHT_DOWNRIGHT, "belts/TR_BR")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_DOWNLEFT_RIGHT, "belts/BL_R")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_DOWNLEFT_UPLEFT, "belts/BL_TL")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_DOWNRIGHT_LEFT, "belts/BR_L")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_DOWNRIGHT_UPRIGHT, "belts/BR_TR")

	// Side joins
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_LEFT_RIGHT_UPLEFT, "belts/L_TL_R")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_LEFT_RIGHT_DOWNLEFT, "belts/L_BL_R")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_RIGHT_LEFT_UPRIGHT, "belts/R_TR_L")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_RIGHT_LEFT_DOWNRIGHT, "belts/R_BR_L")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_UPLEFT_DOWNRIGHT_LEFT, "belts/TL_L_BR")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_UPLEFT_DOWNRIGHT_UPRIGHT, "belts/TL_TR_BR")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_UPRIGHT_DOWNLEFT_RIGHT, "belts/TR_R_BL")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_UPRIGHT_DOWNLEFT_UPLEFT, "belts/TR_TL_BL")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_DOWNRIGHT_UPLEFT_DOWNLEFT, "belts/BR_BL_TL")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_DOWNRIGHT_UPLEFT_RIGHT, "belts/BR_R_TL")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_DOWNLEFT_UPRIGHT_DOWNRIGHT, "belts/BL_BR_TR")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_DOWNLEFT_UPRIGHT_LEFT, "belts/BL_L_TR")

	// Merges
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_DOWNRIGHT_UPRIGHT_LEFT, "belts/BR_TR_L")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_UPLEFT_DOWNLEFT_RIGHT, "belts/TL_BL_R")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_UPRIGHT_LEFT_DOWNRIGHT, "belts/TR_L_BR")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_LEFT_DOWNRIGHT_UPRIGHT, "belts/L_BR_TR")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_DOWNLEFT_RIGHT_UPLEFT, "belts/BL_R_TL")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_RIGHT_UPLEFT_DOWNLEFT, "belts/R_TL_BL")

	// 3-2-1
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_ALL_LEFT, "belts/3_L")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_ALL_RIGHT, "belts/3_R")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_ALL_UPLEFT, "belts/3_TL")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_ALL_DOWNLEFT, "belts/3_BL")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_ALL_UPRIGHT, "belts/3_TR")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_ALL_DOWNRIGHT, "belts/3_BR")

	// Splitters
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_SPLITTER_LEFTUPLEFT_RIGHTDOWNRIGHT, "belts/SP_LTL_RBR")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_SPLITTER_UPLEFTRIGHT_DOWNLEFTRIGHT, "belts/SP_TLR_BRL")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_SPLITTER_RIGHTUPRIGHT_LEFTDOWNLEFT, "belts/SP_TRR_BLL")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_SPLITTER_RIGHTDOWNRIGHT_LEFTUPLEFT, "belts/SP_RBR_LTL")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_SPLITTER_DOWNLEFTRIGHT_UPLEFTRIGHT, "belts/SP_BRL_TLR")
	r.LoadBeltAnimationTexture(ss.BELT_TYPE_SPLITTER_LEFTDOWNLEFT_RIGHTUPRIGHT, "belts/SP_BLL_TRR")
}

func (r *GameRenderer) LoadItemTextures() {
	r.itemTextures[ss.ITEM_TYPE_IRON_PLATE] = r.loadCachedTexture("items/iron_plate")
}

func (r *GameRenderer) LoadStructureGroundTextures() {
	r.objectTextures[ss.OBJECT_TYPE_CHESTBOX_SMALL] = r.loadCachedTexture("chests/chest_small")
	r.objectTextures[ss.OBJECT_TYPE_CHESTBOX_MEDIUM] = r.loadCachedTexture("chests/chest_medium")
	r.objectTextures[ss.OBJECT_TYPE_CHESTBOX_LARGE] = r.loadCachedTexture("chests/chest_large")

	r.objectDirTextures[ss.OBJECT_TYPE_INSERTER1] = [utils.DIR_COUNT]*sdl.Texture{
		utils.DIR_LEFT:       r.loadCachedTexture("inserter/base_l"),
		utils.DIR_RIGHT:      r.loadCachedTexture("inserter/base_r"),
		utils.DIR_UP_LEFT:    r.loadCachedTexture("inserter/base_tl"),
		utils.DIR_DOWN_LEFT:  r.loadCachedTexture("inserter/base_bl"),
		utils.DIR_UP_RIGHT:   r.loadCachedTexture("inserter/base_tr"),
		utils.DIR_DOWN_RIGHT: r.loadCachedTexture("inserter/base_br"),
	}

	r.objectDirTextures[ss.OBJECT_TYPE_FURNACE_STONE] = [utils.DIR_COUNT]*sdl.Texture{
		utils.DIR_LEFT:     r.loadCachedTexture("shape_diamond_lr"),
		utils.DIR_UP_LEFT:  r.loadCachedTexture("shape_diamond_tl_br"),
		utils.DIR_UP_RIGHT: r.loadCachedTexture("shape_diamond_tr_bl"),
	}
}

func (r *GameRenderer) LoadBeltTexture(beltType ss.BeltType, filename string) {
	r.beltTextures[beltType] = r.loadCachedTexture(filename)
}

func (r *GameRenderer) LoadBeltAnimationTexture(beltType ss.BeltType, filename string) {
	r.beltAnimationTextures[beltType] = r.loadCachedTexture(filename)
}

func (r *GameRenderer) LoadBeltOnGroundTexture(beltType ss.BeltType, filename string) {
	r.beltOnGroundTextures[beltType] = r.loadCachedTexture(filename)
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

// func (r *GameRenderer) DrawText(text string, x, y int32, color sdl.Color) {
// 	surface, err := r.font.RenderUTF8Blended(text, color)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer surface.Free()
// 	texture, err := r.renderer.CreateTextureFromSurface(surface)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer texture.Destroy()
// 	r.renderer.Copy(texture, nil, &sdl.Rect{X: x, Y: y, W: surface.W, H: surface.H})
// }

func (r *GameRenderer) DrawString(stringID strings.StringID, pctX, pctY float32) {
	x, y := pctX*RES_X, pctY*RES_Y
	r.stringManager.Render(r.renderer, stringID, x, y)
}

func (r *GameRenderer) DrawWorldCoords(coord utils.WorldCoord, precision int, pctX, pctY float32) {
	x, y := pctX*RES_X, pctY*RES_Y
	x = r.stringManager.RenderFloat(r.renderer, coord.X, precision, x, y)
	x = r.stringManager.Render(r.renderer, strings.STRING_COMMASPACE, x, y)
	r.stringManager.RenderFloat(r.renderer, coord.Y, precision, x, y)
}
func (r *GameRenderer) DrawHexCoords(hex utils.HexCoord, pctX, pctY float32) {
	x, y := pctX*RES_X, pctY*RES_Y
	x = r.stringManager.RenderInt(r.renderer, int(hex.X), 1, x, y)
	x = r.stringManager.Render(r.renderer, strings.STRING_COMMASPACE, x, y)
	r.stringManager.RenderInt(r.renderer, int(hex.Y), 1, x, y)
}

func (r *GameRenderer) DrawFpsTps(fps, tps float64, pctX, pctY float32) {
	x, y := pctX*RES_X, pctY*RES_Y
	x = r.stringManager.Render(r.renderer, strings.STRING_FPS, x, y)
	x = r.stringManager.RenderFloat(r.renderer, fps, 1, x, y)

	x = r.stringManager.Render(r.renderer, strings.STRING_SPACE, x, y)

	x = r.stringManager.Render(r.renderer, strings.STRING_TPS, x, y)
	r.stringManager.RenderFloat(r.renderer, tps, 1, x, y)
}

func (r *GameRenderer) DrawPlayerCoords(coord utils.WorldCoord, pctX, pctY float32) {
	x, y := pctX*RES_X, pctY*RES_Y
	x = r.stringManager.Render(r.renderer, strings.STRING_PLAYER_COORDS, x, y)

	x = r.stringManager.RenderFloat(r.renderer, coord.X, 2, x, y)
	x = r.stringManager.Render(r.renderer, strings.STRING_COMMASPACE, x, y)
	r.stringManager.RenderFloat(r.renderer, coord.Y, 2, x, y)
}

func (r *GameRenderer) DrawObjectDetails(
	name strings.StringID, hex utils.HexCoord, items []utils.ItemInfo, pctX, pctY float32,
) {
	x, y := pctX*RES_X, pctY*RES_Y

	r.stringManager.Render(r.renderer, name, x, y)
	y += 22

	x2 := r.stringManager.RenderInt(r.renderer, int(hex.X), 1, x, y)
	x2 = r.stringManager.Render(r.renderer, strings.STRING_COMMASPACE, x2, y)
	r.stringManager.RenderInt(r.renderer, int(hex.Y), 1, x2, y)
	y += 22

	for _, itemInfo := range items {
		tex := r.itemTextures[itemInfo.Type]
		if tex == nil {
			r.renderer.SetDrawColor(255, 0, 255, 255)
			r.renderer.FillRectF(&sdl.FRect{X: x, Y: y, W: 25, H: 25})
			return
		}
		r.renderer.CopyF(tex, nil, &sdl.FRect{X: x, Y: y, W: 25, H: 25})

		r.stringManager.RenderInt(r.renderer, itemInfo.Count, 1, x+15, y+10)
		x += 30
	}
}
