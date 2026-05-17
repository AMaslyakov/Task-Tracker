<template>
  <main class="page-shell">
    <section class="workspace">
      <!-- Шапка приложения с выбором команды и профилем -->
      <AppHeader
        :commands="commands"
        :selected-command-id="selectedCommandId"
        :current-user="currentUser"
        @update:selected-command-id="selectedCommandId = $event"
        @logout="handleLogout"
      />

      <!-- Исходный неизмененный контейнер сетки дашборда -->
      <div class="content-grid">
        <p v-if="isLoading" class="state-message">Загрузка задач...</p>
        <p v-else-if="errorMessage" class="state-message state-message-error">{{ errorMessage }}</p>

        <template v-else>
          <!--
            Блок управления кнопкой. Размещен прямо над доской,
            выровнен по правому краю без сжатия колонок
          -->
          <div class="dashboard-action-bar">
            <button
              class="create-task-button-dashboard"
              type="button"
              @click="openCreateModal"
            >
              Новая задача
            </button>
          </div>

          <!-- Канбан-доска занимает всю изначальную ширину родителя -->
          <StatusColumnList
            v-model:tasks="tasks"
            :tasks="filteredTasks"
            :command="selectedCommand"
            @edit-task="openEditModal"
            @task-status-changed="handleStatusDragAndDrop"
          />
        </template>
      </div>
    </section>

    <!-- Универсальное модальное окно управления задачами -->
    <TaskModal
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
import TaskModal from '../components/TaskModal.vue'
import { fetchCurrentUser, logout } from '../api/auth'
import { fetchTasks, fetchTeams, createTask, updateTask, deleteTask, updateTaskStatus } from '../api/tasks'

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

const isModalOpen = ref(false)
const selectedTask = ref(null)

const selectedCommand = computed(() => {
  const found = commands.value.find((command) => command.id === selectedCommandId.value)
  return found || commands.value[0] || DEFAULT_COMMAND
})

const filteredTasks = computed(() => {
  return tasks.value.filter((task) => {
    const teamId = task.command?.id || task.team_id || task.team?.id;
    return teamId === selectedCommand.value.id;
  })
})

watch(selectedCommandId, (commandId) => {
  if (commandId) {
    localStorage.setItem(CURRENT_TEAM_STORAGE_KEY, String(commandId))
  }
})

onMounted(loadDashboard)

function openCreateModal() {
  selectedTask.value = null;
  isModalOpen.value = true;
}

function openEditModal(task) {
  selectedTask.value = task;
  isModalOpen.value = true;
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
      priority: formData.priority,
      deadline: formData.deadline || null,
      status: formData.status,
      team_id: selectedCommand.value.id,
      created_by: currentUser.value?.id,
      assigned_to_name: formData.assigned_to_name
    };

    const createdTask = await createTask(taskPayload);
    tasks.value.push(createdTask);
    closeModal();
  } catch (error) {
    console.error('Ошибка при создании задачи:', error);
    alert('Не удалось создать задачу на сервере');
  }
}

async function handleUpdateTask(formData) {
  try {
    const updatedTaskFromServer = await updateTask(formData.id, {
      title: formData.title,
      description: formData.description,
      priority: formData.priority,
      deadline: formData.deadline || null,
      assigned_to_name: formData.assigned_to_name
    });

    tasks.value = tasks.value.map(task => {
      if (task.id === formData.id) {
        return {
          ...task,
          ...updatedTaskFromServer
        };
      }
      return task;
    });

    closeModal();
  } catch (error) {
    console.error('Ошибка при обновлении задачи:', error);
    alert('Не удалось сохранить изменения');
  }
}

async function handleDeleteTask(taskId) {
  try {
    await deleteTask(taskId);
    tasks.value = tasks.value.filter(task => task.id !== taskId);
    closeModal();
  } catch (error) {
    console.error('Ошибка при удалении задачи:', error);
    alert('Не удалось удалить задачу с сервера');
  }
}

async function handleStatusDragAndDrop({ taskId, newStatus }) {
  try {
    await updateTaskStatus(taskId, newStatus);
  } catch (error) {
    console.error('Ошибка при сохранении положения карточки:', error);
    tasks.value = await fetchTasks(commands.value);
  }
}

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
  width: min(1120px, 82%);
  margin: 0 auto;
}

/* ОРИГИНАЛЬНЫЙ ДИЗАЙН: Сетка возвращена в исходное состояние без сжатия */
.content-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 320px;
  gap: 20px;
  align-items: start;
}

/* Строка панели действий прижимает кнопку вправо над доской */
.dashboard-action-bar {
  grid-column: 1 / -1; /* Растягивается на всю сетку */
  display: flex;
  justify-content: flex-end; /* Кнопка строго по правую сторону */
  margin-bottom: 4px; /* Небольшой аккуратный зазор до колонок */
}

/* Стилизация кнопки */
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

  .dashboard-action-bar {
    justify-content: fill;
  }

  .create-task-button-dashboard {
    width: 100%;
  }
}
</style>
