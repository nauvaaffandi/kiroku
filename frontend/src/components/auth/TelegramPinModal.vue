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
  (e: 'submit', token: string, newPassword: string): void
}>()

const token = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const error = ref('')

const handleSubmit = (e: Event) => {
  e.preventDefault()
  if (newPassword.value !== confirmPassword.value) {
    error.value = 'New passwords do not match'
    return
  }
  if (!token.value.trim()) {
    error.value = 'Token is required'
    return
  }
  error.value = ''
  emit('submit', token.value.trim(), newPassword.value)
}

const close = () => {
  token.value = ''
  newPassword.value = ''
  confirmPassword.value = ''
  error.value = ''
  emit('update:modelValue', false)
}
</script>

<template>
  <Transition name="fade">
    <div v-if="modelValue" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
      <div class="w-full max-w-md bg-white dark:bg-gray-900 rounded-2xl shadow-xl p-8">
        <div class="flex justify-between items-center mb-6">
          <h2 class="text-xl font-bold text-gray-900 dark:text-white">{{ title }}</h2>
          <button @click="close" class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <form @submit="handleSubmit" class="space-y-4">
          <div v-if="error" class="p-3 rounded-lg bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 text-sm">
            {{ error }}
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {{ mode === 'change' ? 'Telegram Token' : 'Reset Token from Telegram' }}
            </label>
            <input
              v-model="token"
              type="text"
              required
              placeholder="Paste token from @nauva_kiroku_bot"
              class="w-full px-4 py-2.5 rounded-xl border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500 font-mono text-sm"
            />
            <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
              {{ mode === 'change' 
                ? 'Request token via /changepassword command in Telegram' 
                : 'Request token via /forgotpassword command in Telegram' }}
            </p>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">New Password</label>
            <input
              v-model="newPassword"
              type="password"
              required
              class="w-full px-4 py-2.5 rounded-xl border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Confirm New Password</label>
            <input
              v-model="confirmPassword"
              type="password"
              required
              class="w-full px-4 py-2.5 rounded-xl border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>

          <div class="pt-4 flex justify-end gap-3">
            <button
              type="button"
              @click="close"
              class="px-4 py-2 rounded-xl border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 font-medium"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="loading"
              class="px-4 py-2 rounded-xl bg-indigo-600 text-white font-medium hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 disabled:opacity-50"
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
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>