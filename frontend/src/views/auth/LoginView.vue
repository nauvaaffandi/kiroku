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
            type="number"
            inputmode="numeric"
            pattern="[0-9]{6}"
            maxlength="6"
            min="0"
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

        <!-- Only show fingerprint login button if biometric is set/enabled in settings -->
        <button
          v-if="isBiometricConfigured"
          type="button"
          class="login-fingerprint-btn"
          @click="showFingerprintOverlay = true"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 10a2 2 0 0 0-2 2c0 1.02-.1 2.51-.81 3.5"></path>
            <path d="M15.5 13.5c0-1.8-1.5-3.5-3.5-3.5s-3.5 1.7-3.5 3.5"></path>
            <path d="M19 11c0-3.86-3.14-7-7-7s-7 3.14-7 7c0 1.55.3 3 .83 4.33"></path>
            <path d="M21.5 14C22 12.8 22.5 11.5 22.5 10c0-5.52-4.48-10-10-10S2.5 4.48 2.5 10c0 2.21.73 4.25 1.97 5.88"></path>
            <path d="M12 18v2"></path>
            <path d="M9 19.5c0 .83.67 1.5 1.5 1.5h3c.83 0 1.5-.67 1.5-1.5V18"></path>
          </svg>
          Masuk dengan Sidik Jari
        </button>
      </form>

      <p class="login-footer">
        Single-user mode · Device remembered for 24h of activity
      </p>
    </div>

    <!-- Fingerprint Overlay -->
    <FingerprintOverlay
      :is-open="showFingerprintOverlay"
      @close="showFingerprintOverlay = false"
      @success="handleFingerprintSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useUIStore } from '@/stores/ui'
import { useRouter } from 'vue-router'
import FingerprintOverlay from '@/components/auth/FingerprintOverlay.vue'

const authStore = useAuthStore()
const uiStore = useUIStore()
const router = useRouter()

const pin = ref('')
const showFingerprintOverlay = ref(false)
const isBiometricConfigured = ref(false)

const submitting = ref(false)

const handleVerify = async () => {
  if (submitting.value) return
  submitting.value = true
  const success = await authStore.login(pin.value)
  if (success) {
    uiStore.success('Welcome back!')
    router.push('/dashboard')
  }
  pin.value = ''
  submitting.value = false
}

watch(pin, (newVal) => {
  const strVal = String(newVal || '')
  if (strVal.length >= 6) {
    pin.value = strVal.slice(0, 6)
    handleVerify()
  }
})

const handleFingerprintSuccess = async () => {
  showFingerprintOverlay.value = false
  authStore.user = { user_id: 'authenticated' }
  uiStore.success('Autentikasi sidik jari berhasil!')
  router.push('/dashboard')
}

onMounted(() => {
  // Check if biometric is enabled/configured in localStorage (set via settings)
  const configured = localStorage.getItem('kiroku_biometric_enabled') === 'true'
  isBiometricConfigured.value = configured

  // Only auto-open if biometric is explicitly configured by user
  if (configured) {
    setTimeout(() => {
      showFingerprintOverlay.value = true
    }, 400)
  }
})
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
  gap: 1.25rem;
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

.login-fingerprint-btn {
  width: 100%;
  padding: 0.75rem 1.5rem;
  background-color: var(--bg-secondary, #f1f3f5);
  color: var(--text-primary, #212529);
  font-weight: 500;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-primary, #dee2e6);
  cursor: pointer;
  transition: background-color 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.login-fingerprint-btn:hover {
  background-color: var(--bg-hover, #e9ecef);
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
