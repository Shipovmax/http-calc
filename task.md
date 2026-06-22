# Task #2 — HTTP Calculator

## Цель

Написать REST API-сервер на Go с React-фронтендом. Сервер принимает JSON-запрос с двумя числами и оператором, выполняет вычисление и возвращает JSON-ответ. Фронтенд предоставляет UI для взаимодействия с API через браузер.

Главная учебная цель — освоить `net/http`: роутинг, JSON decode/encode, HTTP статус-коды, middleware. Это фундамент любого Go backend-сервиса.

---

## Acceptance Criteria

- [x] `POST /calculate {"a":10,"op":"+","b":5}` → HTTP 200, `{"result":15}`
- [x] `POST /calculate {"a":10,"op":"/","b":3}` → HTTP 200, `{"result":3.3333333333333335}`
- [x] `POST /calculate {"a":10,"op":"/","b":0}` → HTTP 400, `{"error":"деление на ноль"}`
- [x] `POST /calculate {"a":10,"op":"^","b":2}` → HTTP 400, `{"error":"неизвестный оператор: ^"}`
- [x] `POST /calculate <невалидный JSON>` → HTTP 400, `{"error":"некорректный JSON"}`
- [x] `GET /calculate` → HTTP 405
- [x] Каждый запрос логируется в stdout: метод, путь, статус, время выполнения
- [x] CORS middleware для React dev-сервера (localhost:5173)
- [x] `go vet ./...` проходит без предупреждений
- [x] `go.mod` содержит только `module` и `go` директивы
- [x] React-фронтенд с UI калькулятора, отправляет запросы к API

---

## Технические требования

### Backend (Go)

| Требование | Детали |
|---|---|
| HTTP-сервер | `net/http`, `http.ListenAndServe(":8080", mux)` |
| Роутер | `http.NewServeMux()` — не `DefaultServeMux` |
| Метод | проверка через паттерн `"POST /calculate"` в Go 1.22+ |
| JSON запрос | структура `Request{A float64, Op string, B float64}` |
| JSON ответ (успех) | `{"result": <float64>}` |
| JSON ответ (ошибка) | `{"error": "<сообщение>"}` |
| HTTP статусы | 200 OK, 400 Bad Request, 405 Method Not Allowed |
| Разбивка по файлам | `main.go`, `handler.go`, `calculator.go`, `middleware.go` |
| Middleware цепочка | `loggingMiddleware(corsMiddleware(mux))` |

### Frontend (React)

| Требование | Детали |
|---|---|
| Стек | Vite + React |
| Поля ввода | два числовых инпута + выбор оператора |
| Отправка | fetch POST на `http://localhost:8080/calculate` |
| Отображение | результат или сообщение об ошибке |
| Preview | живой preview JSON-запроса |
| Дизайн | тёмная/светлая тема (prefers-color-scheme) |

### Запрещено

- `panic` для обработки ошибок — только `error` return + HTTP 400/500
- Сторонние Go-пакеты (`gin`, `echo`, `chi`) — только `net/http`
- Глобальный `http.DefaultServeMux`
- `fmt.Println` в хендлерах — только в middleware

---

## Структура файлов

```
http-calc/
├── main.go          # NewServeMux, регистрация хендлеров, ListenAndServe
├── handler.go       # calculateHandler: decode → calculate → encode
├── calculator.go    # func calculate(a float64, op string, b float64) (float64, error)
├── middleware.go    # loggingMiddleware + corsMiddleware
├── go.mod           # module github.com/Shipovmax/http-calc
├── frontend/        # React-приложение
│   └── src/
│       ├── App.jsx
│       └── App.css
└── README.md
```

---

## Definition of Done

- [x] Все Acceptance Criteria выполнены
- [x] Код запушен на GitHub, релиз v1.0.0 создан
- [x] README.md актуален
- [ ] Ты можешь объяснить каждую строку кода вслух без подглядывания

---

## Следующий шаг после сдачи

После ревью переходим к **Task #3 — TODO CLI с файлами**: работа с файловой системой, `encoding/json` для персистентности, операции над слайсами структур.
