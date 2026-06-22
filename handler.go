package main

import (
	"encoding/json"
	"net/http"
)

type Request struct {
	A  float64 `json:"a"`
	Op string  `json:"op"`
	B  float64 `json:"b"`
}

type Response struct {
	Result *float64 `json:"result,omitempty"`
	Error  string   `json:"error,omitempty"`
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// читаем JSON из тела запроса в структуру Request
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Error: "некорректный JSON"})
		return
	}

	// передаём данные в бизнес-логику
	result, err := calculate(req.A, req.Op, req.B)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Error: err.Error()})
		return
	}

	// &result — указатель, чтобы omitempty не выкинул 0.0 из JSON
	json.NewEncoder(w).Encode(Response{Result: &result})
}
