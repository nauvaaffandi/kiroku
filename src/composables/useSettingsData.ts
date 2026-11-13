import { useSettingsStore } from '@/stores/settings'
import { useUIStore } from '@/stores/ui'

export function useSettingsData() {
  const settingsStore = useSettingsStore()
  const uiStore = useUIStore()

  const createBackup = async () => {
    try {
      const response = await settingsStore.createBackup() as { data: { filename?: string } }
      const blob = new Blob([JSON.stringify(response.data, null, 2)], { type: 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = response.data.filename || 'backup.json'
      a.click()
      URL.revokeObjectURL(url)
      uiStore.success('Backup downloaded')
    } catch (e) { uiStore.error('Failed to create backup') }
  }

  const restoreBackup = async (file: File) => {
    const password = prompt('Enter backup password:')
    if (!password) return
    try {
      await settingsStore.restoreBackup(file, password)
      uiStore.success('Backup restored. Reloading...')
      window.location.reload()
    } catch (e) { uiStore.error('Failed to restore backup') }
  }

  const exportCSV = async (_type: string) => {
    uiStore.error('CSV export not available offline')
  }

  const importCSV = async (_file: File, _type: string) => {
    uiStore.error('CSV import not available offline')
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

  return { createBackup, restoreBackup, exportCSV, importCSV, confirmClearData }
}
