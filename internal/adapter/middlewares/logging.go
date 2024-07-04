package middlewares

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Request URI: %s", r.RequestURI)
		next.ServeHTTP(w, r)
		log.Printf("Request processed in %s, Method: %s, Url: %s, Status Code: %v", time.Since(start), r.Method, r.URL.Path, http.StatusOK)
	})
}
