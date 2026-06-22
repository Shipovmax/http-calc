package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Go 1.22+: паттерн "POST /calculate" проверяет метод автоматически, GET вернёт 405
	mux.HandleFunc("POST /calculate", calculateHandler)

	log.Println("сервер запущен на :8080")
	if err := http.ListenAndServe(":8080", loggingMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}
