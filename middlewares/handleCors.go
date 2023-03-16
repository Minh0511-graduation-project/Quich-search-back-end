package middlewares

import (
	"net/http"
	"os"
)

func HandleCors(rw http.ResponseWriter) {
	rw.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS_ORIGIN_WHITELIST"))
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")
}
