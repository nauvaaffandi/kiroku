// Biometric Types
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
