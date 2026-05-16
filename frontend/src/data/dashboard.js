export const tasks = [
  {
    id: "abc123123",
    status: 'В работе',
    title: 'Подготовить backend API',
    description: 'Описать первые маршруты для задач и health-check.',
    priority: 'Высокий',
    deadline: '22-06-2026 14:00',
    tags: ["backend", "frontend"],
  },
  {
    id: "90112abbc",
    status: 'Ожидает',
    title: 'Согласовать модель данных',
    description: 'Зафиксировать поля задачи, пользователя и событий.',
    priority: 'Средний',
    deadline: '22-06-2026 14:00',
    tags: ["general", "users"],
  },
  {
    id: "fabc4567",
    status: 'Готово',
    title: 'Поднять frontend',
    description: 'Собрать Vue-заглушку и раздать ее через Nginx.',
    priority: 'Средний',
    deadline: '22-06-2026 14:00',
    tags: ["vue", "nginx"],
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
