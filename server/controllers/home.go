package controllers

import "net/http"

func HomePageController(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("Welcome to ghostenv! 👻"))
}
