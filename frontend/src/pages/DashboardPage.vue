<template>
  <main class="page-shell">
    <section class="workspace">
      <AppHeader
        :commands="commands"
        :selected-command-id="selectedCommandId"
        :current-user="currentUser"
        @update:selected-command-id="handleSelectedCommandChange"
        @logout="handleLogout"
      />

      <div class="content-grid">
        <p v-if="isLoading" class="state-message">Загрузка задач...</p>
        <p v-else-if="errorMessage" class="state-message state-message-error">{{ errorMessage }}</p>
        <template v-else>
          <p v-if="statusUpdateError" class="state-message state-message-error">{{ statusUpdateError }}</p>
          <StatusColumnList
            :tasks="filteredTasks"
            :command="selectedCommand"
            :is-status-updating="isStatusUpdating"
            @task-status-changed="handleTaskStatusChanged"
          />
        </template>
      </div>
    </section>
  </main>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import AppHeader from '../components/AppHeader.vue'
import StatusColumnList from '../components/StatusColumnList.vue';
import { fetchCurrentUser, logout } from '../api/auth'
import { fetchTasks, fetchTeams, mapTask, updateTaskStatus } from '../api/tasks'

const CURRENT_TEAM_STORAGE_KEY = 'task-tracker-current-team-id'
const DEFAULT_COMMAND = {
  id: 0,
  name: 'Команда',
  config_dashboard: {
    statuses: ['TODO', 'IN PROGRESS', 'REVIEW', 'DONE']
  }
}

const router = useRouter()
const commands = ref([])
const tasks = ref([])
const currentUser = ref(null)
const selectedCommandId = ref(readSavedCommandId())
const isLoading = ref(true)
const isStatusUpdating = ref(false)
const errorMessage = ref('')
const statusUpdateError = ref('')

const selectedCommand = computed(() => {
  return commands.value.find((command) => isSameCommandId(command.id, selectedCommandId.value)) ?? commands.value[0] ?? DEFAULT_COMMAND
})

const filteredTasks = computed(() => {
  return tasks.value.filter((task) => isSameCommandId(task.command.id, selectedCommand.value.id))
})

watch(selectedCommandId, saveCommandId)

onMounted(loadDashboard)

async function loadDashboard() {
  isLoading.value = true
  errorMessage.value = ''

  try {
    const currentUserResponse = await fetchCurrentUser()
    currentUser.value = currentUserResponse.user

    const loadedCommands = await fetchTeams()
    commands.value = loadedCommands

    if (!commands.value.some((command) => isSameCommandId(command.id, selectedCommandId.value))) {
      selectedCommandId.value = commands.value[0]?.id ?? 0
    }

    saveCommandId(selectedCommandId.value)

    tasks.value = await fetchTasks(commands.value)
  } catch (error) {
    console.error(error)
    if (error.status === 401) {
      router.push('/login')
      return
    }

    errorMessage.value = 'Не удалось загрузить задачи из backend'
  } finally {
    isLoading.value = false
  }
}

async function handleLogout() {
  try {
    await logout()
  } catch (error) {
    console.error(error)
  } finally {
    router.push('/login')
  }
}

function handleSelectedCommandChange(commandId) {
  selectedCommandId.value = commandId
  saveCommandId(commandId)
  statusUpdateError.value = ''
}

async function handleTaskStatusChanged({ taskId, oldStatus, newStatus }) {
  if (isStatusUpdating.value || isSameCommandId(oldStatus, newStatus)) {
    return
  }

  isStatusUpdating.value = true
  statusUpdateError.value = ''
  setTaskStatus(taskId, newStatus)

  try {
    const updatedTask = await updateTaskStatus(taskId, newStatus)
    replaceTask(mapTask(updatedTask, commands.value))
  } catch (error) {
    console.error(error)
    setTaskStatus(taskId, oldStatus)
    if (error.status === 401) {
      router.push('/login')
      return
    }

    statusUpdateError.value = 'Не удалось сохранить статус задачи'
  } finally {
    isStatusUpdating.value = false
  }
}

function readSavedCommandId() {
  const savedId = Number(localStorage.getItem(CURRENT_TEAM_STORAGE_KEY))
  return Number.isFinite(savedId) ? savedId : 0
}

function saveCommandId(commandId) {
  if (commandId) {
    localStorage.setItem(CURRENT_TEAM_STORAGE_KEY, String(commandId))
  } else {
    localStorage.removeItem(CURRENT_TEAM_STORAGE_KEY)
  }
}

function isSameCommandId(leftId, rightId) {
  return String(leftId) === String(rightId)
}

function setTaskStatus(taskId, status) {
  tasks.value = tasks.value.map((task) => {
    if (task.id === taskId) {
      return { ...task, status }
    }

    return task
  })
}

function replaceTask(updatedTask) {
  tasks.value = tasks.value.map((task) => {
    if (task.id === updatedTask.id) {
      return updatedTask
    }

    return task
  })
}
</script>

<style scoped>
.page-shell {
  min-height: 100vh;
  padding: 32px;
}

.workspace {
  width: min(1120px, 100%);
  margin: 0 auto;
}

.content-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 320px;
  gap: 20px;
  align-items: start;
}

.state-message {
  grid-column: 1 / -1;
  margin: 0;
  color: #475569;
  font-size: 18px;
}

.state-message-error {
  color: #b91c1c;
}

@media (max-width: 760px) {
  .page-shell {
    padding: 20px;
  }

  .content-grid {
    grid-template-columns: 1fr;
  }
}
</style>
