import { ref } from 'vue'
import { useSettingsStore } from '@/stores/settings'
import { useUIStore } from '@/stores/ui'

export function useSettingsAuth() {
  const settingsStore = useSettingsStore()
  const uiStore = useUIStore()

  const totpLoading = ref(false)
  const totpSetupData = ref<{ qr_code: string; secret: string } | null>(null)
  const showTotpSetup = ref(false)
  const showBackupModal = ref(false)
  const showTelegramPinModal = ref(false)
  const showTelegramRegister = ref(false)
  const telegramLoading = ref(false)
  const pinChangeLoading = ref(false)
  const isTelegramLinked = ref(false)

  const checkTelegramStatus = async () => {
    isTelegramLinked.value = false
  }

  const openTelegramRegister = async () => {
    telegramLoading.value = false
  }

  const onTelegramRegistered = () => {
    isTelegramLinked.value = true
    uiStore.success('Telegram linked successfully')
  }

  const requestPinChange = async () => {
    pinChangeLoading.value = false
  }

  const handleVerifyPinChange = async (_token: string, _newPin: string) => {
    uiStore.error('PIN change not available offline')
  }

  const enableTotp = async () => {
    totpLoading.value = false
  }

  const disableTotp = async () => {
    if (confirm('Are you sure you want to disable 2FA?')) {
      settingsStore.settings.totp_enabled = false
      uiStore.success('2FA disabled')
    }
    totpLoading.value = false
  }

  const verifyTotp = async (_code: string) => {
    settingsStore.settings.totp_enabled = true
    showTotpSetup.value = false
    totpSetupData.value = null
    uiStore.success('2FA enabled successfully')
  }

  return {
    totpLoading, totpSetupData, showTotpSetup, showBackupModal,
    showTelegramPinModal, showTelegramRegister, telegramLoading, pinChangeLoading, isTelegramLinked,
    checkTelegramStatus, openTelegramRegister, onTelegramRegistered, requestPinChange,
    handleVerifyPinChange, enableTotp, disableTotp, verifyTotp
  }
}