package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
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

func normalizeTime(t string) string {
	var hour, min, sec int
	_, err := fmt.Sscanf(t, "%d:%d:%d", &hour, &min, &sec)
	if err != nil {
		return t
	}
	hour = hour % 24
	return fmt.Sprintf("%02d:%02d:%02d", hour, min, sec)
}

func (wh *DatabaseHandler) StopNames(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	stops, err := wh.databaseStore.GetStopsNames()
	if err != nil {
		http.Error(w, "There was an error getting the names", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stops); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (wh *DatabaseHandler) FindRoute(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for all requests
	w.Header().Set("Access-Control-Allow-Origin", "*") // Use specific origin if credentials/cookies are needed
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Credentials", "true") // Only if not using '*'

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
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
	routesChan := make(chan map[string]models.TripHash, 1)
	tempStopChan := make(chan []models.TempStop, 1)
	fromStopChan := make(chan models.Stop, 1)
	dateChan := make(chan string, 1)
	toStopChan := make(chan models.Stop, 1)
	errorChan := make(chan error, 1)

	var service_date string

	handleError := func(err error, msg string) {
		if err != nil {
			select {
			case errorChan <- fmt.Errorf("%s: %w", msg, err):
			default:
			}
			cancel()
		}
	}

	wg.Add(3)
	go func() {
		defer wg.Done()
		err := wh.databaseStore.GetStopsID(from, fromIdChan)
		handleError(err, "Failed to get stops for 'from'")
		close(fromIdChan)
	}()
	go func() {
		defer wg.Done()
		err := wh.databaseStore.GetCalendarType(date, dateChan)
		handleError(err, "Failed to get calendarDate")
		service_date = <-dateChan
		close(dateChan)
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
		err := wh.databaseStore.GetStopTimesInfo(fromIDs, toIDs, service_date, tempStopChan)
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
	var routes map[string]models.TripHash
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
			TripId:        temp.TripID,
			TripName:      routes[temp.TripID].Headsign,
			FromStopId:    temp.FromStopID,
			FromStopName:  fromStop.StopName,
			ToStopId:      temp.ToStopID,
			ToStopName:    toStop.StopName,
			DepartureTime: normalizeTime(temp.FromDepartureTime),
			ArrivalTime:   normalizeTime(temp.ToDepartureTime),
			// ServiceId:          routes[temp.TripID].ServiceID,
			DepartureDayOffset: 0,
			ArrivalDayOffset:   0,
			SearchDate:         date,
		}
		finalRoutes = append(finalRoutes, route)
	}
	sort.Slice(finalRoutes, func(i, j int) bool {
		return finalRoutes[i].DepartureTime < finalRoutes[j].DepartureTime
	})

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(finalRoutes); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}
