-- 1. Статусы (колонки доски)
CREATE TABLE IF NOT EXISTS statuses (
    id SERIAL PRIMARY KEY,
    status_name TEXT NOT NULL UNIQUE
);
INSERT INTO statuses (status_name) VALUES
    ('TODO'),
    ('IN PROGRESS'),
    ('DONE'),
    ('REVIEW'),
    ('BLOCKED'),
    ('PLANNED'),
    ('TESTING'),
    ('ARCHIVED'),
    ('CANCELLED'),
    ('REOPENED');

-- 1 статус -> много задач; 1 задача -> 1 статус → 1:N
-- отношение: statuses (1) ← statuses_id (N) в tasks

-- 2. Приоритеты
CREATE TABLE IF NOT EXISTS priorities (
    id SMALLINT PRIMARY KEY,
    priority_name TEXT NOT NULL UNIQUE,
    weight SMALLINT NOT NULL,
    keywords TEXT[] NOT NULL
);

INSERT INTO priorities (id, priority_name, weight, keywords) VALUES
(1, 'Critical',   1, ARRAY['critical', 'crit', 'urgent']),
(2, 'High',       2, ARRAY['high', 'important', 'major']),
(3, 'Medium',     3, ARRAY['medium', 'normal', 'standard']),
(4, 'Low',        4, ARRAY['low', 'minor', 'trivial']),
(5, 'Backlog',    5, ARRAY['backlog', 'later', 'someday', 'plan']),
(6, 'Blocked',    6, ARRAY['block', 'blocked']);


-- 1 приоритет -> много задач; 1 задача -> 1 приоритет → 1:N

-- 3. Пользователи
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    user_name TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    CONSTRAINT email_fmt CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
);

-- 1:N → один пользователь может:
--   - иметь много задач (как автор/assignee),
--   - быть в нескольких командах (см. team_members)

-- 4. Логин (многие поля связаны с users)
CREATE TABLE IF NOT EXISTS login (
    id SERIAL PRIMARY KEY,
    user_name TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    last_login TIMESTAMP WITH TIME ZONE,
    failed_login INTEGER NOT NULL DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT email_fmt_login CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
);

-- 1:1 → каждый login «привязан» к одному users; один пользователь — один логин

-- 5. Команды
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    team_name TEXT NOT NULL UNIQUE,
    description TEXT,
    config_dashboard JSONB NOT NULL DEFAULT '["TODO", "IN PROGRESS", "DONE"]'::JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- 1 командa -> много членов (team_members); много команд → один пользователь → M:N
-- (через team_members)

-- 6. Связь команды ↔️ пользователи (многие‑ко‑многим)
CREATE TABLE IF NOT EXISTS team_members (
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    role TEXT,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (team_id, user_id)
);

-- M:N: один пользователь — несколько команд, одна команда — несколько пользователей

-- 7. Задачи (основная сущность)
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status_id INTEGER NOT NULL REFERENCES statuses(id) ON DELETE RESTRICT,
    priority_id SMALLINT NOT NULL REFERENCES priorities(id) ON DELETE RESTRICT,
    deadline TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    team_id INTEGER REFERENCES teams(id) ON DELETE SET NULL,         -- 1:N: команда → много задач
    created_by INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT, -- 1:N: автор
    assigned_to INTEGER REFERENCES users(id) ON DELETE SET NULL,       -- 1:N: исполнитель
    CONSTRAINT title_not_blank CHECK (length(trim(title)) > 0)
);

-- 8. Теги/метки
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    tag_name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS task_tags (
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, tag_id)
);

-- внешние ключи и частые фильтры
CREATE INDEX idx_team_members_user ON team_members (user_id);
CREATE INDEX idx_team_members_team ON team_members (team_id);

CREATE INDEX idx_tasks_team     ON tasks (team_id);
CREATE INDEX idx_tasks_status   ON tasks (status_id);
CREATE INDEX idx_tasks_priority ON tasks (priority_id);
CREATE INDEX idx_tasks_created  ON tasks (created_at);  -- для сортировки "сначала новые"
CREATE INDEX idx_tasks_assignee ON tasks (assigned_to);

CREATE INDEX idx_users_email    ON users (email);
CREATE INDEX idx_login_user_id  ON login (user_id);
