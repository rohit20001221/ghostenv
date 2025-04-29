package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	engine "github.com/rabbit-backend/template"
	"github.com/rohit20001221/ghostenv-server/types"
)

func CreateApplicationController(db *sql.DB, e *engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			f, _ := os.Open("templates/html/create_application.html")
			f.WriteTo(w)

			return
		}

		if r.Method == http.MethodPost {
			r.ParseForm()

			query, args := e.Execute(
				"templates/sql/create_application.sql.tmpl",
				map[string]string{
					"app_name":    r.FormValue("app_name"),
					"description": r.FormValue("description"),
					"user":        r.Context().Value(types.SessionKey("user")).(string),
				},
			)

			if _, err := db.Exec(query, args...); err != nil {
				log.Println(err)

				http.Error(w, "unexpected error", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
