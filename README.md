# http-calc — HTTP Калькулятор

> REST API на чистом `net/http` + React-фронтенд. Учебный проект #2 в рамках подготовки к позиции Go Backend Developer.

---

## Для рекрутера

### Что это и зачем

Второй проект в roadmap — переход от CLI к HTTP. Та же логика калькулятора из проекта #1, но теперь обёрнутая в REST API: сервер принимает JSON-запрос, считает, возвращает JSON-ответ. Это фундаментальный паттерн любого backend-сервиса.

Главная цель — освоить стандартный `net/http` без фреймворков: роутинг, парсинг JSON, формирование ответа, HTTP-статусы, middleware. Именно так работает большинство внутренних сервисов в Ozon и WB.

React-фронтенд добавлен как UI-слой поверх API — демонстрирует реальное взаимодействие браузер ↔ Go-сервер через CORS.

### Что демонстрирует этот проект

| Навык | Реализация |
|---|---|
| HTTP-сервер | `net/http`, `http.ListenAndServe`, `http.ServeMux` |
| Роутинг с методом | `"POST /calculate"` через Go 1.22+ паттерн |
| JSON decode/encode | `json.NewDecoder`, `json.NewEncoder` |
| HTTP статус-коды | 200, 400, 405 |
| Middleware цепочка | logging → CORS → handler |
| CORS | preflight OPTIONS + заголовки для React dev-сервера |
| Разделение ответственности | handler / calculator / middleware в отдельных файлах |
| React-фронтенд | Vite + React, fetch API, live preview запроса |

### Стек

- **Backend:** Go 1.22+, только стандартная библиотека
- **Frontend:** React + Vite (без UI-фреймворков)
- **Зависимости:** нет внешних Go-пакетов

---

## Для разработчика

### Структура

```
http-calc/
├── main.go          # точка входа: mux, регистрация хендлеров, запуск сервера
├── handler.go       # HTTP-хендлер: decode → calculate → encode
├── calculator.go    # бизнес-логика: calculate(a, op, b), ничего про HTTP
├── middleware.go    # loggingMiddleware + corsMiddleware
├── go.mod           # только module и go директивы
├── frontend/        # React-приложение (Vite)
│   └── src/
│       ├── App.jsx  # калькулятор UI
│       └── App.css  # стили
└── README.md
```

### Архитектурные решения

#### Почему `net/http` без фреймворка?

`gin` и `echo` — тонкие обёртки над `net/http`. Зная стандартную библиотеку, понимаешь что делает фреймворк. На собеседовании в BigTech спросят именно про `net/http`.

#### Почему middleware цепочка, а не один обработчик?

```go
loggingMiddleware(corsMiddleware(mux))
```

Каждый middleware отвечает за одно — принцип единственной ответственности. Легко добавить auth, rate-limit или tracing без изменения бизнес-логики.

#### Почему `*float64` в Response?

```go
type Response struct {
    Result *float64 `json:"result,omitempty"`
    Error  string   `json:"error,omitempty"`
}
```

`omitempty` со значением `0.0` выкинет поле из JSON. Указатель позволяет отличить "результат равен нулю" от "поле не задано".

### Запуск

**Сервер:**
```bash
go run .
# → сервер на :8080
```

**Фронтенд:**
```bash
cd frontend
npm install
npm run dev
# → открой http://localhost:5173
```

### API

```
POST /calculate
Content-Type: application/json

{"a": <число>, "op": "<оператор>", "b": <число>}
```

**Операторы:** `+` `-` `*` `/`

### Примеры

```bash
# Сложение
curl -s -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"a": 10, "op": "+", "b": 5}'
# {"result":15}

# Деление
curl -s -X POST http://localhost:8080/calculate \
  -d '{"a": 10, "op": "/", "b": 3}'
# {"result":3.3333333333333335}
```

### Обработка ошибок

```bash
# Деление на ноль → HTTP 400
curl -s -X POST http://localhost:8080/calculate -d '{"a":10,"op":"/","b":0}'
# {"error":"деление на ноль"}

# Неизвестный оператор → HTTP 400
curl -s -X POST http://localhost:8080/calculate -d '{"a":10,"op":"^","b":2}'
# {"error":"неизвестный оператор: ^"}

# Невалидный JSON → HTTP 400
curl -s -X POST http://localhost:8080/calculate -d 'not json'
# {"error":"некорректный JSON"}

# GET запрос → HTTP 405
curl -s -X GET http://localhost:8080/calculate
# Method Not Allowed
```
