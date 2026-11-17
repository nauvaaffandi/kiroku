import type { Settings } from '@/types'

export const SECURE_KEYS: readonly (keyof Settings)[] = ['ocr_api_key', 'ai_api_key']

let cachedSecureStorage:
  | typeof import('capacitor-secure-storage-plugin')
  | null = null

// Lazy dynamic import is intentional: the Capacitor native plugin is unavailable
// in web/dev/test and the package is not resolvable outside a native Android
// build. import() defers resolution until first use, avoids module-load failure,
// and keeps the module interceptable by vitest `vi.mock`.
async function getSecureStorage() {
  if (!cachedSecureStorage) {
    cachedSecureStorage = await import('capacitor-secure-storage-plugin')
  }
  return cachedSecureStorage
}

export function useSecureStorage() {
  const setSecure = async (key: string, value: string): Promise<boolean> => {
    console.log('[SecureStorage] set')
    try {
      const { SecureStorage } = await getSecureStorage()
      await SecureStorage.set({ key, value })
      return true
    } catch {
      console.log('[SecureStorage] set failed')
      return false
    }
  }

  const getSecure = async (key: string): Promise<string | null> => {
    console.log('[SecureStorage] get')
    try {
      const { SecureStorage } = await getSecureStorage()
      const result = await SecureStorage.get({ key })
      return result?.value ?? null
    } catch {
      return null
    }
  }

  const removeSecure = async (key: string): Promise<boolean> => {
    console.log('[SecureStorage] remove')
    try {
      const { SecureStorage } = await getSecureStorage()
      await SecureStorage.remove({ key })
      return true
    } catch {
      console.log('[SecureStorage] remove failed')
      return false
    }
  }

  const hasSecure = async (key: string): Promise<boolean> => {
    console.log('[SecureStorage] hasSecure')
    try {
      const { SecureStorage } = await getSecureStorage()
      const { keys } = await SecureStorage.keys()
      return keys.includes(key)
    } catch {
      return false
    }
  }

  const migratePlaintextKeys = async (
    settingsRow: Settings
  ): Promise<{ migrated: boolean; redactedKeys: string[]; failedKeys: string[] }> => {
    const redactedKeys: string[] = []
    const failedKeys: string[] = []
    let migrated = false
    for (const key of SECURE_KEYS) {
      const value = settingsRow[key]
      if (value) {
        const ok = await setSecure(key as string, value as string)
        if (ok) {
          redactedKeys.push(key as string)
          migrated = true
        } else {
          failedKeys.push(key as string)
        }
      }
    }
    return { migrated, redactedKeys, failedKeys }
  }

  return { setSecure, getSecure, removeSecure, hasSecure, migratePlaintextKeys }
}
