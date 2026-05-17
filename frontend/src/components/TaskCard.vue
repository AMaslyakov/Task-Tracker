<!-- ID, название, описание, статус, приоритет, теги, дата создания, дедлайн -->

<template>
  <article class="task-card">
    <div class="task-header">
      <p class="task-status">{{ task.status }}</p>
      <p class="task-id">id: {{ task.id }}</p>
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
      <!-- Email с автоуменьшением шрифта -->
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
defineProps({
  task: {
    type: Object,
    required: true
  }
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
}

.task-card + .task-card {
  margin-top: 14px;
}

.task-header {
  width: 100%;
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
  max-width: 100%;        /* Не дает вылезти за границы карточки */
  overflow: hidden;       /* Скрывает то, что не влезло */
  text-overflow: ellipsis; /* Добавляет три точки ... */
  box-sizing: border-box;
}
</style>
