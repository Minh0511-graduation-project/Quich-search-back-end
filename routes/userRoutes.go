package routes

import (
	"Quick-search-back-end/controllers/lazada"
	"Quick-search-back-end/controllers/shopee"
	"Quick-search-back-end/controllers/tiki"

	"github.com/gorilla/mux"
)

func UserRoute(router *mux.Router) {
	// Tiki routes
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

	// lazada routes
	router.HandleFunc("/lazada/product", lazada.GetProductsBySearchTerm()).
		Queries(
			"searchTerm", "{searchTerm}",
		).Methods("GET")

	router.HandleFunc("/lazada/product", lazada.GetAllProducts()).Methods("GET")
	router.HandleFunc("/lazada/suggestion", lazada.GetSuggestionsByKeyword()).
		Queries(
			"keyword", "{keyword}",
		).Methods("GET")

	router.HandleFunc("/lazada/suggestion", lazada.GetAllSuggestions()).Methods("GET")

	// shopee routes
	router.HandleFunc("/shopee/product", shopee.GetProductsBySearchTerm()).
		Queries(
			"searchTerm", "{searchTerm}",
		).Methods("GET")

	router.HandleFunc("/shopee/product", shopee.GetAllProducts()).Methods("GET")
	router.HandleFunc("/shopee/suggestion", shopee.GetSuggestionsByKeyword()).
		Queries(
			"keyword", "{keyword}",
		).Methods("GET")

	router.HandleFunc("/shopee/suggestion", shopee.GetAllSuggestions()).Methods("GET")
}
