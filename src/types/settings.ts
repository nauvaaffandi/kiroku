// Settings Types
export type Theme = 'light' | 'dark' | 'system'
export type Language = 'en' | 'id' | 'ja'
export type DateFormat = 'YYYY-MM-DD' | 'DD/MM/YYYY' | 'MM/DD/YYYY'
export type NumberFormat = 'ja-JP' | 'en-US' | 'id-ID'
export type WeekStart = 0 | 1 | 2 | 3 | 4 | 5 | 6 // 0 = Sunday, 1 = Monday

export interface Settings {
  theme: Theme
  language: Language
  base_currency: string
  date_format: DateFormat
  number_format: NumberFormat
  week_start: WeekStart
  totp_enabled: boolean
  hide_amounts: boolean
  [key: string]: string | boolean | number | undefined
}

export interface SettingsState {
  settings: Settings
  loading: boolean
  theme: Theme
  language: Language
  baseCurrency: string
  dateFormat: DateFormat
  numberFormat: NumberFormat
  weekStart: WeekStart
  isTotpEnabled: boolean
}

export interface UpdateSettingsRequest {
  theme?: Theme
  language?: Language
  base_currency?: string
  date_format?: DateFormat
  number_format?: NumberFormat
  week_start?: WeekStart
  totp_enabled?: boolean
  [key: string]: string | boolean | number | undefined
}

export interface BackupCreateResponse {
  filename: string
  size: number
  created_at: string
}

export interface BackupRestoreRequest {
  backup: File
  password: string
}