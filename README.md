# 💳 FinStat — Personal Finance REST API & CLI

Легковесный, быстрый и надежный бэкенд-сервис для учета личных финансов, фиксации доходов и расходов, расчета баланса и анализа транзакций. Проект реализован на языке **Go** с использованием **PostgreSQL** и предоставляет два независимых интерфейса: **REST API** и интерактивный **CLI**.

---

## ✨ Особенности

* **Чистая слоистая архитектура (Clean / Layered Architecture):**
  * Полное разделение слоев: Транспорт (`api`, `cli`) $\rightarrow$ Бизнес-логика (`transaction.Service`) $\rightarrow$ Абстракция данных (`transaction.Repository`) $\rightarrow$ Реализация БД (`PostgresRepository` / `MemoryRepository`).
  * Благодаря интерфейсу `Repository` юнит-тесты выполняются мгновенно на in-memory хранилище без необходимости поднимать Postgres.
* **Безопасность и изоляция данных:**
  * Все операции (чтение, создание, редактирование, удаление) строго изолированы по `user_id` (защита от IDOR).
  * Исключительно параметризованные SQL-запросы (`$1, $2, ...`) — полная защита от SQL-инъекций.
* **Производительность:**
  * Cursor-based пагинация (`limit`, `after_id`) для эффективной выборки больших объемов данных.
  * Составные индексы в PostgreSQL под все сценарии фильтрации (`(user_id, id)`, `(user_id, type, id)`, `(user_id, category, id)`).
  * Настроенный Connection Pool к базе данных и явные таймауты HTTP-сервера (`ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, `IdleTimeout`).
* **Две точки входа:**
  * **REST API:** стандартный веб-сервер на базе стандартной библиотеки `net/http`.
  * **CLI:** интерактивная консольная утилита с поддержкой команд и сессий.
* **Встроенные миграции:**
  * SQL-схема автоматически накатывается при старте приложения с помощью механизма `//go:embed`.

---

## 🏗 Архитектура проекта

```
                  ┌──────────────────────┐   ┌──────────────────────┐
                  │    HTTP REST API     │   │      Console CLI     │
                  │       (api/)         │   │        (cli/)        │
                  └──────────┬───────────┘   └──────────┬───────────┘
                             │                          │
                             ▼                          ▼
                  ┌─────────────────────────────────────────────────┐
                  │          Service (Бизнес-логика)                │
                  │             (transaction/service)               │
                  └─────────────────────────┬───────────────────────┘
                                            │
                                            ▼
                  ┌─────────────────────────────────────────────────┐
                  │          Repository (Интерфейс)                 │
                  └─────────────┬─────────────────────┬─────────────┘
                                │                     │
                                ▼                     ▼
                  ┌──────────────────────┐   ┌──────────────────────┐
                  │  PostgresRepository  │   │   MemoryRepository   │
                  │    (PostgreSQL 16)   │   │  (Для быстрых тестов)│
                  └──────────────────────┘   └──────────────────────┘
```

---

## 📁 Структура репозитория

```text
.
├── api/                    # Слой HTTP-транспорта (обработчики, валидация, DTO, тесты)
│   ├── http.go             # Роутер, хэндлеры и middleware
│   └── http_test.go        # HTTP-тесты с использованием httptest
├── cli/                    # Логика интерактивного консольного интерфейса
│   ├── cli.go              # Цикл чтения и обработки потока команд
│   ├── executor.go         # Диспетчеризация команд
│   ├── parser.go           # Парсер пользовательского ввода
│   ├── session.go          # Состояние текущей пользовательской сессии
│   ├── transaction_read.go # Обработка команд чтения
│   ├── transaction_write.go# Обработка команд записи/редактирования
│   └── cli_test.go         # Тесты консольного интерфейса
├── cmd/
│   └── cli/
│       └── main.go         # Точка входа для запуска CLI
├── migrations/             # SQL-миграции схемы БД
│   ├── 001_init.sql        # Инициализация таблиц и индексов
│   └── migrations.go       # Встраивание SQL через //go:embed
├── transaction/            # Доменная модель, интерфейсы и бизнес-логика
│   ├── filter.go           # Структура фильтрации и пагинации
│   ├── memory_repository.go# Реализация in-memory репозитория для тестов
│   ├── model.go            # Доменная модель транзакции
│   ├── postgres_repository.go # Реализация репозитория на PostgreSQL
│   ├── repository.go       # Интерфейс Repository
│   ├── service.go          # Сервисный слой и валидация
│   └── service_test.go     # Тесты бизнес-логики
├── docker-compose.yml      # Спецификация для запуска PostgreSQL в Docker
├── go.mod                  # Определение Go-модуля и зависимостей
├── go.sum                  # Контрольные суммы зависимостей
└── main.go                 # Точка входа REST API сервера
```

---

## 🚀 Быстрый старт

### Требования
* **Go 1.22+**
* **Docker** и **Docker Compose**

### 1. Клонирование репозитория
```bash
git clone https://github.com/<ваш-аккаунт>/fin-stat.git
cd fin-stat
```

### 2. Запуск PostgreSQL
Запустите базу данных в фоновом режиме:
```bash
docker compose up -d
```
База будет доступна по адресу `localhost:5432` со стандартными учетными данными (`finstat` / `finstat`).

### 3. Запуск REST API
```bash
go run .
```
Сервер запустится на порту `8080`:
```text
REST API listening on http://localhost:8080
```

### 4. Запуск CLI (опционально)
В отдельном терминале можно запустить консольный интерфейс:
```bash
go run ./cmd/cli
```

### 5. Запуск тестов
Все тесты изолированы и не требуют запущенной БД:
```bash
go test -v ./...
```

---

## ⚙️ Переменные окружения

| Переменная | По умолчанию | Описание |
|---|---|---|
| `PORT` | `8080` | Порт, на котором слушает HTTP-сервер |
| `DATABASE_URL` | `postgres://finstat:finstat@localhost:5432/finstat?sslmode=disable` | Строка подключения к PostgreSQL |

---

## 📖 Спецификация REST API

> ℹ️ **Концепция авторизации (`X-User-ID`):**  
> Каждый запрос к API требует обязательный заголовок `X-User-ID: <integer>`. В рамках сервиса этот заголовок идентифицирует текущего пользователя. Такой подход является стандартом для **внутренних микросервисов, расположенных за API Gateway** (Gateway проверяет клиентский JWT/OAuth-токен и проксирует запросы во внутренний контур, подставляя проверенный `X-User-ID`).

### Список эндпоинтов

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/balance` | Получение текущего финансового баланса пользователя |
| `GET` | `/transactions` | Получение списка транзакций с поддержкой пагинации и фильтров |
| `POST` | `/transactions` | Создание новой транзакции (доход / расход) |
| `GET` | `/transactions/{id}` | Получение конкретной транзакции по ID |
| `PUT` | `/transactions/{id}` | Обновление существующей транзакции |
| `DELETE` | `/transactions/{id}` | Удаление транзакции |

---

### Примеры запросов c `curl`

#### 1. Создание доходной транзакции
```bash
curl -i -X POST http://localhost:8080/transactions \
  -H "X-User-ID: 1" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "incoming",
    "amount": 75000.00,
    "category": "salary"
  }'
```
**Ответ `201 Created`:**
```json
{
  "id": 1,
  "user_id": 1,
  "type": "incoming",
  "amount": 75000,
  "category": "salary"
}
```

#### 2. Создание расходной транзакции
```bash
curl -i -X POST http://localhost:8080/transactions \
  -H "X-User-ID: 1" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "outgoing",
    "amount": 1500.50,
    "category": "groceries"
  }'
```

#### 3. Проверка текущего баланса
```bash
curl -i -X GET http://localhost:8080/balance \
  -H "X-User-ID: 1"
```
**Ответ `200 OK`:**
```json
{
  "user_id": 1,
  "balance": 73499.5
}
```

#### 4. Получение транзакций с фильтрацией и пагинацией
Параметры запроса (Query Params):
* `type` — фильтр по типу (`incoming` или `outgoing`).
* `category` — фильтр по названию категории.
* `limit` — размер страницы (по умолчанию `100`).
* `after_id` — ID транзакции, после которой запрашивать следующую порцию (курсор).

*Пример (доходы по категории "salary"):*
```bash
curl -i -X GET "http://localhost:8080/transactions?type=incoming&category=salary&limit=10" \
  -H "X-User-ID: 1"
```
**Ответ `200 OK`:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "type": "incoming",
    "amount": 75000,
    "category": "salary"
  }
]
```

#### 5. Обновление транзакции
```bash
curl -i -X PUT http://localhost:8080/transactions/1 \
  -H "X-User-ID: 1" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "incoming",
    "amount": 80000.00,
    "category": "salary"
  }'
```

#### 6. Удаление транзакции
```bash
curl -i -X DELETE http://localhost:8080/transactions/1 \
  -H "X-User-ID: 1"
```
**Ответ `204 No Content`**.

---

## 💻 Работа с интерактивным CLI

Консольный интерфейс поддерживает сессии и следующие команды:

```text
Enter commands (type 'exit' to quit):
Login 1
Logged in as user 1

Enter a command:
Add_Incoming 50000 salary
Added transaction #1

Enter a command:
Add_Outgoing 1200 groceries
Added transaction #2

Enter a command:
Check_Balance
Balance: 48800.00

Enter a command:
All_By_Type incoming
#1  incoming  50000.00  salary

Enter a command:
All_By_Category groceries
#2  outgoing  1200.00   groceries

Enter a command:
Edit_Transaction 2 1500 food
Updated transaction #2

Enter a command:
Delete_Transaction 2
Deleted transaction #2

Enter a command:
exit
```

---

## 🛠 Стек технологий

* **Язык программирования:** Go
* **База данных:** PostgreSQL 16
* **Драйвер БД:** `github.com/lib/pq`
* **Контейнеризация:** Docker, Docker Compose
* **Тестирование:** Стандартный пакет `testing`, `net/http/httptest`
