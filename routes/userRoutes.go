package routes

import (
	"Quick-search-back-end/controllers"

	"github.com/gorilla/mux"
)

func UserRoute(router *mux.Router) {
	router.HandleFunc("/tiki/product", controllers.GetAProduct()).
		Queries(
			"searchTerm", "{searchTerm}",
		).Methods("GET")
	router.HandleFunc("/tiki/product", controllers.GetAllProducts()).Methods("GET")
}
