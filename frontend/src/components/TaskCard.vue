

<template>
  
  <article class="task-card" @click="emit('task-click', props.task)">

    <div class="task-header">
      <div class="task-meta-line">
        <div>
          <p class="task-status">{{ task.status }}</p>
          <p class="task-id">id: {{ task.id }}</p>
        </div>

        
        
        <div class="task-actions-menu" @click.stop>
          <button
            type="button"
            class="pencil-btn"
            @click="toggleMenu"
            aria-label="Действия с задачей"
            title="Действия с задачей"
          >
            <svg viewBox="0 0 20 20" aria-hidden="true" focusable="false">
              <path d="M13.7 2.6a1.7 1.7 0 0 1 2.4 0l1.3 1.3a1.7 1.7 0 0 1 0 2.4L7.2 16.5 3 17l.5-4.2L13.7 2.6Z" />
              <path d="m12.5 4 3.5 3.5" />
            </svg>
          </button>

          
          <div v-if="isMenuOpen" class="actions-dropdown" ref="menuRef">
            <button type="button" class="dropdown-item" @click="triggerEdit">
              Редактировать
            </button>
            <button type="button" class="dropdown-item delete-item" @click="triggerDelete">
              Удалить задачу
            </button>
          </div>
        </div>
      </div>

      
      <h3 class="task-title">{{ task.title }}</h3>
    </div>

    
    <p class="task-desc">{{ task.description }}</p>

    
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

const emit = defineEmits(['task-click', 'task-edit', 'task-delete'])

const isMenuOpen = ref(false)
const menuRef = ref(null)

function toggleMenu() {
  isMenuOpen.value = !isMenuOpen.value
}

function triggerEdit() {
  isMenuOpen.value = false
  emit('task-edit', props.task)
}

function triggerDelete() {
  isMenuOpen.value = false
  if (confirm('Вы уверены, что хотите безвозвратно удалить эту задачу?')) {
    emit('task-delete', props.task.id)
  }
}

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
  position: relative; 
  cursor: pointer;
}

.task-card + .task-card {
  margin-top: 14px;
}

.task-header {
  width: 100%;
}


.task-meta-line {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  width: 100%;
}


.task-actions-menu {
  position: relative;
  display: inline-block;
}


.pencil-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background-color: #f8fafc;
  border: 1px solid #cbd5e1;
  color: #334155;
  cursor: pointer;
  padding: 0;
  border-radius: 8px;
  transition: background-color 0.2s, border-color 0.2s, color 0.2s;
  flex: 0 0 32px;
}

.pencil-btn svg {
  width: 16px;
  height: 16px;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.8;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.pencil-btn:hover,
.pencil-btn:focus-visible {
  background-color: #e2e8f0;
  border-color: #94a3b8;
  color: #0f172a;
  outline: none;
}


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
