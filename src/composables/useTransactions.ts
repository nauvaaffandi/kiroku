import { useI18n } from '@/locales'
import { useUIStore } from '@/stores/ui'
import { transactionRepository } from '@/repository'

export function useTransactions() {
  const { t } = useI18n()
  const uiStore = useUIStore()
  const entity = 'transaction'

  const create = async (data: unknown) => {
    try {
      console.log('[useTransactions] create start', data)
      const id = await transactionRepository.create(data as any)
      console.log('[useTransactions] create success', { id })
      return id
    } catch (err: any) {
      const msg = err.message || t('repository.error.save', { entity })
      uiStore.error(msg)
      console.error('[useTransactions] create error', msg)
      throw err
    }
  }

  const update = async (id: string, data: unknown) => {
    try {
      console.log('[useTransactions] update start', { id, data })
      await transactionRepository.update(id, data as any)
      console.log('[useTransactions] update success', { id })
    } catch (err: any) {
      const msg = err.message || t('repository.error.update', { entity })
      uiStore.error(msg)
      console.error('[useTransactions] update error', msg)
      throw err
    }
  }

  const deleteTransaction = async (id: string) => {
    try {
      console.log('[useTransactions] delete start', { id })
      await transactionRepository.delete(id)
      console.log('[useTransactions] delete success', { id })
    } catch (err: any) {
      const msg = err.message || t('repository.error.delete', { entity })
      uiStore.error(msg)
      console.error('[useTransactions] delete error', msg)
      throw err
    }
  }

  const findById = async (id: string) => {
    try {
      console.log('[useTransactions] findById start', { id })
      const result = await transactionRepository.findById(id)
      console.log('[useTransactions] findById success', { id, found: !!result })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useTransactions] findById error', msg)
      throw err
    }
  }

  const findAll = async (filters?: Record<string, unknown>) => {
    try {
      console.log('[useTransactions] findAll start', { filters })
      const result = await transactionRepository.findAll(filters)
      console.log('[useTransactions] findAll success', { count: result.length })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useTransactions] findAll error', msg)
      throw err
    }
  }

  const count = async (filters?: Record<string, unknown>) => {
    try {
      console.log('[useTransactions] count start', { filters })
      const result = await transactionRepository.count(filters)
      console.log('[useTransactions] count success', { count: result })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useTransactions] count error', msg)
      throw err
    }
  }

  return { create, update, delete: deleteTransaction, findById, findAll, count }
}
