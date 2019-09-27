package middleware

import (
	"log"
	"net/http"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		enableCors(&w)
		log.Println(r.RequestURI + r.Method)
		next.ServeHTTP(w, r)
	})
}
