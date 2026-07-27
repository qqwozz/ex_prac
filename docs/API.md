# Exam Trainer — API Reference (Go Backend)

Актуальная документация Go API. Все эндпоинты, форматы запросов/ответов, ошибки и примеры.

---

## Общая информация

| Параметр | Значение |
|----------|----------|
| Базовый URL | `http://localhost:8080` |
| Формат | JSON |
| Версия API | v1 |
| Кодировка | UTF-8 |

### Ограничения

| Параметр | Значение |
|----------|----------|
| Максимум заданий за запрос | 100 (`limit` capped) |
| Максимальный размер тела | 1 MB |
| Таймаут подключения | 10s (read header), 30s (read/write), 120s (idle) |
| Retry на Supabase | Макс. 2 ретрая, exponential backoff (200ms → 400ms → 800ms) |
| Retry на 5xx | Да (сетевые ошибки + серверные ошибки Supabase) |
| Retry на 4xx | Нет |

### CORS

Сервер разрешает запросы с origins:

- `http://localhost:5500`
- `http://localhost:5080`
- `http://localhost:5081`
- `http://localhost:3000`
- `http://127.0.0.1:5500`
- `http://127.0.0.1:5080`

Разрешённые методы: `GET`, `POST`, `PUT`, `OPTIONS`

---

## Эндпоинты

### GET /health

Health-check эндпоинт. Возвращает статус сервера.

**Ответ (200):**

```json
{
  "status": "ok"
}
```

---

### GET /api/v1/tasks

Получить задание (или список заданий). Задания выдаются в **случайном порядке**.

**Параметры (query):**

| Параметр | Тип | Обязательный | Описание |
|----------|-----|--------------|----------|
| `subject` | string | Да | Предмет: `math`, `informatics` |
| `exam` | string | Нет | Тип экзамена: `ege`, `oge` |
| `type` | string | Нет | Тип задания: `choice`, `number`, `string`, `multi`, `code`, `text` |
| `topic` | string | Нет | Тема/модуль: `Производная`, `Программирование` и т.д. |
| `difficulty` | integer | Нет | Максимальная сложность (1-5). Фильтр: `<=` указанному значению |
| `task_number` | integer | Нет | Номер задания в экзамене (8, 16 и т.д.) |
| `limit` | integer | Нет | Количество заданий (по умолчанию 1, максимум 100) |

**Ответ (200):**

```json
{
  "tasks": [
    {
      "id": "83b25796-49c5-4159-ae46-9e1196719288",
      "content": "Напишите программу, которая...",
      "answer": "-",
      "solution": null,
      "subject": "informatics",
      "exam_type": "oge",
      "level": "medium",
      "topic": "Программирование",
      "task_type": "fipi",
      "target_type": "all",
      "target_id": null,
      "created_by": "0f339b36-c20e-47d0-97f4-51ae9837333b",
      "created_at": "2026-07-11T20:48:49.992368+00:00",
      "task_number": 16,
      "source": "Открытый банк ФИПИ",
      "display_id": "#000002"
    }
  ],
  "count": 1
}
```

**Поля задачи (tasks):**

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | UUID | Уникальный идентификатор |
| `content` | string | Условие задания (текст + LaTeX + ссылки на изображения) |
| `answer` | string | Правильный ответ (для code/text — `"-"`) |
| `solution` | string/null | Разбор решения |
| `subject` | string | Предмет: `math`, `informatics` |
| `exam_type` | string | Тип экзамена: `ege`, `oge` |
| `level` | string | Уровень сложности: `base`, `medium`, `hard` |
| `topic` | string | Тема/модуль |
| `task_type` | string | Тип: `fipi` (банк ФИПИ), `ai` (сгенерировано) |
| `target_type` | string | Для кого: `all`, `plus` |
| `target_id` | UUID/null | Ссылка на конкретного пользователя |
| `created_by` | UUID | Кто создал |
| `created_at` | timestamp | Дата создания |
| `task_number` | integer | Номер задания в экзамене |
| `source` | string | Источник: `Открытый банк ФИПИ`, `ai`, `teacher` |
| `display_id` | string | Читаемый ID: `#000001`, `#000002` |

**Примеры:**

```bash
# Случайное задание по математике
curl "http://localhost:8080/api/v1/tasks?subject=math"

# 5 заданий по информатике, ОГЭ, номер 16
curl "http://localhost:8080/api/v1/tasks?subject=informatics&exam=oge&task_number=16&limit=5"

# Задания по теме "Программирование" с сложностью <= 3
curl "http://localhost:8080/api/v1/tasks?subject=informatics&topic=Программирование&difficulty=3&limit=10"
```

---

### GET /api/v1/tasks/:id

Получить одно задание по UUID.

**Параметры (path):**

| Параметр | Тип | Описание |
|----------|-----|----------|
| `id` | UUID | Уникальный идентификатор задания |

**Ответ (200):**

```json
{
  "task": {
    "id": "83b25796-49c5-4159-ae46-9e1196719288",
    "content": "Напишите программу...",
    "answer": "-",
    "subject": "informatics",
    "exam_type": "oge",
    "level": "medium",
    "topic": "Программирование",
    "task_type": "fipi",
    "task_number": 16,
    "source": "Открытый банк ФИПИ",
    "display_id": "#000002"
  }
}
```

**Пример:**

```bash
curl "http://localhost:8080/api/v1/tasks/83b25796-49c5-4159-ae46-9e1196719288"
```

---

### PUT /api/v1/tasks/:id

Обновить задание по UUID. Обновляются только переданные поля (частичное обновление).

**Параметры (path):**

| Параметр | Тип | Описание |
|----------|-----|----------|
| `id` | UUID | Уникальный идентификатор задания |

**Тело запроса (все поля опциональны):**

```json
{
  "content": "Новое условие задания",
  "answer": "42",
  "solution": "Разбор решения",
  "subject": "math",
  "exam_type": "ege",
  "level": "hard",
  "topic": "Алгебра",
  "task_type": "ai",
  "task_number": 5,
  "source": "ai"
}
```

| Поле | Тип | Описание |
|------|-----|----------|
| `content` | string | Условие задания |
| `answer` | string | Правильный ответ |
| `solution` | string | Разбор решения |
| `subject` | string | Предмет: `math`, `informatics` |
| `exam_type` | string | Тип экзамена: `ege`, `oge` |
| `level` | string | Уровень: `base`, `medium`, `hard` |
| `topic` | string | Тема/модуль |
| `task_type` | string | Тип: `fipi`, `ai` |
| `task_number` | integer | Номер задания |
| `source` | string | Источник |

**Ответ (200):**

```json
{
  "message": "задание обновлено"
}
```

**Пример:**

```bash
curl -X PUT "http://localhost:8080/api/v1/tasks/83b25796-49c5-4159-ae46-9e1196719288" \
  -H "Content-Type: application/json" \
  -d '{"answer": "42", "level": "hard"}'
```

---

### DELETE /api/v1/tasks/:id

Удалить задание по UUID.

**Параметры (path):**

| Параметр | Тип | Описание |
|----------|-----|----------|
| `id` | UUID | Уникальный идентификатор задания |

**Ответ (200):**

```json
{
  "message": "задание удалено"
}
```

**Пример:**

```bash
curl -X DELETE "http://localhost:8080/api/v1/tasks/83b25796-49c5-4159-ae46-9e1196719288"
```

---

### POST /api/v1/check

Проверить ответ ученика. Сервер сам определяет тип задания и выбирает метод проверки.

**Тело запроса:**

```json
{
  "task_id": "83b25796-49c5-4159-ae46-9e1196719288",
  "answer": "4"
}
```

| Поле | Тип | Обязательный | Описание |
|------|-----|--------------|----------|
| `task_id` | string | Да | UUID задания |
| `answer` | string | Да | Ответ ученика |

**Ответ (200):**

```json
{
  "correct": true,
  "correct_answer": "4",
  "explanation": "Производная x^2 = 2x, при x=2: 2*2 = 4"
}
```

| Поле | Тип | Описание |
|------|-----|----------|
| `correct` | bool | Правильный ли ответ |
| `correct_answer` | string | Правильный ответ |
| `explanation` | string/null | Разбор решения (если есть) |

**Как работает проверка:**

| Тип задания | Метод проверки |
|-------------|---------------|
| `choice` | Сравнение строк (регистр-независимое, обрезка пробелов) |
| `number` | Сравнение чисел с допуском ±0.01 (запятая = точка) |
| `string` | Нормализация (нижний регистр, схлопывание пробелов) + сравнение |
| `multi` | Сравнение множеств (порядок неважен, разделители: `,`, `;`, `|`) |
| `code` | Пересылка в Python для запуска кода |
| `text` | Пересылка в Python для AI-анализа |

**Определение типа:**
- Если `task_type` задан в БД — используется он
- Если `task_type` пуст и `answer` равен `"-"` или пуст — тип = `code`
- Иначе — тип = `choice`

**Примеры:**

```bash
# Проверка математического задания
curl -X POST "http://localhost:8080/api/v1/check" \
  -H "Content-Type: application/json" \
  -d '{"task_id": "83b25796-49c5-4159-ae46-9e1196719288", "answer": "4"}'

# Проверка ответа с запятой (нормализуется)
curl -X POST "http://localhost:8080/api/v1/check" \
  -H "Content-Type: application/json" \
  -d '{"task_id": "...", "answer": "3,14"}'
```

---

## Ошибки

| Код | Описание | Пример |
|-----|----------|--------|
| `400` | Неверный запрос | Отсутствует `task_id` или `answer`, невалидный UUID, пустой payload |
| `404` | Задание не найдено | Несуществующий `task_id` |
| `500` | Внутренняя ошибка | Ошибка Supabase (сеть, таймаут), ошибка Python |

**Формат ошибки:**

```json
{
  "error": "нужны task_id и answer"
}
```

```json
{
  "error": "задание не найдено"
}
```

```json
{
  "error": "ошибка проверки через Python: Python вернул 500: ..."
}
```

---

## Типы заданий

### choice (выбор ответа)

Задание с вариантами ответа. Проверяется точное совпадение строк с учётом регистра.

```json
{
  "content": "Чему равна производная f(x) = x^2?",
  "answer": "2x",
  "task_type": "choice"
}
```

### number (числовой ответ)

Числовой ответ с допуском ±0.01. Запятая автоматически заменяется на точку.

```json
{
  "content": "Вычислите интеграл от 0 до 1: x^2 dx",
  "answer": "0.33",
  "task_type": "number"
}
```

### string (текстовый ответ)

Строковый ответ. Нормализация: нижний регистр, схлопывание пробелов.

```json
{
  "content": "Какой язык программирования является интерпретируемым?",
  "answer": "Python",
  "task_type": "string"
}
```

### multi (множественный выбор)

Несколько правильных ответов. Порядок неважен. Разделители: `,`, `;`, `|`.

```json
{
  "content": "Какие из перечисленных чисел являются простыми?",
  "answer": "2,3,5,7",
  "task_type": "multi"
}
```

### code (программирование)

Задание на программирование. Ответ — код ученика. Проверяется через Python.

```json
{
  "content": "Напишите программу, которая выводит сумму чисел...",
  "answer": "-",
  "task_type": "fipi"
}
```

### text (текстовое задание)

Развёрнутый текстовый ответ. Проверяется через Python (AI-анализ).

```json
{
  "content": "Объясните, что такое полиморфизм в ООП",
  "answer": "-",
  "task_type": "text"
}
```

---

## Интеграция с Python

Для заданий типа `code` и `text` Go-сервер пересылает запрос в Python.

**POST** `http://localhost:5080/ai/v1/check`

```json
{
  "task_id": "uuid",
  "task_type": "code",
  "content": "условие задания",
  "answer": "эталонный вывод",
  "user_answer": "код ученика"
}
```

**Ответ Python:**

```json
{
  "correct": true
}
```

**Ограничения:**
- Таймаут: 15 секунд
- Максимальный размер ответа: 1 МБ

---

## Конфигурация

### Переменные окружения (.env)

| Переменная | Описание | Обязательная |
|-----------|----------|--------------|
| `PORT` | Порт Go-сервера (по умолчанию 8080) | Нет |
| `SUPABASE_URL` | URL проекта Supabase | Да |
| `SUPABASE_ANON_KEY` | Анонимный ключ Supabase | Да |
| `SUPABASE_SERVICE_KEY` | Сервисный ключ Supabase (нужен для записи) | Да |
| `PYTHON_URL` | URL Python-сервера (по умолчанию http://localhost:5080) | Нет |

### config.yaml

```yaml
supabase:
  url: "https://your-project.supabase.co"
  anon_key: "${SUPABASE_ANON_KEY}"
  service_key: "${SUPABASE_SERVICE_KEY}"

server:
  go_port: 8080
```

---

## Запуск

```bash
cd api
go run cmd/server/main.go
```

Сервер автоматически запускает проверки при старте:
- Конфигурация (env vars заданы, ключи различаются)
- Подключение к БД (anon + service key)
- Checker (все типы заданий)
- HTTP-эндпоинты

Если хотя бы одна проверка не пройдена — сервер не запустится.

---

## Схема данных (Supabase)

### Таблица tasks

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | UUID | Уникальный идентификатор |
| `content` | TEXT | Условие задания |
| `answer` | TEXT | Правильный ответ |
| `solution` | TEXT | Разбор (nullable) |
| `subject` | TEXT | Предмет: `math`, `informatics` |
| `exam_type` | TEXT | Тип экзамена: `ege`, `oge` |
| `level` | TEXT | Уровень: `base`, `medium`, `hard` |
| `topic` | TEXT | Тема/модуль |
| `task_type` | TEXT | Тип: `fipi`, `ai` |
| `target_type` | TEXT | Для кого: `all`, `plus` |
| `target_id` | UUID | Ссылка на пользователя (nullable) |
| `created_by` | UUID | Кто создал |
| `created_at` | TIMESTAMPTZ | Дата создания |
| `task_number` | INTEGER | Номер задания |
| `source` | TEXT | Источник |
| `display_id` | TEXT | Читаемый ID |
