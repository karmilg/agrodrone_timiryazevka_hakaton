package algo

import "github.com/karmilg/agrodrone_timiryazevka_hakaton/internal/models"

func FindProblemCell(f *models.Field) []*models.Cell {
	var res []*models.Cell
	for i := 0; i < f.Height; i++ {
		for j := 0; j < f.Width; j++ {
			s := f.Grid[i][j]
			if s.Scanned && s.IsProblem() {
				res = append(res, s)
			}
		}
	}

	return res
}