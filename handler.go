package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Request is the JSON payload accepted by POST /calculate.
type Request struct {
	A  float64 `json:"a"`
	Op string  `json:"op"`
	B  float64 `json:"b"`
}

// Response is the JSON payload returned by POST /calculate.
// Exactly one of Result or Error is populated.
type Response struct {
	Result *float64 `json:"result,omitempty"`
	Error  string   `json:"error,omitempty"`
}

// writeJSON encodes v as JSON with the given status code and logs
// (rather than silently drops) any encoding failure — the response
// has already started, so at that point there's nothing left to do
// but record it for diagnostics.
func writeJSON(w http.ResponseWriter, status int, v Response) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// calculateHandler implements POST /calculate: decode the request body,
// run the calculation, and encode the result or error as JSON.
func calculateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// decode the request body into Request; malformed JSON must
	// produce a 400 response, never a panic
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, Response{Error: "invalid JSON"})
		return
	}

	// delegate to the business logic
	result, err := calculate(req.A, req.Op, req.B)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, Response{Error: err.Error()})
		return
	}

	// result is passed as a pointer so that omitempty does not drop
	// a legitimate 0.0 result from the JSON payload
	writeJSON(w, http.StatusOK, Response{Result: &result})
}
