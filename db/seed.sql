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
    '$2a$12$aH67meqOtARuLUoiUHTS4OqkUHJAnNPiNf3Sz8Soh964sKlyzWhpa',
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

-- 6. Создаем задачи так, чтобы у каждой команды были видимые карточки на dashboard.
-- Статусы задач должны входить в config_dashboard соответствующей команды.
INSERT INTO tasks (title, description, status_id, priority_id, deadline, team_id, created_by, assigned_to)
SELECT
    task_data.title,
    task_data.description,
    (SELECT id FROM statuses WHERE status_name = task_data.status_name),
    (SELECT id FROM priorities WHERE priority_name = task_data.priority_name),
    now() + task_data.deadline_offset,
    (SELECT id FROM teams WHERE team_name = task_data.team_name),
    (SELECT id FROM users WHERE user_name = task_data.created_by),
    (SELECT id FROM users WHERE user_name = task_data.assigned_to)
FROM (
    VALUES
        ('Backend Core', 'TODO', 'High', 'Егор', 'Людмила', 'Fix connection pool timeout', 'Поймать таймауты подключения к базе, пока они не поймали нас на демо', interval '1 day'),
        ('Backend Core', 'IN PROGRESS', 'Critical', 'Егор', 'Александр', 'Implement JWT Auth', 'Сделать вход так, чтобы пароль не жил в коде и не смотрел на нас с укором', interval '2 days'),
        ('Backend Core', 'REVIEW', 'Medium', 'Егор', 'Людмила', 'Optimize SQL queries for dashboard', 'Уговорить запросы работать быстрее, чем открывается вкладка с кофе', interval '3 days'),
        ('Backend Core', 'DONE', 'Low', 'Егор', 'Александр', 'Write API documentation', 'Описать API до того, как frontend начнет гадать по логам', interval '4 days'),

        ('Frontend UI', 'TODO', 'High', 'Ольга', 'Дмитрий', 'Create Kanban Board View', 'Собрать доску, где карточки выглядят как задачи, а не как потерянные div', interval '1 day'),
        ('Frontend UI', 'IN PROGRESS', 'Medium', 'Ольга', 'РешательПроблем', 'Fix responsive layout on mobile', 'Сделать так, чтобы кнопки не устраивали давку на маленьком экране', interval '2 days'),
        ('Frontend UI', 'TESTING', 'High', 'Ольга', 'Дмитрий', 'Add input validation for login form', 'Не пускать пустой пароль с видом важного пользователя', interval '3 days'),
        ('Frontend UI', 'DONE', 'Low', 'Ольга', 'РешательПроблем', 'Add dark mode support', 'Добавить темную тему для тех, кто деплоит после полуночи', interval '4 days'),

        ('Мстители перходников', 'TODO', 'Medium', 'StivJobs', 'Сергей', 'Сделать так, чтобы новый софт запускался на старых компьютерах', 'Старые компьютеры не виноваты, что видели больше релизов, чем мы', interval '2 days'),
        ('Мстители перходников', 'IN PROGRESS', 'Critical', 'StivJobs', 'Кристика', 'Сократить количество экранов и кнопок до абсолютного минимума', 'Оставить только то, что нельзя случайно объяснить словом "потом"', interval '3 days'),
        ('Мстители перходников', 'TESTING', 'High', 'StivJobs', 'Сергей', 'Устранить баги и выпустить пакеты обновлений', 'Сначала устранить баги, потом сделать вид, что так и планировалось', interval '4 days'),
        ('Мстители перходников', 'BLOCKED', 'Critical', 'StivJobs', 'Кристика', 'Проанализировать заполнение таблицы', 'Выяснить, почему таблица выглядит так, будто ее заполняли в пятницу вечером', interval '5 days'),

        ('BSOD Survivors', 'TODO', 'High', 'BillGates', 'Marianna', 'Setup CI/CD pipeline', 'Попросить pipeline собираться сам и не требовать личного внимания каждый раз', interval '1 day'),
        ('BSOD Survivors', 'IN PROGRESS', 'Medium', 'BillGates', 'Настя', 'Handle 500 errors gracefully', 'Падать красиво, информативно и без драматического молчания', interval '2 days'),
        ('BSOD Survivors', 'DONE', 'Low', 'BillGates', 'Пиццерия-для-всех', 'Cover endpoints with unit tests', 'Проверить endpoint до того, как endpoint проверит нашу выдержку', interval '3 days'),

        ('To the Moon', 'PLANNED', 'Medium', 'ElonMusk', 'PiterPen', 'Разработать план анимации и вау-эффекты для презентации продукта клиенту', 'Главное, чтобы вау-эффект не оказался ошибкой 500 на большом экране', interval '4 days'),
        ('To the Moon', 'TODO', 'High', 'ElonMusk', 'PiterPen', 'Собрать работающую демо-версию сложного продукта в нереально короткие сроки', 'Демо должно работать хотя бы в той вселенной, где его показывают', interval '5 days'),
        ('To the Moon', 'IN PROGRESS', 'Critical', 'ElonMusk', 'ElonMusk', 'Внедрить ии-помощников в рутинные задачи', 'Пусть рутина автоматизируется, а не расширяется до отдельного проекта', interval '6 days'),
        ('To the Moon', 'DONE', 'Low', 'ElonMusk', 'PiterPen', 'Полировка интерфейса', 'Натереть интерфейс до состояния "не стыдно открыть при людях"', interval '7 days'),

        ('Algorithm King', 'TODO', 'High', 'MarkZukerberg', 'MarkZukerberg', 'Собрать и обрабатывать терабайты данных, чтобы показывать максимально точную таргетную рекламу', 'Сначала собрать данные, потом понять, зачем их было так много', interval '1 day'),
        ('Algorithm King', 'REVIEW', 'Medium', 'MarkZukerberg', 'MarkZukerberg', 'Настроить ленту новостей так, чтобы человек не мог оторваться от ней часами', 'Проверить, что лента не отрывает разработчика от дедлайна первой', interval '2 days'),
        ('Algorithm King', 'TESTING', 'High', 'MarkZukerberg', 'MarkZukerberg', 'Заставить тяжелые алгоритмы работать быстрее', 'Объяснить алгоритмам, что дедлайн тоже имеет вычислительную сложность', interval '3 days'),
        ('Algorithm King', 'ARCHIVED', 'Low', 'MarkZukerberg', 'MarkZukerberg', 'Переписать весь код проекта с нуля за одни выходные', 'Заархивировать идею до тех пор, пока все не выспятся', interval '4 days')
) AS task_data(team_name, status_name, priority_name, created_by, assigned_to, title, description, deadline_offset);

-- 7. Привязываем случайные теги к созданным задачам (по 1-2 тега на задачу)
INSERT INTO task_tags (task_id, tag_id)
SELECT t.id, tag.id
FROM tasks t
CROSS JOIN LATERAL (
    SELECT id FROM tags ORDER BY random() LIMIT 2
) tag
ON CONFLICT DO NOTHING;
