package settings

type DrawingShape int

const (
	DRAWING_SHAPE_SINGLE     DrawingShape = iota // single hex
	DRAWING_SHAPE_DIAMOND_LR DrawingShape = iota // diamond left-right (4 hexes)
	DRAWING_SHAPE_DIAMOND_UL DrawingShape = iota // diamond upleft (4 hexes)
	DRAWING_SHAPE_DIAMOND_UR DrawingShape = iota // diamond upright (4 hexes)

	DRAWING_SHAPE_COUNT
)
