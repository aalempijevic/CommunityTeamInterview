package api

import (
	"net/http"
)

//Route represents an api route that can be handled by our Router.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	Queries		[]string
	HandlerFunction http.HandlerFunc
}

//Routes is a slice of Route
type Routes []Route

//Config for all of our api routes
var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		nil,
		Index,
	},
	Route{
		"WordsByTag",
		"GET",
		"/wordsByTag",
		[]string{"tag", "{tag:[0-9]+}"},
		WordsByTag,
	},
	Route{
		"CommentsByWord",
		"GET",
		"/commentsByWord",
		[]string{"word", "{word}", "skip", "{skip:[0-9]+}", "limit", "{limit:[0-9]+}"},
		CommentsByWord,
	},
}