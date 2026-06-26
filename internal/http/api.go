package http

import (
	"encoding/json"
	"net/http"

	"github.com/karmilg/agrodrone_timiryazevka_hakaton/internal/sim"
	"github.com/karmilg/agrodrone_timiryazevka_hakaton/internal/storage"
)

type API struct {
	Fleet *storage.Fleet
	Fields *storage.FieldCatalog
	Tasks *storage.TaskQueue
	Engine *sim.Engine
}

func NewAPI(fleet *storage.Fleet, fields *storage.FieldCatalog, tasks *storage.TaskQueue, engine *sim.Engine) *API {
	return &API{
		Fleet: fleet,
		Fields: fields,
		Tasks: tasks,
		Engine: engine,
	}
}

func (a *API) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/state", a.getState)

	mux.HandleFunc("POST /api/fields", a.createField)
	mux.HandleFunc("GET /api/fields/{id}", a.getField)
	mux.HandleFunc("GET /api/fields/{id}/problems", a.getProblems)

	mux.HandleFunc("POST /api/drones", a.createDrone)

	mux.HandleFunc("POST /api/tasks", a.createTask)

	mux.HandleFunc("POST /api/tick", a.tick)

	mux.HandleFunc("GET /api/search/field", a.searchField)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-type", "application-json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func readJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}