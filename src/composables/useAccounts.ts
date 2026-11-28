import { useI18n } from '@/locales'
import { useUIStore } from '@/stores/ui'
import { accountRepository } from '@/repository'

export function useAccounts() {
  const { t } = useI18n()
  const uiStore = useUIStore()
  const entity = 'account'

  const create = async (data: unknown) => {
    try {
      console.log('[useAccounts] create start', data)
      const id = await accountRepository.create(data as any)
      console.log('[useAccounts] create success', { id })
      return id
    } catch (err: any) {
      const msg = err.message || t('repository.error.save', { entity })
      uiStore.error(msg)
      console.error('[useAccounts] create error', msg)
      throw err
    }
  }

  const update = async (id: string, data: unknown) => {
    try {
      console.log('[useAccounts] update start', { id, data })
      await accountRepository.update(id, data as any)
      console.log('[useAccounts] update success', { id })
    } catch (err: any) {
      const msg = err.message || t('repository.error.update', { entity })
      uiStore.error(msg)
      console.error('[useAccounts] update error', msg)
      throw err
    }
  }

  const deleteAccount = async (id: string) => {
    try {
      console.log('[useAccounts] delete start', { id })
      await accountRepository.delete(id)
      console.log('[useAccounts] delete success', { id })
    } catch (err: any) {
      const msg = err.message || t('repository.error.delete', { entity })
      uiStore.error(msg)
      console.error('[useAccounts] delete error', msg)
      throw err
    }
  }

  const findById = async (id: string) => {
    try {
      console.log('[useAccounts] findById start', { id })
      const result = await accountRepository.findById(id)
      console.log('[useAccounts] findById success', { id, found: !!result })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useAccounts] findById error', msg)
      throw err
    }
  }

  const findAll = async (filters?: Record<string, unknown>) => {
    try {
      console.log('[useAccounts] findAll start', { filters })
      const result = await accountRepository.findAll(filters)
      console.log('[useAccounts] findAll success', { count: result.length })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useAccounts] findAll error', msg)
      throw err
    }
  }

  const count = async (filters?: Record<string, unknown>) => {
    try {
      console.log('[useAccounts] count start', { filters })
      const result = await accountRepository.count(filters)
      console.log('[useAccounts] count success', { count: result })
      return result
    } catch (err: any) {
      const msg = err.message || t('repository.error.fetch', { entity })
      uiStore.error(msg)
      console.error('[useAccounts] count error', msg)
      throw err
    }
  }

  return { create, update, delete: deleteAccount, findById, findAll, count }
}
