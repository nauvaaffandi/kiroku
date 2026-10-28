<template>
  <div class="settings-page">
    <header class="page-header">
      <h1>{{ t('settings.title') }}</h1>
    </header>

    <SettingSection :title="t('settings.theme')">
      <SettingItem
        :label="t('settings.theme')"
        :description="t('settings.themeSystem')"
      >
        <BaseSelect
          v-model="theme"
          :options="themeOptions"
          @change="onThemeChange"
        />
      </SettingItem>
    </SettingSection>

    <SettingSection :title="t('settings.currency') + ' & ' + t('settings.language')">
      <SettingItem
        :label="t('settings.currency')"
        :description="t('settings.currencyDesc')"
      >
        <div class="flex flex-col gap-2 w-full">
          <BaseSelect
            v-model="settings.base_currency"
            :options="currencyOptions"
            @change="saveSettings"
          />
          <div class="flex gap-2 mt-2">
            <input
              v-model="newCurrencyCode"
              placeholder="Add Currency (e.g. IDR)"
              class="p-2 border border-[var(--border-primary)] rounded-lg bg-[var(--bg-primary)] text-[var(--text-primary)] text-sm flex-1"
              maxlength="5"
            />
            <BaseButton variant="secondary" @click="handleAddCurrency">Add</BaseButton>
          </div>
        </div>
      </SettingItem>

      <SettingItem
        :label="t('settings.language')"
        :description="t('settings.languageDesc')"
      >
        <LanguageSwitcher />
      </SettingItem>

      <SettingItem
        :label="t('settings.date_format')"
        :description="t('settings.dateFormatDesc')"
      >
        <BaseSelect
          v-model="settings.date_format"
          :options="dateFormatOptions"
          @change="saveSettings"
        />
      </SettingItem>

      <SettingItem
        :label="t('settings.number_format')"
        :description="t('settings.numberFormatDesc')"
      >
        <BaseSelect
          v-model="settings.number_format"
          :options="numberFormatOptions"
          @change="saveSettings"
        />
      </SettingItem>

      <SettingItem
        :label="t('settings.week_start')"
        :description="t('settings.weekStartDesc')"
      >
        <BaseSelect
          v-model="weekStartStr"
          :options="weekStartOptions"
          @change="onWeekStartChange"
        />
      </SettingItem>
    </SettingSection>

    <SettingSection :title="t('settings.telegramTitle')">
      <SettingItem
        :label="t('settings.telegram')"
        :description="isTelegramLinked ? t('settings.telegramLinked') : t('settings.telegramDesc')"
      >
        <template #default>
          <BaseButton 
            :variant="isTelegramLinked ? 'secondary' : 'primary'"
            @click="openTelegramRegister"
            :loading="telegramLoading"
          >
            {{ isTelegramLinked ? 'Linked (Re-link)' : t('settings.telegramRegister') }}
          </BaseButton>
        </template>
      </SettingItem>
    </SettingSection>

    <SettingSection :title="t('settings.security')">
      <SettingItem
        :label="t('settings.changePin')"
        :description="t('settings.changePinDesc')"
      >
        <template #default>
          <BaseButton variant="primary" @click="requestPinChange" :loading="pinChangeLoading">
            {{ t('settings.requestToken') }}
          </BaseButton>
        </template>
      </SettingItem>

      <SettingItem
        :label="t('settings.twoFactorAuth')"
        :description="settingsStore.isTotpEnabled ? t('settings.twoFactorEnabled') : t('settings.twoFactorDesc')"
      >
        <template #default>
          <BaseButton
            :variant="totpButtonVariant"
            @click="settingsStore.isTotpEnabled ? disableTotp : enableTotp"
            :loading="totpLoading"
          >
            {{ settingsStore.isTotpEnabled ? t('settings.disable2FA') : t('settings.enable2FA') }}
          </BaseButton>
        </template>
      </SettingItem>

      <SettingItem
        :label="t('settings.biometricAuth')"
        :description="status.enabled ? t('settings.biometricEnabledDesc') : t('settings.biometricDisabledDesc')"
      >
        <template #default>
          <button
            @click="toggleBiometric"
            class="biometric-toggle"
            :class="{ active: status.enabled }"
            :disabled="!status.supported || biometricLoading"
          >
            {{ status.enabled ? t('settings.disableBiometric') : t('settings.enableBiometric') }}
          </button>
          <span v-if="!status.supported" class="text-xs text-[#dc3545] ml-2">{{ t('settings.biometricNotSupported') }}</span>
        </template>
      </SettingItem>
    </SettingSection>

    <SettingSection :title="t('settings.data')">
      <SettingItem
        :label="t('settings.backup')"
        :description="t('settings.backupDesc')"
      >
        <BaseButton variant="primary" @click="showBackupModal = true">
          {{ t('settings.manageBackup') }}
        </BaseButton>
      </SettingItem>

      <SettingItem
        :label="t('settings.clearData')"
        :description="t('settings.clearDataDesc')"
      >
        <BaseButton variant="danger" @click="confirmClearData">
          {{ t('settings.clearAllData') }}
        </BaseButton>
      </SettingItem>
    </SettingSection>

    <TelegramPinModal
      v-model="showTelegramPinModal"
      :title="t('settings.changePin')"
      :loading="pinChangeLoading"
      mode="change"
      @submit="handleVerifyPinChange"
    />

    <TelegramRegisterModal
      v-model="showTelegramRegister"
      :title="t('settings.telegramTitle')"
      :loading="telegramLoading"
      @register="onTelegramRegistered"
    />

    <TotpSetupModal
      v-model="showTotpSetup"
      :title="t('settings.setup2FA')"
      :qrCode="totpSetupData?.qr_code || ''"
      :secret="totpSetupData?.secret || ''"
      :codeLabel="t('settings.enterCode')"
      :submitLabel="t('settings.verify')"
      @submit="verifyTotp"
    />

    <SettingsBackupModal
      v-model="showBackupModal"
      :title="t('settings.backup')"
      @create-backup="createBackup"
      @restore-backup="restoreBackup"
      @export-csv="exportCSV"
      @import-csv="importCSV"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useSettingsStore } from '@/stores/settings'
import { useUIStore } from '@/stores/ui'
import { useI18n } from '@/locales'
import { api } from '@/api/client'
import LanguageSwitcher from '@/components/forms/LanguageSwitcher.vue'
import type { Theme, WeekStart } from '@/types'
import { SettingItem, SettingSection, TotpSetupModal, SettingsBackupModal } from '@/widgets/settings'
import { BaseSelect, BaseButton } from '@/components/ui'
import TelegramPinModal from '@/components/auth/TelegramPinModal.vue'
import TelegramRegisterModal from '@/components/auth/TelegramRegisterModal.vue'
import { useCurrencies } from '@/composables/useCurrencies'
import { useBiometric } from '@/composables/useBiometric'

const settingsStore = useSettingsStore()
const uiStore = useUIStore()
const { t } = useI18n()
const { status, enable, disable, loading: biometricLoading, refreshStatus } = useBiometric()

const settings = computed(() => settingsStore.settings)

const theme = computed<Theme>({
  get: () => settingsStore.theme,
  set: (val) => {
    settingsStore.setSetting('theme', val)
  }
})

const weekStartStr = computed<string>({
  get: () => String(settings.value.week_start ?? 1),
  set: (val) => {
    settingsStore.setSetting('week_start', parseInt(val, 10) as WeekStart)
  }
})

const totpLoading = ref(false)
const totpSetupData = ref<{ qr_code: string; secret: string } | null>(null)
const showTotpSetup = ref(false)
const showBackupModal = ref(false)
const showTelegramPinModal = ref(false)
const showTelegramRegister = ref(false)
const telegramLoading = ref(false)
const pinChangeLoading = ref(false)
const isTelegramLinked = ref(false)

const toggleBiometric = async () => {
  if (status.value.enabled) {
    const success = await disable()
    if (success) {
      uiStore.success(t('settings.biometricDisabled'))
    }
  } else {
    const success = await enable()
    if (success) {
      uiStore.success(t('settings.biometricEnabled'))
    }
  }
}

const totpButtonVariant = computed<'primary' | 'danger'>(() =>
  settingsStore.isTotpEnabled ? 'danger' : 'primary'
)

const themeOptions = [
  { value: 'system', label: 'System' },
  { value: 'light', label: 'Light' },
  { value: 'dark', label: 'Dark' }
]

const { currencies, addCurrency } = useCurrencies()

const currencyOptions = computed(() => 
  currencies.value.map(c => ({ value: c, label: c }))
)

const newCurrencyCode = ref('')
const handleAddCurrency = async () => {
  const code = newCurrencyCode.value.trim().toUpperCase()
  if (!code) return
  const success = await addCurrency(code)
  if (success) {
    uiStore.success(`Currency ${code} added`)
    newCurrencyCode.value = ''
  } else {
    uiStore.error('Failed to add currency')
  }
}

const dateFormatOptions = [
  { value: 'YYYY-MM-DD', label: 'YYYY-MM-DD' },
  { value: 'DD/MM/YYYY', label: 'DD/MM/YYYY' },
  { value: 'MM/DD/YYYY', label: 'MM/DD/YYYY' },
  { value: 'YYYY年M月D日', label: 'YYYY年M月D日 (Japanese)' }
]

const numberFormatOptions = [
  { value: 'ja-JP', label: 'Japanese (1,234,567)' },
  { value: 'en-US', label: 'US English (1,234,567)' },
  { value: 'en-GB', label: 'UK English (1,234,567)' },
  { value: 'de-DE', label: 'German (1.234.567)' }
]

const weekStartOptions = [
  { value: '0', label: 'Sunday' },
  { value: '1', label: 'Monday' },
  { value: '6', label: 'Saturday' }
]

onMounted(async () => {
  settingsStore.fetchSettings()
  await refreshStatus()
  await checkTelegramStatus()
})

const checkTelegramStatus = async () => {
  try {
    const res = await api.get('/api/v1/auth/telegram/status')
    isTelegramLinked.value = res.data?.linked || false
  } catch {
    isTelegramLinked.value = false
  }
}

const openTelegramRegister = async () => {
  telegramLoading.value = true
  try {
    const res = await api.post('/api/v1/auth/telegram/register-token')
    if (res.data?.telegram_url) {
      window.open(res.data.telegram_url, '_blank')
      showTelegramRegister.value = true
    }
  } catch (e: any) {
    uiStore.error(e.response?.data || 'Failed to generate registration token')
  } finally {
    telegramLoading.value = false
  }
}

const onTelegramRegistered = () => {
  isTelegramLinked.value = true
  uiStore.success('Telegram linked successfully')
}

const requestPinChange = async () => {
  pinChangeLoading.value = true
  try {
    // First prompt for new PIN
    const newPin = prompt('Enter your new PIN:')
    if (!newPin) {
      pinChangeLoading.value = false
      return
    }

    const res = await api.post('/api/v1/auth/pin/change-request', { new_pin: newPin })
    if (res.status === 200) {
      uiStore.success(t('settings.tokenRequested'))
      // Store temporary new PIN in sessionStorage or state
      sessionStorage.setItem('kiroku_temp_new_pin', newPin)
      showTelegramPinModal.value = true
    }
  } catch (e: any) {
    uiStore.error(e.response?.data || 'Failed to request PIN change')
  } finally {
    pinChangeLoading.value = false
  }
}

const handleVerifyPinChange = async (token: string, newPin: string) => {
  pinChangeLoading.value = true
  try {
    await api.post('/api/v1/auth/pin/verify', { token, new_pin: newPin })
    uiStore.success(t('settings.pinChangeSuccess'))
    showTelegramPinModal.value = false
  } catch (e: any) {
    uiStore.error(e.response?.data || 'Invalid or expired token')
  } finally {
    pinChangeLoading.value = false
  }
}

const onThemeChange = async (value: string) => {
  await settingsStore.setSetting('theme', value as Theme)
}

const onWeekStartChange = async (value: string) => {
  await settingsStore.setSetting('week_start', parseInt(value, 10) as WeekStart)
}

const saveSettings = async () => {
  await settingsStore.updateSettings(settings.value)
  uiStore.success('Settings saved')
}

const enableTotp = async () => {
  totpLoading.value = true
  try {
    const response = await api.post('/api/v1/settings/totp/enable')
    if (response.data) {
      totpSetupData.value = response.data
      showTotpSetup.value = true
    }
  } catch (e) {
    uiStore.error('Failed to enable 2FA')
  } finally {
    totpLoading.value = false
  }
}

const disableTotp = async () => {
  if (confirm('Are you sure you want to disable 2FA?')) {
    totpLoading.value = true
    try {
      await api.post('/api/v1/settings/totp/disable')
      settingsStore.settings.totp_enabled = false
      uiStore.success('2FA disabled')
    } catch (e) {
      uiStore.error('Failed to disable 2FA')
    } finally {
      totpLoading.value = false
    }
  }
}

const verifyTotp = async (code: string) => {
  totpLoading.value = true
  try {
    await api.post('/api/v1/settings/totp/verify', { code })
    settingsStore.settings.totp_enabled = true
    showTotpSetup.value = false
    totpSetupData.value = null
    uiStore.success('2FA enabled successfully')
  } catch (e) {
    uiStore.error('Invalid code')
  } finally {
    totpLoading.value = false
  }
}

const createBackup = async () => {
  try {
    const response = await settingsStore.createBackup()
    const blob = new Blob([JSON.stringify(response.data, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = response.data.filename || 'backup.json'
    a.click()
    URL.revokeObjectURL(url)
    uiStore.success('Backup downloaded')
  } catch (e) {
    uiStore.error('Failed to create backup')
  }
}

const restoreBackup = async (file: File) => {
  const password = prompt('Enter backup password:')
  if (password) {
    try {
      await settingsStore.restoreBackup(file, password)
      uiStore.success('Backup restored. Reloading...')
      window.location.reload()
    } catch (e) {
      uiStore.error('Failed to restore backup')
    }
  }
}

const exportCSV = async (type: string) => {
  try {
    const response = await api.get(`/api/v1/export/${type}?format=csv`)
    const blob = new Blob([response.data], { type: 'text/csv' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${type}.csv`
    a.click()
    URL.revokeObjectURL(url)
  } catch (e) {
    uiStore.error(`Failed to export ${type}`)
  }
}

const importCSV = async (file: File, type: string) => {
  try {
    const formData = new FormData()
    formData.append('file', file)
    await api.post(`/api/v1/import/${type}`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    uiStore.success('Import completed')
  } catch (e) {
    uiStore.error('Failed to import')
  }
}

const confirmClearData = async () => {
  const confirm1 = confirm('This will delete ALL your data. Are you sure?')
  if (!confirm1) return
  
  const confirm2 = confirm('This action CANNOT be undone. Type "DELETE" to confirm.')
  if (!confirm2) return
  
  try {
    await settingsStore.clearAllData()
    uiStore.success('All data cleared. Reloading...')
    window.location.reload()
  } catch (e) {
    uiStore.error('Failed to clear data')
  }
}
</script>

<style scoped>
@reference "tailwindcss";
.page-header {
  @apply mb-8;
}

.page-header h1 {
  @apply text-2xl font-bold m-0;
  color: var(--text-primary);
}

.biometric-toggle {
  @apply px-4 py-2 font-medium text-sm cursor-pointer transition-all duration-200 border;
  border-radius: var(--radius-sm);
  background-color: var(--bg-secondary);
  color: var(--text-primary);
  border-color: var(--border-primary);
}

.biometric-toggle.active {
  @apply text-white;
  background-color: var(--success);
  border-color: var(--success);
}
</style>