package main

import (
	"net/http"

	"ops-storage/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/update/", handlers.Update)

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		panic(err)
	}

}
