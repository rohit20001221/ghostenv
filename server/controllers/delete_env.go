package controllers

import (
	"database/sql"
	"net/http"

	engine "github.com/rabbit-backend/template"
)

func DeleteEnvVariable(db *sql.DB, e *engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
