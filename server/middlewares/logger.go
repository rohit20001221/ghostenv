package middlewares

import (
	"log"
	"net/http"
)

func LoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("[x] incomming request [user]:", r.Context().Value(reqkey("user")))

		h.ServeHTTP(w, r)
	})
}
