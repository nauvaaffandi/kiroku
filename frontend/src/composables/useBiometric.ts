import { ref, computed } from 'vue'

// Storage keys
const BIOMETRIC_CREDENTIAL_ID_KEY = 'kiroku_biometric_credential_id'
const BIOMETRIC_ENABLED_KEY = 'kiroku_biometric_enabled'

// Types
export interface BiometricCredential {
  id: string
  rawId: ArrayBuffer
  type: 'public-key'
}

export interface BiometricSupportResult {
  supported: boolean
  reason?: string
}

export interface BiometricStatus {
  supported: boolean
  enabled: boolean
  credentialRegistered: boolean
  reason?: string
}

/**
 * Check if the device/browser supports WebAuthn (biometric authentication)
 */
export function checkBiometricSupport(): BiometricSupportResult {
  if (typeof window === 'undefined') {
    return { supported: false, reason: 'Not in browser environment' }
  }

  if (!window.PublicKeyCredential) {
    return { supported: false, reason: 'WebAuthn not supported in this browser' }
  }

  if (!navigator.credentials) {
    return { supported: false, reason: 'Credential Management API not available' }
  }

  // Check if conditional UI is available (for auto-fill style prompts)
  // This is optional - basic WebAuthn works without it
  return { supported: true }
}

/**
 * Check if a biometric credential is already registered for this user
 */
export async function checkCredentialRegistered(): Promise<boolean> {
  const credentialId = localStorage.getItem(BIOMETRIC_CREDENTIAL_ID_KEY)
  if (!credentialId) {
    return false
  }

  try {
    // We don't actually call navigator.credentials.get() here because it would
    // prompt the user. Instead, we just check if the credential ID exists in storage.
    // The actual verification happens during authentication.
    return true
  } catch {
    return false
  }
}

/**
 * Get the full biometric status (support + enabled + credential registered)
 */
export async function getBiometricStatus(): Promise<BiometricStatus> {
  const support = checkBiometricSupport()
  if (!support.supported) {
    return {
      supported: false,
      enabled: false,
      credentialRegistered: false,
      reason: support.reason
    }
  }

  const enabled = localStorage.getItem(BIOMETRIC_ENABLED_KEY) === 'true'
  const credentialRegistered = await checkCredentialRegistered()

  return {
    supported: true,
    enabled,
    credentialRegistered,
    reason: credentialRegistered ? undefined : 'No biometric credential registered'
  }
}

/**
 * Register a new biometric credential (WebAuthn registration)
 * Returns the credential ID if successful
 */
export async function registerBiometricCredential(): Promise<string | null> {
  const support = checkBiometricSupport()
  if (!support.supported) {
    throw new Error(support.reason || 'Biometric not supported')
  }

  try {
    const challenge = new Uint8Array(32)
    crypto.getRandomValues(challenge)

    const credential = await navigator.credentials.create({
      publicKey: {
        challenge,
        rp: { name: 'Kiroku Finance', id: window.location.hostname },
        user: {
          id: new TextEncoder().encode('kiroku_user_1'),
          name: 'user@kiroku.local',
          displayName: 'Kiroku User'
        },
        pubKeyCredParams: [
          { alg: -7, type: 'public-key' },   // ES256
          { alg: -257, type: 'public-key' }  // RS256
        ],
        authenticatorSelection: {
          authenticatorAttachment: 'platform', // Prefer platform authenticator (biometric)
          userVerification: 'required',
          residentKey: 'preferred'
        },
        timeout: 60000,
        attestation: 'none'
      }
    }) as PublicKeyCredential

    // Store the credential ID for later verification
    const credentialId = bufferToBase64Url(credential.rawId)
    localStorage.setItem(BIOMETRIC_CREDENTIAL_ID_KEY, credentialId)
    localStorage.setItem(BIOMETRIC_ENABLED_KEY, 'true')

    return credentialId
  } catch (error) {
    // User may have cancelled the prompt
    console.warn('Biometric registration cancelled or failed:', error)
    throw error
  }
}

/**
 * Authenticate using an existing biometric credential (WebAuthn authentication)
 * Returns true if authentication succeeds
 */
export async function authenticateBiometric(): Promise<boolean> {
  const support = checkBiometricSupport()
  if (!support.supported) {
    throw new Error(support.reason || 'Biometric not supported')
  }

  const credentialId = localStorage.getItem(BIOMETRIC_CREDENTIAL_ID_KEY)
  if (!credentialId) {
    throw new Error('No biometric credential registered')
  }

  try {
    const challenge = new Uint8Array(32)
    crypto.getRandomValues(challenge)

    const assertion = await navigator.credentials.get({
      publicKey: {
        challenge,
        allowCredentials: [{
          id: base64UrlToBuffer(credentialId),
          type: 'public-key',
          transports: ['internal', 'hybrid', 'usb', 'ble', 'nfc']
        }],
        userVerification: 'required',
        timeout: 30000
      }
    })

    return !!assertion
  } catch (error) {
    console.warn('Biometric authentication failed:', error)
    throw error
  }
}

/**
 * Disable biometric authentication (remove stored credential)
 */
export function disableBiometric(): void {
  localStorage.removeItem(BIOMETRIC_CREDENTIAL_ID_KEY)
  localStorage.removeItem(BIOMETRIC_ENABLED_KEY)
}

/**
 * Check if biometric is enabled in settings (regardless of credential)
 */
export function isBiometricEnabled(): boolean {
  return localStorage.getItem(BIOMETRIC_ENABLED_KEY) === 'true'
}

/**
 * Set biometric enabled state
 */
export function setBiometricEnabled(enabled: boolean): void {
  localStorage.setItem(BIOMETRIC_ENABLED_KEY, String(enabled))
}

/**
 * Convert ArrayBuffer to base64url string (for storage)
 */
function bufferToBase64Url(buffer: ArrayBuffer): string {
  const bytes = new Uint8Array(buffer)
  let binary = ''
  for (let i = 0; i < bytes.byteLength; i++) {
    binary += String.fromCharCode(bytes[i])
  }
  return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '')
}

/**
 * Convert base64url string to ArrayBuffer (for WebAuthn)
 */
function base64UrlToBuffer(base64url: string): ArrayBuffer {
  const base64 = base64url.replace(/-/g, '+').replace(/_/g, '/')
  const padded = base64.padEnd(base64.length + (4 - (base64.length % 4)) % 4, '=')
  const binary = atob(padded)
  const bytes = new Uint8Array(binary.length)
  for (let i = 0; i < binary.length; i++) {
    bytes[i] = binary.charCodeAt(i)
  }
  return bytes.buffer
}

/**
 * Vue composable for biometric authentication
 */
export function useBiometric() {
  const status = ref<BiometricStatus>({
    supported: false,
    enabled: false,
    credentialRegistered: false
  })
  const loading = ref(false)
  const error = ref<string | null>(null)

  const isAvailable = computed(() => status.value.supported && status.value.enabled && status.value.credentialRegistered)

  async function refreshStatus() {
    status.value = await getBiometricStatus()
  }

  async function enable() {
    loading.value = true
    error.value = null
    try {
      await registerBiometricCredential()
      await refreshStatus()
      return true
    } catch (e: any) {
      error.value = e.message || 'Failed to enable biometric'
      return false
    } finally {
      loading.value = false
    }
  }

  async function disable() {
    loading.value = true
    error.value = null
    try {
      disableBiometric()
      await refreshStatus()
      return true
    } catch (e: any) {
      error.value = e.message || 'Failed to disable biometric'
      return false
    } finally {
      loading.value = false
    }
  }

  async function authenticate() {
    loading.value = true
    error.value = null
    try {
      const result = await authenticateBiometric()
      return result
    } catch (e: any) {
      error.value = e.message || 'Biometric authentication failed'
      return false
    } finally {
      loading.value = false
    }
  }

  return {
    status,
    isAvailable,
    loading,
    error,
    refreshStatus,
    enable,
    disable,
    authenticate
  }
}