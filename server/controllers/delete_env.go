package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	engine "github.com/rabbit-backend/template"
	"github.com/rohit20001221/ghostenv-server/types"
)

type DeleteEnvReqBody struct {
	Key   string `json:"key"`
	AppId string `json:"application"`
}

func DeleteEnvVariable(db *sql.DB, e *engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var body DeleteEnvReqBody
			json.NewDecoder(r.Body).Decode(&body)

			query, args := e.Execute(
				"templates/sql/get_application_by_app_name.sql.tmpl",
				map[string]string{
					"app_name": body.AppId,
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

			query, args = e.Execute(
				"templates/sql/delete_env_variable.sql",
				body,
			)

			log.Println(query, args)

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
