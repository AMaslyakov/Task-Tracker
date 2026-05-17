const API_BASE = '/api'

const FALLBACK_STATUSES = ['TODO', 'IN PROGRESS', 'REVIEW', 'DONE']

async function request(path) {
  const response = await fetch(`${API_BASE}${path}`)

  if (!response.ok) {
    throw new Error(`API request failed: ${response.status}`)
  }

  return response.json()
}

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

export async function fetchTasks(teams = []) {
  const tasks = await request('/tasks')

  return tasks.map((task) => {
    const command = teams.find((team) => team.id === task.team_id)

    return {
      id: task.id,
      status: task.status_name,
      title: task.title,
      description: task.description,
      priority: task.priority_name,
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
