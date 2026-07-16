package main

import (
	"fmt"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture the status code
// written by the downstream handler, since http.ResponseWriter itself
// exposes no way to read it back.
type responseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader records the status code and forwards it to the
// underlying http.ResponseWriter.
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// loggingMiddleware logs the method, path, status code, and duration
// of every request.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// wrap the writer to observe the status code after the handler runs
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		fmt.Printf("%s %s %d %s\n", r.Method, r.URL.Path, rw.status, time.Since(start))
	})
}

// corsMiddleware allows requests from the React dev server's origin.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// the browser sends an OPTIONS preflight before POST — answer 204 and stop
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
