package middlewares

import (
	"net/http"

	"github.com/rohit20001221/ghostenv-server/types"
)

func LoginRequiredMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(types.SessionKey("user"))
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		h.ServeHTTP(w, r)
	})
}
