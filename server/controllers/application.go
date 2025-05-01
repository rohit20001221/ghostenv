package controllers

import (
	"database/sql"
	"log"
	"net/http"

	engine "github.com/rabbit-backend/template"
	"github.com/rohit20001221/ghostenv-server/lib"
	"github.com/rohit20001221/ghostenv-server/types"
)

func ApplicationController(db *sql.DB, e *engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appName := r.PathValue("app_name")

		query, args := e.Execute(
			"templates/sql/get_application_by_app_name.sql.tmpl",
			map[string]string{
				"app_name": appName,
			},
		)

		row := db.QueryRow(query, args...)

		var userId, app, descritpion string
		if err := row.Scan(&app, &descritpion, &userId); err != nil {
			log.Println(err)
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
			return
		}

		if r.Context().Value(types.SessionKey("user")) != userId {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		html, err := lib.RenderTemplate("templates/html/app.html.tmpl", map[string]string{
			"app_name":    app,
			"description": descritpion,
		})

		if err != nil {
			log.Println(err)

			http.Error(w, "Unexpected error", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(html))
	}
}
