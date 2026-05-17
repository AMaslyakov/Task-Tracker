-- 0. Сброс прерванной транзакции (если ещё не выполнено)
ROLLBACK;

-- 1. Исправление CHECK-констрейнтов для email (выполняется вне транзакции)
DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'email_fmt') THEN
    ALTER TABLE users DROP CONSTRAINT email_fmt;
    ALTER TABLE users ADD CONSTRAINT email_fmt 
      CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');
  END IF;
  IF EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'email_fmt_login') THEN
    ALTER TABLE login DROP CONSTRAINT email_fmt_login;
    ALTER TABLE login ADD CONSTRAINT email_fmt_login 
      CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');
  END IF;
END $$;

-- 2. Основная вставка данных
BEGIN;

-- Статусы
INSERT INTO statuses (id, status_name) VALUES
(1, 'TODO'), (2, 'IN PROGRESS'), (3, 'IN REVIEW'), (4, 'DONE')
ON CONFLICT (id) DO NOTHING;
SELECT setval('statuses_id_seq', (SELECT COALESCE(MAX(id),1) FROM statuses));

-- Приоритеты
INSERT INTO priorities (id, priority_name, weight, keywords) VALUES
(1, 'Low', 1, ARRAY['minor','cosmetic']),
(2, 'Medium', 2, ARRAY['normal','standard']),
(3, 'High', 3, ARRAY['urgent','important']),
(4, 'Critical', 4, ARRAY['blocker','emergency'])
ON CONFLICT (id) DO NOTHING;

-- Пользователи
INSERT INTO users (id, user_name, email) VALUES
(1, 'alice_dev', 'alice@techcorp.com'),
(2, 'bob_lead',  'bob@techcorp.com'),
(3, 'charlie_ui','charlie@techcorp.com'),
(4, 'diana_qa',  'diana@techcorp.com'),
(5, 'evan_ops',  'evan@techcorp.com')
ON CONFLICT (id) DO NOTHING;
SELECT setval('users_id_seq', (SELECT COALESCE(MAX(id),1) FROM users));

-- Логины
INSERT INTO login (id, user_name, email, password_hash, is_admin, user_id) VALUES
(1, 'alice_dev', 'alice@techcorp.com', '$2b$10$dummyhash_alice...', FALSE, 1),
(2, 'bob_lead',  'bob@techcorp.com',   '$2b$10$dummyhash_bob...',    TRUE,  2),
(3, 'charlie_ui','charlie@techcorp.com','$2b$10$dummyhash_charlie...',FALSE, 3),
(4, 'diana_qa',  'diana@techcorp.com', '$2b$10$dummyhash_diana...',  FALSE, 4),
(5, 'evan_ops',  'evan@techcorp.com',  '$2b$10$dummyhash_evan...',   TRUE,  5)
ON CONFLICT (id) DO NOTHING;
SELECT setval('login_id_seq', (SELECT COALESCE(MAX(id),1) FROM login));

-- Команды
INSERT INTO teams (id, team_name, description) VALUES
(1, 'Frontend', 'React, UI-компоненты'),
(2, 'Backend',  'API, БД, микросервисы'),
(3, 'DevOps',   'CI/CD, инфраструктура')
ON CONFLICT (id) DO NOTHING;
SELECT setval('teams_id_seq', (SELECT COALESCE(MAX(id),1) FROM teams));

-- Участники команд (M:N)
INSERT INTO team_members (team_id, user_id, is_admin, role) VALUES
(1, 1, FALSE, 'Senior Frontend Dev'), (1, 3, FALSE, 'UI/UX Designer'), (1, 2, TRUE, 'PM'),
(2, 2, TRUE,  'Tech Lead'),           (2, 4, FALSE, 'QA Engineer'),
(3, 5, TRUE,  'DevOps Lead'),         (3, 1, FALSE, 'Infra Consultant')
ON CONFLICT ON CONSTRAINT team_members_pkey DO NOTHING;

-- Теги
INSERT INTO tags (id, tag_name) VALUES
(1, 'bug'), (2, 'feature'), (3, 'devops'), 
(4, 'docs'), (5, 'ui/ux'), (6, 'performance')
ON CONFLICT (id) DO NOTHING;
SELECT setval('tags_id_seq', (SELECT COALESCE(MAX(id),1) FROM tags));

-- Задачи
INSERT INTO tasks (id, title, description, status_id, priority_id, team_id, created_by, assigned_to, deadline) VALUES
(1, 'Настроить CI/CD пайплайн', 'GitHub Actions + Docker', 1, 3, 3, 5, 5, NOW() + INTERVAL '3 days'),
(2, 'Авторизация OAuth2',       'Google & GitHub',        1, 4, 2, 2, 2, NOW() + INTERVAL '5 days'),
(3, 'Обновить зависимости',     'React 18, Tailwind',     1, 2, 1, 1, 1, NOW() + INTERVAL '10 days'),
(4, 'Реализовать дашборд',      'Графики, виджеты',       2, 2, 1, 2, 3, NOW() + INTERVAL '7 days'),
(5, 'Оптимизация SQL',          'Индексы для аналитики',  2, 3, 2, 2, 2, NOW() + INTERVAL '2 days'),
(6, 'Настроить алерты',         'PagerDuty для 5xx',      2, 3, 3, 5, 5, NOW() + INTERVAL '3 days'),
(7, 'Фикс мобильной навигации', 'Sidebar collapse',       3, 1, 1, 1, 3, NOW() - INTERVAL '1 day'),
(8, 'Документация API v2',      'Swagger + Postman',      3, 2, 2, 2, 4, NOW() + INTERVAL '6 days'),
(9, 'Исправить утечку памяти',  'WebSocket cleanup',      4, 4, 2, 2, 2, NOW() - INTERVAL '5 days'),
(10,'Онбординг пользователей',  'Пошаговый туториал',     4, 2, 1, 2, 1, NOW() - INTERVAL '12 days')
ON CONFLICT (id) DO NOTHING;
SELECT setval('tasks_id_seq', (SELECT COALESCE(MAX(id),1) FROM tasks));

-- Связь задач ↔️ теги
INSERT INTO task_tags (task_id, tag_id) VALUES
(1, 3), (2, 2), (2, 6), (3, 2), (3, 4), (4, 2), (4, 5),
(5, 6), (5, 2), (6, 3), (6, 2), (7, 1), (7, 5),
(8, 4), (9, 1), (9, 6), (10, 5), (10, 2)
ON CONFLICT ON CONSTRAINT task_tags_pkey DO NOTHING;

COMMIT;