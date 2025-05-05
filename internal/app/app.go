package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Hajdudev/ecoDatabase/internal/api"
)

type Application struct {
	Logger          *log.Logger
	DatabaseHandler *api.DatabaseHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	dbHandler := api.NewDatabaseHandler()
	app := &Application{
		Logger:          logger,
		DatabaseHandler: dbHandler,
	}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "THe app is healthy \n")
}
