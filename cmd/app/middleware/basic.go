package middleware

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"strings"
)

// Basic function
func Basic(checkAuth func(string, string) bool) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			login, password, err := getLoginPassword(r)
			if err != nil {
				log.Print(err)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if !checkAuth(login, password) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			handler.ServeHTTP(w, r)
		})
	}
}

func getLoginPassword(r *http.Request) (string, string, error) {

	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Basic" {
		return "", "", errors.New("invalid method")
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)

	if len(pair) != 2 {
		return "", "", errors.New("invalid data")
	}
	return pair[0], pair[1], nil
}
