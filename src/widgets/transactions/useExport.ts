import { exportToExcel } from '@/utils/exportExcel'
import { useTransactionsStore } from '@/stores/transactions'
import { useUIStore } from '@/stores/ui'
import { formatAmountForDisplay } from '@/utils/amountDisplay'

export function useExport() {
  const transactionsStore = useTransactionsStore()
  const uiStore = useUIStore()

  const exportCSV = async () => {
    try {
      const transactions = transactionsStore.transactions
      const headers = ['Date', 'Type', 'Amount', 'Category', 'Account', 'Note', 'Status']
      const rows = transactions.map((tx: any) => [
        tx.date_time,
        tx.type,
        formatAmountForDisplay(tx.amount, 'transfer'),
        tx.category_name || '',
        tx.from_account_name || '',
        tx.note || '',
        tx.status
      ])
      const csv = [headers.join(','), ...rows.map(r => r.map(String).join(','))].join('\n')
      const url = window.URL.createObjectURL(new Blob([csv], { type: 'text/csv' }))
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
      const transactions = transactionsStore.transactions
      const headers = ['Date', 'Type', 'Amount', 'Category', 'Account', 'Note', 'Status']
      const rows = transactions.map((tx: any) => [
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