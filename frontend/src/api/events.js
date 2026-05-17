const EVENTS_URL = '/api/events'

const TASK_EVENT_TYPES = [
  'task.created',
  'task.updated',
  'task.status_changed',
  'task.deleted'
]

export function subscribeToTaskEvents(onTaskEvent, onError) {
  const source = new EventSource(EVENTS_URL)

  TASK_EVENT_TYPES.forEach((eventType) => {
    source.addEventListener(eventType, (message) => {
      onTaskEvent(parseEventData(message.data))
    })
  })

  source.onerror = (error) => {
    if (typeof onError === 'function') {
      onError(error)
    }
  }

  return () => {
    source.close()
  }
}

function parseEventData(data) {
  try {
    return JSON.parse(data)
  } catch {
    return null
  }
}
