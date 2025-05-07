package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Hajdudev/ecoDatabase/internal/store"
)

type DatabaseHandler struct {
	databaseStore store.DatabaseStore
	logger        *log.Logger
}

func NewDatabaseHandler(databaseStore store.DatabaseStore, logger *log.Logger) *DatabaseHandler {
	return &DatabaseHandler{
		databaseStore: databaseStore,
		logger:        logger,
	}
}

func (wh *DatabaseHandler) FindRoute(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	from := query.Get("from")
	to := query.Get("to")
	date := query.Get("date")

	if from == "" || to == "" {
		http.Error(w, "Missing required parameters 'from' and 'to'", http.StatusBadRequest)
		return
	}

	wh.databaseStore.GetRoutes()
	fmt.Fprintf(w, "Route search parameters:\n")
	fmt.Fprintf(w, "From: %s\n", from)
	fmt.Fprintf(w, "To: %s\n", to)
	fmt.Fprintf(w, "Date: %s\n", date)
}
