package middlewares

import (
	"net/http"
)

func HandleCors(rw http.ResponseWriter) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
