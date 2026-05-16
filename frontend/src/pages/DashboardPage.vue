<template>
  <main class="page-shell">
    <section class="workspace">
      <AppHeader
        :commands="commands"
        :selected-command-id="selectedCommandId"
        @update:selected-command-id="selectedCommandId = $event"
      />

      <div class="content-grid">
        <StatusColumnList :tasks="filteredTasks" :command="selectedCommand"/>
      </div>
    </section>
  </main>
</template>

<script setup>
import { computed, ref } from 'vue'
import AppHeader from '../components/AppHeader.vue'
import StatusColumnList from '../components/StatusColumnList.vue';
import { commands, tasks } from '../data/dashboard'

const selectedCommandId = ref(commands[0].id)

const selectedCommand = computed(() => {
  return commands.find((command) => command.id === selectedCommandId.value) ?? commands[0]
})

const filteredTasks = computed(() => {
  return tasks.filter((task) => task.command.id === selectedCommand.value.id)
})
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
