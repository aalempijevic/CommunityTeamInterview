package api

import "github.com/gorilla/mux"

//NewRouter creates a Router decorated with a HttpLogger which handles routes
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
