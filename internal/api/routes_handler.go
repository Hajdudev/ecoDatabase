package api

import (
	"fmt"
	"log"
	"net/http"
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

	routesChan := make(chan []models.Trip, 1)
	fromStopChan := make(chan models.Stop, 1)
	toStopChan := make(chan models.Stop, 1)
	errorChan := make(chan error, 3)

	wg.Add(3)

	go func() {
		err := wh.databaseStore.GetRoutesById(from, to, routesChan, &wg)
		if err != nil {
			errorChan <- err
		}
		close(routesChan)
	}()
	go func() {
		err := wh.databaseStore.GetStopInfo(from, fromStopChan, &wg)
		if err != nil {
			errorChan <- err
		}
		close(fromStopChan)
	}()
	go func() {
		err := wh.databaseStore.GetStopInfo(to, toStopChan, &wg)
		if err != nil {
			errorChan <- err
		}
		close(toStopChan)
	}()
	wg.Wait()
	close(errorChan)

	for err := range errorChan {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	routes := <-routesChan
	fromStop := <-fromStopChan
	toStop := <-toStopChan

	fmt.Fprintf(w, "The routes %+v \n", routes)
	fmt.Fprintf(w, "To data %+v\n", toStop)
	fmt.Fprintf(w, "From data %+v\n", fromStop)
	fmt.Fprintf(w, "Date: %s\n", date)
}
