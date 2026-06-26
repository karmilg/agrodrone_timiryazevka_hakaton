package http

import (
	"net/http"
	"strconv"

	"github.com/karmilg/agrodrone_timiryazevka_hakaton/internal/algo"
	"github.com/karmilg/agrodrone_timiryazevka_hakaton/internal/models"
)

func (a *API) getState(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"now":    a.Engine.Now,
		"tick":   a.Engine.Tick,
		"drones": a.Fleet.All(),
		"fields": fieldsSummary(a.Fields.All()),
		"tasks":  a.Tasks.All(),
	})
}

func fieldsSummary(fs []*models.Field) []map[string]any {
	out := make([]map[string]any, 0, len(fs))
	for _, f := range fs {
		out = append(out, map[string]any{
			"id":            f.ID,
			"name":          f.Name,
			"width":         f.Width,
			"height":        f.Height,
			"avg_ndvi":      f.AVGNdvi(),
			"problem_index": f.ProblemIndex(),
		})
	}
	return out
}

func (a *API) createField(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name   string `json:"name"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	}
	if err := readJSON(r, &body); err != nil {
		writeErr(w, http.StatusBadRequest, "bad json")
		return
	}
	if body.Width <= 0 || body.Height <= 0 || body.Name == "" {
		writeErr(w, http.StatusBadRequest, "name/width/height required")
		return
	}
	f := models.NewField(0, body.Name, body.Width, body.Height)
	a.Fields.Add(f) 

	writeJSON(w, http.StatusCreated, f)
}

func (a *API) getField(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	f, err := a.Fields.Get(id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "field not found")
		return
	}
	writeJSON(w, http.StatusOK, f)
}

func (a *API) getProblems(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	f, err := a.Fields.Get(id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "field not found")
		return
	}

	cells := algo.FindProblemCell(f)

	coords := make([][2]int, 0, len(cells))
	for _, c := range cells {
		coords = append(coords, [2]int{c.X, c.Y})
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"count": len(coords),
		"cells": coords,
	})
}

func (a *API) createDrone(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
		Kind string `json:"kind"`
	}
	if err := readJSON(r, &body); err != nil {
		writeErr(w, http.StatusBadRequest, "bad json")
		return
	}

	base := models.Drone{
		Name:        body.Name,
		Battery:     100,
		StatusDrone: models.StatusAfk,
		X:           0, Y: 0, BaseX: 0, BaseY: 0,
	}

	var d *models.Drone
	switch body.Kind {
	case "spray":
		sp := &models.SprayDrone{Drone: base, Cap: 10}
		d = &sp.Drone 
	default: 
		sc := &models.ScoutDrone{Drone: base, ScanRadius: 1}
		d = &sc.Drone
	}

	a.Fleet.Add(d)
	writeJSON(w, http.StatusCreated, d)
}


func (a *API) createTask(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Type     string `json:"type"`
		FieldID  int    `json:"field_id"`
		X        int    `json:"x"`
		Y        int    `json:"y"`
		Priority int    `json:"priority"`
	}
	if err := readJSON(r, &body); err != nil {
		writeErr(w, http.StatusBadRequest, "bad json")
		return
	}

	t := &models.Task{
		Type:    models.TaskType(body.Type),
		FieldID: body.FieldID,
		X:       body.X,
		Y:       body.Y,
		Prior:   body.Priority,
	}
	a.Tasks.Add(t)
	writeJSON(w, http.StatusCreated, t)
}


func (a *API) tick(w http.ResponseWriter, r *http.Request) {
	logs := a.Engine.Ticks()
	writeJSON(w, http.StatusOK, map[string]any{
		"tick": a.Engine.Tick,
		"now":  a.Engine.Now,
		"log":  logs,
	})
}


func (a *API) searchField(w http.ResponseWriter, r *http.Request) {
	target, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	ids := a.Fields.SortedIdx()
	idx, ok := algo.BinarySearch(ids, target)
	writeJSON(w, http.StatusOK, map[string]any{
		"found":    ok,
		"index":    idx,
		"compared": len(ids), 
	})
}
