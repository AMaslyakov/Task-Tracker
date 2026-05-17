<template>
  <main class="page-shell">
    <section class="workspace">
      <!-- Шапка приложения с выбором команды -->
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

          <!-- Блок кнопок действий над доской справа -->
          <div class="dashboard-action-bar">
            <button
              v-if="currentUser?.is_admin === true"
              class="admin-panel-link-btn-dashboard"
              type="button"
              @click="router.push('/admin')"
            >
              ⚙️ Создать команду
            </button>

            <button
              class="create-task-button-dashboard"
              type="button"
              @click="openCreateModal"
            >
              Новая задача
            </button>
          </div>

          <!-- Канбан-доска -->
          <StatusColumnList
            v-model:tasks="tasks"
            :tasks="filteredTasks"
            :command="selectedCommand"
            :is-status-updating="isStatusUpdating"
            @edit-task="openEditModal"
            @delete-task="handleDeleteTask"
            @task-status-changed="handleTaskStatusChanged"
          />
        </template>
      </div>
    </section>

    <!-- Модальное окно создания задач -->
    <TaskForm
      v-if="isModalOpen"
      :task="selectedTask"
      :command="selectedCommand"
      @close="closeModal"
      @create="handleCreateTask"
      @update="handleUpdateTask"
      @delete="handleDeleteTask"
    />
  </main>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import AppHeader from '../components/AppHeader.vue'
import StatusColumnList from '../components/StatusColumnList.vue'
import TaskForm from '../components/TaskForm.vue'
import { fetchCurrentUser, logout } from '../api/auth'
import { fetchTasks, fetchTeams, createTask, updateTask, mapTask, deleteTask, updateTaskStatus, priorityNameToId } from '../api/tasks'

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

const isModalOpen = ref(false)
const selectedTask = ref(null)

const selectedCommand = computed(() => {
  const found = commands.value.find((command) => isSameCommandId(command.id, selectedCommandId.value))
  return found || commands.value || DEFAULT_COMMAND
})

const filteredTasks = computed(() => {
  return tasks.value.filter((task) => {
    const teamId = task.command?.id || task.team_id || task.team?.id;
    return isSameCommandId(teamId, selectedCommand.value.id);
  })
})

watch(selectedCommandId, saveCommandId)

onMounted(loadDashboard)

function openCreateModal() {
  selectedTask.value = null;
  isModalOpen.value = true;
}

function openEditModal(task) {
  router.push(`/task/${task.id}`)
}

function closeModal() {
  isModalOpen.value = false;
  selectedTask.value = null;
}

async function handleCreateTask(formData) {
  try {
    const taskPayload = {
      title: formData.title,
      description: formData.description,
      priority_id: priorityNameToId(formData.priority),
      deadline: toApiDeadline(formData.deadline),
      status_name: formData.status,
      team_id: selectedCommand.value.id,
      created_by: currentUser.value?.id,
      assigned_to: formData.assigned_to
    };

    const createdTask = await createTask(taskPayload);
    tasks.value.push(mapTask(createdTask, commands.value));
    closeModal();
  } catch (error) {
    console.error(error);
    alert('Не удалось создать задачу на сервере');
  }
}

async function handleUpdateTask(formData) {
  try {
    const updatedTaskFromServer = await updateTask(formData.id, {
      title: formData.title,
      description: formData.description,
      priority_id: priorityNameToId(formData.priority),
      deadline: toApiDeadline(formData.deadline),
      assigned_to: formData.assigned_to
    });

    replaceTask(mapTask(updatedTaskFromServer, commands.value));
    closeModal();
  } catch (error) {
    console.error(error);
    alert('Не удалось сохранить изменения');
  }
}

async function handleDeleteTask(taskId) {
  try {
    await deleteTask(taskId);
    tasks.value = tasks.value.filter(task => !isSameCommandId(task.id, taskId));
    closeModal();
  } catch (error) {
    console.error(error);
    alert('Не удалось удалить задачу с сервера');
  }
}

async function loadDashboard() {
  isLoading.value = true
  errorMessage.value = ''

  try {
    // В момент входа запрашиваем данные текущей сессии пользователя
    const currentUserResponse = await fetchCurrentUser()
    currentUser.value = currentUserResponse.user // Сюда падает объект пользователя с флагом is_admin

    const loadedCommands = await fetchTeams()
    commands.value = loadedCommands

    if (!commands.value.some((command) => isSameCommandId(command.id, selectedCommandId.value))) {
      selectedCommandId.value = commands.value?.id ?? 0
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
  if (leftId == null || rightId == null) return false
  return String(leftId).trim() === String(rightId).trim()
}

function setTaskStatus(taskId, status) {
  tasks.value = tasks.value.map((task) => {
    if (isSameCommandId(task.id, taskId)) {
      return { ...task, status }
    }
    return task
  })
}

function replaceTask(updatedTask) {
  tasks.value = tasks.value.map((task) => {
    if (isSameCommandId(task.id, updatedTask.id)) {
      return updatedTask
    }
    return task
  })
}

function toApiDeadline(value) {
  if (!value) return null
  return new Date(`${value}T00:00:00`).toISOString()
}
</script>

<style scoped>
.page-shell {
  min-height: 100vh;
  padding: 32px;
  background-color: #f8fafc;
}

.workspace {
  width: min(1120px, 82%);
  margin: 0 auto;
}

.content-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 20px;
  align-items: start;
}

/* Контейнер кнопок управления */
.dashboard-action-bar {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}


.admin-panel-link-btn-dashboard {
  min-height: 42px;
  border: 1px solid #cbd5e1;
  border-radius: 10px;
  padding: 0 20px;
  background-color: #ffffff;
  color: #334155;
  font-family: system-ui, -apple-system, sans-serif;
  font-size: 14px;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 1px 2px rgba(0,0,0,0.05);
  transition: all 0.2s ease;
}

.admin-panel-link-btn-dashboard:hover {
  border-color: #a855f7;
  color: #a855f7;
  box-shadow: 0 4px 12px rgba(168, 85, 247, 0.1);
}

.create-task-button-dashboard {
  min-height: 42px;
  border: none;
  border-radius: 10px;
  padding: 0 24px;
  background: linear-gradient(90deg, #a855f7 0%, #ec4899 100%);
  color: #ffffff;
  font-family: system-ui, -apple-system, sans-serif;
  font-size: 14px;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.02em;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(168, 85, 247, 0.25);
  transition: opacity 0.2s, transform 0.1s, box-shadow 0.2s;
}

.create-task-button-dashboard:hover {
  opacity: 0.95;
  box-shadow: 0 4px 16px rgba(168, 85, 247, 0.45);
}

.create-task-button-dashboard:active {
  transform: scale(0.98);
}

.state-message {
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
  .dashboard-action-bar {
    flex-direction: column;
    width: 100%;
  }
  .admin-panel-link-btn-dashboard, .create-task-button-dashboard {
    width: 100%;
  }
}
</style>
