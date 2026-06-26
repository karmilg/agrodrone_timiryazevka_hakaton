package storage

import (
	"sync"

	"github.com/karmilg/agrodrone_timiryazevka_hakaton/internal/models"
)

type TaskQueue struct {
	mu sync.Mutex
	items []*models.Task
	nextID int
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		nextID: 1,
	}
}

func (tq *TaskQueue) Add(t *models.Task) *models.Task {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	t.ID = tq.nextID
	tq.nextID++
	tq.items = append(tq.items, t)

	return t
}


func (tq *TaskQueue) Pop() *models.Task {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	if len(tq.items) == 0 {
		return nil
	}

	f := 0
	for i := 1; i < len(tq.items); i++ {
		if tq.items[i].Prior > tq.items[f].Prior {
			f = i
		}
	}

	best := tq.items[f]
	tq.items = append(tq.items[:f], tq.items[:f+1]...)
	return best
}


func (tq *TaskQueue) All() []*models.Task {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	res := make([]*models.Task, len(tq.items))
	copy(res, tq.items)
	
	return res
}


func (tq *TaskQueue) Len() int {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	return len(tq.items)
}