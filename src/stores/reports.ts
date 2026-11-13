import { defineStore } from 'pinia'
import { ref } from 'vue'
import { localDb } from '@/storage/repository'
import type {
  SpendingReport,
  TrendReport,
  BalanceReport,
  BudgetReport,
  CashflowProjectionReport,
  GoalsReport,
  OutlookReport,
  Transaction,
  Category,
  Account
} from '@/types'

export const useReportsStore = defineStore('reports', () => {
  const loading = ref(false)

  const listTransactions = () => localDb.list<Transaction>('transaction')

  const getSpending = async (params: Record<string, string>): Promise<SpendingReport> => {
    loading.value = true
    try {
      const all = await listTransactions()
      const categories = await localDb.list<Category>('category')
      const byCategory = new Map<string, { category_name: string; category_icon: string; category_color: string; amount: number; count: number }>()
      let totalExpense = 0
      let totalIncome = 0
      for (const t of all) {
        if (t.type === 'income') totalIncome += t.amount
        if (t.type !== 'expense') continue
        const catId = t.category_id ?? 'uncategorized'
        const cat = categories.find((c) => c.id === catId)
        const entry = byCategory.get(catId) ?? {
          category_name: cat?.name ?? 'Uncategorized',
          category_icon: cat?.icon ?? '',
          category_color: cat?.color ?? '',
          amount: 0,
          count: 0
        }
        entry.amount += t.amount
        entry.count += 1
        byCategory.set(catId, entry)
        totalExpense += t.amount
      }
      const groups = [...byCategory.entries()].map(([category_id, entry]) => ({
        category_id,
        ...entry,
        percentage: totalExpense ? (entry.amount / totalExpense) * 100 : 0,
        children: []
      }))
      return { period: params.period ?? '', groups, total_income: totalIncome, total_expense: totalExpense, net: totalIncome - totalExpense }
    } catch (e) {
      console.error('Failed to compute spending report:', e)
      return { period: params.period, groups: [], total_income: 0, total_expense: 0, net: 0 }
    } finally {
      loading.value = false
    }
  }

  const getTrend = async (params: { months: number; account_id?: string; category_id?: string }): Promise<TrendReport> => {
    loading.value = true
    try {
      const all = await listTransactions()
      const now = new Date()
      const months: string[] = []
      for (let i = params.months - 1; i >= 0; i--) {
        const d = new Date(now.getFullYear(), now.getMonth() - i, 1)
        months.push(`${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`)
      }
      const series = months.map((period) => {
        const inMonth = all.filter((t) => t.date_time.slice(0, 7) === period)
        const income = inMonth.filter((t) => t.type === 'income').reduce((s, t) => s + t.amount, 0)
        const expense = inMonth.filter((t) => t.type === 'expense').reduce((s, t) => s + t.amount, 0)
        return { period, income, expense, net: income - expense }
      })
      return { series }
    } catch (e) {
      console.error('Failed to compute trend report:', e)
      return { series: [] }
    } finally {
      loading.value = false
    }
  }

  const getBalance = async (_params: Record<string, string>): Promise<BalanceReport> => {
    loading.value = true
    try {
      const accounts = await localDb.list<Account>('account')
      const series = accounts.map((a) => ({
        date: new Date().toISOString().slice(0, 10),
        balance: a.balance,
        account_id: a.id,
        account_name: a.name
      }))
      return { series }
    } finally {
      loading.value = false
    }
  }

  const getBudget = async (_params: Record<string, string>): Promise<BudgetReport> => {
    return { items: [], total_budget: 0, total_spent: 0, total_remaining: 0 }
  }

  const getCashflowProjection = async (_params: Record<string, string>): Promise<CashflowProjectionReport> => {
    return { series: [], accounts: [] }
  }

  const getGoals = async (): Promise<GoalsReport> => {
    return { goals: [] }
  }

  const getOutlook = async (_params: Record<string, string>): Promise<OutlookReport> => {
    return { items: [], summary: { total_income: 0, total_expense: 0, net: 0 } }
  }

  return {
    loading,
    getSpending,
    getTrend,
    getBalance,
    getBudget,
    getCashflowProjection,
    getGoals,
    getOutlook
  }
})