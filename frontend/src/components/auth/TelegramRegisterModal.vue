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
    <div v-if="modelValue" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
      <div class="w-full max-w-md bg-white dark:bg-gray-900 rounded-2xl shadow-xl p-8">
        <div class="flex justify-between items-center mb-6">
          <h2 class="text-xl font-bold text-gray-900 dark:text-white">{{ title }}</h2>
          <button @click="$emit('update:modelValue', false)" class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div v-if="loading" class="text-center py-8">
          <div class="inline-block animate-spin rounded-full h-10 w-10 border-4 border-indigo-500 border-t-transparent mb-4"></div>
          <p class="text-gray-600 dark:text-gray-300 font-medium">Generating registration token...</p>
        </div>

        <div v-else class="space-y-4">
          <div v-if="error" class="p-3 rounded-lg bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 text-sm">
            {{ error }}
          </div>

          <div class="text-center">
            <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-blue-100 dark:bg-blue-900/30 mb-4">
              <svg class="w-8 h-8 text-blue-600 dark:text-blue-400" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm-1-13h2v6h-2zm0 8h2v2h-2z"/>
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">Link Your Telegram</h3>
            <p class="text-gray-600 dark:text-gray-300 text-sm">
              Click the button below to open Telegram and link your account using the registration token.
            </p>
          </div>

          <div v-if="telegramUrl" class="space-y-3">
            <button
              @click="openTelegram"
              class="w-full flex items-center justify-center gap-2 px-4 py-3 rounded-xl bg-blue-600 text-white font-medium hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm-1-13h2v6h-2zm0 8h2v2h-2z"/>
              </svg>
              Open Telegram (@nauva_kiroku_bot)
            </button>

            <div class="p-3 rounded-lg bg-gray-50 dark:bg-gray-800 text-sm">
              <p class="font-medium text-gray-900 dark:text-white mb-1">Registration Token:</p>
              <div class="flex items-center gap-2">
                <code class="flex-1 font-mono text-xs bg-white dark:bg-gray-700 px-2 py-1.5 rounded">{{ token }}</code>
                <button
                  @click="copyToken"
                  class="px-3 py-1.5 text-xs rounded-lg border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
                >
                  Copy
                </button>
              </div>
            </div>

            <p class="text-xs text-gray-500 dark:text-gray-400 text-center">
              After opening Telegram, send the token with: <code class="bg-gray-100 dark:bg-gray-800 px-1 rounded">/register &lt;token&gt;</code>
            </p>
          </div>

          <div class="pt-4 flex justify-end gap-3">
            <button
              type="button"
              @click="$emit('update:modelValue', false)"
              class="px-4 py-2 rounded-xl border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 font-medium"
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
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>