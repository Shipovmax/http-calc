# Task #2 — HTTP Calculator

## Goal

Write a REST API server in Go with a React frontend. The server accepts a JSON request with two numbers and an operator, performs the calculation, and returns a JSON response. The frontend provides a UI for interacting with the API through the browser.

The main learning goal is to master `net/http`: routing, JSON decode/encode, HTTP status codes, middleware. This is the foundation of any Go backend service.

---

## Acceptance Criteria

- [x] `POST /calculate {"a":10,"op":"+","b":5}` → HTTP 200, `{"result":15}`
- [x] `POST /calculate {"a":10,"op":"/","b":3}` → HTTP 200, `{"result":3.3333333333333335}`
- [x] `POST /calculate {"a":10,"op":"/","b":0}` → HTTP 400, `{"error":"division by zero"}`
- [x] `POST /calculate {"a":10,"op":"^","b":2}` → HTTP 400, `{"error":"unknown operator: ^"}`
- [x] `POST /calculate <invalid JSON>` → HTTP 400, `{"error":"invalid JSON"}`
- [x] `GET /calculate` → HTTP 405
- [x] Every request is logged to stdout: method, path, status, execution time
- [x] CORS middleware for the React dev server (localhost:5173)
- [x] `go vet ./...` passes without warnings
- [x] `go.mod` contains only the `module` and `go` directives
- [x] React frontend with a calculator UI that sends requests to the API

---

## Technical Requirements

### Backend (Go)

| Requirement | Details |
|---|---|
| HTTP server | `net/http`, `http.ListenAndServe(":8080", mux)` |
| Router | `http.NewServeMux()` — not `DefaultServeMux` |
| Method | checked via the `"POST /calculate"` pattern in Go 1.22+ |
| JSON request | `Request{A float64, Op string, B float64}` struct |
| JSON response (success) | `{"result": <float64>}` |
| JSON response (error) | `{"error": "<message>"}` |
| HTTP statuses | 200 OK, 400 Bad Request, 405 Method Not Allowed |
| File split | `main.go`, `handler.go`, `calculator.go`, `middleware.go` |
| Middleware chain | `loggingMiddleware(corsMiddleware(mux))` |

### Frontend (React)

| Requirement | Details |
|---|---|
| Stack | Vite + React |
| Input fields | two numeric inputs + operator selector |
| Submission | fetch POST to `http://localhost:8080/calculate` |
| Display | result or error message |
| Preview | live preview of the JSON request |
| Design | dark/light theme (prefers-color-scheme) |

### Forbidden

- `panic` for error handling — only `error` return values + HTTP 400/500
- Third-party Go packages (`gin`, `echo`, `chi`) — `net/http` only
- The global `http.DefaultServeMux`
- `fmt.Println` in handlers — only in middleware

---

## File Structure

```
http-calc/
├── main.go          # NewServeMux, handler registration, ListenAndServe
├── handler.go       # calculateHandler: decode → calculate → encode
├── calculator.go    # func calculate(a float64, op string, b float64) (float64, error)
├── middleware.go     # loggingMiddleware + corsMiddleware
├── go.mod           # module github.com/Shipovmax/http-calc
├── frontend/        # React application
│   └── src/
│       ├── App.jsx
│       └── App.css
└── README.md
```

---

## Definition of Done

- [x] All acceptance criteria met
- [x] Code pushed to GitHub, release v1.0.0 created
- [x] README.md up to date
- [ ] You can explain every line of code out loud without looking

---

## Next Step After Submission

After review, moving on to **Task #3 — TODO CLI with files**: filesystem operations, `encoding/json` for persistence, operations on slices of structs.
