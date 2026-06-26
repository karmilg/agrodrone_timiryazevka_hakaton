package storage

import (
	"errors"
	"sync"

	"github.com/karmilg/agrodrone_timiryazevka_hakaton/internal/models"
)

type Fleet struct {
	mu sync.RWMutex
	drones map[int]*models.Drone
	nextID int
}

func NewFleet() *Fleet {
	return &Fleet{
		drones: make(map[int]*models.Drone),
		nextID: 1,
	}
}

func (f *Fleet) Add(d *models.Drone) *models.Drone {
	f.mu.Lock()
	defer f.mu.Unlock()

	d.ID = f.nextID
	f.nextID++
	f.drones[d.ID] = d

	return d
}

func (f *Fleet) Get(id int) (*models.Drone, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	d, ok := f.drones[id]
	if !ok {
		return nil, errors.New("Дрон не найден")
	}

	return d, nil
}

func (f *Fleet) All() []*models.Drone {
	f.mu.RLock()
	defer f.mu.RUnlock()

	res := make([]*models.Drone, 0, len(f.drones))
	for _, r := range f.drones {
		res = append(res, r)
	}

	return res
}


func (f *Fleet) Free() *models.Drone {
	f.mu.RLock()
	defer f.mu.RUnlock()

	for _, r := range f.drones {
		if r.StatusDrone == models.StatusAfk && r.Battery > 30 {
			return r
		}
	}

	return nil
}