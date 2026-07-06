package controllers

import (
	"log"
	"net/http"

	"github.com/rohit20001221/ghostenv-server/lib"
)

func IndexPageController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		html, err := lib.RenderTemplate("templates/html/index.html", map[string]any{})
		if err != nil {
			log.Println(err)

			http.Error(w, "Unknown Error", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(html))
	}
}
