# Backend API

Backend слушает порт `8080` внутри Docker-сети. Frontend обращается к backend через Nginx-прокси по пути `/api/`.

## Health

### `GET /health`

Проверка доступности backend.

Ответ `200`:

```json
{
  "status": "ok"
}
```

## Tasks

### `GET /api/tasks`

Возвращает список задач.

Если параметр команды не передан, endpoint возвращает все задачи. Это поведение оставлено для совместимости с текущим frontend и быстрой разработки MVP.

Query-параметры:

- `team_id` - необязательный положительный integer. Если передан, backend вернет задачи только этой команды.

Примеры:

```http
GET /api/tasks
GET /api/tasks?team_id=1
```

Ответ `200`:

```json
[
  {
    "id": 1,
    "title": "Подготовить dashboard",
    "description": "Собрать колонки задач",
    "priority_name": "High",
    "status_name": "TODO",
    "deadline": "2026-05-17T18:00:00Z",
    "created_at": "2026-05-17T10:00:00Z",
    "updated_at": "2026-05-17T10:00:00Z",
    "team_id": 1,
    "team_name": "Frontend",
    "created_by": "user1",
    "assigned_to": "user2"
  }
]
```

Ошибки:

- `400` - `team_id` передан, но не является положительным integer.
- `500` - ошибка чтения задач из PostgreSQL.

### `GET /api/task/:id`

Возвращает одну задачу по `id`.

Пример:

```http
GET /api/task/1
```

Ответ `200`:

```json
{
  "id": 1,
  "title": "Подготовить dashboard",
  "description": "Собрать колонки задач",
  "priority_name": "High",
  "status_name": "TODO",
  "deadline": "2026-05-17T18:00:00Z",
  "created_at": "2026-05-17T10:00:00Z",
  "updated_at": "2026-05-17T10:00:00Z",
  "team_id": 1,
  "team_name": "Frontend",
  "created_by": "user1",
  "assigned_to": "user2"
}
```

Ошибки:

- `400` - `id` не является положительным integer.
- `404` - задача не найдена.
- `500` - ошибка чтения задачи из PostgreSQL.

## Teams

### `GET /api/teams`

Возвращает список команд.

Пример:

```http
GET /api/teams
```

Ответ `200`:

```json
[
  {
    "id": 1,
    "name": "Frontend",
    "description": "Команда frontend-разработки",
    "config_dashboard": ["TODO", "IN PROGRESS", "DONE"],
    "created_at": "2026-05-17 10:00:00+00",
    "updated_at": "2026-05-17 10:00:00+00",
    "tasks": ["Подготовить dashboard"],
    "members": ["user1", "user2"]
  }
]
```

Ошибки:

- `500` - ошибка чтения команд из PostgreSQL.

### `GET /api/team/:id`

Возвращает данные одной команды и настройки dashboard.

Пример:

```http
GET /api/team/1
```

Ответ `200`:

```json
{
  "id": 1,
  "name": "Frontend",
  "description": "Команда frontend-разработки",
  "config_dashboard": ["TODO", "IN PROGRESS", "DONE"],
  "created_at": "2026-05-17 10:00:00+00",
  "updated_at": "2026-05-17 10:00:00+00",
  "tasks": ["Подготовить dashboard"],
  "members": ["user1", "user2"]
}
```

Ошибки:

- `400` - `id` не является положительным integer.
- `404` - команда не найдена.
- `500` - ошибка чтения команды из PostgreSQL.

