package renderer

type UIElement uint8

const (
	UI_ELEMENT_WINDOW UIElement = iota
	UI_ELEMENT_BUTTON
	UI_ELEMENT_ITEM_SLOT

	UI_ELEMENT_COUNT
)

var uiColors = [UI_ELEMENT_COUNT][4]uint8{
	UI_ELEMENT_WINDOW:    {38, 28, 22, 255},
	UI_ELEMENT_BUTTON:    {51, 37, 29, 255},
	UI_ELEMENT_ITEM_SLOT: {172, 176, 191, 255},
}

var uiColorsBorder = [UI_ELEMENT_COUNT][4]uint8{
	UI_ELEMENT_WINDOW:    {76, 61, 53, 255},
	UI_ELEMENT_BUTTON:    {76, 61, 53, 255},
	UI_ELEMENT_ITEM_SLOT: {135, 110, 94, 255},
}

var uiColorsHlight = [UI_ELEMENT_COUNT][4]uint8{
	UI_ELEMENT_BUTTON:    {63, 46, 36, 255},
	UI_ELEMENT_ITEM_SLOT: {232, 237, 255, 255},
}

var buttonDownColor = [4]uint8{30, 22, 17, 255}
var windowHeaderColor = [4]uint8{63, 46, 36, 255}
