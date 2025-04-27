package middlewares

import (
	"log"
	"net/http"

	"github.com/rohit20001221/ghostenv-server/types"
)

func LoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("[x] incomming request [user]:", r.Context().Value(types.SessionKey("user")))

		h.ServeHTTP(w, r)
	})
}
