<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from '@/locales'

defineProps<{
  modelValue: boolean
  title: string
  loading: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'submit', oldPin: string, newPin: string): void
}>()

const { t } = useI18n()

const oldPin = ref('')
const newPin = ref('')
const confirmPin = ref('')
const error = ref('')

const handleSubmit = (e: Event) => {
  e.preventDefault()
  if (newPin.value !== confirmPin.value) {
    error.value = t('auth.changePin.pinsMismatch')
    return
  }
  error.value = ''
  emit('submit', oldPin.value, newPin.value)
}

const close = () => {
  oldPin.value = ''
  newPin.value = ''
  confirmPin.value = ''
  error.value = ''
  emit('update:modelValue', false)
}
</script>

<template>
  <Transition name="fade">
      <div v-if="modelValue" class="cp-backdrop">
        <div class="cp-modal">
          <div class="cp-header">
            <h2 class="cp-title">{{ title }}</h2>
            <button @click="close" class="cp-close-btn">
              <svg class="cp-close-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <div v-if="loading" class="cp-loading">
            <div class="cp-spinner"></div>
            <p class="cp-loading-text">{{ t('auth.changePin.waitingApproval') }}</p>
            <p class="cp-loading-hint">{{ t('auth.changePin.approvalHint') }}</p>
          </div>

          <form v-else @submit="handleSubmit" class="cp-form">
            <div v-if="error" class="cp-error">
              {{ error }}
            </div>

            <div>
              <label class="cp-label">{{ t('auth.changePin.currentPin') }}</label>
              <input
                v-model="oldPin"
                type="number"
                inputmode="numeric"
                pattern="[0-9]{6}"
                maxlength="6"
                min="0"
                required
                class="cp-input"
              />
            </div>

            <div>
              <label class="cp-label">{{ t('auth.changePin.newPin') }}</label>
              <input
                v-model="newPin"
                type="number"
                inputmode="numeric"
                pattern="[0-9]{6}"
                maxlength="6"
                min="0"
                required
                class="cp-input"
              />
            </div>

            <div>
              <label class="cp-label">{{ t('auth.changePin.confirmNewPin') }}</label>
              <input
                v-model="confirmPin"
                type="number"
                inputmode="numeric"
                pattern="[0-9]{6}"
                maxlength="6"
                min="0"
                required
                class="cp-input"
              />
            </div>

            <div class="cp-actions">
              <button
                type="button"
                @click="close"
                class="cp-btn-cancel"
              >
                {{ t('common.cancel') }}
              </button>
              <button
                type="submit"
                class="cp-btn-submit"
              >
                {{ t('auth.changePin.submitRequest') }}
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

.cp-backdrop {
  @apply fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4;
}
.cp-modal {
  @apply w-full max-w-md bg-white dark:bg-zinc-900 rounded-2xl shadow-xl p-8 border border-transparent dark:border-zinc-800;
}
.cp-header {
  @apply flex justify-between items-center mb-6;
}
.cp-title {
  @apply text-xl font-black tracking-tight text-zinc-900 dark:text-zinc-50;
}
.cp-close-btn {
  @apply text-zinc-500 hover:text-zinc-700 dark:text-zinc-400 dark:hover:text-zinc-200;
}
.cp-close-icon {
  @apply w-6 h-6;
}
.cp-loading {
  @apply text-center py-8;
}
.cp-spinner {
  @apply inline-block animate-spin rounded-full h-10 w-10 border-4 border-indigo-500 border-t-transparent mb-4;
}
.cp-loading-text {
  @apply text-zinc-600 dark:text-zinc-300 font-bold;
}
.cp-loading-hint {
  @apply text-sm text-zinc-400 dark:text-zinc-400 mt-1 font-medium;
}
.cp-form {
  @apply space-y-4;
}
.cp-error {
  @apply p-3 rounded-lg bg-red-50 dark:bg-red-900/30 text-red-600 dark:text-red-400 text-sm font-medium;
}
.cp-label {
  @apply block text-sm font-bold text-zinc-700 dark:text-zinc-300 mb-1;
}
.cp-input {
  @apply w-full px-4 py-2.5 rounded-xl border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-950 text-zinc-900 dark:text-zinc-50 focus:outline-none focus:ring-2 focus:ring-indigo-500;
}
.cp-actions {
  @apply pt-4 flex justify-end gap-3;
}
.cp-btn-cancel {
  @apply px-4 py-2 rounded-xl border border-zinc-300 dark:border-zinc-700 text-zinc-700 dark:text-zinc-300 hover:bg-zinc-50 dark:hover:bg-zinc-800 font-bold;
}
.cp-btn-submit {
  @apply px-4 py-2 rounded-xl bg-indigo-600 text-white font-bold hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500;
}
</style>
