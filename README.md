# JobTrack API

Backend-сервис для отслеживания заявок на работу. Построен на Go + Gin + PostgreSQL.

## Стек

- **Go** + **Gin** — HTTP-сервер
- **PostgreSQL** + **GORM** — база данных с авто-миграциями
- **JWT** — аутентификация
- **Docker / Docker Compose** — контейнеризация

## Быстрый старт

### Docker (рекомендуется)

```bash
cp .env.example .env
docker compose up
```

API будет доступен на `http://localhost:8080`.

### Локально

Нужен запущенный PostgreSQL. Создайте БД `jobtrack`, затем:

```bash
cp .env.example .env
# отредактируйте .env под свои настройки
go run ./cmd/api
```

## Переменные окружения

| Переменная | По умолчанию | Описание |
|---|---|---|
| `PORT` | `8080` | Порт сервера |
| `DATABASE_URL` | `host=localhost ...` | DSN строка PostgreSQL |
| `JWT_SECRET` | `change-me-in-production` | Секрет для подписи JWT |

## API

### Здоровье

```
GET /health
```

### Аутентификация

```
POST /api/v1/auth/register
POST /api/v1/auth/login
```

**Регистрация:**
```json
{
  "email": "user@example.com",
  "password": "secret123",
  "name": "Иван"
}
```

**Вход** — возвращает `token`, который нужно передавать в заголовке:
```
Authorization: Bearer <token>
```

---

Все эндпоинты ниже требуют авторизации.

### Заявки на работу

```
GET    /api/v1/applications           — список заявок
POST   /api/v1/applications           — создать заявку
GET    /api/v1/applications/:id       — получить заявку (с интервью)
PUT    /api/v1/applications/:id       — обновить заявку
DELETE /api/v1/applications/:id       — удалить заявку
```

**Создание заявки:**
```json
{
  "company": "Yandex",
  "position": "Backend Developer",
  "status": "applied",
  "applied_at": "2024-01-15T10:00:00Z",
  "notes": "Откликнулся через hh.ru",
  "url": "https://hh.ru/vacancy/123"
}
```

**Статусы заявки:** `applied`, `phone_screen`, `interview`, `offer`, `rejected`, `accepted`, `withdrawn`

### Интервью

```
GET    /api/v1/applications/:id/interviews
POST   /api/v1/applications/:id/interviews
PUT    /api/v1/applications/:id/interviews/:interview_id
DELETE /api/v1/applications/:id/interviews/:interview_id
```

**Создание интервью:**
```json
{
  "scheduled_at": "2024-01-20T14:00:00Z",
  "type": "technical",
  "notes": "Алгоритмы и системный дизайн"
}
```

**Типы интервью:** `phone`, `technical`, `hr`, `onsite`, `other`

**Обновление** (можно добавить результат):
```json
{
  "result": "прошёл успешно",
  "notes": "спросили про горутины"
}
```

### Статистика

```
GET /api/v1/stats
```

**Ответ:**
```json
{
  "data": {
    "total": 15,
    "by_status": {
      "applied": 5,
      "interview": 4,
      "offer": 2,
      "rejected": 4
    },
    "this_month": 6,
    "interviews": 8
  }
}
```

## Структура проекта

```
cmd/api/          — точка входа
internal/
  config/         — загрузка конфигурации
  model/          — модели данных (User, JobApplication, Interview)
  repository/     — слой работы с БД (GORM)
  service/        — бизнес-логика
  handler/        — HTTP-хендлеры Gin
  middleware/     — JWT middleware
Dockerfile
docker-compose.yml
.env.example
```
