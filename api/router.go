package api

import "github.com/gorilla/mux"

//NewRouter creates a Router with handlers (decorated with a HttpLogger) for routes
// configured in routes
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		handler := HttpLogger(route.HandlerFunction, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler).
			Queries(route.Queries...)
	}
	return router
}
