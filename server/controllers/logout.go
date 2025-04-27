package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	engine "github.com/rabbit-backend/template"
)

func LogoutController(db *sql.DB, e *engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_id")
		if err != nil {
			log.Println(err)
			http.Error(w, "Unkown error", http.StatusInternalServerError)
			return
		}

		sessionId := sessionCookie.Value
		query, args := e.Execute(
			"templates/sql/delete_user_session.sql.tmpl",
			map[string]string{
				"session_id": sessionId,
			},
		)

		if _, err := db.Exec(query, args...); err != nil {
			log.Println(err)
			http.Error(w, "Unkown error", http.StatusInternalServerError)
			return
		}

		cookie := &http.Cookie{
			Name:    "session_id",
			Expires: time.Unix(0, 0),
		}

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
