import { useSettingsStore } from '@/stores/settings'
import { useUIStore } from '@/stores/ui'
import { useI18n } from '@/locales'

export function useSettingsData() {
  const settingsStore = useSettingsStore()
  const uiStore = useUIStore()
  const { t } = useI18n()

  const createBackup = async () => {
    try {
      console.log('[SettingsData] Initiating backup download')
      const blob = await settingsStore.createBackup()
      const timestamp = new Date().toISOString().replace(/[:.]/g, '-').slice(0, 19)
      const filename = `kiroku-backup-${timestamp}.zip`

      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = filename
      try {
        document.body.appendChild(a)
        a.click()
        document.body.removeChild(a)
      } finally {
        URL.revokeObjectURL(url)
      }

      console.log(`[SettingsData] Backup downloaded as ${filename}`)
      uiStore.success(t('settings.backupDownloaded'))
    } catch (e) {
      console.error('[SettingsData] Failed to create backup:', e)
      uiStore.error(t('settings.backupFailed'))
    }
  }

  const confirmClearData = async () => {
    if (!confirm('This will delete ALL your data. Are you sure?')) return
    if (!confirm('This action CANNOT be undone. Type "DELETE" to confirm.')) return
    try {
      await settingsStore.clearAllData()
      uiStore.success('All data cleared. Reloading...')
      window.location.reload()
    } catch (e) { uiStore.error('Failed to clear data') }
  }

  return { createBackup, confirmClearData }
}
