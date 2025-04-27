package main

import (
	"database/sql"
	"fmt"
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
	DB.Exec(query, args...)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/home", controllers.HomePageController)
	mux.HandleFunc("/login", controllers.LoginPageController(DB, sqlEngine))
	mux.HandleFunc("/register", controllers.RegisterController(DB, sqlEngine))

	handler := middlewares.NewSessionMiddleware(
		middlewares.LoggerMiddleware(
			mux,
		),
		DB,
		sqlEngine,
	)

	fmt.Println("👻 server started on http://localhost:3000")
	http.ListenAndServe(":3000", handler)
}
