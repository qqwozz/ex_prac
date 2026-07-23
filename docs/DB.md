
# Exam Trainer — База данных

Актуальная схема Supabase (PostgreSQL). Данные получены напрямую из проекта.

---

## Таблица tasks

Задания для тренажёра (ЕГЭ/ОГЭ по математике и информатике).

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | UUID | Уникальный идентификатор задания |
| `content` | TEXT | Условие задания (текст + LaTeX + ссылки на изображения) |
| `answer` | TEXT | Правильный ответ |
| `solution` | TEXT | Разбор решения (nullable) |
| `subject` | TEXT | Предмет: `math`, `informatics` |
| `exam_type` | TEXT | Тип экзамена: `ege`, `oge` |
| `level` | TEXT | Уровень сложности: `base`, `medium`, `hard` |
| `topic` | TEXT | Тема/модуль: `Производная`, `Программирование` и т.д. |
| `task_type` | TEXT | Тип задания: `fipi` (из банка ФИПИ), `ai` (сгенерировано) |
| `target_type` | TEXT | Для кого: `all` (все), `plus` (только подписчики) |
| `target_id` | UUID | Ссылка на конкретного пользователя (nullable) |
| `created_by` | UUID | Кто создал задание (ссылка на users.id) |
| `created_at` | TIMESTAMPTZ | Дата создания |
| `task_number` | INTEGER | Номер задания в экзамене (8, 16 и т.д.) |
| `source` | TEXT | Источник: `Открытый банк ФИПИ`, `ai`, `teacher` |
| `display_id` | TEXT | Читаемый ID для отображения: `#000001`, `#000002` |

### Пример данных

```json
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
```

### Запросы через PostgREST

| Что нужно | Запрос |
|-----------|--------|
| Все задания | `GET /rest/v1/tasks?select=*` |
| По предмету | `?subject=eq.math` |
| По теме | `?topic=eq.Программирование` |
| По экзамену | `?exam_type=eq.oge` |
| По уровню | `?level=eq.medium` |
| С лимитом | `?limit=5` |
| Комбинированный | `?subject=eq.informatics&exam_type=eq.oge&limit=3` |

---

## Таблица users

Пользователи системы (ученики и преподаватели).

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | UUID | Уникальный идентификатор (совпадает с auth.users.id) |
| `email` | TEXT | Email пользователя |
| `role` | TEXT | Роль: `admin`, `teacher`, `student` |
| `auth_provider` | TEXT | Провайдер авторизации: `email`, `google` и т.д. |
| `created_at` | TIMESTAMPTZ | Дата регистрации |
| `first_name` | TEXT | Имя |
| `last_name` | TEXT | Фамилия |
| `avatar_url` | TEXT | Ссылка на аватар (в Supabase Storage) |
| `photo_url` | TEXT | Путь к фото |
| `description` | TEXT | Описание / биография |
| `subjects` | TEXT | Предметы (для преподавателей) |
| `university` | TEXT | Университет |
| `experience` | TEXT | Опыт преподавания |

### Пример данных

```json
{
  "id": "d41bc28c-df35-47c4-a508-cd46c394e980",
  "email": "chabata33@yandex.ru",
  "role": "admin",
  "auth_provider": "email",
  "created_at": "2026-06-30T09:18:53.280736+00:00",
  "first_name": "Максим",
  "last_name": "Ковалёв",
  "avatar_url": "https://...supabase.co/storage/v1/object/public/avatars/...",
  "description": "Студент МГТУ им. Баумана...",
  "subjects": "Математика, Физика (ЕГЭ и ОГЭ)",
  "university": "МГТУ им. Баумана",
  "experience": "2 года преподавания"
}
```

---

## Таблица attempts (planned)

Попытки решения заданий учениками. **Пока не создана.**

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | UUID | Уникальный идентификатор попытки |
| `user_id` | UUID | Ссылка на пользователя |
| `task_id` | UUID | Ссылка на задание (tasks.id) |
| `answer` | TEXT | Ответ ученика |
| `correct` | BOOLEAN | Правильный ли ответ |
| `created_at` | TIMESTAMPTZ | Время попытки |

---

## Связи между таблицами

```
users (1) ──→ (N) tasks        — created_by
users (1) ──→ (N) attempts     — user_id
tasks (1) ──→ (N) attempts     — task_id
```

---

## Supabase Storage

| Бакет | Назначение |
|-------|------------|
| `avatars` | Аватары пользователей |
| `tasks-images` | Изображения к заданиям |

---

## RLS (Row Level Security)

| Таблица | Политика |
|---------|----------|
| `tasks` | Чтение — всем (anon), запись — только admin/teacher |
| `users` | Чтение — свои данные, запись — свои данные |
| `attempts` | Чтение — свои попытки, запись — все авторизованные |

---

## Статистика по данным

| Таблица | Записей | Последнее обновление |
|---------|---------|---------------------|
| `tasks` | 2 | 2026-07-11 |
| `users` | 1 | 2026-06-30 |
| `attempts` | — | таблица не создана |
