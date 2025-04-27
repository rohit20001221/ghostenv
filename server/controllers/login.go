package controllers

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	engine "github.com/rabbit-backend/template"
)

func LoginPageController(db *sql.DB, e *engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			f, _ := os.Open("templates/html/login.html")

			w.WriteHeader(http.StatusOK)
			f.WriteTo(w)
			return
		}

		if r.Method == http.MethodPost {
			r.ParseForm()

			query, args := e.Execute(
				"templates/sql/get_user.sql.tmpl",
				map[string]string{
					"username": r.FormValue("username"),
				},
			)

			var username string
			var passwordHashFromDB string

			if err := db.QueryRow(query, args...).Scan(&username, &passwordHashFromDB); err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			passwordHash := fmt.Sprintf("%x", sha256.Sum256([]byte(r.FormValue("password"))))
			passwordMatch := passwordHash == passwordHashFromDB

			if !passwordMatch {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			sessionId := uuid.New()

			query, args = e.Execute(
				"templates/sql/create_session.sql.tmpl",
				map[string]string{
					"sessionid": sessionId.String(),
					"username":  username,
				},
			)

			log.Println(query, args)

			if _, err := db.Exec(query, args...); err != nil {
				log.Println(err)

				http.Error(w, "Login Failed!! try after sometime", http.StatusInternalServerError)
				return
			}

			cookie := http.Cookie{}
			cookie.Name = "session_id"

			cookie.Value = sessionId.String()

			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}

		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
