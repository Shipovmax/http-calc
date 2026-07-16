package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Go 1.22+: the "POST /calculate" pattern checks the method automatically,
	// any other method on this path returns 405
	mux.HandleFunc("POST /calculate", calculateHandler)

	log.Println("server listening on :8080")
	if err := http.ListenAndServe(":8080", loggingMiddleware(corsMiddleware(mux))); err != nil {
		log.Fatal(err)
	}
}
