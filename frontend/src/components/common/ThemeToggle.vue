<template>
  <div
    @click="toggleTheme"
    class="theme-toggle"
    role="button"
    tabindex="0"
    aria-label="Toggle theme"
  >
    <span class="theme-toggle-label">{{ isDark ? 'Dark Mode' : 'Light Mode' }}</span>
    <div class="theme-switch">
      <div
        class="theme-switch-thumb"
        :class="{ 'translate-x-5': isDark }"
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

<style scoped>
@reference "tailwindcss";

.theme-toggle {
  @apply flex items-center justify-between gap-3 font-semibold px-3 py-2 cursor-pointer w-full select-none rounded-lg;
  background-color: var(--bg-card);
  color: var(--text-primary);
  border: 1px solid var(--border-primary);
}

.theme-toggle:hover {
  background-color: var(--bg-hover);
}

.theme-toggle-label {
  @apply text-sm pointer-events-none;
}

.theme-switch {
  @apply relative inline-flex items-center h-6 w-11 shrink-0 rounded-full;
  background-color: var(--border-primary);
  transition: background-color 0.2s ease;
}

.theme-switch-thumb {
  @apply absolute left-0.5 h-5 w-5 rounded-full bg-white shadow-md transition-transform duration-200;
}
</style>
