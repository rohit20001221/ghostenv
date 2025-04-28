package controllers

import (
	"net/http"
	"os"
)

func HomePageController(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	f, _ := os.Open("templates/html/home.html")
	f.WriteTo(w)
}
