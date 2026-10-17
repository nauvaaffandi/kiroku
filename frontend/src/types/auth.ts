// Auth Types
export interface User {
  id: string
  email: string
  totpEnabled: boolean
}

export interface AuthState {
  user: User | null
  loading: boolean
  error: string | null
  isAuthenticated: boolean
}

export interface BootstrapRequest {
  email: string
  password: string
}

export interface LoginRequest {
  email: string
  password: string
  totp_code?: string
}

export interface TotpSetupResponse {
  secret: string
  qr_code: string
}

export interface TotpVerifyRequest {
  code: string
}

export interface AuthStatusResponse {
  has_users: boolean
}