package settings

type Shape int

const (
	SHAPE_SINGLE  Shape = iota // single hex
	SHAPE_DIAMOND Shape = iota // diamond (4 hexes)
	SHAPE_BIGHEX  Shape = iota // big hex

	SHAPE_COUNT Shape = iota
)
