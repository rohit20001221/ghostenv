package main

import (
	"database/sql"
	"fmt"
	"net/http"

	engine "github.com/rabbit-backend/template"
	"github.com/rohit20001221/ghostenv-server/templates/db"
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("Welcome to ghostenv!"))
	})

	fmt.Println("👻 server started on http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}
