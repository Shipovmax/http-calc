# http-calc — HTTP Calculator

> A REST API using pure `net/http` + a React frontend. Learning Project #2 as part of preparation for a Go Backend Developer role.

---

## For the Recruiter

### What It Is and Why

The second project in the roadmap marks the transition from CLI to HTTP. It utilizes the same calculator logic from project #1, but now wrapped in a REST API: the server accepts a JSON request, performs the calculation, and returns a JSON response. This is a fundamental pattern for any backend service.

The main goal is to master the standard `net/http` library without external frameworks: routing, JSON parsing, response building, HTTP status codes, and middleware. This is exactly how most internal services at major tech companies operate.

A React frontend is added as a UI layer on top of the API, demonstrating real-world browser-to-Go-server interaction via CORS.

### What This Project Demonstrates

| Skill | Implementation |
|---|---|
| HTTP Server | `net/http`, `http.ListenAndServe`, `http.ServeMux` |
| Method-based Routing | `"POST /calculate"` using the Go 1.22+ pattern |
| JSON Decode/Encode | `json.NewDecoder`, `json.NewEncoder` |
| HTTP Status Codes | 200, 400, 405 |
| Middleware Chain | `logging` → `CORS` → `handler` |
| CORS | preflight `OPTIONS` + headers for the React dev server |
| Separation of Concerns | handler / calculator / middleware split into separate files |
| React Frontend | Vite + React, Fetch API, live request preview |

### Stack

- **Backend:** Go 1.22+, standard library only
- **Frontend:** React + Vite (no UI frameworks)
- **Dependencies:** No external Go packages

---

## For the Developer

### Structure


```

http-calc/
├── main.go          # entry point: mux, handler registration, starting the server
├── handler.go       # HTTP handler: decode → calculate → encode
├── calculator.go    # business logic: calculate(a, op, b), completely decoupled from HTTP
├── middleware.go    # loggingMiddleware + corsMiddleware
├── go.mod           # contains only module and go directives
├── frontend/        # React application (Vite)
│   └── src/
│       ├── App.jsx  # calculator UI
│       └── App.css  # styles
└── README.md

```

### Architectural Decisions

#### Why `net/http` without a framework?

`gin` and `echo` are thin wrappers around `net/http`. By knowing the standard library, you deeply understand what a framework does under the hood. Tech interviews at BigTech companies often focus explicitly on `net/http`.

#### Why a middleware chain instead of a single handler?

```go
loggingMiddleware(corsMiddleware(mux))

```

Each middleware handles a single task, adhering to the Single Responsibility Principle. This makes it effortless to plug in authentication, rate limiting, or tracing without altering the core business logic.

#### Why `*float64` in the Response struct?

```go
type Response struct {
    Result *float64 `json:"result,omitempty"`
    Error  string   `json:"error,omitempty"`
}

```

The `omitempty` tag will omit the field from the JSON payload if the value is `0.0`. Using a pointer allows the application to distinguish between "the result is zero" and "the field was not provided/unset".

### Running the Project

**Server:**

```bash
go run .
# → server running on :8080

```

**Frontend:**

```bash
cd frontend
npm install
npm run dev
# → open http://localhost:5173

```

### API

```
POST /calculate
Content-Type: application/json

{"a": <number>, "op": "<operator>", "b": <number>}

```

**Operators:** `+`, `-`, `*`, `/`

### Examples

```bash
# Addition
curl -s -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"a": 10, "op": "+", "b": 5}'
# {"result":15}

# Division
curl -s -X POST http://localhost:8080/calculate \
  -d '{"a": 10, "op": "/", "b": 3}'
# {"result":3.3333333333333335}

```

### Error Handling

```bash
# Division by zero → HTTP 400
curl -s -X POST http://localhost:8080/calculate -d '{"a":10,"op":"/","b":0}'
# {"error":"division by zero"}

# Unknown operator → HTTP 400
curl -s -X POST http://localhost:8080/calculate -d '{"a":10,"op":"^","b":2}'
# {"error":"unknown operator: ^"}

# Invalid JSON → HTTP 400
curl -s -X POST http://localhost:8080/calculate -d 'not json'
# {"error":"invalid JSON"}

# GET request → HTTP 405
curl -s -X GET http://localhost:8080/calculate
# Method Not Allowed

```

