package routes

import (
	"Quick-search-back-end/controllers/tiki"

	"github.com/gorilla/mux"
)

func UserRoute(router *mux.Router) {
	router.HandleFunc("/tiki/product", tiki.GetProductsBySearchTerm()).
		Queries(
			"searchTerm", "{searchTerm}",
		).Methods("GET")

	router.HandleFunc("/tiki/product", tiki.GetAllProducts()).Methods("GET")
	router.HandleFunc("/tiki/suggestion", tiki.GetSuggestionsByKeyword()).
		Queries(
			"keyword", "{keyword}",
		).Methods("GET")

	router.HandleFunc("/tiki/suggestion", tiki.GetAllSuggestions()).Methods("GET")
}
