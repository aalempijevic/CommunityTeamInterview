package api

import "github.com/gorilla/mux"

//NewRouter creates a Router decorated with a Logger which handles routes
// configured in the routes file
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		handler := Logger(route.HandlerFunction, route.Name)

		router.
		Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler).
			Queries(route.Queries...)
	}
	return router
}
