package controllers

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
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

			passwordHash := sha256.Sum256([]byte(r.FormValue("password")))

			query, args := e.Execute(
				"templates/sql/create_user.sql.tmpl",
				map[string]string{
					"username": r.FormValue("username"),
					"password": fmt.Sprintf("%x", passwordHash),
				},
			)

			log.Println(query, args)
			if _, err := db.Exec(query, args...); err != nil {
				log.Println(err)
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
				return
			}

			w.Write([]byte("ok"))
			return
		}

		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
