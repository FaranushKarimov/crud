package middleware

import "net/http"

//CheckHeader func
func CheckHeader(header, value string) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if value != r.Header.Get(header) {
				http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			handler.ServeHTTP(rw, r)
		})
	}
}
