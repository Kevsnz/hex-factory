package utils

type Dir int8

func (d Dir) Reverse() Dir {
	return (d + DIR_COUNT/2) % DIR_COUNT
}

func (d Dir) IsAcute(d2 Dir) bool {
	dd := (d2 + DIR_COUNT - d) % DIR_COUNT
	return dd < 2 || dd > DIR_COUNT-2
}

func (d Dir) NextCW() Dir {
	return (d + 1) % DIR_COUNT
}

func (d Dir) NextCCW() Dir {
	return (d + DIR_COUNT - 1) % DIR_COUNT
}

func (d Dir) Next(cw bool) Dir {
	if cw {
		return d.NextCW()
	}
	return d.NextCCW()
}

const (
	DIR_LEFT       Dir = iota // 0
	DIR_UP_LEFT    Dir = iota // 1
	DIR_UP_RIGHT   Dir = iota // 2
	DIR_RIGHT      Dir = iota // 3
	DIR_DOWN_RIGHT Dir = iota // 4
	DIR_DOWN_LEFT  Dir = iota // 5
	DIR_COUNT      Dir = iota // 6 total
)

var adjacency_dirs = [3][3]Dir{ // [dy][dx] !!!!!!!!!!!!!!!!!!!
	{DIR_COUNT, DIR_UP_LEFT, DIR_UP_RIGHT},
	{DIR_LEFT, DIR_COUNT, DIR_RIGHT},
	{DIR_DOWN_LEFT, DIR_DOWN_RIGHT, DIR_COUNT},
}

var AllDirs = [DIR_COUNT]Dir{DIR_LEFT, DIR_UP_LEFT, DIR_UP_RIGHT, DIR_RIGHT, DIR_DOWN_RIGHT, DIR_DOWN_LEFT}
