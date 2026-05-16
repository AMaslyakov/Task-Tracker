<template>
  <section class="task-column">
    <div class="section-heading">
      <h2>{{ title }}</h2>
      <span>{{ tasks.length }} задачи</span>
    </div>

    <TaskCard v-for="task in tasks" :key="task.title" :task="task" />
  </section>
</template>

<script setup>
import TaskCard from './TaskCard.vue'

defineProps({
  title: {
    type: String,
    required: true
  },
  tasks: {
    type: Array,
    required: true
  }
})
</script>

<style scoped>
.task-column {
  border: 1px solid #0955cf;
  border-radius: 8px;
  padding: 20px;
  background: #ffffff;
}

.section-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 16px;
}

.section-heading h2 {
  margin-bottom: 0;
  font-size: 22px;
}

.section-heading span {
  color: #64748b;
  font-size: 14px;
}
</style>
