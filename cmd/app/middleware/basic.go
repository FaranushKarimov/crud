package middleware

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
)

//Basic func
func Basic(auth func(ctx context.Context, login string, pass string) bool) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			value := r.Header.Get("Authorization")
			if value == "" {
				http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			slitted := strings.Split(value, " ")
			if len(slitted) != 2 {
				http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			credentialsPart := slitted[1]
			data, err := base64.StdEncoding.DecodeString(credentialsPart)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			credentials := strings.Split(string(data), ":")
			if len(credentials) != 2 {
				http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			login := credentials[0]
			password := credentials[1]

			if !auth(r.Context(), login, password) {
				http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			handler.ServeHTTP(rw, r)
		})
	}
}
