<template>
  <main class="page-shell">
    <section class="workspace">
      <AppHeader
        :commands="commands"
        :selected-command-id="selectedCommandId"
        :current-user="currentUser"
        @update:selected-command-id="selectedCommandId = $event"
        @logout="handleLogout"
      />

      <div class="content-grid">
        <p v-if="isLoading" class="state-message">Загрузка задач...</p>
        <p v-else-if="errorMessage" class="state-message state-message-error">{{ errorMessage }}</p>
        <StatusColumnList v-else :tasks="filteredTasks" :command="selectedCommand"/>
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
import { fetchTasks, fetchTeams } from '../api/tasks'

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
const errorMessage = ref('')

const selectedCommand = computed(() => {
  return commands.value.find((command) => command.id === selectedCommandId.value) ?? commands.value[0] ?? DEFAULT_COMMAND
})

const filteredTasks = computed(() => {
  return tasks.value.filter((task) => task.command.id === selectedCommand.value.id)
})

watch(selectedCommandId, (commandId) => {
  if (commandId) {
    localStorage.setItem(CURRENT_TEAM_STORAGE_KEY, String(commandId))
  }
})

onMounted(loadDashboard)

async function loadDashboard() {
  isLoading.value = true
  errorMessage.value = ''

  try {
    const currentUserResponse = await fetchCurrentUser()
    currentUser.value = currentUserResponse.user

    const loadedCommands = await fetchTeams()
    commands.value = loadedCommands

    if (!commands.value.some((command) => command.id === selectedCommandId.value)) {
      selectedCommandId.value = commands.value[0]?.id ?? 0
    }

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

function readSavedCommandId() {
  const savedId = Number(localStorage.getItem(CURRENT_TEAM_STORAGE_KEY))
  return Number.isFinite(savedId) ? savedId : 0
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
