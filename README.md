# http-calc — HTTP-калькулятор на Go

> REST API-сервер на чистом `net/http`. Учебный проект #2 в рамках подготовки к позиции Go Backend Developer.

---

## Для рекрутера

### Что это и зачем

Второй проект в roadmap — переход от CLI к HTTP. Та же логика калькулятора из проекта #1, но теперь обёрнутая в REST API: сервер принимает JSON-запрос, считает, возвращает JSON-ответ. Это фундаментальный паттерн любого backend-сервиса.

Главная цель — освоить стандартный `net/http` без фреймворков: роутинг, парсинг JSON из тела запроса, формирование JSON-ответа, коды статусов HTTP, и middleware для логирования. Именно так работает большинство внутренних сервисов в Ozon и WB — Go + `net/http` или минимальный роутер поверх него.

Проект намеренно обходится без `gin`, `echo`, `chi` — чтобы понять что происходит под капотом этих фреймворков, прежде чем их использовать.

### Что демонстрирует этот проект

| Навык | Реализация |
|---|---|
| HTTP-сервер | `net/http`, `http.ListenAndServe`, `http.ServeMux` |
| Роутинг | регистрация хендлеров через `mux.HandleFunc` |
| JSON decode/encode | `encoding/json`, `json.NewDecoder`, `json.NewEncoder` |
| HTTP статус-коды | `http.StatusOK`, `http.StatusBadRequest`, `http.StatusMethodNotAllowed` |
| Middleware | логирующий middleware через `http.Handler` обёртку |
| Разделение ответственности | handler / calculator / middleware в отдельных файлах |
| Обработка ошибок | единый формат JSON-ошибки для всех случаев |

### Стек

- **Язык:** Go 1.22+
- **Зависимости:** только стандартная библиотека
- **Платформа:** Linux / macOS / Windows

---

## Для разработчика

### Архитектурные решения

#### Почему `net/http` без фреймворка?

`gin` и `echo` — тонкие обёртки над `net/http`. Зная стандартную библиотеку, ты понимаешь что делает фреймворк. Не зная — используешь магию. На собеседовании в BigTech спросят именно про `net/http`.

#### Почему `http.ServeMux`, а не глобальный `DefaultServeMux`?

```go
// Плохо — глобальное состояние, нельзя тестировать изолированно
http.HandleFunc("/calculate", handler)
http.ListenAndServe(":8080", nil)

// Хорошо — явный mux, можно передавать в httptest.NewServer
mux := http.NewServeMux()
mux.HandleFunc("POST /calculate", handler)
http.ListenAndServe(":8080", mux)
```

#### Почему единый формат JSON-ошибки?

Клиент обязан знать структуру ответа заранее — и для успеха, и для ошибки. Смешивать строку ошибки и объект результата — антипаттерн.

```json
// Успех
{"result": 15}

// Ошибка
{"error": "деление на ноль"}
```

Один и тот же клиентский код обрабатывает оба случая: проверяй поле `error`, если пустое — читай `result`.

#### Почему метод проверяется в хендлере, а не в роутере?

В Go 1.22+ `ServeMux` поддерживает `"POST /calculate"` прямо в паттерне. Используй это — не городи `if r.Method != "POST"` внутри хендлера.

#### Почему middleware для логирования, а не `log.Println` в каждом хендлере?

Сквозная функциональность (логирование, аутентификация, метрики) не должна быть размазана по бизнес-логике. Middleware оборачивает любой `http.Handler` и работает для всех эндпоинтов сразу — это основа production-архитектуры.

### Структура

```
http-calc/
├── main.go          # точка входа: создание mux, регистрация хендлеров, запуск сервера
├── handler.go       # HTTP-хендлер: decode запроса, вызов calculate, encode ответа
├── calculator.go    # бизнес-логика: функция calculate, не знает про HTTP
├── middleware.go    # логирующий middleware
├── go.mod
└── README.md
```

### Установка и запуск

```bash
git clone https://github.com/Shipovmax/http-calc
cd http-calc
go run .
# Сервер запущен на :8080
```

### Использование

```
POST /calculate
Content-Type: application/json

{"a": <число>, "op": "<оператор>", "b": <число>}
```

**Поддерживаемые операторы:** `+` `-` `*` `/`

### Примеры

```bash
curl -s -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"a": 10, "op": "+", "b": 5}'
# {"result":15}

curl -s -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"a": 10, "op": "/", "b": 3}'
# {"result":3.3333333333333335}

curl -s -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"a": -5, "op": "*", "b": 2}'
# {"result":-10}
```

### Обработка ошибок

```bash
# Деление на ноль
curl -s -X POST http://localhost:8080/calculate \
  -d '{"a": 10, "op": "/", "b": 0}'
# HTTP 400: {"error":"деление на ноль"}

# Неизвестный оператор
curl -s -X POST http://localhost:8080/calculate \
  -d '{"a": 10, "op": "^", "b": 2}'
# HTTP 400: {"error":"неизвестный оператор: ^"}

# Некорректный JSON
curl -s -X POST http://localhost:8080/calculate \
  -d 'not json'
# HTTP 400: {"error":"некорректный JSON"}

# Неверный метод
curl -s -X GET http://localhost:8080/calculate
# HTTP 405: Method Not Allowed
```

### Запуск без сборки

```bash
go run .
```
