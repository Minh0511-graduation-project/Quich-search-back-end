package routes

import (
	"Quick-search-back-end/controllers/lazada"
	"Quick-search-back-end/controllers/searchCount"
	"Quick-search-back-end/controllers/shopee"
	"Quick-search-back-end/controllers/tiki"

	"github.com/gorilla/mux"
)

func UserRoute(router *mux.Router) {
	// Tiki routes
	router.HandleFunc("/tiki/product", tiki.GetProductsBySearchTerm()).
		Queries(
			"searchTerm", "{searchTerm}",
		).Methods("GET", "OPTIONS")

	router.HandleFunc("/tiki/product", tiki.GetAllProducts()).Methods("GET", "OPTIONS")
	router.HandleFunc("/tiki/suggestion", tiki.GetSuggestionsByKeyword()).
		Queries(
			"keyword", "{keyword}",
		).Methods("GET", "OPTIONS")

	router.HandleFunc("/tiki/suggestion", tiki.GetAllSuggestions()).Methods("GET", "OPTIONS")

	// lazada routes
	router.HandleFunc("/lazada/product", lazada.GetProductsBySearchTerm()).
		Queries(
			"searchTerm", "{searchTerm}",
		).Methods("GET", "OPTIONS")

	router.HandleFunc("/lazada/product", lazada.GetAllProducts()).Methods("GET", "OPTIONS")
	router.HandleFunc("/lazada/suggestion", lazada.GetSuggestionsByKeyword()).
		Queries(
			"keyword", "{keyword}",
		).Methods("GET", "OPTIONS")

	router.HandleFunc("/lazada/suggestion", lazada.GetAllSuggestions()).Methods("GET", "OPTIONS")

	// shopee routes
	router.HandleFunc("/shopee/product", shopee.GetProductsBySearchTerm()).
		Queries(
			"searchTerm", "{searchTerm}",
		).Methods("GET", "OPTIONS")

	router.HandleFunc("/shopee/product", shopee.GetAllProducts()).Methods("GET", "OPTIONS")
	router.HandleFunc("/shopee/suggestion", shopee.GetSuggestionsByKeyword()).
		Queries(
			"keyword", "{keyword}",
		).Methods("GET", "OPTIONS")

	router.HandleFunc("/shopee/suggestion", shopee.GetAllSuggestions()).Methods("GET", "OPTIONS")

	router.HandleFunc("/suggestionCount", searchCount.GetCountByKeyword()).
		Queries(
			"keyword", "{keyword}",
			"site", "{site}",
		).Methods("GET", "OPTIONS")
}
