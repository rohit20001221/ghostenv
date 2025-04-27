package controllers

import (
	"net/http"
	"os"
)

func LoginPageController(w http.ResponseWriter, r *http.Request) {
	f, _ := os.Open("templates/html/login.html")

	w.WriteHeader(http.StatusOK)
	f.WriteTo(w)
}
