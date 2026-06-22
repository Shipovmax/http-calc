# Task #2 — HTTP Calculator

## Цель

Написать REST API-сервер на Go, который принимает JSON-запрос с двумя числами и оператором, выполняет вычисление и возвращает JSON-ответ. Главная учебная цель — освоить `net/http`: роутинг, JSON decode/encode, HTTP статус-коды, middleware. Это фундамент любого Go backend-сервиса и обязательная тема на собеседованиях в Ozon/WB/Сбер.

---

## Acceptance Criteria

- [ ] `POST /calculate {"a":10,"op":"+","b":5}` → HTTP 200, `{"result":15}`
- [ ] `POST /calculate {"a":10,"op":"/","b":3}` → HTTP 200, `{"result":3.3333333333333335}`
- [ ] `POST /calculate {"a":10,"op":"/","b":0}` → HTTP 400, `{"error":"деление на ноль"}`
- [ ] `POST /calculate {"a":10,"op":"^","b":2}` → HTTP 400, `{"error":"неизвестный оператор: ^"}`
- [ ] `POST /calculate <невалидный JSON>` → HTTP 400, `{"error":"некорректный JSON"}`
- [ ] `GET /calculate` → HTTP 405
- [ ] Каждый запрос логируется в stdout: метод, путь, статус, время выполнения
- [ ] `go vet ./...` проходит без предупреждений
- [ ] `go.mod` содержит только `module` и `go` директивы

---

## Технические требования

### Обязательно

| Требование | Детали |
|---|---|
| HTTP-сервер | `net/http`, `http.ListenAndServe(":8080", mux)` |
| Роутер | `http.NewServeMux()` — не `DefaultServeMux` |
| Метод | проверка через паттерн `"POST /calculate"` в Go 1.22+ |
| JSON запрос | структура `Request{A float64, Op string, B float64}` |
| JSON ответ (успех) | `{"result": <float64>}` |
| JSON ответ (ошибка) | `{"error": "<сообщение>"}` |
| HTTP статусы | 200 OK, 400 Bad Request, 405 Method Not Allowed |
| Разбивка по файлам | минимум 3 файла: `main.go`, `handler.go`, `calculator.go` |
| Middleware | логирующий: метод + путь + статус + duration |

### Запрещено

- `panic` для обработки ошибок — только `error` return + HTTP 400/500
- Сторонние пакеты (`gin`, `echo`, `chi`, `gorilla/mux`) — только `net/http`
- Глобальный `http.DefaultServeMux` — использовать явный `http.NewServeMux()`
- `fmt.Println` в хендлерах для логирования — выносить в middleware
- Смешивать бизнес-логику и HTTP в одном файле

---

## Темы Go, которые ты прокачиваешь

> Это не просто список — это checklist того, что **обязан использовать** в реализации.

- **`net/http`** — `http.ListenAndServe`, `http.NewServeMux`, `http.HandlerFunc`, `http.ResponseWriter`, `*http.Request`
- **`encoding/json`** — `json.NewDecoder(r.Body).Decode(&req)` для чтения запроса, `json.NewEncoder(w).Encode(resp)` для ответа
- **`http.Handler` интерфейс** — обёртка для middleware: `func(next http.Handler) http.Handler`
- **HTTP статус-коды** — `w.WriteHeader(http.StatusBadRequest)` до записи тела
- **Структуры с json-тегами** — `json:"a"`, `json:"op"`, `json:"b"`, `json:"result"`, `json:"error"`
- **`time.Since`** — замер времени выполнения запроса в middleware
- **Множественный возврат + именованные типы ошибок** — `calculate(a float64, op string, b float64) (float64, error)`

---

## Структура файлов

```
http-calc/
├── main.go          # NewServeMux, регистрация хендлеров, ListenAndServe
├── handler.go       # calculateHandler: decode → calculate → encode
├── calculator.go    # func calculate(a, b float64, op string) (float64, error)
├── middleware.go    # func loggingMiddleware(next http.Handler) http.Handler
├── go.mod           # module github.com/Shipovmax/http-calc
└── README.md
```

---

## Подсказки по архитектуре

```go
// calculator.go — чистая функция, ничего не знает про HTTP
func calculate(a float64, op string, b float64) (float64, error)

// handler.go — только HTTP-слой
type Request struct {
    A  float64 `json:"a"`
    Op string  `json:"op"`
    B  float64 `json:"b"`
}

type Response struct {
    Result *float64 `json:"result,omitempty"`
    Error  string   `json:"error,omitempty"`
}

func calculateHandler(w http.ResponseWriter, r *http.Request)

// middleware.go — оборачивает любой хендлер
func loggingMiddleware(next http.Handler) http.Handler

// main.go
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("POST /calculate", calculateHandler)
    http.ListenAndServe(":8080", loggingMiddleware(mux))
}
```

> Обрати внимание на `*float64` в Response — указатель нужен чтобы `omitempty` работал корректно с нулевым значением `0.0`.

---

## Definition of Done

1. Все Acceptance Criteria выполнены
2. Код запушен на GitHub в репозиторий `http-calc`
3. README.md в репозитории соответствует шаблону проекта
4. Ты можешь объяснить каждую строку кода вслух без подглядывания

---

## Следующий шаг после сдачи

После ревью переходим к **Task #3 — TODO CLI с файлами**: работа с файловой системой, `encoding/json` для персистентности, операции над слайсами структур.
