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
	sqlEngine = engine.NewEngineWithPlaceHolder(engine.NewSqlitePlaceholder())
	DB = db.CreateConnection()

	query, args := sqlEngine.Execute("templates/sql/create_db.sql.tmpl", nil)
	log.Println("[x] Initializing the application...")
	DB.Exec(query, args...)
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/home", middlewares.LoginRequiredMiddleware(http.HandlerFunc(
		controllers.HomePageController(DB, sqlEngine),
	)))

	mux.HandleFunc("/login", controllers.LoginPageController(DB, sqlEngine))
	mux.HandleFunc("/register", controllers.RegisterController(DB, sqlEngine))
	mux.HandleFunc("/logout", controllers.LogoutController(DB, sqlEngine))

	mux.Handle("/create_application", middlewares.LoginRequiredMiddleware(http.HandlerFunc(controllers.CreateApplicationController(DB, sqlEngine))))

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
