# Exam Trainer Backend

REST API на Go для тренировки заданий ЕГЭ и ОГЭ.

## Суть

Бэкенд хранит базу заданий из ЕГЭ и ОГЭ по разным предметам. Пользователь берёт задание, отвечает, бэкенд проверяет ответ и возвращает результат.

## Как работает

```
Клиент                    Backend (Go)                SQLite
  │                           │                          │
  │  GET /tasks?subject=math  │                          │
  │ ─────────────────────────>│  SELECT * FROM tasks     │
  │                           │ ────────────────────────>│
  │  { task, options }        │                          │
  │ <─────────────────────────│                          │
  │                           │                          │
  │  POST /check              │                          │
  │  { task_id, answer }      │                          │
  │ ─────────────────────────>│  SELECT answer FROM tasks│
  │                           │ ────────────────────────>│
  │  { correct: true/false }  │                          │
  │ <─────────────────────────│                          │
```

## Форматы заданий

Задания бывают разных типов:

| Тип | Пример | Проверка |
|-----|--------|----------|
| Выбор ответа | 4 варианта, один правильный | Точное сравнение |
| Ввод числа | Ответ — число (42, 3.14) | Сравнение с допуском |
| Ввод строки | Ответ — слово или выражение | Нормализация + сравнение |
| Множественный выбор | Несколько правильных из N | Совпадение множества |

## API

### Получить задание

```
GET /api/v1/tasks?subject=math&exam=ege&type=choice&limit=1
```

Ответ:
```json
{
  "id": "task_001",
  "subject": "math",
  "exam": "ege",
  "type": "choice",
  "number": 1,
  "text": "Вычислите: 2³ + 3²",
  "options": ["11", "15", "17", "22"],
  "image_url": null
}
```

### Проверить ответ

```
POST /api/v1/check
```

```json
{
  "task_id": "task_001",
  "answer": "17"
}
```

Ответ:
```json
{
  "correct": true,
  "correct_answer": "17",
  "explanation": "2³ = 8, 3² = 9, 8 + 9 = 17"
}
```

### Получить статистику

```
GET /api/v1/stats?user_id=123
```

```json
{
  "total": 150,
  "correct": 112,
  "accuracy": 74.6,
  "by_subject": {
    "math": { "total": 50, "correct": 42 },
    "russian": { "total": 40, "correct": 30 }
  }
}
```

## Структура

```
.
├── cmd/
│   └── server/
│       └── main.go           # точка входа
├── internal/
│   ├── api/
│   │   ├── router.go         # маршруты
│   │   ├── handlers.go       # обработчики запросов
│   │   └── middleware.go     # логирование, CORS
│   ├── models/
│   │   ├── task.go           # Task, Answer
│   │   └── stats.go          # Stats
│   ├── storage/
│   │   ├── db.go             # подключение к БД
│   │   └── tasks.go          # запросы к таблице tasks
│   └── checker/
│       └── checker.go        # логика проверки ответов
├── migrations/
│   └── 001_init.sql          # создание таблиц
├── seeds/
│   └── math_ege.json         # тестовые данные
├── go.mod
└── README.md
```

## База данных

SQLite — просто, ноль настроек, файл в одной папке.

```sql
CREATE TABLE tasks (
    id          TEXT PRIMARY KEY,
    subject     TEXT NOT NULL,        -- math, russian, physics, ...
    exam        TEXT NOT NULL,        -- ege, oge
    type        TEXT NOT NULL,        -- choice, number, string, multi
    number      INTEGER,             -- номер задания в экзамене
    text        TEXT NOT NULL,        -- условие
    options     TEXT,                 -- JSON с вариантами (для choice)
    answer      TEXT NOT NULL,        -- правильный ответ
    explanation TEXT,                 -- разбор
    image_url   TEXT,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE attempts (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     TEXT NOT NULL,
    task_id     TEXT NOT NULL,
    answer      TEXT NOT NULL,
    correct     BOOLEAN NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks(id)
);

CREATE INDEX idx_tasks_subject ON tasks(subject, exam);
CREATE INDEX idx_attempts_user ON attempts(user_id, created_at);
```

## Стек

| Компонент | Решение |
|-----------|---------|
| Язык | Go 1.22+ |
| HTTP | `net/http` + `chi` или `gin` |
| БД | SQLite через `modernc.org/sqlite` (чистый Go, CGO не нужен) |
| Миграции | `golang-migrate` |
| Тесты | `testing` + `httptest` |

## Запуск

```bash
# сборка
go build -o server ./cmd/server

# запуск
./server -port 8080 -db exam.db
```

```bash
# тесты
go test ./...
```

## Что дальше

- Авторизация (JWT)
- Режим экзамена (таймер, количество заданий)
- Импорт заданий из PDF/Excel
- Рейтинг пользователей
- Повторение ошибок (adaptive)
