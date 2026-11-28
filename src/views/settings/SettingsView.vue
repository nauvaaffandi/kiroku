<template>
  <div class="settings-page">
    <header class="page-header"><h1>{{ t('settings.title') }}</h1></header>

    <SettingsThemeSection
      :theme="theme" :theme-options="themeOptions" @theme-change="onThemeChange" />

    <SettingsCurrencyLanguageSection
      :settings="settings" :base-currency="settings.base_currency" :currency-options="currencyOptions" :date-format-options="dateFormatOptions"
      :number-format-options="numberFormatOptions" :week-start-options="weekStartOptions" :week-start-str="weekStartStr"
      :new-currency-code="newCurrencyCode" @update:new-currency-code="newCurrencyCode = $event" @add-currency="handleAddCurrency"
      @date-format-change="(v: any) => (settingsStore.settings.date_format = v)"
      @number-format-change="(v: any) => (settingsStore.settings.number_format = v)"
      @week-start-change="(v: any) => settingsStore.setSetting('week_start', parseInt(v, 10) as WeekStart)"
      @base-currency-change="(v: any) => (settingsStore.settings.base_currency = v)" />

    <SettingsDataSection @download-backup="createBackup" @clear-data="confirmClearData" />

    <SettingsOCRAISection />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useI18n } from '@/locales'
import { useSettingsStore } from '@/stores/settings'
import { useSettingsPage } from '@/composables/useSettingsPage'
import { useSettingsData } from '@/composables/useSettingsData'
import type { WeekStart } from '@/types'
import SettingsThemeSection from '@/views/settings/SettingsThemeSection.vue'
import SettingsCurrencyLanguageSection from '@/views/settings/SettingsCurrencyLanguageSection.vue'
import SettingsDataSection from '@/views/settings/SettingsDataSection.vue'
import SettingsOCRAISection from '@/views/settings/SettingsOCRAISection.vue'

const settingsStore = useSettingsStore()
const { t } = useI18n()
const {
  theme, weekStartStr, newCurrencyCode, currencyOptions, themeOptions,
  dateFormatOptions, numberFormatOptions, weekStartOptions, handleAddCurrency, onThemeChange,
} = useSettingsPage()
const { createBackup, confirmClearData } = useSettingsData()

const settings = settingsStore.settings

onMounted(async () => {
  settingsStore.fetchSettings()
})
</script>

<style scoped>
@reference "tailwindcss";
.page-header { @apply mb-8; }
.page-header h1 { @apply text-2xl font-semibold m-0 text-[var(--text-primary)]; }
</style>