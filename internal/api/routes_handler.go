package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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

	routes, err := wh.databaseStore.GetRoutesById(from, to)
	if err != nil {
		http.Error(w, "Error fetching routes", http.StatusInternalServerError)
		return
	}
	fromStop, err := wh.databaseStore.GetStopInfo(from)
	if err != nil {
		http.Error(w, "Invalid 'from' id ", http.StatusBadRequest)
		return
	}

	toStop, err := wh.databaseStore.GetStopInfo(to)
	if err != nil {
		http.Error(w, "Invalid 'to' id", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "The routes %s \n", strings.Join(routes, ", "))
	fmt.Fprintf(w, "To data %+v\n", toStop)
	fmt.Fprintf(w, "From data %+v\n", fromStop)
	fmt.Fprintf(w, "From: %s\n", from)
	fmt.Fprintf(w, "To: %s\n", to)
	fmt.Fprintf(w, "Date: %s\n", date)
}
