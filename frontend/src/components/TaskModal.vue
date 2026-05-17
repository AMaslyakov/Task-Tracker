<template>
  
  <div class="modal-overlay" @click.self="emit('close')">

    
    <div class="modal-card">

      
      <div class="modal-header">
        <h2 class="modal-title">
          {{ isEditMode ? 'Редактирование задачи' : 'Создание задачи' }}
        </h2>
        <button type="button" class="close-btn" @click="emit('close')">&times;</button>
      </div>

      
      <form @submit.prevent="handleSubmit" class="modal-form">

        
        <label class="form-field">
          <span class="field-label">Название задачи *</span>
          <input
            v-model="form.title"
            type="text"
            placeholder="Введите название..."
            required
            class="cyber-input"
          />
        </label>

        
        <label class="form-field">
          <span class="field-label">Описание</span>
          <textarea
            v-model="form.description"
            placeholder="Добавьте детали задачи..."
            rows="3"
            class="cyber-textarea"
          ></textarea>
        </label>

        
        <label class="form-field">
          <span class="field-label">Исполнитель</span>
          <select v-model="form.assigned_to_name" class="cyber-select">
            <option :value="null">Не назначен</option>
            <option
              v-for="memberName in props.command?.members || []"
              :key="memberName"
              :value="memberName"
            >
              {{ memberName }}
            </option>
          </select>
        </label>

        
        <label class="form-field">
          <span class="field-label">Приоритет</span>
          <select v-model="form.priority" class="cyber-select">
            <option value="Critical">🔥 Critical (1)</option>
            <option value="High">⚠️ High (2)</option>
            <option value="Medium">⚡ Medium (3)</option>
            <option value="Low">🌱 Low (4)</option>
            <option value="Backlog">📂 Backlog (5)</option>
            <option value="Blocked">🚫 Blocked (6)</option>
          </select>
        </label>

        
        <label class="form-field">
          <span class="field-label">Выполнить до (Дедлайн)</span>
          <input
            v-model="form.deadline"
            type="date"
            class="cyber-input"
          />
        </label>

        
        <div class="modal-actions">
          
          <button
            v-if="isEditMode"
            type="button"
            class="action-btn delete-btn"
            @click="handleDelete"
          >
            Удалить задачу
          </button>

          <div class="right-actions">
            <button type="button" class="action-btn cancel-btn" @click="emit('close')">
              Отмена
            </button>
            <button type="submit" class="action-btn save-btn">
              {{ isEditMode ? 'Сохранить изменения' : 'Создать задачу' }}
            </button>
          </div>
        </div>

      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';

const props = defineProps({
  task: {
    type: Object,
    default: null
  },
  command: {
    type: Object,
    required: true
  }
});

const emit = defineEmits(['close', 'create', 'update', 'delete']);

const isEditMode = computed(() => props.task !== null);

const form = ref({
  title: '',
  description: '',
  assigned_to_name: null,
  priority: 'Medium',
  deadline: ''
});

onMounted(() => {
  if (isEditMode.value && props.task) {
    form.value = {
      title: props.task.title || '',
      description: props.task.description || '',
      assigned_to_name: props.task.assigned_to?.name || props.task.asigned_to?.name || props.task.assigned_to || null,
      priority: props.task.priority || 'Medium',
      deadline: props.task.deadline ? props.task.deadline.split('T')[0] : ''
    };
  }
});

function handleSubmit() {
  if (!form.value.title.trim()) return;

  if (isEditMode.value) {
    emit('update', { id: props.task.id, ...form.value });
  } else {
    emit('create', { ...form.value, status: 'TODO' });
  }
}

function handleDelete() {
  if (confirm('Вы уверены, что хотите безвозвратно удалить эту задачу?')) {
    emit('delete', props.task.id);
  }
}
</script>

<style scoped>

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 16px;
  box-sizing: border-box;
}


.modal-card {
  font-family: system-ui, -apple-system, sans-serif;
  background-color: #1e293b; 
  border: 1px solid #334155;
  border-radius: 14px;
  width: 100%;
  max-width: 520px;
  box-shadow: 0 20px 25px -5px rgba(168, 85, 247, 0.15), 0 10px 10px -5px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}


.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  border-bottom: 1px solid #334155;
}

.modal-title {
  font-size: 20px;
  font-weight: 800;
  color: #f8fafc;
  margin: 0;
  text-transform: uppercase;
  letter-spacing: -0.01em;
}

.close-btn {
  background: none;
  border: none;
  color: #94a3b8;
  font-size: 28px;
  cursor: pointer;
  line-height: 1;
  padding: 0;
}

.close-btn:hover {
  color: #f43f5e;
}


.modal-form {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field-label {
  font-size: 12px;
  font-weight: 700;
  color: #cbd5e1;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}


.cyber-input, .cyber-textarea, .cyber-select {
  font-family: inherit;
  width: 100%;
  border: 1px solid #475569;
  border-radius: 8px;
  padding: 10px 14px;
  background-color: #f1f5f9; 
  color: #0f172a;           
  font-size: 14px;
  font-weight: 600;
  outline: none;
  box-sizing: border-box;
  transition: border-color 0.2s;
}

.cyber-textarea {
  resize: vertical;
}

.cyber-input:focus, .cyber-textarea:focus, .cyber-select:focus {
  border-color: #a855f7;
  box-shadow: 0 0 0 3px rgba(168, 85, 247, 0.25);
}


.modal-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 12px;
  border-top: 1px solid #334155;
  padding-top: 16px;
  gap: 12px;
}

.right-actions {
  display: flex;
  gap: 12px;
  margin-left: auto;
}

.action-btn {
  font-family: inherit;
  font-size: 14px;
  font-weight: 700;
  padding: 10px 18px;
  border-radius: 8px;
  cursor: pointer;
  border: none;
  transition: opacity 0.2s, transform 0.1s;
}

.action-btn:active {
  transform: scale(0.98);
}


.save-btn {
  background: linear-gradient(90deg, #7c3aed 0%, #db2777 100%);
  color: #ffffff;
}

.save-btn:hover {
  opacity: 0.9;
}


.cancel-btn {
  background-color: #334155;
  color: #cbd5e1;
  border: 1px solid #475569;
}

.cancel-btn:hover {
  background-color: #475569;
}


.delete-btn {
  background-color: #mx-auto;
  background-color: #1e293b;
  border: 1px solid #f43f5e;
  color: #f43f5e;
}

.delete-btn:hover {
  background-color: #f43f5e;
  color: #ffffff;
}
</style>
