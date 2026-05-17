# Backend API

Backend слушает порт `8080` внутри Docker-сети. Frontend обращается к backend через Nginx-прокси по пути `/api/`.

В основном `docker-compose.yml` backend не публикует порт на хост. В деплойной сборке backend отдает только рабочие API-роуты приложения.

Endpoints задач и команд требуют авторизации через cookie `session_id`. Получить cookie можно через `POST /api/login`.

## Health

### `GET /health`

Проверка доступности backend.

Ответ `200`:

```json
{
  "status": "ok"
}
```

## Auth

Для локальных seed-пользователей используется единый демо-пароль:

```text
password123
```

### `POST /api/login`

Проверяет email/password, создает сессию в PostgreSQL и выставляет `HttpOnly` cookie `session_id` на 30 минут.

Payload:

```json
{
  "email": "egor@codiki.com",
  "password": "password123"
}
```

Ответ `200`:

```json
{
  "user": {
    "id": 2,
    "user_name": "Егор",
    "email": "egor@codiki.com",
    "is_admin": true
  }
}
```

Ошибки:

- `400` - некорректный payload или пустой email/password.
- `401` - неверный email или пароль.
- `500` - ошибка чтения пользователя или создания сессии.

### `GET /api/me`

Возвращает текущего пользователя по cookie `session_id`.

Ответ `200`:

```json
{
  "user": {
    "id": 2,
    "user_name": "Егор",
    "email": "egor@codiki.com",
    "is_admin": true
  }
}
```

Ошибки:

- `401` - cookie отсутствует, сессия не найдена или истекла.
- `500` - ошибка чтения сессии.

### `POST /api/logout`

Удаляет текущую сессию из PostgreSQL и очищает cookie `session_id`.

Ответ `204` - тело ответа пустое.

## Tasks

### `GET /api/tasks`

Возвращает список задач.

Требуется cookie `session_id`.

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

- `401` - пользователь не авторизован.
- `400` - `team_id` передан, но не является положительным integer.
- `500` - ошибка чтения задач из PostgreSQL.

### `GET /api/task/:id`

Возвращает одну задачу по `id`.

Требуется cookie `session_id`.

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

- `401` - пользователь не авторизован.
- `400` - `id` не является положительным integer.
- `404` - задача не найдена.
- `500` - ошибка чтения задачи из PostgreSQL.

### `POST /api/task`

Создает задачу.

Требуется cookie `session_id`.

Payload:

```json
{
  "title": "Подготовить CRUD",
  "description": "Добавить create/update/delete для задач",
  "status_name": "TODO",
  "priority_id": 2,
  "deadline": "2026-05-18T18:00:00Z",
  "team_id": 1,
  "created_by": 1,
  "assigned_to": 2
}
```

Ответ `201` - созданная задача в формате `GET /api/task/:id`.

Ошибки:

- `401` - пользователь не авторизован.
- `400` - payload некорректен или обязательные поля пустые.
- `404` - статус, приоритет, команда, автор или исполнитель не найдены.
- `500` - ошибка записи задачи в PostgreSQL.

### `PATCH /api/task/:id`

Обновляет поля задачи. Можно передавать любое непустое подмножество полей.
Чтобы очистить дедлайн или исполнителя, передайте `null` в `deadline` или `assigned_to`.

Требуется cookie `session_id`.

Payload:

```json
{
  "title": "Подготовить CRUD задач",
  "description": "Обновленное описание",
  "status_name": "IN PROGRESS",
  "priority_id": 1,
  "deadline": "2026-05-19T12:00:00Z",
  "team_id": 1,
  "assigned_to": 3
}
```

Ответ `200` - обновленная задача в формате `GET /api/task/:id`.

Ошибки:

- `401` - пользователь не авторизован.
- `400` - `id` не является положительным integer, payload пустой или поля невалидны.
- `404` - задача или связанная запись не найдены.
- `500` - ошибка обновления задачи в PostgreSQL.

### `PATCH /api/task/:id/status`

Меняет статус задачи. Endpoint предназначен для drag-and-drop между колонками.

Требуется cookie `session_id`.

Payload:

```json
{
  "status_name": "DONE"
}
```

Ответ `200` - обновленная задача в формате `GET /api/task/:id`.

Ошибки:

- `401` - пользователь не авторизован.
- `400` - `id` не является положительным integer или `status_name` пустой.
- `404` - задача или статус не найдены.
- `500` - ошибка обновления статуса в PostgreSQL.

### `DELETE /api/task/:id`

Удаляет задачу.

Требуется cookie `session_id`.

Ответ `204` - задача удалена, тело ответа пустое.

Ошибки:

- `401` - пользователь не авторизован.
- `400` - `id` не является положительным integer.
- `404` - задача не найдена.
- `500` - ошибка удаления задачи из PostgreSQL.

## Events

### `GET /api/events`

Открывает SSE-поток realtime-событий.

Требуется cookie `session_id`. Frontend подключается через `EventSource('/api/events')`; браузер отправляет cookie автоматически для same-origin запроса.

Backend отправляет именованные SSE-события:

- `task.created`;
- `task.updated`;
- `task.status_changed`;
- `task.deleted`.

Формат `data`:

```json
{
  "event_type": "task.status_changed",
  "payload": {
    "task_id": 1,
    "team_id": 1,
    "new_status": "DONE"
  }
}
```

Ошибки:

- `401` - пользователь не авторизован.

## Teams

### `GET /api/teams`

Возвращает список команд.

Требуется cookie `session_id`.

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
    "members": ["user1", "user2"],
    "member_details": [
      {
        "id": 1,
        "name": "user1",
        "email": "user1@example.com",
        "role": "Lead",
        "is_admin": true
      }
    ]
  }
]
```

Ошибки:

- `401` - пользователь не авторизован.
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
  "members": ["user1", "user2"],
  "member_details": [
    {
      "id": 1,
      "name": "user1",
      "email": "user1@example.com",
      "role": "Lead",
      "is_admin": true
    }
  ]
}
```

Ошибки:

- `400` - `id` не является положительным integer.
- `404` - команда не найдена.
- `500` - ошибка чтения команды из PostgreSQL.

Users
GET /api/users
Возвращает список пользователей.
Пример:
http
1
Ответ 200:
json
12345678910111213141516
Ошибки:
500 - ошибка чтения пользователей из PostgreSQL.
GET /api/user/:id
Возвращает данные одного пользователя по id.
Пример:
http
1
Ответ 200:
json
1234567
Ошибки:
400 - id не является положительным integer.
404 - пользователь не найден.
500 - ошибка чтения пользователя из PostgreSQL.
POST /api/user
Создает нового пользователя. Пароль хешируется перед сохранением.
Payload:
json
12345
Ответ 201:
json
1234
Ошибки:
400 - payload некорректен, обязательные поля пустые или не проходят валидацию (username < 3 символов, невалидный email, пароль < 6 символов).
409 - пользователь с таким username или email уже существует.
500 - ошибка записи пользователя в PostgreSQL.
PATCH /api/user/:id
Обновляет поля пользователя. Можно передавать любое непустое подмножество полей. Не переданные поля остаются без изменений.
Payload:
json
1234
Ответ 200:
json
123
Ошибки:
400 - id не является положительным integer, payload пустой или поля не проходят валидацию.
404 - пользователь не найден.
409 - новый username или email уже занят другим пользователем.
500 - ошибка обновления пользователя в PostgreSQL.
DELETE /api/user/:id
Удаляет пользователя. Связанная запись в таблице login удаляется автоматически (ON DELETE CASCADE).
Ответ 200:
json
123
Ошибки:
400 - id не является положительным integer.
404 - пользователь не найден.
500 - ошибка удаления пользователя из PostgreSQL.
