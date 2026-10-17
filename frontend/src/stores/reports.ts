import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/api/client'
import type { SpendingReport, TrendReport, BalanceReport, BudgetReport, CashflowProjectionReport, GoalsReport, OutlookReport } from '@/types'

export const useReportsStore = defineStore('reports', () => {
  const loading = ref(false)

  const getSpending = async (params: Record<string, string>): Promise<SpendingReport> => {
    loading.value = true
    try {
      const response = await api.get('/api/v1/reports/spending', { params })
      return response.data
    } catch (e) {
      console.error('Failed to fetch spending report:', e)
      return { period: params.period, groups: [], total_income: 0, total_expense: 0, net: 0 }
    } finally {
      loading.value = false
    }
  }

  const getTrend = async (params: { months: number; account_id?: string; category_id?: string }): Promise<TrendReport> => {
    loading.value = true
    try {
      const response = await api.get('/api/v1/reports/trend', { params })
      return response.data
    } catch (e) {
      console.error('Failed to fetch trend report:', e)
      return { series: [] }
    } finally {
      loading.value = false
    }
  }

  const getBalance = async (params: Record<string, string>): Promise<BalanceReport> => {
    loading.value = true
    try {
      const response = await api.get('/api/v1/reports/balance', { params })
      return response.data
    } catch (e) {
      console.error('Failed to fetch balance report:', e)
      return { series: [] }
    } finally {
      loading.value = false
    }
  }

  const getBudget = async (params: Record<string, string>): Promise<BudgetReport> => {
    loading.value = true
    try {
      const response = await api.get('/api/v1/reports/budget', { params })
      return response.data
    } catch (e) {
      console.error('Failed to fetch budget report:', e)
      return { items: [], total_budget: 0, total_spent: 0, total_remaining: 0 }
    } finally {
      loading.value = false
    }
  }

  const getCashflowProjection = async (params: Record<string, string>): Promise<CashflowProjectionReport> => {
    loading.value = true
    try {
      const response = await api.get('/api/v1/reports/cashflow-projection', { params })
      return response.data
    } catch (e) {
      console.error('Failed to fetch cashflow projection:', e)
      return { series: [], accounts: [] }
    } finally {
      loading.value = false
    }
  }

  const getGoals = async (): Promise<GoalsReport> => {
    loading.value = true
    try {
      const response = await api.get('/api/v1/reports/goals')
      return response.data
    } catch (e) {
      console.error('Failed to fetch goals report:', e)
      return { goals: [] }
    } finally {
      loading.value = false
    }
  }

  const getOutlook = async (params: Record<string, string>): Promise<OutlookReport> => {
    loading.value = true
    try {
      const response = await api.get('/api/v1/reports/outlook', { params })
      return response.data
    } catch (e) {
      console.error('Failed to fetch outlook report:', e)
      return { items: [], summary: { total_income: 0, total_expense: 0, net: 0 } }
    } finally {
      loading.value = false
    }
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