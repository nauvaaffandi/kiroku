<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

defineProps<{
  modelValue: boolean
  title: string
  loading: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'register', token: string, telegramUrl: string): void
}>()

const telegramUrl = ref<string>('')
const token = ref<string>('')
const error = ref<string>('')
let checkInterval: ReturnType<typeof setInterval> | null = null

const openTelegram = () => {
  if (telegramUrl.value) {
    window.open(telegramUrl.value, '_blank')
  }
}

const copyToken = async () => {
  if (token.value) {
    await navigator.clipboard.writeText(token.value)
  }
}

const checkRegistration = async () => {
  try {
    const response = await fetch('/api/v1/auth/telegram/status', {
      credentials: 'include'
    })
    if (response.ok) {
      const data = await response.json()
      if (data.linked) {
        emit('update:modelValue', false)
        emit('register', '', '')
      }
    }
  } catch {
    // ignore
  }
}

onMounted(() => {
  // Start polling for registration completion
  checkInterval = setInterval(checkRegistration, 3000)
  checkRegistration()
})

onUnmounted(() => {
  if (checkInterval) clearInterval(checkInterval)
})
</script>

<template>
  <Transition name="fade">
<div v-if="modelValue" class="trm-backdrop">
        <div class="trm-modal">
          <div class="trm-header">
            <h2 class="trm-title">{{ title }}</h2>
            <button @click="$emit('update:modelValue', false)" class="trm-close-btn">
              <svg class="trm-close-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <div v-if="loading" class="trm-loading">
            <div class="trm-spinner"></div>
            <p class="trm-loading-text">Generating registration token...</p>
          </div>

          <div v-else class="trm-form">
            <div v-if="error" class="trm-error">
              {{ error }}
            </div>

            <div class="trm-icon-wrap">
              <div class="trm-icon-circle">
                <svg class="trm-icon" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm-1-13h2v6h-2zm0 8h2v2h-2z"/>
                </svg>
              </div>
              <h3 class="trm-heading">Link Your Telegram</h3>
              <p class="trm-subtext">
                Click the button below to open Telegram and link your account using the registration token.
              </p>
            </div>

            <div v-if="telegramUrl" class="trm-actions">
              <button
                @click="openTelegram"
                class="trm-btn-primary"
              >
                <svg class="trm-btn-icon" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm-1-13h2v6h-2zm0 8h2v2h-2z"/>
                </svg>
                Open Telegram (@nauva_kiroku_bot)
              </button>

              <div class="trm-token-box">
                <p class="trm-token-label">Registration Token:</p>
                <div class="trm-token-row">
                  <code class="trm-token-code">{{ token }}</code>
                  <button
                    @click="copyToken"
                    class="trm-copy-btn"
                  >
                    Copy
                  </button>
                </div>
              </div>

              <p class="trm-hint">
                After opening Telegram, send the token with: <code class="trm-hint-code">/register &lt;token&gt;</code>
              </p>
            </div>

            <div class="trm-footer">
              <button
                type="button"
                @click="$emit('update:modelValue', false)"
                class="trm-btn-cancel"
              >
                Close
              </button>
            </div>
          </div>
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

.trm-backdrop {
  @apply fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4;
}
.trm-modal {
  @apply w-full max-w-md bg-white dark:bg-zinc-900 rounded-2xl shadow-xl p-8 border border-transparent dark:border-zinc-800;
}
.trm-header {
  @apply flex justify-between items-center mb-6;
}
.trm-title {
  @apply text-xl font-black tracking-tight text-zinc-900 dark:text-zinc-50;
}
.trm-close-btn {
  @apply text-zinc-500 hover:text-zinc-700 dark:text-zinc-400 dark:hover:text-zinc-200;
}
.trm-close-icon {
  @apply w-6 h-6;
}
.trm-loading {
  @apply text-center py-8;
}
.trm-spinner {
  @apply inline-block animate-spin rounded-full h-10 w-10 border-4 border-indigo-500 border-t-transparent mb-4;
}
.trm-loading-text {
  @apply text-zinc-600 dark:text-zinc-300 font-bold;
}
.trm-form {
  @apply space-y-4;
}
.trm-error {
  @apply p-3 rounded-lg bg-red-50 dark:bg-red-900/30 text-red-600 dark:text-red-400 text-sm font-medium;
}
.trm-icon-wrap {
  @apply text-center;
}
.trm-icon-circle {
  @apply inline-flex items-center justify-center w-16 h-16 rounded-full bg-blue-100 dark:bg-blue-950 text-blue-600 dark:text-blue-400 mb-4;
}
.trm-icon {
  @apply w-8 h-8;
}
.trm-heading {
  @apply text-lg font-black tracking-tight text-zinc-900 dark:text-zinc-50 mb-2;
}
.trm-subtext {
  @apply text-zinc-600 dark:text-zinc-300 text-sm font-medium;
}
.trm-actions {
  @apply space-y-3;
}
.trm-btn-primary {
  @apply w-full flex items-center justify-center gap-2 px-4 py-3 rounded-xl bg-blue-600 text-white font-bold hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500;
}
.trm-btn-icon {
  @apply w-5 h-5;
}
.trm-token-box {
  @apply p-3 rounded-lg bg-zinc-50 dark:bg-zinc-950 text-sm border border-zinc-200 dark:border-zinc-800;
}
.trm-token-label {
  @apply font-bold text-zinc-900 dark:text-zinc-100 mb-1;
}
.trm-token-row {
  @apply flex items-center gap-2;
}
.trm-token-code {
  @apply flex-1 font-mono text-xs bg-white dark:bg-zinc-900 text-zinc-900 dark:text-zinc-100 px-2 py-1.5 rounded border border-zinc-200 dark:border-zinc-800;
}
.trm-copy-btn {
  @apply px-3 py-1.5 text-xs rounded-lg border border-zinc-300 dark:border-zinc-700 text-zinc-700 dark:text-zinc-300 hover:bg-zinc-100 dark:hover:bg-zinc-800 font-bold;
}
.trm-hint {
  @apply text-xs text-zinc-500 dark:text-zinc-400 text-center font-medium;
}
.trm-hint-code {
  @apply bg-zinc-100 dark:bg-zinc-950 px-1 rounded text-zinc-900 dark:text-zinc-200;
}
.trm-footer {
  @apply pt-4 flex justify-end gap-3;
}
.trm-btn-cancel {
  @apply px-4 py-2 rounded-xl border border-zinc-300 dark:border-zinc-700 text-zinc-700 dark:text-zinc-300 hover:bg-zinc-50 dark:hover:bg-zinc-800 font-bold;
}
</style>