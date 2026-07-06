package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	engine "github.com/rabbit-backend/template"
)

func PullEnvironmentVariables(db *sql.DB, e *engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appName := r.PathValue("app_name")

		query, args := e.Execute(
			"templates/sql/get_env_variables.sql.tmpl",
			map[string]string{
				"app_name": appName,
			},
		)

		log.Println(query, args)
		envVariables := make([]map[string]string, 0)

		rows, err := db.Query(query, args...)
		if err != nil {
			log.Println(err)
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			var key, value string

			err := rows.Scan(&key, &value)
			if err != nil {
				log.Println(err)
				continue
			}

			envVariables = append(envVariables, map[string]string{"key": key, "value": value})
		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(envVariables)
	}
}
