package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Hajdudev/ecoDatabase/internal/api"
	"github.com/Hajdudev/ecoDatabase/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Logger          *log.Logger
	DatabaseHandler *api.DatabaseHandler
	Database        *pgxpool.Pool
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := store.Open()
	if err != nil {
		return nil, err
	}

	databaseStore := store.NewPostgresStore(db)
	dbHandler := api.NewDatabaseHandler(databaseStore, logger)

	app := &Application{
		Logger:          logger,
		DatabaseHandler: dbHandler,
		Database:        db,
	}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "The app is healthy\n")
}
