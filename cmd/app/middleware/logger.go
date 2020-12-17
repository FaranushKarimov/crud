package middleware

import (
	"log"
	"net/http"
)

//Logger func
func Logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("START: %s %s", r.Method, r.URL.Path)

		handler.ServeHTTP(rw, r)

		log.Printf("END: %s %s", r.Method, r.URL.Path)
	})
}
