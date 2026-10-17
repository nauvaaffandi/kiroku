<template>
  <div class="language-switcher">
    <label class="language-label">{{ t('settings.language') }}</label>
    <select
      v-model="localeStore.currentLocale"
      @change="changeLocale"
      class="language-select"
      :aria-label="t('settings.language')"
    >
      <option v-for="l in localeStore.availableLocales" :key="l.code" :value="l.code">
        {{ l.name }}
      </option>
    </select>
  </div>
</template>

<script setup lang="ts">
import { useI18nStore } from '@/locales'

const localeStore = useI18nStore()
const { t } = localeStore

function changeLocale() {
  localeStore.setLocale(localeStore.currentLocale)
}
</script>

<style scoped>
.language-switcher {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}

.language-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-secondary);
}

.language-select {
  padding: 0.5rem 0.75rem;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-sm);
  background-color: var(--bg-card);
  color: var(--text-primary);
  font-size: 0.875rem;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.language-select:focus {
  outline: none;
  border-color: var(--accent);
  box-shadow: 0 0 0 2px var(--accent-light);
}
</style>