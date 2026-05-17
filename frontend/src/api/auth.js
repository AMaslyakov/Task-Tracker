const API_BASE = '/api'

async function request(path, options = {}) {
  const response = await fetch(`${API_BASE}${path}`, {
    credentials: 'same-origin',
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers
    }
  })

  if (!response.ok) {
    let message = 'Ошибка авторизации'

    try {
      const payload = await response.json()
      if (payload?.error) {
        message = payload.error
      }
    } catch {
      message = `Ошибка запроса: ${response.status}`
    }

    const error = new Error(message)
    error.status = response.status
    throw error
  }

  if (response.status === 204) {
    return null
  }

  return response.json()
}

export async function login(email, password) {
  return request('/login', {
    method: 'POST',
    body: JSON.stringify({ email, password })
  })
}

export async function fetchCurrentUser() {
  return request('/me')
}

export async function logout() {
  return request('/logout', {
    method: 'POST'
  })
}
