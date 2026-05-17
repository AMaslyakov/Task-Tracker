<template>
  <div class="toolbar cyberpunk-graphite">
    <!-- Стилизованный блок заголовка с неоном -->
    <div class="logo-area-cyber">
      <p class="eyebrow-neon">System v2.6</p>
      <!-- ИСПРАВЛЕНО: Заголовок стал крупнее и выразительнее -->
      <h1 class="main-title-neon">Task Tracker</h1>
    </div>

    <!-- Селектор выбора команды -->
    <label class="command-select">
      <!-- ИСПРАВЛЕНО: Более читаемый, крупный и чистый шрифт -->
      <span class="select-label">Команда</span>
      <div class="select-wrapper">
        <select :value="selectedCommandId" @change="handleCommandChange">
          <option
            v-for="command in commands"
            :key="command.id"
            :value="command.id"
          >
            {{ command.name }}
          </option>
        </select>
      </div>
    </label>

    <!-- Профиль пользователя -->
    <div class="user auth-profile-block">
      <img class="user-icon" src="../assets/user.png" :alt="user2.name">
      <div class="user-data">
          <span class="user-name">{{ user2.name }}</span>
          <span class="user-email">{{ user2.email }}</span>
          <span class="user-id">ID: {{ user2.id }}</span>
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
/* Главный контейнер тулбара */
.toolbar {
  /* Комбинируем: заголовок в стиле Courier, а интерфейс — в читаемом без засечек */
  font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 32px;
  margin-bottom: 32px;
  background-color: #1e293b;
  padding: 18px 24px;
  border-radius: 14px;
  border: 1px solid #334155;
  box-shadow: 0 10px 25px -5px rgba(168, 85, 247, 0.12);
}

/* Блок заголовка с неоном */
.logo-area-cyber {
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.eyebrow-neon {
  font-family: 'Courier New', monospace;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  color: #c084fc;
  letter-spacing: 0.15em;
  margin: 0 0 2px 0;
}

/* ИСПРАВЛЕНО: Шрифт увеличен до 36px, вес 900 (максимально жирный) */
.main-title-neon {
  font-family: 'Courier New', monospace;
  font-size: 36px;
  font-weight: 900;
  margin: 0;
  letter-spacing: -0.03em;
  text-transform: uppercase;
  background: linear-gradient(90deg, #a855f7 0%, #ec4899 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  filter: drop-shadow(0 0 8px rgba(168, 85, 247, 0.4));
}

/* Стилизация выпадающего списка */
.command-select {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 200px;
  flex-grow: 1;
  max-width: 280px;
}

/* ИСПРАВЛЕНО: Шрифт переключен на более читаемый system-ui */
.select-label {
  font-family: system-ui, sans-serif;
  font-size: 13px;
  font-weight: 800;
  color: #ffffff;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

/* ИСПРАВЛЕНО: Чистый шрифт для выпадающего списка вариантов */
.command-select select {
  font-family: system-ui, sans-serif;
  width: 100%;
  min-height: 42px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 0 14px;
  background-color: #f1f5f9;
  color: #0f172a;
  font-size: 14px;
  font-weight: 700;
  outline: none;
  cursor: pointer;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.command-select select:focus {
  border-color: #a855f7;
  box-shadow: 0 0 0 3px rgba(168, 85, 247, 0.3);
}

/* Блок профиля пользователя */
.user {
  display: flex;
  align-items: center;
  gap: 14px;
  background-color: #334155;
  padding: 10px 20px;
  border-radius: 10px;
  border: 1px solid #475569;
  min-width: 260px;
  flex-shrink: 0;
}

.user-icon {
  width: 46px;
  height: 46px;
  object-fit: cover;
  border-radius: 50%;
  border: 2px solid #475569;
}

.user-data {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  width: 100%;
}

.user-name {
  font-size: 15px;
  font-weight: 700;
  color: #f8fafc;
  line-height: 1.2;
  white-space: nowrap;
}

.user-email {
  font-size: 12px;
  color: #cbd5e1;
  margin-top: 1px;
  white-space: nowrap;
}

.user-id {
  font-size: 10px;
  font-weight: 700;
  color: #09090b;
  background-color: #a855f7;
  padding: 1px 6px;
  border-radius: 4px;
  margin-top: 4px;
}

/* Адаптивность */
@media (max-width: 760px) {
  .toolbar {
    display: grid;
    grid-template-columns: 1fr;
    gap: 18px;
    padding: 16px;
  }
  .command-select {
    width: 100%;
    min-width: 0;
    max-width: 100%;
  }
  .user {
    width: 100%;
    min-width: 0;
    box-sizing: border-box;
  }
}
</style>
