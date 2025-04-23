package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("Welcome to ghostenv!"))
	})

	fmt.Println("👻 server started on http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}
