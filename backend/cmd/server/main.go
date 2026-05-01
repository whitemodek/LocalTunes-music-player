package main

import (
	"log"
	"net/http"

	"localtunes/internal/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/tracks", handler.TrackHandler)

	log.Println("LocalTunes backend listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
