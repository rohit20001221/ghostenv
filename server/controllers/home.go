package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	engine "github.com/rabbit-backend/template"
	"github.com/rohit20001221/ghostenv-server/types"
)

func HomePageController(db *sql.DB, e *engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		// get list of all the applications
		user := r.Context().Value(types.SessionKey("user")).(string)

		query, args := e.Execute(
			"templates/sql/get_applications_by_user.sql.tmpl",
			map[string]string{
				"user_id": user,
			},
		)

		rows, err := db.Query(query, args...)
		if err != nil {
			log.Println(err)

			http.Error(w, "Unknown Error", http.StatusInternalServerError)
			return
		}

		applications := make([]map[string]string, 0)

		for rows.Next() {
			var appName, description string
			rows.Scan(&appName, &description)

			app := map[string]string{
				"appName":     appName,
				"description": description,
			}

			applications = append(applications, app)
		}

		log.Println(applications)

		f, _ := os.Open("templates/html/home.html")
		f.WriteTo(w)
	}
}
