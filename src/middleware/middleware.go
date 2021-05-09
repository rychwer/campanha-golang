package middleware

import (
	"log"
	"net/http"
)

func Logger(proximaFuncao http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\n%s %s %s", r.Method, r.RequestURI, r.Host)
		proximaFuncao.ServeHTTP(w, r)
	})
}