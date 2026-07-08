package main

import (
	"database/sql"
	"log"
	"net/http"

	engine "github.com/rabbit-backend/template"
	"github.com/rohit20001221/ghostenv-server/controllers"
	"github.com/rohit20001221/ghostenv-server/db"
	"github.com/rohit20001221/ghostenv-server/middlewares"
)

var sqlEngine *engine.Engine
var DB *sql.DB

func init() {
	sqlEngine = engine.NewEngineWithPlaceHolder(engine.NewPostgresPlaceHolder())
	DB = db.CreateConnection()

	query, args := sqlEngine.Execute("templates/sql/create_db.sql.tmpl", nil)
	log.Println("[x] Initializing the application...")
	if _, err := DB.Exec(query, args...); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/dashboard",
		middlewares.LoginRequiredMiddleware(
			http.HandlerFunc(
				controllers.DashboardController(DB, sqlEngine),
			),
		),
	)

	mux.Handle("/",
		http.HandlerFunc(
			controllers.HomeController(),
		),
	)

	mux.HandleFunc("/auth/login",
		controllers.LoginController(
			DB,
			sqlEngine,
		),
	)

	mux.HandleFunc("/auth/register",
		controllers.RegisterController(
			DB,
			sqlEngine,
		),
	)

	mux.HandleFunc("/auth/logout",
		controllers.LogoutController(
			DB,
			sqlEngine,
		),
	)

	mux.Handle("/create_application", middlewares.LoginRequiredMiddleware(http.HandlerFunc(controllers.CreateApplicationController(DB, sqlEngine))))
	mux.Handle("/application/{app_name}", middlewares.LoginRequiredMiddleware(http.HandlerFunc(controllers.ApplicationController(DB, sqlEngine))))
	mux.Handle("DELETE /application/{app_name}/delete", middlewares.LoginRequiredMiddleware(http.HandlerFunc(controllers.DeleteApplicationController(DB, sqlEngine))))
	mux.Handle("/create_env", middlewares.LoginRequiredMiddleware(http.HandlerFunc(controllers.CreateEnvVariable(DB, sqlEngine))))
	mux.Handle("/delete_env", middlewares.LoginRequiredMiddleware(http.HandlerFunc(controllers.DeleteEnvVariable(DB, sqlEngine))))
	mux.Handle("/env/{app_name}", middlewares.NewBasicAuthMiddleware(http.HandlerFunc(controllers.PullEnvironmentVariables(DB, sqlEngine))))

	handler := middlewares.NewSessionMiddleware(
		middlewares.LoggerMiddleware(
			mux,
		),
		DB,
		sqlEngine,
	)

	log.Println("👻 server started on http://localhost:3000")
	http.ListenAndServe(":3000", handler)
}
