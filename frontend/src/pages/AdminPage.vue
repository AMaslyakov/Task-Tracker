<template>
  <main class="page-admin">
    <div class="admin-nav">
      <button type="button" class="back-btn" @click="forceRedirect">
        ← Вернуться к задачам
      </button>
      <h1 class="admin-title">Панель администратора</h1>
    </div>

    <div class="admin-grid">
      <section class="admin-card">
        <h2 class="card-subtitle">📁 Создание новой команды</h2>

        <form @submit.prevent="handleCreateTeam" class="admin-form">
          <label class="form-field">
            <span>Название команды *</span>
            <input v-model="newTeam.name" type="text" placeholder="Например: Backend API" required />
          </label>

          <label class="form-field">
            <span>Описание</span>
            <textarea v-model="newTeam.description" placeholder="Чем занимается команда..." rows="2"></textarea>
          </label>

          <label class="form-field">
            <span>Статусы канбан-доски (через запятую)</span>
            <input v-model="newTeam.statusesStr" type="text" placeholder="TODO, IN PROGRESS, DONE" />
          </label>

          <button type="submit" class="action-btn save-btn">Создать команду</button>
        </form>

        <div class="teams-list-section">
          <h3 class="inner-title">Существующие команды ({{ teams.length }})</h3>
          <div class="teams-container">
            <div
              v-for="team in teams"
              :key="team.id"
              class="team-item-row"
              :class="{ 'is-selected': selectedTeam?.id === team.id }"
              @click="selectTeam(team)"
            >
              <div class="team-meta-info">
                <span class="team-row-name">{{ team.name }}</span>
                <span class="team-row-desc">{{ team.description || 'Нет описания' }}</span>
              </div>
              <span class="members-count-badge">👥 {{ team.members?.length || 0 }}</span>
            </div>
          </div>
        </div>
      </section>

      <aside class="admin-card text-right-panel">
        <h2 class="card-subtitle">👥 Управление участниками</h2>

        <div v-if="selectedTeam" class="active-management-zone">
          <div class="selected-team-banner">
            Выбрана команда: <strong>{{ selectedTeam.name }}</strong>
          </div>

          <form @submit.prevent="handleAddUserToTeam" class="admin-form add-user-inline">
            <label class="form-field">
              <span>Выбрать сотрудника из базы данных *</span>
              <div class="inline-input-group">
                <select v-model="newUserLog" class="admin-select-field" required>
                  <option :value="''" disabled selected>-- Выберите пользователя --</option>
                  <option
                    v-for="user in allSystemUsers"
                    :key="user"
                    :value="user"
                    :disabled="selectedTeam.members?.includes(user)"
                  >
                    👤 {{ user }} {{ selectedTeam.members?.includes(user) ? '(Уже в команде)' : '' }}
                  </option>
                </select>
                <button type="submit" class="action-btn add-btn" :disabled="!newUserLog">Добавить</button>
              </div>
            </label>
          </form>

          <div class="current-members-list">
            <h3 class="inner-title">Состав команды:</h3>
            <div v-if="selectedTeam.members && selectedTeam.members.length" class="members-chips-grid">
              <div v-for="member in selectedTeam.members" :key="member" class="admin-member-pill">
                <span class="member-pill-name">{{ member }}</span>
                <button type="button" class="remove-member-btn" @click="handleRemoveUser(member)">&times;</button>
              </div>
            </div>
            <p v-else class="empty-text">В этой команде ещё нет участников. Добавьте первого!</p>
          </div>
        </div>

        <div v-else class="placeholder-state">
          <span class="placeholder-icon">👈</span>
          <p>Выберите команду из списка слева, чтобы управлять её участниками и добавлять людей.</p>
        </div>
      </aside>
    </div>
  </main>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { fetchTeams, fetchAllUsers } from '../api/tasks'

const teams = ref([])
const selectedTeam = ref(null)
const newUserLog = ref('')
const allSystemUsers = ref([])

const newTeam = ref({
  name: '',
  description: '',
  statusesStr: 'TODO, IN PROGRESS, REVIEW, DONE'
})

onMounted(loadAdminData)

async function loadAdminData() {
  try {
    const loadedTeams = await fetchTeams()
    teams.value = loadedTeams
    if (teams.value.length > 0) {
      selectedTeam.value = teams.value[0]
    }
  } catch (error) {
    console.error('Ошибка загрузки команд с бэкенда:', error)
  }

  try {
    const usersResponse = await fetchAllUsers()
    allSystemUsers.value = usersResponse.map(u => u.user_name).filter(Boolean)
  } catch (e) {
    console.error('Ошибка загрузки списка юзеров из таблицы user:', e)
  }
}

function selectTeam(team) {
  selectedTeam.value = team
  newUserLog.value = ''
}

async function handleCreateTeam() {
  if (!newTeam.value.name.trim()) return

  const statusesArray = newTeam.value.statusesStr
    .split(',')
    .map(s => s.trim().toUpperCase())
    .filter(s => s.length > 0)

  const teamPayload = {
    name: newTeam.value.name,
    description: newTeam.value.description,
    config_dashboard: { statuses: statusesArray },
    members: []
  }

  try {
    const response = await fetch('/api/teams', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(teamPayload)
    })

    if (response.ok) {
      const created = await response.json()
      teams.value.push(created)
    }
  } catch (error) {
    console.error(error)
  } finally {
    newTeam.value.name = ''
    newTeam.value.description = ''
    forceRedirect()
  }
}

async function handleAddUserToTeam() {
  if (!newUserLog.value || !selectedTeam.value) return

  try {
    if (!selectedTeam.value.members) {
      selectedTeam.value.members = []
    }

    selectedTeam.value.members.push(newUserLog.value)
    newUserLog.value = ''
  } catch (error) {
    console.error(error)
  }
}

async function handleRemoveUser(memberName) {
  if (!selectedTeam.value) return

  if (confirm(`Исключить сотрудника ${memberName}?`)) {
    try {
      selectedTeam.value.members = selectedTeam.value.members.filter(m => m !== memberName)
    } catch (error) {
      console.error(error)
    }
  }
}

function forceRedirect() {
  window.location.href = '/tasks'
}
</script>

<style scoped>
.page-admin {
  font-family: system-ui, -apple-system, sans-serif;
  min-height: 100vh;
  padding: 32px;
  background-color: #f8fafc;
  max-width: 1120px;
  margin: 0 auto;
  box-sizing: border-box;
}

.admin-nav {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 28px;
}

.back-btn {
  background: none;
  border: none;
  color: #6366f1;
  font-size: 15px;
  font-weight: 700;
  cursor: pointer;
  padding: 0;
  align-self: flex-start;
}

.admin-title {
  font-size: 28px;
  font-weight: 800;
  color: #0f172a;
  margin: 0;
}

.admin-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
  align-items: start;
}

.admin-card {
  background-color: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 24px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.02);
}

.card-subtitle {
  font-size: 16px;
  font-weight: 800;
  color: #1e293b;
  margin: 0 0 20px 0;
  border-bottom: 1px solid #f1f5f9;
  padding-bottom: 12px;
}

.admin-form {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-field span {
  font-size: 12px;
  font-weight: 700;
  color: #64748b;
  text-transform: uppercase;
}

.admin-form input, .admin-form textarea, .admin-select-field {
  font-family: inherit;
  width: 100%;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 10px 14px;
  background-color: #f8fafc;
  color: #0f172a;
  font-size: 14px;
  font-weight: 600;
  outline: none;
  box-sizing: border-box;
}

.admin-select-field {
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://w3.org' viewBox='0 0 24 24' fill='none' stroke='%23475569' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3e%3cpolyline points='6 9 12 15 18 9'%3e%3c/polyline%3e%3c/svg%3e");
  background-repeat: no-repeat;
  background-position: right 14px center;
  background-size: 16px;
  padding-right: 40px;
}

.admin-form input:focus, .admin-form textarea:focus, .admin-select-field:focus {
  border-color: #6366f1;
}

.action-btn {
  font-family: inherit;
  font-size: 14px;
  font-weight: 700;
  padding: 10px 16px;
  border-radius: 8px;
  cursor: pointer;
  border: none;
  text-transform: uppercase;
  letter-spacing: 0.02em;
}

.save-btn {
  background: linear-gradient(90deg, #7c3aed 0%, #db2777 100%);
  color: #ffffff;
  margin-top: 6px;
}

.teams-list-section {
  margin-top: 24px;
  border-top: 1px solid #f1f5f9;
  padding-top: 20px;
}

.inner-title {
  font-size: 13px;
  font-weight: 700;
  color: #94a3b8;
  text-transform: uppercase;
  margin: 0 0 12px 0;
}

.teams-container {
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-height: 260px;
  overflow-y: auto;
}

.team-item-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background-color: #f8fafc;
  border: 1px solid #e2e8f0;
  padding: 12px 16px;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.15s;
}

.team-item-row:hover {
  background-color: #f1f5f9;
  border-color: #cbd5e1;
}

.team-item-row.is-selected {
  background-color: #eef2ff;
  border-color: #6366f1;
}

.team-meta-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.team-row-name {
  font-size: 14px;
  font-weight: 700;
  color: #1e293b;
}

.team-row-desc {
  font-size: 12px;
  color: #64748b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.members-count-badge {
  font-size: 12px;
  font-weight: 600;
  color: #475569;
  background-color: #e2e8f0;
  padding: 4px 8px;
  border-radius: 6px;
  flex-shrink: 0;
}

.selected-team-banner {
  background-color: #f1f5f9;
  border-left: 4px solid #6366f1;
  padding: 10px 14px;
  border-radius: 4px;
  font-size: 14px;
  color: #334155;
  margin-bottom: 16px;
}

.inline-input-group {
  display: flex;
  gap: 10px;
}

.add-btn {
  background-color: #1e293b;
  color: white;
  flex-shrink: 0;
}

.add-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.current-members-list {
  margin-top: 20px;
}

.members-chips-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.admin-member-pill {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background-color: #f1f5f9;
  border: 1px solid #cbd5e1;
  border-radius: 20px;
  padding: 4px 12px;
  font-size: 13px;
  font-weight: 600;
  color: #334155;
}

.remove-member-btn {
  background: none;
  border: none;
  color: #94a3b8;
  font-size: 16px;
  cursor: pointer;
  padding: 0;
  line-height: 1;
}

.remove-member-btn:hover {
  color: #f43f5e;
}

.placeholder-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #94a3b8;
  text-align: center;
}

.placeholder-icon {
  font-size: 32px;
  margin-bottom: 12px;
}

.empty-text {
  color: #94a3b8;
  font-style: italic;
  font-size: 13px;
}

@media (max-width: 760px) {
  .admin-grid { grid-template-columns: 1fr; }
  .inline-input-group { flex-direction: column; }
}
</style>

