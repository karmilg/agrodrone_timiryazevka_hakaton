package storage

import (
	"errors"
	"sort"
	"sync"

	"github.com/karmilg/agrodrone_timiryazevka_hakaton/internal/models"
)

type FieldCatalog struct {
	mu sync.RWMutex
	fields map[int]*models.Field
	nextID int
}

func NewFieldCatalog() *FieldCatalog {
	return &FieldCatalog{
		fields: make(map[int]*models.Field),
		nextID: 1,
	}
}

func (fc *FieldCatalog) Add(f *models.Field) *models.Field {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	f.ID = fc.nextID
	fc.nextID++
	fc.fields[f.ID] = f

	return f
}

func (fc *FieldCatalog) Get(id int) (*models.Field, error) {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	f, ok := fc.fields[id]
	if !ok {
		return nil, errors.New("Поле не найдено")
	}

	return f, nil
}

func (fc *FieldCatalog) All() []*models.Field {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	res := make([]*models.Field, 0, len(fc.fields))
	for _, r := range fc.fields {
		res = append(res, r)
	}

	return res
}

func (fc *FieldCatalog) SortedIdx() []int {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	idx := make([]int, 0, len(fc.fields))
	for id := range fc.fields {
		idx = append(idx, id)
	}

	sort.Ints(idx)
	return idx
}
