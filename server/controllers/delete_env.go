package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	engine "github.com/rabbit-backend/template"
)

type DeleteEnvReqBody struct {
	Key   string `json:"key"`
	AppId string `json:"app_id"`
}

func DeleteEnvVariable(db *sql.DB, e *engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var body DeleteEnvReqBody
			json.NewDecoder(r.Body).Decode(&body)

			query, args := e.Execute(
				"templates/sql/delete_env_variable.sql",
				body,
			)

			if _, err := db.Exec(query, args...); err != nil {
				log.Println(err)

				http.Error(w, "unexpected error", http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
