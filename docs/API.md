
# Exam Trainer — API Reference

Все эндпоинты, форматы запросов/ответов и ошибки. Для фронтенда и интеграций.

---

## Общая информация

| Параметр | Значение |
|----------|----------|
| Базовый URL (Go) | `http://localhost:8080` |
| Базовый URL (Python) | `http://localhost:5080` |
| Формат | JSON |
| Версия API | v1 |

Python проксирует запросы к Go, поэтому фронтенд обращается к Python (5080).

---

## Go → Supabase

| Запрос | Метод | Путь Supabase |
|--------|-------|---------------|
| Задания | `GET` | `/rest/v1/tasks?subject=eq.{s}&topic=eq.{t}&difficulty=lte.{d}&limit={n}` |
| Проверка | `POST` | `/rest/v1/attempts` |
| Статистика | `GET` | `/rest/v1/attempts?user_id=eq.{id}&select=*` |
| Пользователь | `GET` | `/rest/v1/users?id=eq.{id}` |

---

## Python → Go

| Запрос | Метод | Путь Go |
|--------|-------|---------|
| Задания | `GET` | `http://localhost:8080/api/v1/tasks?...` |
| Проверка | `POST` | `http://localhost:8080/api/v1/check` |
| Статистика | `GET` | `http://localhost:8080/api/v1/stats?...` |
| Пользователь | `GET` | `http://localhost:8080/api/v1/user/{id}` |

---

## Python → AI

| Запрос | Провайдер | Модель |
|--------|-----------|--------|
| Подсказка | DeepSeek | `deepseek-chat` |
| Подсказка (фолбэк) | OpenAI | `gpt-4o-mini` |
| Генерация заданий | DeepSeek | `deepseek-chat` |

---

## Эндпоинты

### GET /api/v1/tasks

Получить задание (или список заданий).

**Параметры (query):**

| Параметр | Тип | Обязательный | Описание |
|----------|-----|--------------|----------|
| `subject` | string | Да | Предмет: `math`, `informatics` |
| `exam` | string | Нет | Тип экзамена: `ege`, `oge` |
| `type` | string | Нет | Тип задания: `choice`, `number`, `string`, `multi` |
| `topic` | string | Нет | Модуль: `geometry`, `equations`, `probability` |
| `difficulty` | integer | Нет | Максимальная сложность (1–5) |
| `task_number` | integer | Нет | Номер задания в экзамене (1–27 для ОГЭ, 1–19 для ЕГЭ) |
| `mode` | string | Нет | Режим: `improve` (улучшить) / `weak` (подтянуть) |
| `limit` | integer | Нет | Количество заданий (по умолчанию 1) |

Задания выдаются в **случайном порядке**. Если фильтры не указаны — вернётся любое случайное задание из всей базы.

**Ответ (200):**

```json
{
  "id": "uuid",
  "subject": "math",
  "exam": "ege",
  "type": "choice",
  "number": 1,
  "difficulty": 3,
  "topic": "geometry",
  "text": "Вычислите: 2³ + 3²",
  "options": ["11", "15", "17", "22"],
  "image_url": null
}
```

---

### POST /api/v1/check

Проверить ответ ученика.

**Тело запроса:**

```json
{
  "task_id": "uuid",
  "answer": "17"
}
```

**Ответ (200):**

```json
{
  "correct": true,
  "correct_answer": "17",
  "explanation": "2³ = 8, 3² = 9, 8 + 9 = 17"
}
```

**Форматы проверки:**

| Тип задания | Кто проверяет | Как проверяет |
|-------------|---------------|---------------|
| `choice` | Go | Точное сравнение строки |
| `number` | Go | Сравнение с допуском (±0.01) |
| `string` | Go | Нормализация (нижний регистр, пробелы) + сравнение |
| `multi` | Go | Совпадение множества |
| `code` | Python | Запуск кода + сравнение вывода |
| `text` | Python | AI-анализ или эталонное сравнение |

**Как работает пересылка в Python:**

Если тип задания `code` или `text`, Go отправляет POST на `PYTHON_URL/ai/v1/check` с телом:

```json
{
  "task_id": "uuid",
  "task_type": "code",
  "content": "условие задания",
  "answer": "эталонный вывод",
  "user_answer": "код ученика"
}
```

---

### GET /api/v1/stats

Получить статистику пользователя.

**Параметры (query):**

| Параметр | Тип | Обязательный | Описание |
|----------|-----|--------------|----------|
| `user_id` | string | Да | UUID пользователя |

**Ответ (200):**

```json
{
  "total": 150,
  "correct": 112,
  "accuracy": 74.6,
  "by_topic": {
    "geometry": { "total": 50, "correct": 42 },
    "equations": { "total": 40, "correct": 30 }
  }
}
```

---

### GET /api/v1/user/{id}

Получить данные пользователя (подписка, лимиты).

**Параметры (path):**

| Параметр | Тип | Описание |
|----------|-----|----------|
| `id` | string | UUID пользователя |

**Ответ (200):**

```json
{
  "id": "uuid",
  "subscription_type": "plus",
  "hints_used_today": 3,
  "hints_limit": 10,
  "referred_by": "uuid"
}
```

---

### POST /api/v1/tasks

Добавить задание (для преподавателей/ботов).

**Тело запроса:**

```json
{
  "subject": "math",
  "exam": "ege",
  "type": "choice",
  "number": 1,
  "difficulty": 3,
  "topic": "geometry",
  "text": "Найдите площадь треугольника...",
  "options": ["12", "24", "36", "48"],
  "answer": "24",
  "explanation": "S = ½ · a · b",
  "source": "teacher"
}
```

**Ответ (201):**

```json
{
  "id": "uuid",
  "status": "created"
}
```

---

## AI-эндпоинты

### POST /ai/v1/hint

Получить AI-подсказку к заданию.

**Тело запроса:**

```json
{
  "task_id": "uuid",
  "user_answer": "15"
}
```

**Ответ (200):**

```json
{
  "hint": "Подумайте о том, как возводить числа в степень. Что такое 2³?",
  "provider": "deepseek"
}
```

**Ограничения:**

- Лимит: 10 подсказок в день
- Тип: направляющая мысль, НЕ ответ
- Счётчик хранится в `users.hints_used_today`
- Фолбэк: DeepSeek → OpenAI

---

### POST /ai/v1/generate

Сгенерировать задание с помощью AI.

**Тело запроса:**

```json
{
  "subject": "math",
  "topic": "geometry",
  "type": "choice",
  "difficulty": 3
}
```

**Ответ (200):**

```json
{
  "text": "Найдите площадь треугольника со сторонами 6 и 8 и углом 90°",
  "options": ["12", "24", "36", "48"],
  "answer": "24",
  "explanation": "S = ½ · a · b = ½ · 6 · 8 = 24"
}
```

---

## Промты

### Подсказка (hint_system.md)

```
Ты — AI-ассистент для подготовки к ЕГЭ/ОГЭ.

Ученик решил задание неправильно. Твоя задача — дать направляющую подсказку,
НО НЕ ОТВЕТ.

Задание: {task_text}
Правильный ответ: {correct_answer}
Ответ ученика: {user_answer}

Дай одну короткую подсказку (1-2 предложения), которая направит ученика
к правильному ответу, но не раскроет его.
```

### Генерация заданий (generate_system.md)

```
Ты — генератор заданий для тренажёра ЕГЭ/ОГЭ.

Сгенерируй задание по предмету: {subject}
Тема: {topic}
Тип: {type}
Сложность: {difficulty}/5

Формат ответа (JSON):
{
  "text": "условие задания",
  "options": ["вариант1", "вариант2", "вариант3", "вариант4"],
  "answer": "правильный ответ",
  "explanation": "разбор"
}
```

---

## Ошибки

| Код | Описание | Пример |
|-----|----------|--------|
| `400` | Неверный запрос | Отсутствует обязательный параметр |
| `404` | Задание не найдено | Неверный `task_id` |
| `429` | Превышен лимит подсказок | `hints_used_today >= 10` |
| `500` | Внутренняя ошибка сервера | Ошибка Supabase / AI |

---

## Примеры запросов

### Получить задание по математике

```bash
curl "http://localhost:5080/api/v1/tasks?subject=math&topic=geometry&limit=1"
```

### Получить случайное задание номер 16 из ОГЭ по информатике

```bash
curl "http://localhost:5080/api/v1/tasks?subject=informatics&exam=oge&task_number=16&limit=1"
```

### Проверить ответ

```bash
curl -X POST "http://localhost:5080/api/v1/check" \
  -H "Content-Type: application/json" \
  -d '{"task_id": "uuid", "answer": "24"}'
```

### Получить подсказку

```bash
curl -X POST "http://localhost:5080/ai/v1/hint" \
  -H "Content-Type: application/json" \
  -d '{"task_id": "uuid", "user_answer": "15"}'
```
