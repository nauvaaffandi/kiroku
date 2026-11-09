import { api } from '@/api/client'
import { exportToExcel } from '@/utils/exportExcel'
import { useTransactionsStore } from '@/stores/transactions'
import { useUIStore } from '@/stores/ui'

export function useExport() {
  const transactionsStore = useTransactionsStore()
  const uiStore = useUIStore()

  const exportCSV = async () => {
    try {
      const params = { ...transactionsStore.filters }
      const response = await api.get('/api/v1/transactions/export', {
        params,
        responseType: 'blob'
      })
      const url = window.URL.createObjectURL(new Blob([response.data]))
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', `transactions-${new Date().toISOString().split('T')[0]}.csv`)
      document.body.appendChild(link)
      link.click()
      link.remove()
      uiStore.success('CSV exported successfully')
    } catch (e) {
      uiStore.error('Failed to export CSV')
    }
  }

  const exportExcel = async () => {
    try {
      const params = { page: 1, limit: 10000, ...transactionsStore.filters }
      const response = await api.get('/api/v1/transactions', { params })
      const data = Array.isArray(response.data) ? response.data : []
      const headers = ['Date', 'Type', 'Amount', 'Category', 'Account', 'Note', 'Status']
      const rows = data.map((tx: any) => [
        tx.date_time,
        tx.type,
        tx.amount,
        tx.category_name || '',
        tx.from_account_name || '',
        tx.note || '',
        tx.status
      ])
      exportToExcel(headers, rows, `transactions-${new Date().toISOString().split('T')[0]}`)
      uiStore.success('Excel exported successfully')
    } catch (e) {
      uiStore.error('Failed to export Excel')
    }
  }

  return { exportCSV, exportExcel }
}
