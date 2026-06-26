package models

type Field struct {
	ID     int
	Name   string
	Width  int
	Height int
	Grid   [][]*Cell
}

func NewField(id int, name string, width, height int) *Field {
	grid := make([][]*Cell, height)
	for i := 0; i < height; i++ {
		grid[i] = make([]*Cell, width)
		for j := 0; j < width; j++ {
			grid[i][j] = &Cell{
				X: j,
				Y: i,
			}
		}
	}
	return &Field{
		ID: id,
		Name: name,
		Width: width,
		Height: height,
		Grid: grid,
	}
}



func (f *Field) AVGNdvi() float64 {
	sum := 0.0
	c := 0
	for _, row := range f.Grid {
		for _, r := range row {
			if r.Scanned {
				sum += r.NDVI
				c++
			}
		}
	}
	if c == 0 {
		return 0
	}

	return sum / float64(c)
}


func (f *Field) ProblemIndex() float64 {
	return 1 - f.AVGNdvi()
}