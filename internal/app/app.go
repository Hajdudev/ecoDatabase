package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Hajdudev/ecoDatabase/internal/api"
	"github.com/Hajdudev/ecoDatabase/internal/store"
)

type Application struct {
	Logger          *log.Logger
	DatabaseHandler *api.DatabaseHandler
	Database        *sql.DB
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := store.Open()
	if err != nil {
		return nil, err
	}
	postgreStore := store.NewPostgresStore(db)

	defer db.Close()
	dbHandler := api.NewDatabaseHandler(*postgreStore, logger)
	app := &Application{
		Logger:          logger,
		DatabaseHandler: dbHandler,
		Database:        db,
	}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "THe app is healthy \n")
}
