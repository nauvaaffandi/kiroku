<template>
  <div
    @click="toggleTheme"
    class="flex items-center justify-between gap-3 font-semibold text-zinc-600 dark:text-zinc-200 hover:bg-zinc-100 dark:hover:bg-zinc-800 rounded-lg px-3 py-2 cursor-pointer w-full select-none"
    role="button"
    tabindex="0"
    aria-label="Toggle theme"
  >
    <span class="text-sm pointer-events-none">{{ isDark ? 'Dark Mode' : 'Light Mode' }}</span>
    <div
      class="relative inline-flex items-center h-6 w-11 shrink-0 rounded-full bg-zinc-300 dark:bg-zinc-700 transition-colors duration-200 pointer-events-none"
    >
      <div
        :class="[
          'absolute left-0.5 h-5 w-5 rounded-full bg-white shadow-md transition-transform duration-200',
          isDark ? 'translate-x-5' : 'translate-x-0'
        ]"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useSettingsStore } from '@/stores/settings'

const settingsStore = useSettingsStore()

const isDark = computed(() => {
  return document.documentElement.classList.contains('dark') || settingsStore.theme === 'dark'
})

const toggleTheme = () => {
  const newTheme = isDark.value ? 'light' : 'dark'
  document.documentElement.classList.toggle('dark', newTheme === 'dark')
  if (newTheme === 'dark') {
    document.documentElement.setAttribute('data-theme', 'dark')
  } else {
    document.documentElement.setAttribute('data-theme', 'light')
  }
  localStorage.setItem('theme', newTheme)
  settingsStore.updateSettings({ theme: newTheme })
}
</script>
