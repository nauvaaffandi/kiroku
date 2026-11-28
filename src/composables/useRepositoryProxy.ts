import { useI18n } from '@/locales'
import { useUIStore } from '@/stores/ui'

export function useRepositoryProxy<T extends object>(repo: T, entity: string): T {
  const { t } = useI18n()
  const uiStore = useUIStore()

  return new Proxy<T>(repo, {
    get(target, prop, receiver) {
      const value = Reflect.get(target, prop, receiver)

      if (typeof value === 'function') {
        return async (...args: unknown[]) => {
          const method = String(prop)
          const targetId = args[0] ?? '*'

          console.log(`[${entity}Repository] ${method} start target=${targetId}`, args)

          try {
            const result = await value.apply(target, args)

            console.log(`[${entity}Repository] ${method} success target=${targetId}`, {
              result: result ?? null,
            })
            return result
          } catch (err: unknown) {
            const errorMsg = err instanceof Error ? err.message : String(err)
            uiStore.error(t('repository.error.fetch', { entity }))

            console.error(`[${entity}Repository] ${method} error target=${targetId}`, errorMsg)
            throw err
          }
        }
      }

      return value
    },
  })
}
