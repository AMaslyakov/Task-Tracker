<!-- ID, название, описание, статус, приоритет, теги, дата создания, дедлайн -->

<template>
  <!-- Клик по карточке открывает форму редактирования, как и раньше -->
  <article class="task-card" @click="emit('task-click', props.task)">

    <div class="task-header">
      <div class="task-meta-line">
        <div>
          <p class="task-status">{{ task.status }}</p>
          <p class="task-id">id: {{ task.id }}</p>
        </div>

        <!-- ИСПРАВЛЕНО: Кнопка-карандашик с контекстным меню -->
        <!-- @click.stop запрещает клику проваливаться в саму карточку -->
        <div class="task-actions-menu" @click.stop>
          <button
            type="button"
            class="pencil-btn"
            @click="toggleMenu"
            title="Действия с задачей"
          >
            ✏️
          </button>

          <!-- Всплывающее меню действий -->
          <div v-if="isMenuOpen" class="actions-dropdown" ref="menuRef">
            <button type="button" class="dropdown-item" @click="triggerEdit">
              📝 Редактировать
            </button>
            <button type="button" class="dropdown-item delete-item" @click="triggerDelete">
              🗑️ Удалить задачу
            </button>
          </div>
        </div>
      </div>

      <!-- Заголовок задачи с адаптивным размером -->
      <h3 class="task-title">{{ task.title }}</h3>
    </div>

    <!-- Описание задачи -->
    <p class="task-desc">{{ task.description }}</p>

    <!-- Данные исполнителя -->
    <div class="task-assigned">
      <div class="task-assigned-name">
        Исполнитель: {{ task.assigned_to?.name || task.asigned_to?.name || 'Не назначен' }}
      </div>
      <div class="task-assigned-email">
        {{ task.assigned_to?.email || task.asigned_to?.email }}
      </div>
    </div>

    <div class="task-priority">
      <span>Приоритет: {{ task.priority }}</span>
    </div>

    <div class="task-deadline">
      <span>Выполнить до: {{ task.deadline }}</span>
    </div>

    <!-- Блок тегов -->
    <div v-if="task.tags && task.tags.length" class="task-tags">
        <span v-for="tag in task.tags" :key="tag" class="tag-badge">{{ tag }}</span>
    </div>
  </article>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';

const props = defineProps({
  task: {
    type: Object,
    required: true
  }
})

// ИСПРАВЛЕНО: Добавлены новые типы событий, чтобы другие разработчики могли их легко поймать
const emit = defineEmits(['task-click', 'task-edit', 'task-delete'])

const isMenuOpen = ref(false)
const menuRef = ref(null)

// Показать/скрыть меню
function toggleMenu() {
  isMenuOpen.value = !isMenuOpen.value
}

// Передаем команду "Редактировать"
function triggerEdit() {
  isMenuOpen.value = false
  emit('task-click', props.task) // Вызывает открытие стандартной формы заполнения
}

// Передаем команду "Удалить" родителю с подтверждением
function triggerDelete() {
  isMenuOpen.value = false
  if (confirm('Вы уверены, что хотите безвозвратно удалить эту задачу?')) {
    // На главной странице сработает функция handleDeleteTask(taskId)
    // Другие разработчики смогут привязать сюда реальный запрос DELETE
    const dashboardElement = document.querySelector('.page-shell');
    if (dashboardElement) {
      // Имитируем клик по удалению из модалки для моментальной связи с DashboardPage.vue
      emit('task-click', props.task);
      setTimeout(() => {
        const modalDeleteBtn = document.querySelector('.delete-btn');
        if (modalDeleteBtn) modalDeleteBtn.click();
      }, 50);
    }
  }
}

// Закрытие выпадающего меню при клике в любое другое место экрана
function handleClickOutside(event) {
  if (menuRef.value && !menuRef.value.contains(event.target)) {
    isMenuOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.task-card {
  display: grid;
  gap: 12px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 14px;
  background-color: #ffffff;
  width: 100%;
  box-sizing: border-box;
  overflow: hidden;
  container-type: inline-size;
  position: relative; /* Важно для позиционирования меню */
  cursor: pointer;
}

.task-card + .task-card {
  margin-top: 14px;
}

.task-header {
  width: 100%;
}

/* Флекс-линия для ID и кнопки-карандашика */
.task-meta-line {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  width: 100%;
}

/* ИСПРАВЛЕНО: Контейнер для выпадающего меню действий */
.task-actions-menu {
  position: relative;
  display: inline-block;
}

/* Маленькая стильная кнопка-карандаш */
.pencil-btn {
  background: none;
  border: none;
  font-size: 13px;
  cursor: pointer;
  padding: 4px 6px;
  border-radius: 4px;
  transition: background-color 0.2s;
  opacity: 0.6;
}

.pencil-btn:hover {
  background-color: #f1f5f9;
  opacity: 1;
}

/* Всплывающее контекстное меню */
.actions-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  background-color: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  z-index: 10;
  min-width: 160px;
  padding: 4px 0;
  margin-top: 4px;
}

/* Пункты меню действий */
.dropdown-item {
  width: 100%;
  border: none;
  background: none;
  text-align: left;
  padding: 8px 12px;
  font-family: system-ui, sans-serif;
  font-size: 13px;
  font-weight: 600;
  color: #334155;
  cursor: pointer;
  transition: background-color 0.15s;
}

.dropdown-item:hover {
  background-color: #f1f5f9;
}

/* Красный пункт удаления */
.dropdown-item.delete-item {
  color: #f43f5e;
  border-top: 1px solid #f1f5f9;
}

.dropdown-item.delete-item:hover {
  background-color: #fff1f2;
}

.task-title {
  margin: 6px 0 0 0;
  font-weight: 700;
  color: #1e293b;
  width: 100%;
  font-size: clamp(13px, 10cqw, 18px);
  white-space: normal;
  word-break: break-word;
}

.task-desc {
  margin: 0;
  color: #475569;
  font-size: 13px;
  line-height: 1.5;
  width: 100%;
  white-space: normal;
  word-break: break-word;
}

.task-status {
  margin: 0;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  color: #64748b;
  letter-spacing: 0.05em;
}

.task-id {
  margin: 2px 0 0 0;
  font-size: 11px;
  color: #94a3b8;
}

.task-assigned {
  font-size: 13px;
  color: #475569;
  width: 100%;
}

.task-assigned-name {
  font-weight: 600;
  color: #334155;
}

.task-assigned-email {
  color: #64748b;
  width: 100%;
  font-size: clamp(10px, 7.5cqw, 12px);
  word-break: break-all;
}

.task-priority {
  font-size: 13px;
  font-weight: 600;
  color: #475569;
}

.task-deadline {
  font-size: 12px;
  color: #64748b;
}

.task-tags {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
  margin-top: 4px;
  width: 100%;
}

.tag-badge {
  font-family: 'Courier New', Courier, monospace;
  font-size: 11px;
  font-weight: 700;
  color: #ffffff;
  padding: 4px 8px;
  background-color: #4f46e5;
  border-radius: 6px;
  white-space: nowrap;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  box-sizing: border-box;
}
</style>
