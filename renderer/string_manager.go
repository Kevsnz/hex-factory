package renderer

import (
	"fmt"
	"hextopdown/settings/strings"
	"math"
	"path"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type TextAlignment int

const (
	TEXT_ALIGN_CENTER TextAlignment = iota
	TEXT_ALIGN_LEFT
	TEXT_ALIGN_RIGHT
)

type CompoundString struct {
	strings []*stringTexture
}

func (s *CompoundString) AddString(sid strings.StringID, sm *StringManager) {
	s.strings = append(s.strings, sm.stringTextures[sid])
}
func (s *CompoundString) AddInt(number int, minDigits int, sm *StringManager) {
	s.strings = append(s.strings, sm.StringListInt(number, minDigits)...)
}
func (s *CompoundString) AddFloat(number float64, precision int, sm *StringManager) {
	s.strings = append(s.strings, sm.StringListFloat(number, precision)...)
}

type stringTexture struct {
	Texture *sdl.Texture
	W, H    int32
}

type StringManager struct {
	stringTextures [strings.STRING_COUNT]*stringTexture
	digitTextures  [10]*stringTexture
	font           *ttf.Font
}

func NewStringManager() *StringManager {
	font, err := ttf.OpenFont(path.Join("resources", "Roboto-Regular.ttf"), 20)
	if err != nil {
		panic(err)
	}
	return &StringManager{
		font: font,
	}
}

func (s *StringManager) Destroy() {
	for _, t := range s.stringTextures {
		if t != nil {
			t.Texture.Destroy()
		}
	}
	s.font.Close()
}

func (s *StringManager) Prerender(r *sdl.Renderer) {
	color := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	for i := 0; i < 10; i++ {
		s.digitTextures[i] = s.createTexture(r, fmt.Sprintf("%d", i), color)
	}
	for i := 0; i < int(strings.STRING_COUNT); i++ {
		s.stringTextures[i] = s.createTexture(r, strings.Strings[i], color)
	}
}

func (s *StringManager) createTexture(r *sdl.Renderer, str string, color sdl.Color) *stringTexture {
	surface, err := s.font.RenderUTF8Blended(str, color)
	if err != nil {
		panic(err)
	}
	defer surface.Free()

	texture, err := r.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	return &stringTexture{texture, surface.W, surface.H}
}

func (s *StringManager) calcWidthHeight(strings []*stringTexture) (int32, int32) {
	width := int32(0)
	maxHeigth := int32(0)

	for _, tex := range strings {
		width += tex.W
		if tex.H > maxHeigth {
			maxHeigth = tex.H
		}
	}
	return width, maxHeigth
}

func (s *StringManager) RenderCompoundString(r *sdl.Renderer, cs *CompoundString, x, y int32, align TextAlignment) int32 {
	w, h := s.calcWidthHeight(cs.strings)

	switch align {
	case TEXT_ALIGN_RIGHT:
		x -= w
	case TEXT_ALIGN_CENTER:
		x -= w / 2
	}
	y -= h / 2

	for _, st := range cs.strings {
		if st == nil {
			panic("string not prerendered")
		}
		r.Copy(st.Texture, nil, &sdl.Rect{X: x, Y: y, W: st.W, H: st.H})
		x += st.W
	}
	return w
}

func (s *StringManager) Render(r *sdl.Renderer, id strings.StringID, x, y int32) int32 {
	st := s.stringTextures[id]
	if st == nil {
		panic("string not prerendered")
	}
	r.Copy(st.Texture, nil, &sdl.Rect{X: x, Y: y, W: st.W, H: st.H})
	return st.W
}

func (s *StringManager) StringListInt(number int, minDigits int) []*stringTexture {
	l := 1
	n2 := number
	if n2 < 0 {
		l++
		n2 = -n2
	}
	n2 /= 10

	for n2 > 0 {
		n2 /= 10
		l++
	}
	if l < minDigits {
		l = minDigits
	}

	strs := make([]*stringTexture, l)
	i := l - 1

	if number < 0 {
		strs[0] = s.stringTextures[strings.STRING_DASH]
		number = -number
	} else {
		strs[0] = s.digitTextures[0]
	}

	for ; number > 0; i-- {
		strs[i] = s.digitTextures[number%10]
		number /= 10
	}
	for ; i >= 1; i-- {
		strs[i] = s.digitTextures[0]
	}

	return strs
}

func (s *StringManager) RenderInt(r *sdl.Renderer, number int, minDigits int, x, y int32) int32 {
	curX := x
	digits := s.StringListInt(number, minDigits)

	for _, s := range digits {
		r.Copy(s.Texture, nil, &sdl.Rect{X: curX, Y: y, W: s.W, H: s.H})
		curX += s.W
	}

	return curX - x
}

func (s *StringManager) StringListFloat(number float64, precision int) []*stringTexture {
	precisionMult := 1
	for i := 0; i < precision; i++ {
		precisionMult *= 10
	}
	numint := int(math.Round(number * float64(precisionMult)))

	l := 2
	n2 := numint
	if n2 < 0 {
		l++
		n2 = -n2
	}
	n2 /= 10
	for n2 > 0 {
		n2 /= 10
		l++
	}

	strs := make([]*stringTexture, l)
	i := l - 1
	if numint < 0 {
		strs[0] = s.stringTextures[strings.STRING_DASH]
		numint = -numint
	} else {
		strs[0] = s.digitTextures[0]
	}

	for numint > 0 {
		strs[i] = s.digitTextures[numint%10]
		numint /= 10
		i--
		if l-i-1 == precision {
			strs[i] = s.stringTextures[strings.STRING_PERIOD]
			i--
		}
	}
	for ; i >= 1; i-- {
		strs[i] = s.digitTextures[0]
	}

	return strs
}

func (s *StringManager) RenderFloat(r *sdl.Renderer, number float64, precision int, x, y int32) int32 {
	strs := s.StringListFloat(number, precision)

	curX := x
	for _, str := range strs {
		r.Copy(str.Texture, nil, &sdl.Rect{X: curX, Y: y, W: str.W, H: str.H})
		curX += str.W
	}

	return curX - x
}
