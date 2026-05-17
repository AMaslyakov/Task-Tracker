<template>
  <main class="page-task">
    <div class="task-page-nav">
      <button type="button" class="back-link-btn" @click="goBack">
        ← Вернуться к доске
      </button>
      <span class="task-page-id-badge">ID: {{ task?.id || $route.params.id }}</span>
    </div>

    <div v-if="isLoading" class="task-page-state">
      <div class="spinner"></div>
      <p>Синхронизация данных с PostgreSQL...</p>
    </div>

    <div v-else-if="error" class="task-page-state error-state">
      <p>❌ {{ error }}</p>
      <button type="button" class="retry-btn" @click="loadTaskDetails">Повторить попытку</button>
    </div>

    <div v-else-if="task" class="task-page-grid">

      <section class="task-main-content">
        <h1 class="task-page-title">{{ task.title }}</h1>

        <div class="task-section">
          <h2 class="section-subtitle">Описание задачи</h2>
          <div class="task-description-box">
            <p v-if="task.description">{{ task.description }}</p>
            <p v-else class="empty-text">Описание для этой задачи отсутствует.</p>
          </div>
        </div>

        <!-- БЛОК 1: Прикрепление и просмотр файлов -->
        <div class="task-section">
          <h2 class="section-subtitle">Вложения и файлы</h2>
          <div class="files-block">
            <label class="file-upload-label">
              <span>📎 Выбрать файлы для загрузки</span>
              <input type="file" multiple class="hidden-file-input" @change="handleFileUpload" />
            </label>

            <div v-if="attachedFiles.length" class="files-list">
              <div v-for="file in attachedFiles" :key="file.id" class="file-item-badge">
                <span class="file-icon">📄</span>
                <div class="file-meta">
                  <span class="file-name">{{ file.name }}</span>
                  <span class="file-size">{{ file.size }}</span>
                </div>
                <button type="button" class="delete-file-btn" @click="removeFile(file.id)">&times;</button>
              </div>
            </div>
            <p v-else class="empty-text">К этой задаче ещё не прикреплено ни одного файла.</p>
          </div>
        </div>

        <!-- БЛОК 2: Комментарии -->
        <div class="task-section">
          <h2 class="section-subtitle">Комментарии ({{ comments.length }})</h2>
          <div class="comments-block">
            <form @submit.prevent="addComment" class="comment-form">
              <textarea
                v-model="newCommentText"
                placeholder="Напишите комментарий или задайте вопрос..."
                rows="3"
                required
                class="cyber-comment-textarea"
              ></textarea>
              <button type="submit" class="send-comment-btn">Отправить</button>
            </form>

            <div v-if="comments.length" class="comments-list">
              <div v-for="comment in comments" :key="comment.id" class="comment-card-item">
                <div class="comment-item-header">
                  <span class="comment-author">{{ comment.author }}</span>
                  <span class="comment-date">{{ comment.date }}</span>
                  <button type="button" class="delete-comment-link" @click="deleteComment(comment.id)">Удалить</button>
                </div>
                <p class="comment-body-text">{{ comment.text }}</p>
              </div>
            </div>
            <p v-else class="empty-text">Здесь пока нет обсуждений. Станьте первым!</p>
          </div>
        </div>
      </section>

      <aside class="task-sidebar-panel">
        <h2 class="sidebar-panel-title">Атрибуты задачи</h2>

        <!-- БЛОК 3: Учёт затраченного времени -->
        <div class="sidebar-widget">
          <label class="widget-label">Затраченное время</label>
          <div class="time-tracking-widget">
            <div class="time-display-counter">⏱️ {{ formattedLoggedTime }}</div>

            <div class="time-controls">
              <button
                type="button"
                class="time-control-btn"
                :class="{ 'stop-active': isTimerRunning }"
                @click="toggleTimer"
              >
                {{ isTimerRunning ? 'Пауза' : 'Запустить таймер' }}
              </button>

              <button type="button" class="time-control-btn manual-btn" @click="openManualTimeInput">
                Вручную
              </button>
            </div>

            <div v-if="showManualTime" class="manual-time-popover">
              <input
                v-model.number="manualMinutes"
                type="number"
                min="1"
                placeholder="Минуты"
                class="manual-time-input"
              />
              <button type="button" class="apply-time-btn" @click="addManualTime">Добавить</button>
            </div>
          </div>
        </div>

        <div class="sidebar-widget">
          <label class="widget-label">Статус доски</label>
          <div class="status-badge-display" :class="statusClass(task.status)">
            {{ task.status }}
          </div>
        </div>

        <div class="sidebar-widget">
          <label class="widget-label">Важность</label>
          <div class="priority-badge-display" :class="priorityClass(task.priority)">
            {{ task.priority }}
          </div>
        </div>

        <div class="sidebar-widget">
          <label class="widget-label">Исполнитель</label>
          <div class="assigned-user-box">
            <div class="user-avatar-placeholder">
              {{ task.asigned_to?.name?.charAt(0).toUpperCase() || '?' }}
            </div>
            <div class="user-info-text">
              <span class="user-display-name">{{ task.asigned_to?.name || 'Не назначен' }}</span>
              <span v-if="task.asigned_to?.email" class="user-display-email">{{ task.asigned_to.email }}</span>
            </div>
          </div>
        </div>

        <div class="sidebar-widget">
          <label class="widget-label">Срок выполнения</label>
          <div class="deadline-box-display">
            📅 {{ task.deadline || 'Без срока' }}
          </div>
        </div>
      </aside>

    </div>
  </main>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { fetchTasks } from '../api/tasks';

const route = useRoute();
const router = useRouter();

const task = ref(null);
const isLoading = ref(true);
const error = ref('');

// Реактивные переменные для новых блоков
const attachedFiles = ref([]);
const comments = ref([]);
const newCommentText = ref('');

const loggedMinutes = ref(0);
const isTimerRunning = ref(false);
const showManualTime = ref(false);
const manualMinutes = ref(null);
let timerInterval = null;

// Форматирование времени в понятный вид: Х ч. ХХ мин.
const formattedLoggedTime = computed(() => {
  const hours = Math.floor(loggedMinutes.value / 60);
  const minutes = loggedMinutes.value % 60;
  if (hours === 0) return `${minutes} мин.`;
  return `${hours} ч. ${minutes < 10 ? '0' + minutes : minutes} мин.`;
});

onMounted(() => {
  loadTaskDetails();
  // Имитируем предзаполненные данные из базы данных для демонстрации интерфейса
  comments.value = [
    { id: 1, author: 'elena_qa', text: 'Проверила эндпоинты, всё работает стабильно.', date: 'Вчера, 16:45' }
  ];
  attachedFiles.value = [
    { id: 1, name: 'api_schema.json', size: '24 KB' }
  ];
  loggedMinutes.value = 45;
});

onUnmounted(() => {
  clearInterval(timerInterval);
});

async function loadTaskDetails() {
  isLoading.value = true;
  error.value = '';
  try {
    const taskId = route.params.id;
    const allTasks = await fetchTasks();
    const foundTask = allTasks.find(t => String(t.id) === String(taskId));
    if (foundTask) {
      task.value = foundTask;
    } else {
      error.value = `Задача с ID ${taskId} не найдена`;
    }
  } catch (err) {
    console.error(err);
    error.value = 'Не удалось загрузить данные задачи';
  } finally {
    isLoading.value = false;
  }
}

// Логика файлов
function handleFileUpload(event) {
  const files = event.target.files;
  for (let i = 0; i < files.length; i++) {
    const sizeInKb = Math.round(files[i].size / 1024);
    attachedFiles.value.push({
      id: Date.now() + i,
      name: files[i].name,
      size: sizeInKb > 1024 ? (sizeInKb / 1024).toFixed(1) + ' MB' : sizeInKb + ' KB'
    });
  }
  // Здесь в будущем будет: await uploadFileToBackend(files)
}

function removeFile(fileId) {
  attachedFiles.value = attachedFiles.value.filter(f => f.id !== fileId);
}

// Логика комментариев
function addComment() {
  if (!newCommentText.value.trim()) return;
  comments.value.unshift({
    id: Date.now(),
    author: 'Вы',
    text: newCommentText.value,
    date: 'Только что'
  });
  newCommentText.value = '';
  // Здесь в будущем будет: await saveCommentToBackend(taskId, text)
}

function deleteComment(commentId) {
  comments.value = comments.value.filter(c => c.id !== commentId);
}

// Логика таймера учета времени
function toggleTimer() {
  if (isTimerRunning.value) {
    clearInterval(timerInterval);
    isTimerRunning.value = false;
    // Здесь будет: await saveTimeProgressToBackend(loggedMinutes.value)
  } else {
    isTimerRunning.value = true;
    // Для симуляции таймер прибавляет 1 минуту каждые несколько секунд на этапе теста,
    // в реальном коде интервал выставляется на 60000 (1 минута)
    timerInterval = setInterval(() => {
      loggedMinutes.value++;
    }, 2000);
  }
}

function openManualTimeInput() {
  showManualTime.value = !showManualTime.value;
}

function addManualTime() {
  if (!manualMinutes.value || manualMinutes.value <= 0) return;
  loggedMinutes.value += manualMinutes.value;
  manualMinutes.value = null;
  showManualTime.value = false;
}

function goBack() {
  router.push('/tasks');
}

function statusClass(status) {
  if (!status) return '';
  return 'status-' + status.toLowerCase().replace(' ', '-');
}

function priorityClass(priority) {
  if (!priority) return '';
  return 'priority-' + priority.toLowerCase();
}
</script>

<style scoped>
.page-task {
  font-family: system-ui, -apple-system, sans-serif;
  min-height: 100vh;
  padding: 32px;
  background-color: #f8fafc;
  box-sizing: border-box;
  width: 100%;
  max-width: 1120px;
  margin: 0 auto;
}

.task-page-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.back-link-btn {
  background: none;
  border: none;
  color: #6366f1;
  font-size: 15px;
  font-weight: 700;
  cursor: pointer;
  padding: 0;
}

.task-page-id-badge {
  font-family: 'Courier New', monospace;
  font-size: 12px;
  font-weight: 800;
  background-color: #e2e8f0;
  color: #475569;
  padding: 4px 10px;
  border-radius: 6px;
}

.task-page-grid {
  display: grid;
  grid-template-columns: 1fr 340px;
  gap: 32px;
  align-items: start;
}

.task-main-content {
  background-color: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 32px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.02);
}

.task-page-title {
  font-size: 32px;
  font-weight: 800;
  color: #0f172a;
  margin: 0 0 24px 0;
}

.task-section {
  margin-top: 28px;
  border-top: 1px solid #f1f5f9;
  padding-top: 24px;
}

.section-subtitle {
  font-size: 13px;
  font-weight: 800;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin: 0 0 14px 0;
}

.task-description-box {
  font-size: 16px;
  color: #334155;
  line-height: 1.6;
}

.empty-text {
  color: #94a3b8;
  font-style: italic;
  font-size: 14px;
  margin: 0;
}

/* Стилизация блока файлов */
.file-upload-label {
  display: inline-block;
  padding: 10px 16px;
  background-color: #f1f5f9;
  border: 1px dashed #cbd5e1;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  color: #475569;
  margin-bottom: 14px;
}

.file-upload-label:hover {
  background-color: #e2e8f0;
}

.hidden-file-input {
  display: none;
}

.files-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 10px;
}

.file-item-badge {
  display: flex;
  align-items: center;
  gap: 10px;
  background-color: #f8fafc;
  border: 1px solid #e2e8f0;
  padding: 8px 12px;
  border-radius: 8px;
  position: relative;
}

.file-meta {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.file-name {
  font-size: 13px;
  font-weight: 600;
  color: #1e293b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-size {
  font-size: 11px;
  color: #94a3b8;
}

.delete-file-btn {
  background: none;
  border: none;
  color: #94a3b8;
  font-size: 18px;
  cursor: pointer;
  margin-left: auto;
}

.delete-file-btn:hover {
  color: #ef4444;
}

/* Стилизация блока комментариев */
.comment-form {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 10px;
  margin-bottom: 20px;
}

.cyber-comment-textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  background-color: #f8fafc;
  font-family: inherit;
  font-size: 14px;
  outline: none;
  box-sizing: border-box;
}

.cyber-comment-textarea:focus {
  border-color: #6366f1;
}

.send-comment-btn {
  padding: 8px 16px;
  background-color: #6366f1;
  color: white;
  border: none;
  border-radius: 6px;
  font-weight: 700;
  cursor: pointer;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.comment-card-item {
  background-color: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 14px;
}

.comment-item-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 6px;
}

.comment-author {
  font-size: 13px;
  font-weight: 700;
  color: #1e293b;
}

.comment-date {
  font-size: 11px;
  color: #94a3b8;
}

.delete-comment-link {
  background: none;
  border: none;
  color: #94a3b8;
  font-size: 11px;
  cursor: pointer;
  margin-left: auto;
}

.delete-comment-link:hover {
  color: #ef4444;
}

.comment-body-text {
  margin: 0;
  font-size: 14px;
  color: #334155;
  line-height: 1.5;
}

/* Правая панель (Сайдбар) */
.task-sidebar-panel {
  background-color: #1e293b;
  border: 1px solid #334155;
  border-radius: 14px;
  padding: 24px;
  box-shadow: 0 10px 25px -5px rgba(168, 85, 247, 0.05);
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.sidebar-panel-title {
  font-size: 13px;
  font-weight: 800;
  color: #c084fc;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  margin: 0;
  border-bottom: 1px solid #334155;
  padding-bottom: 10px;
}

.sidebar-widget {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.widget-label {
  font-size: 11px;
  font-weight: 700;
  color: #94a3b8;
  text-transform: uppercase;
}

/* Виджет времени */
.time-tracking-widget {
  background-color: #334155;
  padding: 14px;
  border-radius: 10px;
  border: 1px solid #475569;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.time-display-counter {
  font-size: 18px;
  font-weight: 800;
  color: #f8fafc;
}

.time-controls {
  display: flex;
  gap: 8px;
}

.time-control-btn {
  flex: 1;
  padding: 6px 10px;
  border-radius: 6px;
  border: none;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  background-color: #22c55e;
  color: white;
}

.time-control-btn.stop-active {
  background-color: #eab308;
}

.time-control-btn.manual-btn {
  background-color: #475569;
  color: #cbd5e1;
}

.manual-time-popover {
  display: flex;
  gap: 8px;
  margin-top: 4px;
}

.manual-time-input {
  width: 100%;
  padding: 6px;
  border: 1px solid #475569;
  border-radius: 4px;
  background-color: #1e293b;
  color: white;
  font-size: 13px;
}

.apply-time-btn {
  padding: 6px 12px;
  background-color: #6366f1;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
}

.status-badge-display {
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 700;
  text-align: center;
  background-color: #334155;
  color: #cbd5e1;
  border: 1px solid #475569;
}

.priority-badge-display {
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 700;
  text-align: center;
}
.priority-badge-display.priority-critical { background-color: #7f1d1d; color: #fca5a5; }
.priority-badge-display.priority-high { background-color: #7c2d12; color: #ffedd5; }
.priority-badge-display.priority-medium { background-color: #713f12; color: #fef08a; }
.priority-badge-display.priority-low { background-color: #064e3b; color: #a7f3d0; }

.assigned-user-box {
  display: flex;
  align-items: center;
  gap: 12px;
  background-color: #334155;
  padding: 10px 14px;
  border-radius: 10px;
  border: 1px solid #475569;
}

.user-avatar-placeholder {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background-color: #6366f1;
  color: #ffffff;
  display: grid;
  place-items: center;
  font-weight: 700;
  font-size: 14px;
}

.user-info-text {
  display: flex;
  flex-direction: column;
}

.user-display-name {
  font-size: 14px;
  font-weight: 700;
  color: #f8fafc;
}

.user-display-email {
  font-size: 11px;
  color: #cbd5e1;
}

.deadline-box-display {
  background-color: #334155;
  padding: 10px 14px;
  border-radius: 10px;
  border: 1px solid #475569;
  color: #f8fafc;
  font-size: 13px;
  font-weight: 600;
}

.task-page-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px;
  background: #ffffff;
  border-radius: 14px;
  border: 1px solid #e2e8f0;
  color: #64748b;
}

.error-state {
  color: #ef4444;
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #cbd5e1;
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 12px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@media (max-width: 900px) {
  .task-page-grid { grid-template-columns: 1fr; }
  .task-sidebar-panel { order: -1; }
}
</style>
