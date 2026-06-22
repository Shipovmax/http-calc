package main

import (
	"fmt"
	"net/http"
	"time"
)

// responseWriter оборачивает http.ResponseWriter чтобы перехватить статус-код
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// loggingMiddleware логирует метод, путь, статус и время выполнения каждого запроса
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// оборачиваем writer чтобы поймать статус после выполнения хендлера
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		fmt.Printf("%s %s %d %s\n", r.Method, r.URL.Path, rw.status, time.Since(start))
	})
}

// corsMiddleware разрешает запросы с localhost (нужно для React dev-сервера)
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// браузер шлёт OPTIONS preflight перед POST — отвечаем 204 и завершаем
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
