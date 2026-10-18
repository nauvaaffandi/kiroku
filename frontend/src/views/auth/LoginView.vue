<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <h1>Kiroku</h1>
        <p>Personal Finance Tracker</p>
      </div>

      <form @submit.prevent="handleVerify" class="login-form">
        <div v-if="authStore.error" class="form-error">
          {{ authStore.error }}
        </div>

        <div class="form-group">
          <label for="pin">PIN</label>
          <input
            id="pin"
            v-model="pin"
            type="password"
            inputmode="numeric"
            pattern="[0-9]*"
            required
            autocomplete="current-password"
            :disabled="authStore.loading"
            autofocus
            @keydown.enter="handleVerify"
          />
        </div>

        <button
          type="submit"
          :disabled="authStore.loading"
          class="login-submit"
        >
          <span v-if="authStore.loading" class="btn-loading">
            <span class="spinner"></span>
            Verifying...
          </span>
          <span v-else>Sign In</span>
        </button>
      </form>

      <p class="login-footer">
        Single-user mode · Device remembered for 24h of activity
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useUIStore } from '@/stores/ui'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const uiStore = useUIStore()
const router = useRouter()

const pin = ref('')

const handleVerify = async () => {
  const success = await authStore.login(pin.value)
  if (success) {
    uiStore.success('Welcome back!')
    router.push('/dashboard')
  }
  pin.value = ''
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--color-neutral-50);
  padding: 1rem;
}

.login-card {
  width: 100%;
  max-width: 28rem;
  background-color: var(--color-neutral-100);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 2rem;
}

.login-header {
  text-align: center;
  margin-bottom: 2rem;
}

.login-header h1 {
  font-size: 1.875rem;
  font-weight: 700;
  color: var(--color-primary);
}

.login-header p {
  color: var(--color-neutral-500);
  margin-top: 0.5rem;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.form-group label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-neutral-700);
}

.form-group input {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 1px solid var(--color-neutral-300);
  border-radius: var(--radius-sm);
  background-color: var(--color-neutral-100);
  color: var(--color-neutral-900);
  font-size: 1rem;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.form-group input:focus {
  outline: none;
  border-color: transparent;
  box-shadow: 0 0 0 2px var(--color-primary);
}

.form-group input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.form-error {
  padding: 0.75rem 1rem;
  background-color: var(--color-red-50);
  color: var(--color-red-600);
  border-radius: var(--radius-sm);
  font-size: 0.875rem;
}

.login-submit {
  width: 100%;
  padding: 0.75rem 1.5rem;
  background-color: var(--color-primary);
  color: white;
  font-weight: 500;
  border-radius: var(--radius-sm);
  border: none;
  cursor: pointer;
  transition: background-color 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.login-submit:hover:not(:disabled) {
  background-color: var(--color-primary-hover);
}

.login-submit:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.spinner {
  width: 1.25rem;
  height: 1.25rem;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.login-footer {
  margin-top: 1.5rem;
  text-align: center;
  font-size: 0.75rem;
  color: var(--color-neutral-400);
}
</style>