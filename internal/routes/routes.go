package routes

import (
	"github.com/Hajdudev/ecoDatabase/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)
	r.Get("/find/route", app.DatabaseHandler.FindRoute)
	r.Get("/names", app.DatabaseHandler.StopNames)
	return r
}
