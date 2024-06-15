package input

import (
	"hextopdown/utils"

	"github.com/veandco/go-sdl2/sdl"
)

type ActionType uint8

const (
	ACTION_PAUSE = ActionType(iota)
	ACTION_QUIT  = ActionType(iota)
	ACTION_FF    = ActionType(iota)

	ACTION_MOVE_LEFT  = ActionType(iota)
	ACTION_MOVE_RIGHT = ActionType(iota)
	ACTION_MOVE_UP    = ActionType(iota)
	ACTION_MOVE_DOWN  = ActionType(iota)

	ACTION_ZOOM_IN  = ActionType(iota)
	ACTION_ZOOM_OUT = ActionType(iota)

	ACTION_ROTATE_CW        = ActionType(iota)
	ACTION_ROTATE_CCW       = ActionType(iota)
	ACTION_PLACE_ITEM       = ActionType(iota)
	ACTION_PLOP_SPLITTER    = ActionType(iota)
	ACTION_PLOP_UNDERGROUND = ActionType(iota)

	ACTION_PLOP_INSERTER        = ActionType(iota)
	ACTION_PLOP_CHESTBOX_SMALL  = ActionType(iota)
	ACTION_PLOP_CHESTBOX_MEDIUM = ActionType(iota)
	ACTION_PLOP_CHESTBOX_LARGE  = ActionType(iota)
	ACTION_PLOP_FURNACE         = ActionType(iota)
	ACTION_PLOP_ASSEMBLER       = ActionType(iota)

	ACTION_COUNT = ActionType(iota)
)

var framedActions = map[ActionType]struct{}{
	ACTION_PAUSE: {},
	ACTION_QUIT:  {},
	ACTION_FF:    {},
}

type ActionEventType uint32

const (
	ACTION_TYPE_DOWN = ActionEventType(sdl.KEYDOWN)
	ACTION_TYPE_UP   = ActionEventType(sdl.KEYUP)
)

type ActionEvent struct {
	Action     ActionType
	Type       ActionEventType
	MouseCoord utils.WorldCoord // in world space
}

type MouseButton uint8

const (
	MOUSE_BUTTON_LEFT   = MouseButton(iota)
	MOUSE_BUTTON_RIGHT  = MouseButton(iota)
	MOUSE_BUTTON_MIDDLE = MouseButton(iota)

	MOUSE_BUTTON_X1 = MouseButton(iota)
	MOUSE_BUTTON_X2 = MouseButton(iota)

	MOUSE_BUTTON_COUNT = MouseButton(iota)
)

var mouseButtonMapping = map[uint8]MouseButton{
	sdl.BUTTON_LEFT:   MOUSE_BUTTON_LEFT,
	sdl.BUTTON_RIGHT:  MOUSE_BUTTON_RIGHT,
	sdl.BUTTON_MIDDLE: MOUSE_BUTTON_MIDDLE,
	sdl.BUTTON_X1:     MOUSE_BUTTON_X1,
	sdl.BUTTON_X2:     MOUSE_BUTTON_X2,
}

type MouseButtonEventType uint32

const (
	MOUSE_BUTTON_DOWN = MouseButtonEventType(sdl.MOUSEBUTTONDOWN)
	MOUSE_BUTTON_UP   = MouseButtonEventType(sdl.MOUSEBUTTONUP)
)

type MouseButtonEvent struct {
	Coord  utils.WorldCoord
	Type   MouseButtonEventType
	Button MouseButton
}

type InputHandler struct {
	KeyMappingFramed      map[sdl.Scancode]ActionType
	KeyMappingTicked      map[sdl.Scancode]ActionType
	ActionState           [ACTION_COUNT]bool
	MouseButtonState      [MOUSE_BUTTON_COUNT]bool
	MousePos              utils.WorldCoord
	KeyboardActionsFramed utils.RingBuffer[ActionEvent]
	KeyboardActionsTicked utils.RingBuffer[ActionEvent]
	MouseActions          utils.RingBuffer[MouseButtonEvent]
}

func NewInputHandler() *InputHandler {
	return &InputHandler{
		KeyMappingFramed:      map[sdl.Scancode]ActionType{},
		KeyMappingTicked:      map[sdl.Scancode]ActionType{},
		ActionState:           [ACTION_COUNT]bool{},
		MouseButtonState:      [MOUSE_BUTTON_COUNT]bool{},
		KeyboardActionsFramed: utils.NewRingBuffer[ActionEvent](10),
		KeyboardActionsTicked: utils.NewRingBuffer[ActionEvent](10),
		MouseActions:          utils.NewRingBuffer[MouseButtonEvent](10),
	}
}

func (ih *InputHandler) HandleKeyboardEvent(event *sdl.KeyboardEvent) {
	if event.Repeat != 0 {
		return
	}

	action, ok := ih.KeyMappingFramed[event.Keysym.Scancode]
	if ok {
		ih.ActionState[action] = event.Type == sdl.KEYDOWN
		ih.KeyboardActionsFramed.Push(ActionEvent{action, ActionEventType(event.Type), ih.MousePos})
	}

	action, ok = ih.KeyMappingTicked[event.Keysym.Scancode]
	if ok {
		ih.ActionState[action] = event.Type == sdl.KEYDOWN
		ih.KeyboardActionsTicked.Push(ActionEvent{action, ActionEventType(event.Type), ih.MousePos})
	}
}

func (ih *InputHandler) HandleMouseButtonEvent(event *sdl.MouseButtonEvent) {
	mouseButton, ok := mouseButtonMapping[event.Button]
	if !ok {
		return
	}

	ih.MouseButtonState[mouseButton] = event.Type == sdl.MOUSEBUTTONDOWN

	coord := utils.ScreenToWorld(float32(event.X), float32(event.Y))
	mouseEvent := MouseButtonEvent{coord, MouseButtonEventType(event.Type), mouseButton}
	_ = ih.MouseActions.Push(mouseEvent)
}

func (ih *InputHandler) HandleMouseMotionEvent(event *sdl.MouseMotionEvent) {
	ih.MousePos = utils.ScreenToWorld(float32(event.X), float32(event.Y))
}

func (ih *InputHandler) SetKeybind(key sdl.Scancode, action ActionType) {
	if _, ok := framedActions[action]; ok {
		ih.KeyMappingFramed[key] = action
	} else {
		ih.KeyMappingTicked[key] = action
	}
}

func (ih *InputHandler) RemoveKeybind(key sdl.Scancode) {
	delete(ih.KeyMappingFramed, key)
	delete(ih.KeyMappingTicked, key)
}

func (ih *InputHandler) GetActionState(action ActionType) bool {
	return ih.ActionState[action]
}

func (ih *InputHandler) GetMouseButtonState(button MouseButton) bool {
	return ih.MouseButtonState[button]
}
