package models

import "time"

type TaskType string

const (
	TaskScan  TaskType = "Сканировать"
	TaskSpray TaskType = "Опрыскивать"
)

type Task struct {
	ID       int
	X        int
	Y        int
	Type     TaskType
	FieldID  int
	Prior    int
	Deadline time.Time
	Done     bool
}

func (t *Task) IsOver(now time.Time) bool {
	return !t.Done && now.After(t.Deadline)
}