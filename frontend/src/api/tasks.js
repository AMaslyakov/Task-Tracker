const API_BASE = '/api'

const FALLBACK_STATUSES = ['TODO', 'IN PROGRESS', 'REVIEW', 'DONE']

// Ваша оригинальная базовая функция запросов
async function request(path, options = {}) {
  const response = await fetch(`${API_BASE}${path}`, {
    credentials: 'same-origin',
    ...options
  })

  if (!response.ok) {
    throw new Error(`API request failed: ${response.status}`)
  }

  return response.json()
}

// Ваша оригинальная функция получения команд
export async function fetchTeams() {
  const teams = await request('/teams')

  return teams.map((team) => ({
    ...team,
    config_dashboard: {
      statuses: Array.isArray(team.config_dashboard) && team.config_dashboard.length > 0
        ? team.config_dashboard
        : FALLBACK_STATUSES
    }
  }))
}

// Ваша оригинальная функция получения задач
export async function fetchTasks(teams = []) {
  const tasks = await request('/tasks')

  return tasks.map((task) => {
    const command = teams.find((team) => team.id === task.team_id)

    return {
      id: task.id,
      status: task.status_name, // Маппинг статуса из БД
      title: task.title,
      description: task.description,
      priority: task.priority_name, // Маппинг приоритета из БД
      deadline: formatDate(task.deadline),
      asigned_to: {
        name: task.assigned_to || 'Не назначен',
        email: ''
      },
      command: command ?? {
        id: task.team_id,
        name: task.team_name,
        config_dashboard: {
          statuses: FALLBACK_STATUSES
        }
      },
      tags: []
    }
  })
}

// ДОБАВЛЕНО 1: Функция создания задачи (заглушка, использующая ваш формат)
export async function createTask(taskPayload) {
  console.log('API -> Запрос на создание задачи отправлен:', taskPayload);

  // Когда бэкенд-эндпоинт POST /api/tasks будет готов, код заменится на:
  // return await request('/tasks', {
  //   method: 'POST',
  //   headers: { 'Content-Type': 'application/json' },
  //   body: JSON.stringify(taskPayload)
  // });

  // Возвращаем mock-объект, адаптированный под маппинг вашей доски:
  return {
    id: Math.floor(Math.random() * 10000),
    title: taskPayload.title,
    description: taskPayload.description,
    priority: taskPayload.priority,
    deadline: formatDate(taskPayload.deadline),
    status: taskPayload.status, // 'TODO'
    team_id: taskPayload.team_id,
    command: { id: taskPayload.team_id },
    asigned_to: { name: taskPayload.assigned_to_name || 'Не назначен', email: '' },
    tags: []
  };
}

// ДОБАВЛЕНО 2: Функция обновления задачи (заглушка)
export async function updateTask(taskId, updateData) {
  console.log(`API -> Запрос на обновление задачи ${taskId} отправлен:`, updateData);

  // Когда эндпоинт PUT /api/tasks/:id будет готов:
  // return await request(`/tasks/${taskId}`, {
  //   method: 'PUT',
  //   headers: { 'Content-Type': 'application/json' },
  //   body: JSON.stringify(updateData)
  // });

  return {
    id: taskId,
    title: updateData.title,
    description: updateData.description,
    priority: updateData.priority,
    deadline: formatDate(updateData.deadline),
    asigned_to: { name: updateData.assigned_to_name || 'Не назначен', email: '' }
  };
}

// ДОБАВЛЕНО 3: Функция удаления задачи (заглушка)
export async function deleteTask(taskId) {
  console.log(`API -> Запрос на удаление задачи ${taskId} отправлен`);

  // Когда эндпоинт DELETE /api/tasks/:id будет готов:
  // return await request(`/tasks/${taskId}`, { method: 'DELETE' });

  return { success: true };
}

// ДОБАВЛЕНО 4: Функция обновления статуса при перетаскивании карточек (заглушка)
export async function updateTaskStatus(taskId, newStatus) {
  console.log(`API -> Перетаскивание: задача ${taskId} сменила статус на ${newStatus}`);

  // Когда эндпоинт PATCH /api/tasks/:id/status будет готов:
  // return await request(`/tasks/${taskId}/status`, {
  //   method: 'PATCH',
  //   headers: { 'Content-Type': 'application/json' },
  //   body: JSON.stringify({ status: newStatus })
  // });

  return { success: true };
}

// Ваша оригинальная функция форматирования даты
function formatDate(value) {
  if (!value) {
    return 'Без срока'
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }

  return new Intl.DateTimeFormat('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

