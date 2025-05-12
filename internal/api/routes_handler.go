package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/Hajdudev/ecoDatabase/internal/store"
	"github.com/Hajdudev/ecoDatabase/models"
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

	var wg sync.WaitGroup

	routesChan := make(chan []string)
	fromStopChan := make(chan models.Stop)
	toStopChan := make(chan models.Stop)

	func() {
		wg.Add(3)

		go wh.databaseStore.GetRoutesById(from, to, routesChan, &wg)
		go wh.databaseStore.GetStopInfo(from, fromStopChan, &wg)
		go wh.databaseStore.GetStopInfo(to, toStopChan, &wg)
		wg.Wait()
		close(routesChan)
		close(fromStopChan)
		close(toStopChan)
	}()

	routes := <-routesChan
	fromStop := <-fromStopChan
	toStop := <-toStopChan

	fmt.Fprintf(w, "The routes %s \n", strings.Join(routes, ", "))
	fmt.Fprintf(w, "To data %+v\n", toStop)
	fmt.Fprintf(w, "From data %+v\n", fromStop)
	fmt.Fprintf(w, "From: %s\n", from)
	fmt.Fprintf(w, "To: %s\n", to)
	fmt.Fprintf(w, "Date: %s\n", date)
}
