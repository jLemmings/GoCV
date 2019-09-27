package middleware

import (
	"context"
	"fmt"
	"github.com/jLemmings/GoCV/models"
	"github.com/jLemmings/GoCV/utils"
	"log"
	"net/http"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI + r.Method)

		uid := r.Header.Get("authorization")
		if uid != "" {
			fmt.Println("UID: ", uid)
			user, err := models.GetAuth().VerifyIDToken(context.Background(), uid)
			utils.HandleErr(err)
			claims := user.Claims
			if admin, ok := claims["admin"]; ok {
				if admin.(bool) {
					log.Println("I AM ADMIN")
					next.ServeHTTP(w, r)
				}
			} else {
				fmt.Println("IN NOT OK: ", ok)
				http.Error(w, "Forbidden", http.StatusForbidden)
			}
		} else {
			fmt.Println("No UID")
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
		// Do stuff here
		log.Println(r.RequestURI + r.Method)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
	})
}
