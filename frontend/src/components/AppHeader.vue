<template>

  <div class="toolbar">


    <div>
      <!-- <p class="eyebrow">Victory Group</p> -->
      <h1>Task Tracker</h1>
    </div>
    <!-- <button type="button">Новая задача</button> -->

    <label class="command-select">
      <span>Команда</span>
      <select :value="selectedCommandId" @change="handleCommandChange">
        <option
          v-for="command in commands"
          :key="command.id"
          :value="command.id"
        >
          {{ command.name }}
        </option>
      </select>
    </label>

    <div class="user">
      <img class="user-icon" src="../assets/user.png" :alt="user2.name">
      <div class="user-data">
          <span class="user-id">Id: {{ user2.id}}</span>
          <span class="user-name">{{ user2.name}}</span>
          <span class="user-email">{{ user2.email}}</span>
      </div>
    </div>

  </div>
  
  
</template>

<script setup>
  import { user2 } from "../data/dashboard"

  defineProps({
    commands: {
      type: Array,
      required: true
    },
    selectedCommandId: {
      type: Number,
      required: true
    }
  })

  const emit = defineEmits(['update:selectedCommandId'])

  function handleCommandChange(event) {
    emit('update:selectedCommandId', Number(event.target.value))
  }
</script>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 28px;
}

.command-select {
  display: grid;
  gap: 6px;
  min-width: 180px;
  color: #334155;
  font-weight: 700;
}

.command-select select {
  min-height: 42px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 0 12px;
  background: #ffffff;
  color: #172033;
  font: inherit;
}

@media (max-width: 760px) {
  .toolbar {
    display: grid;
  }

  button,
  .command-select {
    width: 100%;
  }
}

.user-icon{
  max-width: 40px;
}

.user-data{
  display: flex;
  flex-direction: column;
  align-items:self-start;
  font-size: 22px;
  row-gap: 8px;
}

</style>
