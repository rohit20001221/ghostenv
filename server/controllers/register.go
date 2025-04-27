package controllers

import (
	"database/sql"
	"net/http"
	"os"

	engine "github.com/rabbit-backend/template"
)

func RegisterController(db *sql.DB, e *engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			f, _ := os.Open("templates/html/register.html")

			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", "text/html")
			f.WriteTo(w)
			return
		}

		if r.Method == http.MethodPost {
			r.ParseForm()

			query, args := e.Execute(
				"templates/sql/create_user.sql.tmpl",
				map[string]string{
					"username": r.FormValue("username"),
					"password": r.FormValue("password"),
				},
			)

			if _, err := db.Exec(query, args...); err != nil {
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
				return
			}

			w.Write([]byte("ok"))
			return
		}

		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
