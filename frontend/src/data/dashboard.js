const STATUSES = {
  prog: "In progress",
  todo: "TODO",
  review: "Review",
  done: "Done",
}


export const user1 = {
  id: 2333,
  name: "Vasya",
  email: "sdfsdf@sdf.ru",
  is_admin: false
} 

export const user2 = {
  id: 3,
  name: "petya",
  email: "op@op.ru",
  is_admin: true
}



export const user3 = {
  id: 3,
  name: "user user",
  email: "user@le.ru",
  is_admin: false
}

export const command1 = {
  id: 123132,
  name: "Команда A",
  members: [user1, user2, user3],
  config_dashboard: {
    statuses: [STATUSES.todo, STATUSES.prog, STATUSES.review, STATUSES.done]
  }
}

export const command2 = {
  id: 984221,
  name: "Команда B",
  members: [user1, user3],
  config_dashboard: {
    statuses: [STATUSES.todo, STATUSES.prog, STATUSES.review, STATUSES.done]
  }
}

export const command3 = {
  id: 564322,
  name: "Команда C",
  members: [user2],
  config_dashboard: {
    statuses: [STATUSES.todo, STATUSES.prog, STATUSES.review, STATUSES.done]
  }
}

export const commands = [command1, command2, command3]

export const tasks = [
  {
    id: "tsk001aa",
    status: STATUSES.todo,
    title: 'Настроить авторизацию',
    description: 'Реализовать JWT-авторизацию и refresh token.',

    priority: 'Высокий',
    deadline: '18-06-2026 10:00',
    asigned_to: user1,
    command: command1,
    tags: ["backend", "auth"],
  },
  {
    id: "tsk002bb",
    status: STATUSES.todo,
    title: 'Создать базу данных',
    description: 'Поднять PostgreSQL и настроить миграции.',
    priority: 'Высокий',
    deadline: '19-06-2026 12:00',
    asigned_to: user2,
    command: command1,
    tags: ["database", "backend"],
  },
  {
    id: "tsk003cc",
    status: STATUSES.prog,
    title: 'Сверстать страницу логина',
    description: 'Добавить форму входа и валидацию полей.',
    priority: 'Средний',
    deadline: '20-06-2026 16:00',
    asigned_to: user3,
    command: command2,
    tags: ["frontend", "vue"],
  },
  {
    id: "tsk004dd",
    status: STATUSES.done,
    title: 'Настроить Docker',
    description: 'Создать docker-compose для frontend и backend.',
    priority: 'Средний',
    deadline: '17-06-2026 11:00',
    asigned_to: user1,
    command: command2,
    tags: ["docker", "devops"],
  },
  {
    id: "tsk005ee",
    status: 'Review',
    title: 'Подготовить Swagger',
    description: 'Описать документацию для API маршрутов.',
    priority: 'Средний',
    deadline: '21-06-2026 15:00',
    asigned_to: user2,
    command: command1,
    tags: ["backend", "docs"],
  },
  {
    id: "tsk006ff",
    status: STATUSES.todo,
    title: 'Добавить CI/CD',
    description: 'Настроить GitHub Actions для автодеплоя.',
    priority: 'Высокий',
    deadline: '24-06-2026 18:00',
    asigned_to: user3,
    command: command3,
    tags: ["ci", "github"],
  },
  {
    id: "tsk007gg",
    status: STATUSES.prog,
    title: 'Настроить Nginx',
    description: 'Добавить reverse proxy для frontend и API.',
    priority: 'Средний',
    deadline: '22-06-2026 13:00',
    asigned_to: user1,
    command: command1,
    tags: ["nginx", "server"],
  },
  {
    id: "tsk008hh",
    status: STATUSES.todo,
    title: 'Создать страницу задач',
    description: 'Вывести список задач с фильтрацией.',
    priority: 'Высокий',
    deadline: '23-06-2026 17:00',
    asigned_to: user2,
    command: command2,
    tags: ["frontend", "tasks"],
  },
  {
    id: "tsk009ii",
    status: STATUSES.done,
    title: 'Подключить ESLint',
    description: 'Настроить правила линтинга для проекта.',
    priority: 'Низкий',
    deadline: '16-06-2026 09:00',
    asigned_to: user3,
    command: command3,
    tags: ["eslint", "frontend"],
  },
  {
    id: "tsk010jj",
    status: 'Review',
    title: 'Проверить адаптивность',
    description: 'Протестировать интерфейс на мобильных устройствах.',
    priority: 'Средний',
    deadline: '25-06-2026 14:00',
    asigned_to: user1,
    command: command2,
    tags: ["ui", "mobile"],
  },
  {
    id: "tsk011kk",
    status: STATUSES.todo,
    title: 'Добавить уведомления',
    description: 'Реализовать toast-уведомления для действий пользователя.',
    priority: 'Средний',
    deadline: '26-06-2026 12:00',

    asigned_to: user2,
    command: command1,
    tags: ["frontend", "ux"],
  },
  {
    id: "tsk012ll",
    status: STATUSES.prog,
    title: 'Реализовать поиск',
    description: 'Добавить поиск задач по названию и тегам.',
    priority: 'Высокий',
    deadline: '27-06-2026 11:00',

    asigned_to: user3,
    command: command2,
    tags: ["search", "tasks"],
  },
  {
    id: "tsk013mm",
    status: STATUSES.done,
    title: 'Настроить логирование',
    description: 'Добавить Winston для логирования backend-сервиса.',
    priority: 'Средний',
    deadline: '18-06-2026 08:00',

    asigned_to: user1,
    command: command3,
    tags: ["backend", "logs"],
  },
  {
    id: "tsk014nn",
    status: 'Review',
    title: 'Подготовить релиз',
    description: 'Проверить сборку и подготовить changelog.',
    priority: 'Высокий',
    deadline: '30-06-2026 19:00',

    asigned_to: user2,
    command: command3,
    tags: ["release", "deploy"],
  }
]

export const summaryItems = [
  {
    label: 'В работе',
    value: 7
  },
  {
    label: 'Завершено',
    value: 12
  },
  {
    label: 'Просрочено',
    value: 1
  }
]

