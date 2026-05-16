-- Очищаем таблицы перед заполнением (опционально, если нужно начать с чистого листа)
-- TRUNCATE task_tags, tags, tasks, team_members, teams, login, users RESTART IDENTITY CASCADE;

-- 1. Создаем случайных пользователей (5 человек)
INSERT INTO users (user_name, email)
VALUES
    ('ivan_dev', 'ivan@company.com'),
    ('elena_qa', 'elena@company.com'),
    ('alex_pm', 'alex@company.com'),
    ('dmitry_go', 'dmitry@company.com'),
    ('olga_vue', 'olga@company.com')
ON CONFLICT DO NOTHING;

-- 2. Создаем учетные записи (login) для этих пользователей
-- Все пароли захешированы заглушкой 'mock_hash_123'
INSERT INTO login (user_name, email, password_hash, is_admin, user_id)
SELECT
    user_name,
    email,
    '$2a$12$MockHashForTestingPurposesOnlyDocNotUseInProd', -- Имитация bcrypt хэша
    CASE WHEN user_name = 'alex_pm' THEN TRUE ELSE FALSE END, -- PM сделаем админом системы
    id
FROM users
ON CONFLICT DO NOTHING;

-- 3. Создаем тестовые команды
INSERT INTO teams (team_name, description, config_dashboard)
VALUES
    ('Backend Core', 'Разработка серверной части на Go (Gin + pgx)', '["TODO", "IN PROGRESS", "REVIEW", "DONE"]'::jsonb),
    ('Frontend UI', 'Разработка клиентского интерфейса на Vue.js', '["TODO", "IN PROGRESS", "TESTING", "DONE"]'::jsonb)
ON CONFLICT DO NOTHING;

-- 4. Добавляем пользователей в команды (Участники)
INSERT INTO team_members (team_id, user_id, is_admin, role)
VALUES
    -- Backend Core Team
    ((SELECT id FROM teams WHERE team_name = 'Backend Core'), (SELECT id FROM users WHERE user_name = 'dmitry_go'), TRUE, 'Lead Go Developer'),
    ((SELECT id FROM teams WHERE team_name = 'Backend Core'), (SELECT id FROM users WHERE user_name = 'ivan_dev'), FALSE, 'Junior Go Developer'),
    ((SELECT id FROM teams WHERE team_name = 'Backend Core'), (SELECT id FROM users WHERE user_name = 'alex_pm'), FALSE, 'Project Manager'),
    -- Frontend UI Team
    ((SELECT id FROM teams WHERE team_name = 'Frontend UI'), (SELECT id FROM users WHERE user_name = 'olga_vue'), TRUE, 'Lead Vue Developer'),
    ((SELECT id FROM teams WHERE team_name = 'Frontend UI'), (SELECT id FROM users WHERE user_name = 'elena_qa'), FALSE, 'QA Engineer'),
    ((SELECT id FROM teams WHERE team_name = 'Frontend UI'), (SELECT id FROM users WHERE user_name = 'alex_pm'), FALSE, 'Project Manager')
ON CONFLICT DO NOTHING;

-- 5. Создаем справочник базовых тегов
INSERT INTO tags (tag_name) VALUES
    ('bug'), ('feature'), ('refactoring'), ('api'), ('ui'), ('security'), ('docs')
ON CONFLICT DO NOTHING;

-- 6. Генерируем 15 случайных задач (ситуаций)
INSERT INTO tasks (title, description, status_id, priority_id, deadline, team_id, created_by, assigned_to)
SELECT
    -- Случайные названия и описания задач
    (ARRAY[
        'Fix connection pool timeout', 'Implement JWT Auth', 'Create Kanban Board View',
        'Optimize SQL queries for dashboard', 'Fix responsive layout on mobile', 'Write API documentation',
        'Setup CI/CD pipeline', 'Add input validation for login form', 'Refactor store module in Vue',
        'Handle 500 errors gracefully', 'Cover endpoints with unit tests', 'Add dark mode support'
    ])[floor(random() * 12) + 1] || ' #' || i AS title,

    'Автоматически сгенерированное описание для тестовой задачи под номером ' || i AS description,

    -- Случайный статус (от 1 до 5: TODO, IN PROGRESS, DONE, REVIEW, BLOCKED)
    (floor(random() * 5) + 1)::integer AS status_id,

    -- Случайный приоритет (от 1 до 4: Critical, High, Medium, Low)
    (floor(random() * 4) + 1)::smallint AS priority_id,

    -- Случайный дедлайн (в пределах следующих 14 дней)
    now() + (random() * interval '14 days') AS deadline,

    -- Случайная команда (1 или 2)
    (floor(random() * 2) + 1)::integer AS team_id,

    -- Автор всегда Project Manager (ID: 3)
    3 AS created_by,

    -- Случайный исполнитель из списка пользователей (от 1 до 5) или NULL (задача не назначена)
    CASE WHEN random() > 0.2 THEN (floor(random() * 5) + 1)::integer ELSE NULL END AS assigned_to
FROM generate_series(1, 15) AS i;

-- 7. Привязываем случайные теги к созданным задачам
INSERT INTO task_tags (task_id, tag_id)
SELECT
    t.id AS task_id,
    (floor(random() * 7) + 1)::integer AS tag_id
FROM tasks t
-- Добавим по 1-2 тега для каждой задачи
CROSS JOIN generate_series(1, 2)
ON CONFLICT DO NOTHING;
