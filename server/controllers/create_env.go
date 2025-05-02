package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	engine "github.com/rabbit-backend/template"
)

func CreateEnvVariable(db *sql.DB, e *engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				log.Println(err)

				http.Error(w, "Unexpected error", http.StatusInternalServerError)
				return
			}

			query, args := e.Execute(
				"templates/sql/create_environment_variable.sql.tmpl",
				map[string]string{
					"key":         r.FormValue("key"),
					"value":       r.FormValue("value"),
					"application": r.FormValue("app_name"),
				},
			)

			if _, err := db.Exec(query, args...); err != nil {
				log.Println(err)

				http.Error(w, "Unexpected error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{
				"key":         r.FormValue("key"),
				"value":       r.FormValue("value"),
				"application": r.FormValue("app_name"),
			})

			return
		}

		http.Error(w, "Method Not allowed", http.StatusMethodNotAllowed)
	}
}
