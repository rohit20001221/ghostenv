package middlewares

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	engine "github.com/rabbit-backend/template"
)

type reqkey string

func NewSessionMiddleware(h http.Handler, db *sql.DB, e *engine.Engine) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user any
		if cookie, err := r.Cookie("session_id"); err != nil {
			user = nil
		} else {
			sessionId := cookie.Value
			query, args := e.Execute(
				"templates/sql/get_user_session.sql.tmpl",
				map[string]string{"session_id": sessionId},
			)

			if err := db.QueryRow(query, args...).Scan(&user); err != nil {
				log.Println(err)
				user = nil
			}
		}

		ctx := context.WithValue(r.Context(), reqkey("user"), user)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
