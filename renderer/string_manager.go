package renderer

import (
	"fmt"
	"hextopdown/settings/strings"
	"math"
	"path"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type stringTexture struct {
	Texture *sdl.Texture
	W, H    float32
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
		s.digitTextures[i] = s.createtexture(r, fmt.Sprintf("%d", i), color)
	}
	for i := 0; i < int(strings.STRING_COUNT); i++ {
		s.stringTextures[i] = s.createtexture(r, strings.Strings[i], color)
	}
}

func (s StringManager) createtexture(r *sdl.Renderer, str string, color sdl.Color) *stringTexture {
	surface, err := s.font.RenderUTF8Blended(str, color)
	if err != nil {
		panic(err)
	}
	defer surface.Free()

	texture, err := r.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	return &stringTexture{texture, float32(surface.W), float32(surface.H)}
}

func (s *StringManager) Render(r *sdl.Renderer, id strings.StringID, x, y float32) float32 {
	st := s.stringTextures[id]
	if st == nil {
		panic("string not prerendered")
	}
	r.CopyF(st.Texture, nil, &sdl.FRect{X: x, Y: y, W: st.W, H: st.H})
	return x + st.W
}

func (s *StringManager) RenderInt(r *sdl.Renderer, number int, minDigits int, x, y float32) float32 {
	curX := x
	digits := make([]int, 0)

	if number < 0 {
		curX = s.Render(r, strings.STRING_DASH, curX, y)
		number = -number
	}

	for number > 0 {
		digits = append(digits, number%10)
		number /= 10
	}
	for len(digits) < minDigits {
		digits = append(digits, 0)
	}

	for i := len(digits) - 1; i >= 0; i-- {
		r.CopyF(s.digitTextures[digits[i]].Texture, nil, &sdl.FRect{X: curX, Y: y, W: s.digitTextures[digits[i]].W, H: s.digitTextures[digits[i]].H})
		curX += s.digitTextures[digits[i]].W
	}

	return curX
}

func (s *StringManager) RenderFloat(r *sdl.Renderer, number float64, precision int, x, y float32) float32 {
	precisionMult := 1
	for i := 0; i < precision; i++ {
		precisionMult *= 10
	}
	numint := int(math.Round(number * float64(precisionMult)))

	curX := x
	digits := make([]int, 0)

	if numint < 0 {
		curX = s.Render(r, strings.STRING_DASH, curX, y)
		numint = -numint
	}

	for numint > 0 {
		digits = append(digits, numint%10)
		numint /= 10
	}
	for len(digits) < precision+1 {
		digits = append(digits, 0)
	}

	for i := len(digits) - 1; i >= 0; i-- {
		r.CopyF(s.digitTextures[digits[i]].Texture, nil, &sdl.FRect{X: curX, Y: y, W: s.digitTextures[digits[i]].W, H: s.digitTextures[digits[i]].H})
		curX += s.digitTextures[digits[i]].W
		if i == precision {
			curX = s.Render(r, strings.STRING_PERIOD, curX, y)
		}
	}

	return curX
}
