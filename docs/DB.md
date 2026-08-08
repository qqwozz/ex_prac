# Rubium — База данных (Supabase)

Актуальная схема на 08.08.2026.

---

## Таблица tasks

Задания для тренажёра и ежедневных.

| Поле          | Тип      | Nullable | Описание                                          |
| ----------------- | ----------- | -------- | --------------------------------------------------------- |
| `id`            | UUID        | NO       | PK                                                        |
| `title`         | TEXT        | YES      | Заголовок задания                         |
| `content`       | TEXT        | NO       | Условие (текст + LaTeX)                       |
| `answer`        | TEXT        | NO       | Правильный ответ (LaTeX или текст) |
| `solution`      | TEXT        | YES      | Разбор решения                               |
| `subject`       | TEXT        | NO       | Предмет                                            |
| `topic`         | TEXT        | NO       | Тема                                                  |
| `level`         | TEXT        | NO       | Класс/курс:`9`, `10`, `11`, `bachelor_1` |
| `exam_type`     | TEXT        | YES      | `EGE`, `OGE`, `university`                          |
| `difficulty`    | INTEGER     | YES      | Сложность 1-5                                    |
| `task_number`   | INTEGER     | YES      | Номер задания в экзамене             |
| `code_template` | TEXT        | YES      | Шаблон кода                                     |
| `tags`          | TEXT[]      | YES      | Массив тегов                                   |
| `task_type`     | TEXT        | YES      | `fipi`, `ai`                                          |
| `created_at`    | TIMESTAMPTZ | YES      | Дата создания                                 |

---

## Таблица rubium_users

Пользователи платформы.

| Поле                  | Тип      | Nullable | Описание                                         |
| ------------------------- | ----------- | -------- | -------------------------------------------------------- |
| `id`                    | UUID        | NO       | PK                                                       |
| `auth_id`               | UUID        | NO       | FK → auth.users                                         |
| `email`                 | TEXT        | NO       | Почта                                               |
| `first_name`            | TEXT        | YES      | Имя                                                   |
| `last_name`             | TEXT        | YES      | Фамилия                                           |
| `avatar_url`            | TEXT        | YES      | Аватар                                             |
| `xp`                    | INTEGER     | YES      | Опыт                                                 |
| `level`                 | INTEGER     | YES      | Уровень пользователя                  |
| `streak`                | INTEGER     | YES      | Дней подряд                                    |
| `last_streak_date`      | DATE        | YES      | Последняя активность                  |
| `tasks_solved`          | INTEGER     | YES      | Решено задач                                  |
| `tasks_correct`         | INTEGER     | YES      | Правильно решено                          |
| `accuracy`              | NUMERIC     | YES      | Процент правильных                      |
| `status`                | TEXT        | YES      | `user` / `admin`                                     |
| `created_at`            | TIMESTAMPTZ | YES      | Дата регистрации                          |
| `updated_at`            | TIMESTAMPTZ | YES      | Последнее обновление                  |
| `pinned_notebook_id` 🆕 | UUID        | YES      | FK → notebooks, закреплённая тетрадь |

---

## Таблица rubium_daily_stats

Статистика ежедневных заданий.

---

## Таблица notebooks 🆕

Тетради пользователей. Разделы и страницы хранятся в `content` (JSONB). Рейтинг — в `ratings` (JSONB).

| Поле           | Тип      | Nullable | Default             | Описание                                                                   |
| ------------------ | ----------- | -------- | ------------------- | ---------------------------------------------------------------------------------- |
| `id`             | UUID        | NO       | gen_random_uuid()   | PK                                                                                 |
| `user_id`        | UUID        | NO       | —                  | FK → auth.users                                                                   |
| `title`          | TEXT        | NO       | —                  | Название                                                                   |
| `color`          | TEXT        | YES      | `#A78BFA`         | HEX-цвет                                                                       |
| `tags`           | TEXT[]      | YES      | `[]`              | Теги                                                                           |
| `is_public`      | BOOLEAN     | YES      | `false`           | Публичная                                                                 |
| `is_verified`    | BOOLEAN     | YES      | `false`           | Тетрадь разработчика (ставит админ-владелец) |
| `content`        | JSONB       | YES      | `{"sections":[]}` | Разделы и страницы                                                 |
| `ratings`        | JSONB       | YES      | `{}`              | `{"user_id": rating}` — оценки                                            |
| `average_rating` | NUMERIC     | YES      | `0`               | Средняя оценка                                                        |
| `views_count`    | INTEGER     | YES      | `0`               | Просмотры                                                                 |
| `copies_count`   | INTEGER     | YES      | `0`               | Копирования                                                             |
| `created_at`     | TIMESTAMPTZ | YES      | now()               |                                                                                    |
| `updated_at`     | TIMESTAMPTZ | YES      | now()               |                                                                                    |

**Структура `content`:**

```json
{
  "sections": [
    {
      "id": "uuid",
      "title": "Название раздела",
      "pages": [
        {
          "id": "uuid",
          "title": "Заголовок страницы",
          "content": { "type": "doc", "content": [] },
          "source_task_id": "uuid или null",
          "created_at": "...",
          "updated_at": "..."
        }
      ]
    }
  ]
}
```
