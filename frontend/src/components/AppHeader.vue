<template>
  <div class="toolbar cyberpunk-graphite">
    
    <div class="logo-area-cyber">
      <p class="eyebrow-neon">System v2.6</p>
      <h1 class="main-title-neon">Task Tracker</h1>
    </div>

    
    <label class="command-select">
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

    
    <div class="user auth-profile-block">
      <img class="user-icon" src="../assets/user.png" :alt="currentUserName">
      <div class="user-data">
          <span class="user-name">{{ currentUserName }}</span>
          <span class="user-email">{{ currentUserEmail }}</span>
          <span class="user-id">ID: {{ currentUserId }}</span>
      </div>
      <button class="logout-button" type="button" @click="emit('logout')">
        Выйти
      </button>
    </div>
  </div>
</template>

<script setup>
  import { computed } from 'vue'

  const props = defineProps({
    commands: {
      type: Array,
      required: true
    },
    selectedCommandId: {
      type: Number,
      required: true
    },
    currentUser: {
      type: Object,
      default: null
    }
  })

  const emit = defineEmits(['update:selectedCommandId', 'logout'])

  const currentUserName = computed(() => props.currentUser?.user_name ?? 'Пользователь')
  const currentUserEmail = computed(() => props.currentUser?.email ?? '')
  const currentUserId = computed(() => props.currentUser?.id ?? '-')

  function handleCommandChange(event) {
    emit('update:selectedCommandId', Number(event.target.value))
  }
</script>

<style scoped>

.toolbar {
  font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 32px;
  margin-bottom: 32px;
  background-color: #1e293b;
  padding: 18px 24px;
  border-radius: 14px;
  border: 1px solid #334155;
  box-shadow: 0 10px 25px -5px rgba(168, 85, 247, 0.12);
}


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

.command-select {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 200px;
  flex: 1 1 220px;
  max-width: 280px;
}

.select-label {
  font-family: system-ui, sans-serif;
  font-size: 13px;
  font-weight: 800;
  color: #ffffff;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

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


.user {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  align-items: center;
  gap: 14px;
  background-color: #334155;
  padding: 10px 14px;
  border-radius: 10px;
  border: 1px solid #475569;
  width: min(100%, 420px);
  min-width: 0;
  flex: 1 1 300px;
  box-sizing: border-box;
}

.user-icon {
  width: 46px;
  height: 46px;
  min-width: 46px;
  object-fit: cover;
  border-radius: 50%;
  border: 2px solid #475569;
}

.user-data {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  width: 100%;
  min-width: 0;
}

.user-name {
  font-size: 15px;
  font-weight: 700;
  color: #f8fafc;
  line-height: 1.2;
  white-space: nowrap;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-email {
  font-size: 12px;
  color: #cbd5e1;
  margin-top: 1px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
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

.logout-button {
  min-height: 36px;
  min-width: 68px;
  border: 1px solid #64748b;
  border-radius: 8px;
  padding: 0 12px;
  background-color: #0f172a;
  color: #f8fafc;
  font: inherit;
  font-size: 13px;
  font-weight: 800;
  cursor: pointer;
}

.logout-button:hover {
  border-color: #c084fc;
  background-color: #1e293b;
}


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
  }
}

@media (max-width: 420px) {
  .user {
    grid-template-columns: auto minmax(0, 1fr);
  }

  .logout-button {
    grid-column: 1 / -1;
    width: 100%;
  }
}
</style>
