
<!-- ID, название, описание, статус, приоритет, теги, дата создания, дедлайн -->

<template>
  <article class="task-card">
    <div>
      <p class="task-status">{{ task.status }}</p>
      <p class="task-id">id: {{task.id}}</p>
      <h3>{{ task.title }}</h3>
    </div>

    <p>{{ task.description }}</p>
    
    <div class="task-asigned">
      <div class="task-asigned-name">
        Исполнитель: {{ task.asigned_to.name }}
      </div>
      <div class="task-asigned-email">
      {{ task.asigned_to.email }}
      </div>
    </div>
    <div class="task-priority">
      <span>Приоритет: {{ task.priority }}</span>
    </div>
    <div class="task-deadline">
      <span>Выполнить до: {{ task.deadline }}</span>
    </div>

    <div class="task-tags" >
        <span v-for="tag in task.tags">{{ tag }}</span>
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
  gap: 14px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 18px;
}

.task-card + .task-card {
  margin-top: 14px;
}

.task-card h3 {
  margin-bottom: 0;
  font-size: 18px;
}

.task-card p {
  margin-bottom: 0;
  color: #475569;
}

.task-priority{
  color: #475569;
}

.task-tags{
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 5px;
  }

.task-tags span{
  font-family: 'Courier New', Courier, monospace;
  font-size: 13px;
  color: white;
  padding: 8px;
  background-color: #3730a3;
  border-radius: 10px;
}
</style>
