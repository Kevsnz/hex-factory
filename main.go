package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"hextopdown/game"
	"hextopdown/input"
	"hextopdown/renderer"
	"hextopdown/utils"
)

const MAX_FPS = 100
const MIN_FT = 1000 / MAX_FPS

func main() {
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	err = ttf.Init()
	if err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("Hex Top Down World!", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		renderer.RES_X, renderer.RES_Y, 0)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	gameState := game.NewGame()
	defer gameState.Destroy()

	r := renderer.NewGameRenderer(window)
	defer r.Destroy()

	r.LoadTextures()

	ih := input.NewInputHandler()
	ih.SetKeybind(sdl.SCANCODE_ESCAPE, input.ACTION_QUIT)
	ih.SetKeybind(sdl.SCANCODE_P, input.ACTION_PAUSE)
	ih.SetKeybind(sdl.SCANCODE_F1, input.ACTION_TOGGLE_UI)
	ih.SetKeybind(sdl.SCANCODE_W, input.ACTION_MOVE_UP)
	ih.SetKeybind(sdl.SCANCODE_S, input.ACTION_MOVE_DOWN)
	ih.SetKeybind(sdl.SCANCODE_A, input.ACTION_MOVE_LEFT)
	ih.SetKeybind(sdl.SCANCODE_D, input.ACTION_MOVE_RIGHT)
	ih.SetKeybind(sdl.SCANCODE_KP_PLUS, input.ACTION_ZOOM_IN)
	ih.SetKeybind(sdl.SCANCODE_KP_MINUS, input.ACTION_ZOOM_OUT)
	ih.SetKeybind(sdl.SCANCODE_E, input.ACTION_ROTATE_CW)
	ih.SetKeybind(sdl.SCANCODE_Q, input.ACTION_ROTATE_CCW)
	ih.SetKeybind(sdl.SCANCODE_SPACE, input.ACTION_PLACE_ITEM)
	ih.SetKeybind(sdl.SCANCODE_SLASH, input.ACTION_PLOP_SPLITTER)
	ih.SetKeybind(sdl.SCANCODE_U, input.ACTION_PLOP_UNDERGROUND)
	ih.SetKeybind(sdl.SCANCODE_I, input.ACTION_PLOP_INSERTER)
	ih.SetKeybind(sdl.SCANCODE_C, input.ACTION_PLOP_CHESTBOX_SMALL)
	ih.SetKeybind(sdl.SCANCODE_V, input.ACTION_PLOP_CHESTBOX_MEDIUM)
	ih.SetKeybind(sdl.SCANCODE_B, input.ACTION_PLOP_CHESTBOX_LARGE)
	ih.SetKeybind(sdl.SCANCODE_F, input.ACTION_PLOP_FURNACE)
	ih.SetKeybind(sdl.SCANCODE_M, input.ACTION_PLOP_ASSEMBLER)

	frameTime := 0.0
	tickTime := 0.0
	currentTicks := sdl.GetTicks64()
	lastTicks := currentTicks

	gameState.SetTime(currentTicks)
gameloop:
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch ev := event.(type) {
			case *sdl.QuitEvent:
				break gameloop
			case *sdl.KeyboardEvent:
				ih.HandleKeyboardEvent(ev)
			case *sdl.MouseButtonEvent:
				ih.HandleMouseButtonEvent(ev)
			case *sdl.MouseMotionEvent:
				ih.HandleMouseMotionEvent(ev)
			}
		}

		if ih.GetActionState(input.ACTION_ZOOM_IN) {
			utils.ZoomViewIn(currentTicks - lastTicks)
		}
		if ih.GetActionState(input.ACTION_ZOOM_OUT) {
			utils.ZoomViewOut(currentTicks - lastTicks)
		}
		gameState.ProcessInputFramed(ih, r)
		if !gameState.Running {
			break gameloop
		}

		gameState.Update(currentTicks, ih)

		utils.ShiftView(gameState.GetPlayerPos(), currentTicks-lastTicks)

		r.StartNewFrame(currentTicks)
		gameState.Draw(r)
		r.DrawFpsTps(1000.0/frameTime, 1000.0/tickTime, 0.01, 0.01)
		r.Finish()

		nextTicks := sdl.GetTicks64()
		if nextTicks-currentTicks < MIN_FT {
			sdl.Delay(MIN_FT - uint32(nextTicks-currentTicks))
			nextTicks += MIN_FT - (nextTicks - currentTicks)
		}

		frameTime = frameTime*0.9 + 0.1*float64(nextTicks-currentTicks)
		tickTime = tickTime*0.9 + 0.1*float64(gameState.TickTime)

		lastTicks = currentTicks
		currentTicks = nextTicks
	}
}
