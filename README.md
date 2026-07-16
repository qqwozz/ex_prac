# Exam Trainer

Тренировка заданий ЕГЭ и ОГЭ с геймификацией.

## Структура

```
ex_prac/
├── assets/              # Статические ресурсы
│   └── images/
│       └── logo.png
├── frontend/            # Интерфейс пользователя
│   └── training.html    # Тренировочный режим
├── model/               # AI-интеграция (не в git)
│   └── AI_API_KEY.txt
├── docs/                # Документация
│   ├── API.md           # Описание API
│   └── ARCHITECTURE.md  # Архитектура системы
└── README.md
```

## Компоненты

| Папка | Описание |
|-------|----------|
| `assets/` | Логотипы, изображения, шрифты |
| `frontend/` | HTML/CSS/JS страницы |
| `model/` | API-ключи для AI (в .gitignore) |
| `docs/` | Техническая документация |

## Запуск

Откройте `frontend/training.html` в браузере.

## Документация

Смотрите папку `docs/`:
- [API.md](docs/API.md) — описание эндпоинтов
- [ARCHITECTURE.md](docs/ARCHITECTURE.md) — структура приложения
