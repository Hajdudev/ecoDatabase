package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type DatabaseHandler struct{}

func NewDatabaseHandler() *DatabaseHandler {
	return &DatabaseHandler{}
}

func (wh *DatabaseHandler) FindRoutes(w http.ResponseWriter, r *http.Request) {
	paramsRoutes := chi.URLParam(r, "names")
	if paramsRoutes == "" {
		http.NotFound(w, r)
		return
	}

	routesID, err := strconv.ParseInt(paramsRoutes, 10, 64)
	if err != nil {
		fmt.Fprintln(w, "Error 2")
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Found the routes with a id %d |n", routesID)
}
