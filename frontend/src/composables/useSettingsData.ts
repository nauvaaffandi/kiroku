import { useSettingsStore } from '@/stores/settings'
import { useUIStore } from '@/stores/ui'
import { api } from '@/api/client'

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
    } catch (e) { uiStore.error(`Failed to export ${type}`) }
  }

  const importCSV = async (file: File, type: string) => {
    try {
      const formData = new FormData()
      formData.append('file', file)
      await api.post(`/api/v1/import/${type}`, formData, { headers: { 'Content-Type': 'multipart/form-data' } })
      uiStore.success('Import completed')
    } catch (e) { uiStore.error('Failed to import') }
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
