// Fallback shim for @capacitor-community/secure-storage-plugin.
// Used when the native plugin package is not installed (web / dev / CI).
// Keys are kept in an in-memory Map here — NOT at-rest encrypted.
//
// PRODUCTION NOTE: for real Capacitor Android builds, install the real
// `@capacitor-community/secure-storage-plugin` package, run `npx cap sync android`,
// then remove this shim, the vite alias, and the tsconfig path so the native
// EncryptedSharedPreferences / Keystore plugin provides at-rest encryption.
const storage = new Map<string, string>()

export const SecureStorage = {
  async set({ key, value }: { key: string; value: string }): Promise<void> {
    storage.set(key, value)
  },

  async get({ key }: { key: string }): Promise<{ value: string } | null> {
    const value = storage.get(key)
    return value === undefined ? null : { value }
  },

  async remove({ key }: { key: string }): Promise<void> {
    storage.delete(key)
  },

  async keys(): Promise<{ keys: string[] }> {
    return { keys: Array.from(storage.keys()) }
  }
}
