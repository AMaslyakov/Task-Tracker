<template>
  <main class="auth-page">
    <form class="auth-form" @submit.prevent="handleSubmit">
      <div>
        <p class="eyebrow">Task Tracker</p>
        <h1>Вход</h1>
      </div>

      <label>
        Email
        <input v-model="email" type="email" autocomplete="email" required />
      </label>

      <label>
        Пароль
        <input
          v-model="password"
          type="password"
          autocomplete="current-password"
          required
        />
      </label>

      <button type="submit">Войти</button>

      <p v-if="errorMessage" class="error-message">{{ errorMessage }}</p>
    </form>
  </main>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const email = ref('')
const password = ref('')
const errorMessage = ref('')

function handleSubmit() {
  errorMessage.value = ''

  if (!email.value || !password.value) {
    errorMessage.value = 'Введите email и пароль'
    return
  }

  router.push('/tasks')
}
</script>

<style scoped>
.auth-page {
  display: grid;
  min-height: 100vh;
  place-items: center;
  padding: 24px;
}

.auth-form {
  display: grid;
  gap: 18px;
  width: min(400px, 100%);
  border: 1px solid #dce3ee;
  border-radius: 8px;
  padding: 28px;
  background: #ffffff;
}

.auth-form h1 {
  font-size: 36px;
}

.auth-form label {
  display: grid;
  gap: 8px;
  color: #334155;
  font-weight: 700;
}

.auth-form input {
  min-height: 44px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 0 12px;
  color: #172033;
  font: inherit;
}

.auth-form input:focus {
  border-color: #2563eb;
  outline: 3px solid #dbeafe;
}

.error-message {
  margin-bottom: 0;
  color: #dc2626;
  font-weight: 700;
}
</style>
