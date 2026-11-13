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

    <SettingsTelegramSection :is-telegram-linked="isTelegramLinked" :telegram-loading="telegramLoading" @open-register="openTelegramRegister" />

    <SettingsSecuritySection
      :totp-button-variant="totpButtonVariant" :totp-loading="totpLoading" :biometric-loading="biometricLoading" :status="status"
      @change-pin="requestPinChange" @toggle-totp="settingsStore.isTotpEnabled ? disableTotp : enableTotp" @toggle-biometric="toggleBiometric" />

    <SettingsDataSection @manage-backup="showBackupModal = true" @clear-data="confirmClearData" />

    <TelegramPinModal v-model="showTelegramPinModal" :title="t('settings.changePin')" :loading="pinChangeLoading" mode="change" @submit="handleVerifyPinChange" />
    <TelegramRegisterModal v-model="showTelegramRegister" :title="t('settings.telegramTitle')" :loading="telegramLoading" @register="onTelegramRegistered" />
    <TotpSetupModal v-model="showTotpSetup" :title="t('settings.setup2FA')" :qr-code="totpSetupData?.qr_code || ''" :secret="totpSetupData?.secret || ''" :code-label="t('settings.enterCode')" :submit-label="t('settings.verify')" @submit="verifyTotp" />
    <SettingsBackupModal v-model="showBackupModal" :title="t('settings.backup')" @create-backup="createBackup" @restore-backup="restoreBackup" @export-csv="exportCSV" @import-csv="importCSV" />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useI18n } from '@/locales'
import { useSettingsStore } from '@/stores/settings'
import { useSettingsPage } from '@/composables/useSettingsPage'
import { useSettingsAuth } from '@/composables/useSettingsAuth'
import { useSettingsData } from '@/composables/useSettingsData'
import type { WeekStart } from '@/types'
import SettingsThemeSection from '@/components/settings/SettingsThemeSection.vue'
import SettingsCurrencyLanguageSection from '@/components/settings/SettingsCurrencyLanguageSection.vue'
import SettingsTelegramSection from '@/components/settings/SettingsTelegramSection.vue'
import SettingsSecuritySection from '@/components/settings/SettingsSecuritySection.vue'
import SettingsDataSection from '@/components/settings/SettingsDataSection.vue'
import TelegramPinModal from '@/components/auth/TelegramPinModal.vue'
import TelegramRegisterModal from '@/components/auth/TelegramRegisterModal.vue'
import { TotpSetupModal, SettingsBackupModal } from '@/widgets/settings'

const settingsStore = useSettingsStore()
const { t } = useI18n()
const {
  theme, weekStartStr, newCurrencyCode, totpButtonVariant, currencyOptions, themeOptions,
  dateFormatOptions, numberFormatOptions, weekStartOptions, handleAddCurrency, onThemeChange,
  toggleBiometric, biometricLoading, status, refreshStatus
} = useSettingsPage()
const {
  totpLoading, totpSetupData, showTotpSetup, showBackupModal, showTelegramPinModal,
  showTelegramRegister, telegramLoading, pinChangeLoading, isTelegramLinked,
  openTelegramRegister, onTelegramRegistered, requestPinChange, handleVerifyPinChange,
  enableTotp, disableTotp, verifyTotp, checkTelegramStatus
} = useSettingsAuth()
const { createBackup, restoreBackup, exportCSV, importCSV, confirmClearData } = useSettingsData()

const settings = settingsStore.settings

onMounted(async () => {
  settingsStore.fetchSettings()
  await refreshStatus()
  await checkTelegramStatus()
})
</script>

<style scoped>
@reference "tailwindcss";
.page-header { @apply mb-8; }
.page-header h1 { @apply text-2xl font-bold m-0; color: var(--text-primary); }
.biometric-toggle {
  @apply px-4 py-2 font-medium text-sm cursor-pointer transition-all duration-200 border;
  border-radius: var(--radius-sm);
  background-color: var(--bg-secondary);
  color: var(--text-primary);
  border-color: var(--border-primary);
}
.biometric-toggle.active { @apply text-white; background-color: var(--success); border-color: var(--success); }
</style>
