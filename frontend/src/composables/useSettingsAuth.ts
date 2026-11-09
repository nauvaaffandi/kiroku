import { ref } from 'vue'
import { useSettingsStore } from '@/stores/settings'
import { useUIStore } from '@/stores/ui'
import { api } from '@/api/client'

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
      const newPin = prompt('Enter your new PIN:')
      if (!newPin) { pinChangeLoading.value = false; return }
      const res = await api.post('/api/v1/auth/pin/change-request', { new_pin: newPin })
      if (res.status === 200) {
        uiStore.success('Token requested')
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
      uiStore.success('PIN changed successfully')
      showTelegramPinModal.value = false
    } catch (e: any) {
      uiStore.error(e.response?.data || 'Invalid or expired token')
    } finally {
      pinChangeLoading.value = false
    }
  }

  const enableTotp = async () => {
    totpLoading.value = true
    try {
      const response = await api.post('/api/v1/settings/totp/enable')
      if (response.data) {
        totpSetupData.value = response.data
        showTotpSetup.value = true
      }
    } catch (e) { uiStore.error('Failed to enable 2FA') }
    finally { totpLoading.value = false }
  }

  const disableTotp = async () => {
    if (!confirm('Are you sure you want to disable 2FA?')) return
    totpLoading.value = true
    try {
      await api.post('/api/v1/settings/totp/disable')
      settingsStore.settings.totp_enabled = false
      uiStore.success('2FA disabled')
    } catch (e) { uiStore.error('Failed to disable 2FA') }
    finally { totpLoading.value = false }
  }

  const verifyTotp = async (code: string) => {
    totpLoading.value = true
    try {
      await api.post('/api/v1/settings/totp/verify', { code })
      settingsStore.settings.totp_enabled = true
      showTotpSetup.value = false
      totpSetupData.value = null
      uiStore.success('2FA enabled successfully')
    } catch (e) { uiStore.error('Invalid code') }
    finally { totpLoading.value = false }
  }

  return {
    totpLoading, totpSetupData, showTotpSetup, showBackupModal,
    showTelegramPinModal, showTelegramRegister, telegramLoading, pinChangeLoading, isTelegramLinked,
    checkTelegramStatus, openTelegramRegister, onTelegramRegistered, requestPinChange,
    handleVerifyPinChange, enableTotp, disableTotp, verifyTotp
  }
}
