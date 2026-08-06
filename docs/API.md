
# Rubium API Reference (Go Backend)

> 🆕 — новый | 🔧 — изменился | ✅ — без изменений

## Общая информация

| Параметр      | Значение                                          |
| --------------------- | --------------------------------------------------------- |
| Базовый URL    | `http://localhost:8080`                                 |
| Прокси (Python) | `http://localhost:5080` → проброс `/api/v1/*` |
| Формат          | JSON                                                      |
| Версия API      | v1                                                        |
| Кодировка    | UTF-8                                                     |

### Ограничения

| Параметр                                  | Значение                                                        |
| ------------------------------------------------- | ----------------------------------------------------------------------- |
| Максимум заданий за запрос | 100 (`limit` capped)                                                  |
| Максимальный размер тела    | 1 MB                                                                    |
| Таймаут подключения             | 10s (read header), 30s (read/write), 120s (idle)                        |
| Retry на Supabase                               | Макс. 2 ретрая, exponential backoff (200ms → 400ms → 800ms) |
| Retry на 5xx                                    | Да                                                                    |
| Retry на 4xx                                    | Нет                                                                  |

### CORS

Разрешённые origins: `localhost:5500`, `localhost:5501`, `localhost:5080`, `localhost:3000` (и 127.0.0.1 аналоги)

Методы: `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `OPTIONS`

---

## Эндпоинты

### GET /health ✅

```json
{ "status": "ok" }
```

---

## Задания (Tasks)

### GET /api/v1/tasks 🔧

**Параметры:**

| Параметр | Тип  | Обязательный | Описание                                                                                                                              |
| ---------------- | ------- | ------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------- |
| `subject`      | string  | Да                     | `math`, `russian`, `physics`, `chemistry`, `biology`, `history`, `social_studies`, `english`, `literature`, `informatics` |
| `exam_type`    | string  | Нет                   | `EGE`, `OGE`, `university`                                                                                                              |
| `level`        | string  | Нет                   | `9`, `10`, `11`, `bachelor_1`                                                                                                         |
| `topic`        | string  | Нет                   | Тема                                                                                                                                      |
| `difficulty`   | integer | Нет                   | Сложность 1-5                                                                                                                        |
| `task_number`  | integer | Нет                   | Номер задания                                                                                                                     |
| `tags`         | string  | Нет                   | Теги через запятую                                                                                                            |
| `limit`        | integer | Нет                   | По умолчанию 1, максимум 100                                                                                               |

**Ответ (200):**

```json
{
  "count": 2,
  "tasks": [
    {
      "id": "02be48bf-1638-4c62-98e9-8c953618e790",
      "content": "Найдите $\\sin \\alpha$, если $\\cos \\alpha = \\frac{3}{5}$",
      "answer": "\\frac{4}{5}",
      "solution": "$\\sin \\alpha = \\sqrt{1 - \\cos^2 \\alpha} = \\frac{4}{5}$",
      "subject": "math",
      "topic": "Тригонометрия",
      "level": "11",
      "exam_type": "EGE",
      "difficulty": 2,
      "task_number": 6,
      "task_type": "fipi",
      "tags": ["тригонометрия", "основное_тождество"],
      "created_at": "2026-07-30T12:08:29.48185+00:00"
    }
  ]
}
```

### GET /api/v1/tasks/:id 🔧

**Ответ (200):**

```json
{ "task": { ... } }
```

**404:** `{ "error": "задание не найдено" }`

### POST /api/v1/check ✅

**Тело:**

```json
{ "task_id": "uuid", "answer": "4/5" }
```

**Ответ:**

```json
{ "correct": true, "correct_answer": "\\frac{4}{5}", "explanation": "..." }
```

---

## Тетради (Notebooks) 🆕

Все требуют `Authorization: Bearer <token>`, кроме GET для публичных и community.

### GET /api/v1/notebooks

**Параметры:** `is_public` (bool, опционально)

**Ответ (200):**

```json
{
  "notebooks": [
    {
      "id": "uuid",
      "title": "Тригонометрия ЕГЭ",
      "color": "#A78BFA",
      "tags": ["math", "EGE", "11"],
      "is_public": false,
      "sections_count": 3,
      "pages_count": 12,
      "views_count": 0,
      "copies_count": 0,
      "created_at": "...",
      "updated_at": "..."
    }
  ]
}
```

### POST /api/v1/notebooks

**Тело:**

```json
{
  "title": "Тригонометрия ЕГЭ",
  "color": "#A78BFA",
  "tags": ["math", "EGE", "11"],
  "is_public": false
}
```

| Поле      | Тип   | Обязательный | По умолчанию |
| ------------- | -------- | ------------------------ | ----------------------- |
| `title`     | string   | Да                     | —                      |
| `color`     | string   | Нет                   | `#A78BFA`             |
| `tags`      | []string | Нет                   | `[]`                  |
| `is_public` | bool     | Нет                   | `false`               |

**Ответ (201):** `{ "id": "uuid", "title": "...", ... }`

### GET /api/v1/notebooks/:id

Приватные — только владельцу, публичные — всем.

**Ответ (200):**

```json
{
  "notebook": {
    "id": "uuid",
    "title": "...",
    "color": "#A78BFA",
    "tags": ["math", "EGE"],
    "is_public": true,
    "sections_count": 3,
    "pages_count": 12,
    "views_count": 42,
    "copies_count": 5,
    "author": { "id": "uuid", "first_name": "Илья" },
    "created_at": "...",
    "updated_at": "..."
  }
}
```

### PUT /api/v1/notebooks/:id

Только владелец. Все поля опциональны.

**Тело:** `{ "title": "...", "color": "...", "tags": [...], "is_public": true }`

**Ответ:** `{ "message": "тетрадь обновлена" }`

### DELETE /api/v1/notebooks/:id

Только владелец. Каскадное удаление разделов и страниц.

**Ответ:** `{ "message": "тетрадь удалена" }`

### POST /api/v1/notebooks/:id/copy

Скопировать публичную тетрадь себе. Увеличивает `copies_count` оригинала.

**Ответ (201):** `{ "id": "new-uuid", "title": "Тригонометрия ЕГЭ (копия)", ... }`

---

## Разделы (Sections) 🆕

### GET /api/v1/notebooks/:id/sections

**Ответ (200):**

```json
{
  "sections": [
    {
      "id": "uuid",
      "notebook_id": "uuid",
      "title": "Основные формулы",
      "order_index": 0,
      "pages_count": 4,
      "created_at": "...",
      "updated_at": "..."
    }
  ]
}
```

### POST /api/v1/notebooks/:id/sections

Только владелец. `order_index` — автоматически в конец.

**Тело:** `{ "title": "Основные формулы" }`

**Ответ (201):** `{ "id": "uuid", "title": "...", "order_index": 3, ... }`

### PUT /api/v1/notebooks/sections/:id

Только владелец.

**Тело:** `{ "title": "Новое название" }`

### DELETE /api/v1/notebooks/sections/:id

Только владелец. Каскадное удаление страниц.

### PUT /api/v1/notebooks/:id/sections/reorder

Только владелец.

**Тело:** `{ "order": ["uuid-1", "uuid-3", "uuid-2"] }`

---

## Страницы (Pages) 🆕

### GET /api/v1/notebooks/sections/:id/pages

**Ответ (200):**

```json
{
  "pages": [
    {
      "id": "uuid",
      "section_id": "uuid",
      "title": "Формулы приведения",
      "content": { "type": "doc", "content": [] },
      "source_task_id": null,
      "order_index": 0,
      "created_at": "...",
      "updated_at": "..."
    }
  ]
}
```

| Поле           | Тип    | Описание                                         |
| ------------------ | --------- | -------------------------------------------------------- |
| `content`        | JSONB     | TipTap JSON (формат ProseMirror)                   |
| `source_task_id` | UUID/null | Привязка к задаче из тренажёра |

### POST /api/v1/notebooks/sections/:id/pages

Только владелец. `order_index` — автоматически в конец.

**Тело:**

```json
{
  "title": "Формулы приведения",
  "content": { "type": "doc", "content": [] },
  "source_task_id": null
}
```

**Ответ (201):** `{ "id": "uuid", "section_id": "uuid", "title": "...", "order_index": 0, ... }`

### PUT /api/v1/notebooks/pages/:id

Только владелец. Все поля опциональны.

**Тело:** `{ "title": "...", "content": {...} }`

### DELETE /api/v1/notebooks/pages/:id

Только владелец.

### PUT /api/v1/notebooks/sections/:id/pages/reorder

Только владелец.

**Тело:** `{ "order": ["uuid-1", "uuid-3", "uuid-2"] }`

---

## Рейтинг (Rating) 🆕

### POST /api/v1/notebooks/:id/rate

Только для публичных тетрадей, не своей. Один голос на пользователя (перезаписывается).

**Тело:** `{ "rating": 4 }` (1-5)

**Ответ (200):** `{ "average_rating": 4.3, "total_ratings": 12 }`

### GET /api/v1/notebooks/:id/rating

**Ответ (200):**

```json
{
  "average_rating": 4.3,
  "total_ratings": 12,
  "user_rating": 4
}
```

`user_rating` = null если пользователь не ставил оценку.

---

## Публичный каталог (Community) 🆕

### GET /api/v1/notebooks/community

Без авторизации.

**Параметры:**

| Параметр | Тип  | По умолчанию | Описание                               |
| ---------------- | ------- | ----------------------- | ---------------------------------------------- |
| `search`       | string  | —                      | Поиск по названию и тегам |
| `tags`         | string  | —                      | Теги через запятую             |
| `sort`         | string  | `rating`              | `rating`, `newest`, `popular`            |
| `page`         | integer | 1                       | Страница                               |
| `limit`        | integer | 20                      | На странице (макс. 50)           |

**Ответ (200):**

```json
{
  "notebooks": [
    {
      "id": "uuid",
      "title": "Тригонометрия для ЕГЭ",
      "color": "#A78BFA",
      "tags": ["math", "EGE", "тригонометрия"],
      "average_rating": 4.5,
      "total_ratings": 28,
      "pages_count": 15,
      "author": { "id": "uuid", "first_name": "Илья" },
      "created_at": "..."
    }
  ],
  "total": 142,
  "page": 1,
  "limit": 20
}
```

---

## Профиль (Users) 🆕

### PUT /api/v1/users/me/pinned-notebook

Закрепить тетрадь в профиле.

**Тело:** `{ "notebook_id": "uuid" }`

**Ответ:** `{ "message": "тетрадь закреплена" }`

### DELETE /api/v1/users/me/pinned-notebook

Открепить тетрадь.

**Ответ:** `{ "message": "тетрадь откреплена" }`

---

## Ошибки

| Код  | Описание                                                                                             |
| ------- | ------------------------------------------------------------------------------------------------------------ |
| `400` | Неверный запрос (отсутствуют поля, невалидный UUID)                   |
| `401` | Не авторизован                                                                                  |
| `403` | Нет доступа (чужая приватная тетрадь, попытка оценить свою) |
| `404` | Не найдено                                                                                          |
| `500` | Внутренняя ошибка                                                                            |

Формат: `{ "error": "описание ошибки" }`

---
