package algo

import "github.com/karmilg/agrodrone_timiryazevka_hakaton/internal/models"

func SortByProb(f []*models.Field) {
	quicksort(f, 0, len(f)-1)
}

func quicksort(f []*models.Field, low, high int) {
	if low >= high {
		return
	}

	p := partition(f, low, high)
	quicksort(f, low, p-1)
	quicksort(f, p+1, high)
}

func partition(f []*models.Field, low, high int) int {
	pivot := f[high].ProblemIndex()
	i := low - 1
	for j := low; j < high; j++ {
		if f[j].ProblemIndex() > pivot {
			i++
			f[i], f[j] = f[j], f[i]
		}
	}
	f[i+1], f[high] = f[high], f[i+1]
	return i+1
} 