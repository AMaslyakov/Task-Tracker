<template>
  <main class="page-shell">
    <section class="workspace">
      <AppHeader />

      <div class="content-grid">
        <!-- <TaskList title="Сегодня" :tasks="tasks" /> -->
        <StatusColumnList :tasks="tasks"/>
        <!-- <SummaryPanel :items="summaryItems" /> -->
      </div>
    </section>
  </main>
</template>

<script setup>
import AppHeader from '../components/AppHeader.vue'
// import SummaryPanel from '../components/SummaryPanel.vue'
import TaskList from '../components/TaskList.vue'
import StatusColumnList from '../components/StatusColumnList.vue';
import { tasks } from '../data/dashboard'
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

@media (max-width: 760px) {
  .page-shell {
    padding: 20px;
  }

  .content-grid {
    grid-template-columns: 1fr;
  }
}
</style>
