<script setup lang="ts">
import { ref } from 'vue'

defineProps<{
  modelValue: boolean
  title: string
  loading: boolean
  mode: 'change' | 'reset'
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'submit', token: string, newPin: string): void
}>()

const token = ref('')
const newPin = ref('')
const confirmPin = ref('')
const error = ref('')

const handleSubmit = (e: Event) => {
  e.preventDefault()
  if (newPin.value !== confirmPin.value) {
    error.value = 'New PINs do not match'
    return
  }
  if (!token.value.trim()) {
    error.value = 'Token is required'
    return
  }
  error.value = ''
  emit('submit', token.value.trim(), newPin.value)
}

const close = () => {
  token.value = ''
  newPin.value = ''
  confirmPin.value = ''
  error.value = ''
  emit('update:modelValue', false)
}
</script>

<template>
  <Transition name="fade">
      <div v-if="modelValue" class="tp-backdrop">
        <div class="tp-modal">
          <div class="tp-header">
            <h2 class="tp-title">{{ title }}</h2>
            <button @click="close" class="tp-close-btn">
              <svg class="tp-close-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <form @submit="handleSubmit" class="tp-form">
            <div v-if="error" class="tp-error">
              {{ error }}
            </div>

            <div>
              <label class="tp-label">
                {{ mode === 'change' ? 'Telegram Token' : 'Reset Token from Telegram' }}
              </label>
              <input
                v-model="token"
                type="text"
                required
                placeholder="Paste token from @nauva_kiroku_bot"
                class="tp-input"
              />
              <p class="tp-hint">
                {{ mode === 'change' 
                  ? 'Request token via /changepin command in Telegram' 
                  : 'Request token via /forgotpin command in Telegram' }}
              </p>
            </div>

            <div>
              <label class="tp-label">New PIN</label>
              <input
                v-model="newPin"
                type="number"
                inputmode="numeric"
                pattern="[0-9]{6}"
                maxlength="6"
                min="0"
                required
                class="tp-input"
              />
            </div>

            <div>
              <label class="tp-label">Confirm New PIN</label>
              <input
                v-model="confirmPin"
                type="number"
                inputmode="numeric"
                pattern="[0-9]{6}"
                maxlength="6"
                min="0"
                required
                class="tp-input"
              />
            </div>

            <div class="tp-actions">
              <button
                type="button"
                @click="close"
                class="tp-btn-cancel"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="loading"
                class="tp-btn-submit"
              >
                {{ loading ? 'Verifying...' : 'Verify & Change' }}
              </button>
            </div>
          </form>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
@reference "tailwindcss";
.fade-enter-active,
.fade-leave-active {
  @apply transition-opacity duration-200 ease-out;
}
.fade-enter-from,
.fade-leave-to {
  @apply opacity-0;
}

.tp-backdrop {
  @apply fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4;
}
.tp-modal {
  @apply w-full max-w-md bg-white dark:bg-zinc-900 rounded-2xl shadow-xl p-8 border border-transparent dark:border-zinc-800;
}
.tp-header {
  @apply flex justify-between items-center mb-6;
}
.tp-title {
  @apply text-xl font-black tracking-tight text-zinc-900 dark:text-zinc-50;
}
.tp-close-btn {
  @apply text-zinc-500 hover:text-zinc-700 dark:text-zinc-400 dark:hover:text-zinc-200;
}
.tp-close-icon {
  @apply w-6 h-6;
}
.tp-form {
  @apply space-y-4;
}
.tp-error {
  @apply p-3 rounded-lg bg-red-50 dark:bg-red-900/30 text-red-600 dark:text-red-400 text-sm font-medium;
}
.tp-label {
  @apply block text-sm font-bold text-zinc-700 dark:text-zinc-300 mb-1;
}
.tp-input {
  @apply w-full px-4 py-2.5 rounded-xl border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-950 text-zinc-900 dark:text-zinc-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 font-mono text-sm;
}
.tp-hint {
  @apply mt-1 text-xs text-zinc-500 dark:text-zinc-400 font-medium;
}
.tp-actions {
  @apply pt-4 flex justify-end gap-3;
}
.tp-btn-cancel {
  @apply px-4 py-2 rounded-xl border border-zinc-300 dark:border-zinc-700 text-zinc-700 dark:text-zinc-300 hover:bg-zinc-50 dark:hover:bg-zinc-800 font-bold;
}
.tp-btn-submit {
  @apply px-4 py-2 rounded-xl bg-indigo-600 text-white font-bold hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 disabled:opacity-50;
}
</style>
