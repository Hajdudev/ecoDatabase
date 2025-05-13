package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	toIdChan := make(chan []string, 1)
	fromIdChan := make(chan []string, 1)
	routesChan := make(chan map[string]string, 1)
	tempStopChan := make(chan []models.TempStop, 1)
	fromStopChan := make(chan models.Stop, 1)
	toStopChan := make(chan models.Stop, 1)
	errorChan := make(chan error, 1)

	handleError := func(err error, msg string) {
		if err != nil {
			select {
			case errorChan <- fmt.Errorf("%s: %w", msg, err):
			default:
			}
			cancel()
		}
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		err := wh.databaseStore.GetStopsID(from, fromIdChan)
		handleError(err, "Failed to get stops for 'from'")
		close(fromIdChan)
	}()
	go func() {
		defer wg.Done()
		err := wh.databaseStore.GetStopsID(to, toIdChan)
		handleError(err, "Failed to get stops for 'to'")
		close(toIdChan)
	}()

	wg.Wait()

	fromIDs := <-fromIdChan
	toIDs := <-toIdChan

	if fromIDs == nil || toIDs == nil {
		http.Error(w, "No stops found for given 'from' or 'to' locations", http.StatusNotFound)
		return
	}

	wg.Add(4)
	go func() {
		defer wg.Done()
		err := wh.databaseStore.GetRoutesById(fromIDs, toIDs, routesChan)
		handleError(err, "Failed to get routes by ID")
		close(routesChan)
	}()
	go func() {
		defer wg.Done()
		err := wh.databaseStore.GetStopTimesInfo(fromIDs[1], toIDs[1], tempStopChan)
		handleError(err, "Failed to get stop times info")
		close(tempStopChan)
	}()
	go func() {
		defer wg.Done()
		err := wh.databaseStore.GetStopInfo(fromIDs[1], fromStopChan)
		handleError(err, "Failed to get info for 'from' stop")
		close(fromStopChan)
	}()
	go func() {
		defer wg.Done()
		err := wh.databaseStore.GetStopInfo(toIDs[1], toStopChan)
		handleError(err, "Failed to get info for 'to' stop")
		close(toStopChan)
	}()

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	for err := range errorChan {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var tempStops []models.TempStop
	var routes map[string]string
	var fromStop models.Stop
	var toStop models.Stop

	select {
	case tempStops = <-tempStopChan:
	case <-ctx.Done():
		http.Error(w, "Timeout while fetching temporary stops", http.StatusGatewayTimeout)
		return
	}

	select {
	case routes = <-routesChan:
	case <-ctx.Done():
		http.Error(w, "Timeout while fetching routes", http.StatusGatewayTimeout)
		return
	}

	select {
	case fromStop = <-fromStopChan:
	case <-ctx.Done():
		http.Error(w, "Timeout while fetching 'from' stop info", http.StatusGatewayTimeout)
		return
	}

	select {
	case toStop = <-toStopChan:
	case <-ctx.Done():
		http.Error(w, "Timeout while fetching 'to' stop info", http.StatusGatewayTimeout)
		return
	}

	var finalRoutes []models.RouteResult
	for _, temp := range tempStops {
		route := models.RouteResult{
			TripId:             temp.TripID,
			TripName:           routes[temp.TripID],
			FromStopId:         temp.FromStopID,
			FromStopName:       fromStop.StopName,
			ToStopId:           temp.ToStopID,
			ToStopName:         toStop.StopName,
			DepartureTime:      temp.FromDepartureTime,
			ArrivalTime:        temp.ToDepartureTime,
			ServiceId:          "",
			DepartureDayOffset: 0,
			ArrivalDayOffset:   0,
			SearchDate:         date,
		}
		finalRoutes = append(finalRoutes, route)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(finalRoutes); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}
