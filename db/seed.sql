-- Очищаем таблицы перед заполнением (опционально)
TRUNCATE task_tags, tags, tasks, team_members, teams, login, users RESTART IDENTITY CASCADE;

-- 1. Создаем пользователей
INSERT INTO users (user_name, email)
VALUES
    ('Александр', 'alexandr@codiki.com'),
    ('Егор', 'egor@codiki.com'),
    ('Людмила', 'ludmila@codiki.com'),
    ('Дмитрий', 'dmitry@cobalt.com'),
    ('Ольга', 'olga@cobalt.com'),
    ('Петр', 'petr@cobalt.com'),
    ('Сергей', 'sergty@cobalt.com'),
    ('Настя', 'nastia@mail.ru'),
    ('Кристика', 'kristina@mail.ru'),
    ('Пиццерия-для-всех', 'pizza@niam.com'),
    ('РешательПроблем', 'coolboy@njproblem.info'),
    ('Marianna', 'mari@syperagent.net'),
    ('PiterPen', 'icanfly@good.org'),
    ('StivJobs', 'next_liza@onlymac.com'),
    ('BillGates', 'ctrl_alt_delete@microsoft.com'),
    ('MarkZukerberg', 'reptiloid1984@cobalt.com'),
    ('ElonMusk', 'marsorbust@x.com')
ON CONFLICT DO NOTHING;

-- 2. Создаем учетные записи (login) для этих пользователей
INSERT INTO login (user_name, email, password_hash, is_admin, user_id)
SELECT
    u.user_name,
    u.email,
    '$2a$12$MockHashForTestingPurposesOnlyDocNotUseInProd',
    CASE WHEN u.user_name IN ('MarkZukerberg','Ольга','StivJobs','BillGates','PiterPen','Егор') THEN TRUE ELSE FALSE END,
    u.id
FROM users u
ON CONFLICT DO NOTHING;

-- 3. Создаем тестовые команды
INSERT INTO teams (team_name, description, config_dashboard)
VALUES
    ('Backend Core', 'Разработка серверной части на Go (Gin + pgx)', '["TODO", "IN PROGRESS", "REVIEW", "DONE"]'::jsonb),
    ('Frontend UI', 'Разработка клиентского интерфейса на Vue.js', '["TODO", "IN PROGRESS", "TESTING", "DONE"]'::jsonb),
    ('Мстители перходников', 'Для тех, кто решает сложные проблемы совместимости', '["TODO", "IN PROGRESS", "TESTING", "BLOCKED","DONE"]'::jsonb),
    ('BSOD Survivors', 'Выжившие после синего экрана. Тут только стрессоустойчивые тестировщики и DevOps-инженеры', '["TODO", "IN PROGRESS", "DONE"]'::jsonb),
    ('To the Moon', 'для стартапов с амбициозными планами и взрывным ростом', '["PLANNED","TODO", "IN PROGRESS","DONE"]'::jsonb),
    ('Algorithm King', 'Для специалистов по Big Data которые знают о пользователях всё', '["TODO", "IN PROGRESS", "REVIEW","PLANNED","TESTING", "ARCHIVED","CANCELLED", "DONE"]'::jsonb)
ON CONFLICT DO NOTHING;

-- 4. Добавляем пользователей в команды (team_members)
INSERT INTO team_members (team_id, user_id, is_admin, role)
VALUES
    ((SELECT id FROM teams WHERE team_name = 'Backend Core'), (SELECT id FROM users WHERE user_name = 'Егор'), TRUE, 'Lead Go Developer'),
    ((SELECT id FROM teams WHERE team_name = 'Backend Core'), (SELECT id FROM users WHERE user_name = 'Людмила'), FALSE, 'Junior Go Developer'),
    ((SELECT id FROM teams WHERE team_name = 'Backend Core'), (SELECT id FROM users WHERE user_name = 'Александр'), FALSE, 'Project Manager'),

    ((SELECT id FROM teams WHERE team_name = 'Frontend UI'), (SELECT id FROM users WHERE user_name = 'Ольга'), TRUE, 'Lead Vue Developer'),
    ((SELECT id FROM teams WHERE team_name = 'Frontend UI'), (SELECT id FROM users WHERE user_name = 'Дмитрий'), FALSE, 'QA Engineer'),
    ((SELECT id FROM teams WHERE team_name = 'Frontend UI'), (SELECT id FROM users WHERE user_name = 'РешательПроблем'), FALSE, 'Project Manager'),

    ((SELECT id FROM teams WHERE team_name = 'Мстители перходников'), (SELECT id FROM users WHERE user_name = 'StivJobs'), TRUE, 'Lead'),
    ((SELECT id FROM teams WHERE team_name = 'Мстители перходников'), (SELECT id FROM users WHERE user_name = 'Сергей'), FALSE, 'DevOps'),
    ((SELECT id FROM teams WHERE team_name = 'Мстители перходников'), (SELECT id FROM users WHERE user_name = 'Кристика'), FALSE, 'DataEngineer'),

    ((SELECT id FROM teams WHERE team_name = 'BSOD Survivors'), (SELECT id FROM users WHERE user_name = 'BillGates'), TRUE, 'TeamLead'),
    ((SELECT id FROM teams WHERE team_name = 'BSOD Survivors'), (SELECT id FROM users WHERE user_name = 'Marianna'), FALSE, 'Engineer'),
    ((SELECT id FROM teams WHERE team_name = 'BSOD Survivors'), (SELECT id FROM users WHERE user_name = 'Пиццерия-для-всех'), FALSE, 'Manager'),
    ((SELECT id FROM teams WHERE team_name = 'BSOD Survivors'), (SELECT id FROM users WHERE user_name = 'Настя'), FALSE, 'Manager'),

    ((SELECT id FROM teams WHERE team_name = 'To the Moon'), (SELECT id FROM users WHERE user_name = 'ElonMusk'), TRUE, 'Lead'),
    ((SELECT id FROM teams WHERE team_name = 'To the Moon'), (SELECT id FROM users WHERE user_name = 'PiterPen'), FALSE, 'Contributor'),

    ((SELECT id FROM teams WHERE team_name = 'Algorithm King'), (SELECT id FROM users WHERE user_name = 'MarkZukerberg'), TRUE, 'Lead Vue Developer')
ON CONFLICT DO NOTHING;

-- 5. Создаем справочник базовых тегов
INSERT INTO tags (tag_name) VALUES
    ('bug'), ('feature'), ('refactoring'), ('api'), ('ui'), ('security'), ('docs')
ON CONFLICT DO NOTHING;

-- 6. Генерируем 15 случайных задач
INSERT INTO tasks (title, description, status_id, priority_id, deadline, team_id, created_by, assigned_to)
SELECT
    (titles[array_idx] || ' #' || i) AS title,
    'Автоматически сгенерированное описание для тестовой задачи под номером ' || i AS description,
    (floor(random() * 5) + 1)::integer AS status_id,
    (floor(random() * 4) + 1)::smallint AS priority_id,
    now() + (random() * interval '14 days') AS deadline,
    (SELECT id FROM teams ORDER BY random() LIMIT 1) AS team_id,
    (SELECT id FROM users WHERE user_name = 'Александр' LIMIT 1) AS created_by,
    CASE WHEN random() > 0.2 THEN (SELECT id FROM users ORDER BY random() LIMIT 1) ELSE NULL END AS assigned_to
FROM (
    SELECT
        ARRAY[
        'Fix connection pool timeout', 'Implement JWT Auth', 'Create Kanban Board View',
        'Optimize SQL queries for dashboard', 'Fix responsive layout on mobile', 'Write API documentation',
        'Setup CI/CD pipeline', 'Add input validation for login form', 'Refactor store module in Vue',
        'Handle 500 errors gracefully', 'Cover endpoints with unit tests', 'Add dark mode support',
        'Проанализировать заполнение таблицы', 'Добавить модуль передвижения абонентов', 'Сократить количество экранов и кнопок до абсолютного минимума',
        'Полировка интерфейса', 'Разработать план анимации и вау-эффекты для презентации продукта клиенту','Сделать так, чтобы новый софт запускался на старых компьютерах',
        'Внедрить ии-помощников в рутинные задачи','Устранить баги и выпустить пакеты обновлений', 'Настроить ленту новостей так, чтобы человек не мог оторваться от ней часами',
        'Собрать и обрабатывать терабайты данных, чтобы показывать максимально точную таргетную рекламу',
        'Перепивать весь код проекта с нуля за одни выходные','Заставить тяжелые алгоритмы работать быстрее','Собрать работающую демо-версию сложного продукта в нереально короткие сроки'
        ] AS titles,
        floor(random() * array_length(ARRAY[
        'Fix connection pool timeout', 'Implement JWT Auth', 'Create Kanban Board View',
        'Optimize SQL queries for dashboard', 'Fix responsive layout on mobile', 'Write API documentation',
        'Setup CI/CD pipeline', 'Add input validation for login form', 'Refactor store module in Vue',
        'Handle 500 errors gracefully', 'Cover endpoints with unit tests', 'Add dark mode support',
        'Проанализировать заполнение таблицы', 'Добавить модуль передвижения абонентов', 'Сократить количество экранов и кнопок до абсолютного минимума',
        'Полировка интерфейса', 'Разработать план анимации и вау-эффекты для презентации продукта клиенту','Сделать так, чтобы новый софт запускался на старых компьютерах',
        'Внедрить ии-помощников в рутинные задачи','Устранить баги и выпустить пакеты обновлений', 'Настроить ленту новостей так, чтобы человек не мог оторваться от ней часами',
        'Собрать и обрабатывать терабайты данных, чтобы показывать максимально точную таргетную рекламу',
        'Перепивать весь код проекта с нуля за одни выходные','Заставить тяжелые алгоритмы работать быстрее','Собрать работающую демо-версию сложного продукта в нереально короткие сроки'
        ], 1) + 1)::int AS array_idx
) s
CROSS JOIN generate_series(1,15) AS i;

-- 7. Привязываем случайные теги к созданным задачам (по 1-2 тега на задачу)
INSERT INTO task_tags (task_id, tag_id)
SELECT t.id, tag.id
FROM tasks t
CROSS JOIN LATERAL (
    SELECT id FROM tags ORDER BY random() LIMIT 2
) tag
ON CONFLICT DO NOTHING;


