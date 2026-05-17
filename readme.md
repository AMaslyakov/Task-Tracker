# Архитектура Task-Tracker


## 1. Назначение системы

Task-Tracker - веб-приложение для управления задачами команд. Пользователь входит в систему, выбирает команду и работает с kanban-доской задач, где колонки определяются настройкой команды.

Основные домены:

- пользователи и данные входа;
- сессии авторизации;
- команды и участники команд;
- задачи, статусы, приоритеты и теги;
- планируемые события задач через RabbitMQ и SSE.

## 2. Высокоуровневая схема

```text
Browser
  |
  | HTTP, cookie session_id
  v
frontend container
  Nginx
  - static Vue build
  - /api/* proxy
  |
  | /api/* -> http://backend:8080
  v
backend container
  Go + Gin API
  |
  | SQL via pgxpool
  v
postgres container
  PostgreSQL

broker container
  RabbitMQ management
  Сейчас поднят в Docker Compose, но backend-код к нему еще не подключен.
```

В production-compose внешний вход в приложение - `http://localhost:8080`, где Nginx раздает frontend и проксирует `/api/` в backend. Backend также проброшен наружу как `http://localhost:8081` для прямой проверки API и Swagger.

## 3. Контейнеры и сеть

Система разворачивается через Docker Compose. Имена сервисов и контейнеров фиксированы:

- `frontend` - production-сборка Vue, раздача через Nginx, порт хоста `8080:80`;
- `backend` - Go/Gin HTTP API, порт внутри сети `8080`, порт хоста `8081:8080`;
- `postgres` - PostgreSQL, порт `5432:5432`, volume `postgres_data`;
- `broker` - RabbitMQ management image, порты `5672:5672` и `15672:15672`, volume `broker_data`.

Внутри Docker-сети сервисы обращаются друг к другу по compose-именам. Для frontend это важно: Nginx проксирует API на `backend:8080`.

`docker-compose.dev.yml` отличается тем, что `frontend` запускается как Vite dev server в контейнере `node:22-alpine` на порту `5173`, а Vite проксирует `/api` на target из `VITE_API_PROXY_TARGET`.

## 4. Frontend

Frontend расположен в `frontend/` и построен на Vue 3 + Vite.

Ключевые файлы:

- `frontend/src/main.js` - точка входа Vue;
- `frontend/src/App.vue` - корневой компонент;
- `frontend/src/router/index.js` - маршруты приложения;
- `frontend/src/pages/LoginPage.vue` - экран входа;
- `frontend/src/pages/DashboardPage.vue` - основная доска задач;
- `frontend/src/pages/TaskPage.vue` - заготовка страницы задачи;
- `frontend/src/api/auth.js` - auth-запросы;
- `frontend/src/api/tasks.js` - запросы команд и задач;
- `frontend/src/components/*` - header, kanban-колонки, карточки задач, форма задачи.

Активные маршруты:

- `/` перенаправляет на `/login`;
- `/login` показывает форму входа;
- `/tasks` показывает dashboard;
- `/task/id` подключен к странице задачи-заготовке.

Авторизация на frontend не использует `localStorage` для токенов. После `POST /api/login` браузер хранит `HttpOnly` cookie `session_id`, а fetch-запросы к `/api/...` выполняются с `credentials: 'same-origin'`.

Dashboard при загрузке:

1. вызывает `GET /api/me`;
2. при `401` отправляет пользователя на `/login`;
3. загружает команды через `GET /api/teams`;
4. загружает задачи через `GET /api/tasks`;
5. фильтрует задачи выбранной команды на клиенте;
6. сохраняет выбранную команду в `localStorage`.

Kanban-доска строится через `StatusColumnList`. Список колонок берется из `team.config_dashboard`, который backend отдает как массив статусов. Drag-and-drop реализован через `sortablejs-vue3` и сохраняет изменение статуса через `PATCH /api/task/:id/status`.

Фактическое состояние frontend CRUD:

- создание задачи открывается из dashboard через `TaskForm` и вызывает `POST /api/task`;
- редактирование вызывает `PATCH /api/task/:id`;
- удаление вызывает `DELETE /api/task/:id`;
- задачи после ответов backend маппятся в формат карточек через `mapTask`.

Важные текущие ограничения frontend:

- список задач пока загружается целиком через `GET /api/tasks`, а не через `GET /api/tasks?team_id=...`;
- теги в карточках не загружаются из backend и маппятся как пустой список;
- поле исполнителя в UI исторически называется `asigned_to`, а backend отдает `assigned_to`; маппинг выполняется в `frontend/src/api/tasks.js`;
- realtime через `EventSource` еще не подключен.

## 5. Backend

Backend расположен в `backend/src` и написан на Go + Gin.

Ключевые файлы:

- `backend/src/main.go` - запуск приложения, подключение PostgreSQL, регистрация маршрутов;
- `backend/src/API/models.go` - DTO и модели ответов;
- `backend/src/API/DBwork.go` - SQL-запросы и работа с `pgxpool`;
- `backend/src/API/auth.go` - login, me, logout, session middleware;
- `backend/src/API/get.go` - чтение задач и команд;
- `backend/src/API/post.go` - создание задач и пользователей;
- `backend/src/API/patch.go` - обновление задач, статуса и пользователей;
- `backend/src/API/delete.go` - удаление задач и пользователей;
- `backend/src/docs/` - generated Swagger/OpenAPI docs.

Backend подключается к PostgreSQL через переменную `DB_DSN`. Подключение создается один раз через `pgxpool.New` и хранится в глобальной переменной `API.Pool`.

Роуты:

- `GET /health` - публичная проверка backend;
- `POST /api/login` - публичный вход;
- `GET /api/me` - текущий пользователь по cookie;
- `POST /api/logout` - удаление сессии;
- `POST /api/user`, `PATCH /api/users/:id`, `DELETE /api/users/:id` - пользовательские операции, сейчас зарегистрированы публично;
- защищенная группа `/api` с `RequireAuth()`:
  - `GET /api/tasks`;
  - `GET /api/tasks?team_id=...`;
  - `GET /api/task/:id`;
  - `POST /api/task`;
  - `PATCH /api/task/:id`;
  - `PATCH /api/task/:id/status`;
  - `DELETE /api/task/:id`;
  - `GET /api/teams`;
  - `GET /api/team/:id`;
- `GET /swagger/index.html` - Swagger UI.

Auth-flow:

1. `POST /api/login` принимает email/password.
2. Backend ищет запись в `login` по email.
3. Пароль проверяется через bcrypt.
4. Backend генерирует случайный `session_id`.
5. Сессия сохраняется в таблицу `sessions`.
6. Ответ выставляет `HttpOnly` cookie `session_id` на 30 минут.
7. Защищенные роуты используют `RequireAuth()`, читают cookie и ищут пользователя через `sessions`.

Сессии истекают логически через условие `sessions.created_at >= now() - INTERVAL '30 minutes'`. Дополнительно в БД есть trigger, который удаляет старые сессии перед вставкой новых.

## 6. База данных

Схема лежит в `db/query.sql`, seed-данные (для первичного заполнения выдуманными данными) - в `db/seed.sql`.

Таблицы:

- `statuses` - справочник статусов/колонок;
- `priorities` - справочник приоритетов;
- `users` - профиль пользователя;
- `login` - учетные данные, bcrypt hash, admin-флаг и метаданные входа;
- `teams` - команды и `config_dashboard` как JSONB-массив статусов;
- `team_members` - связь пользователей с командами;
- `tasks` - основная сущность задачи;
- `tags` - справочник тегов;
- `task_tags` - связь задач и тегов;
- `sessions` - cookie-сессии.

Основные связи:

- `login.user_id -> users.id`;
- `team_members.team_id -> teams.id`;
- `team_members.user_id -> users.id`;
- `tasks.status_id -> statuses.id`;
- `tasks.priority_id -> priorities.id`;
- `tasks.team_id -> teams.id`;
- `tasks.created_by -> users.id`;
- `tasks.assigned_to -> users.id`;
- `task_tags.task_id -> tasks.id`;
- `task_tags.tag_id -> tags.id`;
- `sessions.user_id -> users.id`.

PostgreSQL init-файлы подключены в compose:

- `./db/query.sql` -> `/docker-entrypoint-initdb.d/01-schema.sql`;
- `./db/seed.sql` -> `/docker-entrypoint-initdb.d/02-seed.sql`.

Важно: эти SQL-файлы выполняются только при создании нового пустого `postgres_data` volume. Для существующей базы изменения схемы или seed не применяются автоматически.

## 7. API и данные dashboard

Dashboard использует две группы данных:

- команды из `GET /api/teams`;
- задачи из `GET /api/tasks`.

Backend для задач возвращает плоский DTO:

- `id`;
- `title`;
- `description`;
- `priority_name`;
- `status_name`;
- `deadline`;
- `created_at`;
- `updated_at`;
- `team_id`;
- `team_name`;
- `created_by`;
- `assigned_to`.

Frontend преобразует его в карточку:

- `status` берется из `status_name`;
- `priority` берется из `priority_name`;
- `deadline` форматируется для отображения;
- `command` сопоставляется по `team_id`;
- `asigned_to.name` заполняется из `assigned_to`.

Команда возвращает:

- `id`;
- `name`;
- `description`;
- `config_dashboard`;
- `tasks`;
- `members`;
- `member_details`.

`config_dashboard` управляет тем, какие колонки будут показаны. Если у задачи статус отсутствует в настройке выбранной команды, frontend не покажет эту задачу ни в одной колонке.

## 8. RabbitMQ и realtime

RabbitMQ уже присутствует в `docker-compose.yml` и `docker-compose.dev.yml` как сервис `broker`.

Compose передает backend переменные:

- `BROKER_URL`;
- `BROKER_EXCHANGE=task_events`.

Фактическое состояние backend-кода:

- AMQP-клиент в `go.mod` не подключен;
- backend не читает `BROKER_URL`;
- exchange `task_events` не создается;
- события задач не публикуются;
- consumer отсутствует;
- `GET /api/events` отсутствует;
- frontend `EventSource` не использует.

Целевая схема realtime из `PLAN.md`:

```text
Frontend A -> PATCH /api/task/:id/status
Backend -> PostgreSQL update
Backend -> RabbitMQ publish task.status_changed
Backend consumer -> receives event
Backend SSE hub -> broadcasts event to connected browsers
Frontend B/C -> update dashboard without reload
```

Минимальные планируемые события:

- `task.created`;
- `task.updated`;
- `task.status_changed`;
- `task.deleted`.

Frontend не должен подключаться к RabbitMQ напрямую. Единственный realtime-интерфейс для браузера по плану - backend endpoint `GET /api/events` через SSE.

## 9. Сборка и запуск

Production-like запуск:

```bash
docker compose up --build
```

Dev-запуск:

```bash
docker compose -f docker-compose.dev.yml up --build
```

Проверка frontend production build:

```bash
docker build -t task-tracker-frontend-check ./frontend
```

Проверка backend:

```bash
docker run --rm -v "$PWD/backend/src:/app" -w /app golang:1.25 go test ./...
```

Форматирование Go:

```bash
docker run --rm -v "$PWD/backend/src:/app" -w /app golang:1.25 gofmt -w main.go API/*.go
```

По текущей договоренности локально полагаться на Docker, потому что зависимости на хосте могут отсутствовать.

## 10. Текущие архитектурные ограничения

- Нет полноценной системы миграций: используется init SQL PostgreSQL.
- Нет автоматического применения изменений схемы/seed к существующему `postgres_data`.
- Backend пока монолитно организован вокруг пакета `API` и глобального `API.Pool`.
- Пользовательские роуты `POST /api/user`, `PATCH /api/users/:id`, `DELETE /api/users/:id` зарегистрированы вне auth middleware.
- Ролевая модель фактически не enforced, кроме хранения `is_admin`.
- RabbitMQ присутствует инфраструктурно, но не участвует в runtime-потоке.
- SSE endpoint и frontend realtime отсутствуют.
- Автотесты для frontend/backend пока не добавлены.
- Swagger генерируется в Dockerfile через `swag init`, generated docs лежат в репозитории.

## 11. Что можно доделать, но мы не успели


1. Подключить AMQP-клиент backend к RabbitMQ.
2. Создавать exchange `task_events` при старте backend.
3. Публиковать события после успешных операций с задачами.
4. Добавить consumer и внутренний SSE hub.
5. Добавить `GET /api/events`.
6. Подключить frontend к SSE через `EventSource`.
7. Перевести dashboard на загрузку задач выбранной команды через `GET /api/tasks?team_id=...`.
8. Добавить миграции PostgreSQL вместо init-only схемы.
9. Добавить базовые backend/frontend тесты.

## Развернутое приложение

http://80.87.102.162
