package sim

import (
	"math/rand/v2"
	"time"

	"github.com/karmilg/agrodrone_timiryazevka_hakaton/internal/models"
	"github.com/karmilg/agrodrone_timiryazevka_hakaton/internal/storage"
)

type Engine struct {
	Fleet  *storage.Fleet
	Fields *storage.FieldCatalog
	Tasks  *storage.TaskQueue
	Now    time.Time
	Tick   int
}

func NewEngine(fleet *storage.Fleet, fields *storage.FieldCatalog, tasks *storage.TaskQueue) *Engine {
	return &Engine{
		Fleet:  fleet,
		Fields: fields,
		Tasks:  tasks,
		Now:    time.Date(2026, 6, 23, 8, 0, 0, 0, time.UTC),
	}
}

func (e *Engine) Ticks() []string {
	e.Tick++
	e.Now = e.Now.Add(time.Hour)
	log := []string{}

	for {
		drone := e.Fleet.Free()
		if drone == nil {
			break
		}
		task := e.Tasks.Pop()
		if task == nil {
			break
		}

		drone.StatusDrone = models.StatusFly
		e.assign(drone, task)
		log = append(log, fmtAssign(drone, task))
	}

	for _, d := range e.Fleet.All() {
		switch d.StatusDrone {
		case models.StatusFly:
			e.stepFlying(d, &log)
		case models.StatusCharge:
			d.Charge(20)
			if d.StatusDrone == models.StatusAfk {
				log = append(log, d.Name+" заряжен до 100% и готов к работе")
			}
		case models.StatusAfk:
			if d.Battery < 30 && (d.X != d.BaseX || d.Y != d.BaseY) {
				d.Fly(d.BaseX, d.BaseY)
				d.StatusDrone = models.StatusCharge
				log = append(log, d.Name+" вернулся на базу заряжаться")
			} else if d.Battery < 100 && d.X == d.BaseX && d.Y == d.BaseY {
				d.StatusDrone = models.StatusCharge
			}
		}
	}

	return log
}

var assignments = map[int]*models.Task{}

func (e *Engine) assign(d *models.Drone, t *models.Task) {
	assignments[d.ID] = t
}

func (e *Engine) stepFlying(d *models.Drone, log *[]string) {
	t, ok := assignments[d.ID]
	if !ok {
		d.StatusDrone = models.StatusAfk
		return
	}

	if d.X == t.X && d.Y == t.Y {
		field, err := e.Fields.Get(t.FieldID)
		if err == nil {
			cell := field.Grid[t.Y][t.X]
			cell.NDVI = 0.3 + rand.Float64()*0.6
			cell.Soil = 10 + rand.Float64()*60
			cell.Temp = 18 + rand.Float64()*12
			cell.Scanned = true
			cell.LastScanned = e.Now
		}
		t.Done = true
		delete(assignments, d.ID)
		d.Battery -= 2
		d.StatusDrone = models.StatusAfk
		*log = append(*log, d.Name+" выполнил задачу #"+itoa(t.ID))
		return
	}

	switch {
	case d.X < t.X:
		d.X++
	case d.X > t.X:
		d.X--
	case d.Y < t.Y:
		d.Y++
	case d.Y > t.Y:
		d.Y--
	}
	d.Battery--
	if d.Battery <= 5 {
		// аварийный возврат
		delete(assignments, d.ID)
		d.StatusDrone = models.StatusAfk
		*log = append(*log, d.Name+" вернулся: низкий заряд")
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

func fmtAssign(d *models.Drone, t *models.Task) string {
	return d.Name + " получил задачу #" + itoa(t.ID) +
		" (поле " + itoa(t.FieldID) + ", клетка " +
		itoa(t.X) + "," + itoa(t.Y) + ")"
}
