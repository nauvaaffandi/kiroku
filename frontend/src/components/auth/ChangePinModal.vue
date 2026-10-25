<script setup lang="ts">
import { ref } from 'vue'

defineProps<{
  modelValue: boolean
  title: string
  loading: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'submit', oldPin: string, newPin: string): void
}>()

const oldPin = ref('')
const newPin = ref('')
const confirmPin = ref('')
const error = ref('')

const handleSubmit = (e: Event) => {
  e.preventDefault()
  if (newPin.value !== confirmPin.value) {
    error.value = 'New PINs do not match'
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
    <div v-if="modelValue" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
      <div class="w-full max-w-md bg-white dark:bg-zinc-900 rounded-2xl shadow-xl p-8 border border-transparent dark:border-zinc-800">
        <div class="flex justify-between items-center mb-6">
          <h2 class="text-xl font-black tracking-tight text-zinc-900 dark:text-zinc-50">{{ title }}</h2>
          <button @click="close" class="text-zinc-500 hover:text-zinc-700 dark:text-zinc-400 dark:hover:text-zinc-200">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div v-if="loading" class="text-center py-8">
          <div class="inline-block animate-spin rounded-full h-10 w-10 border-4 border-indigo-500 border-t-transparent mb-4"></div>
          <p class="text-zinc-600 dark:text-zinc-300 font-bold">Waiting for Telegram approval...</p>
          <p class="text-sm text-zinc-400 dark:text-zinc-400 mt-1 font-medium">Please reply YES on Telegram to approve PIN change.</p>
        </div>

        <form v-else @submit="handleSubmit" class="space-y-4">
          <div v-if="error" class="p-3 rounded-lg bg-red-50 dark:bg-red-900/30 text-red-600 dark:text-red-400 text-sm font-medium">
            {{ error }}
          </div>

          <div>
            <label class="block text-sm font-bold text-zinc-700 dark:text-zinc-300 mb-1">Current PIN</label>
            <input
              v-model="oldPin"
              type="password"
              inputmode="numeric"
              pattern="[0-9]{6}"
              maxlength="6"
              required
              class="w-full px-4 py-2.5 rounded-xl border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-950 text-zinc-900 dark:text-zinc-50 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>

          <div>
            <label class="block text-sm font-bold text-zinc-700 dark:text-zinc-300 mb-1">New PIN</label>
            <input
              v-model="newPin"
              type="password"
              inputmode="numeric"
              pattern="[0-9]{6}"
              maxlength="6"
              required
              class="w-full px-4 py-2.5 rounded-xl border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-950 text-zinc-900 dark:text-zinc-50 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>

          <div>
            <label class="block text-sm font-bold text-zinc-700 dark:text-zinc-300 mb-1">Confirm New PIN</label>
            <input
              v-model="confirmPin"
              type="password"
              inputmode="numeric"
              pattern="[0-9]{6}"
              maxlength="6"
              required
              class="w-full px-4 py-2.5 rounded-xl border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-950 text-zinc-900 dark:text-zinc-50 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>

          <div class="pt-4 flex justify-end gap-3">
            <button
              type="button"
              @click="close"
              class="px-4 py-2 rounded-xl border border-zinc-300 dark:border-zinc-700 text-zinc-700 dark:text-zinc-300 hover:bg-zinc-50 dark:hover:bg-zinc-800 font-bold"
            >
              Cancel
            </button>
            <button
              type="submit"
              class="px-4 py-2 rounded-xl bg-indigo-600 text-white font-bold hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            >
              Submit & Request Telegram Approval
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
