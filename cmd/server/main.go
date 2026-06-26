package main

import (
	"log"
	"net/http"

	httpapi "github.com/karmilg/agrodrone_timiryazevka_hakaton/internal/http"
	"github.com/karmilg/agrodrone_timiryazevka_hakaton/internal/models"
	"github.com/karmilg/agrodrone_timiryazevka_hakaton/internal/sim"
	"github.com/karmilg/agrodrone_timiryazevka_hakaton/internal/storage"
)

func main() {
	fleet := storage.NewFleet()
	fieds := storage.NewFieldCatalog()
	tasks := storage.NewTaskQueue()
	engine := sim.NewEngine(fleet, fieds, tasks)

	seed(fleet, fieds, tasks)

	api := httpapi.NewAPI(fleet, fieds, tasks, engine)
	mux := http.NewServeMux()
	api.Register(mux)
	mux.Handle("/", http.FileServer(http.Dir("./web")))
	log.Printf("Сервер запущен на http://localhost%s", ":8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func seed(fleet *storage.Fleet, fields *storage.FieldCatalog, tasks *storage.TaskQueue) {
	fields.Add(models.NewField(0, "Северное", 10, 10))
	fields.Add(models.NewField(0, "Южное", 8, 8))

	scout := &models.ScoutDrone{
		Drone: models.Drone{
			Name:        "Пчела-1",
			Battery:     100,
			StatusDrone: models.StatusAfk,
		},
		ScanRadius: 1,
	}

	fleet.Add(&scout.Drone)
	tasks.Add(&models.Task{Type: models.TaskScan, FieldID: 1, X: 3, Y: 4, Prior: 5})
	tasks.Add(&models.Task{Type: models.TaskScan, FieldID: 1, X: 7, Y: 2, Prior: 8})
}
