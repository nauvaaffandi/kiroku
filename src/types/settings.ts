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
  ocr_enabled: boolean
  // Populated from secure storage (useSecureStorage), not the SQLite payload.
  // The SQLite payload stores '' for this field; the live value lives in encrypted at-rest storage.
  ocr_api_key: string
  ai_provider_enabled: boolean
  ai_provider_url: string
  ai_model: string
  // Populated from secure storage (useSecureStorage), not the SQLite payload.
  // The SQLite payload stores '' for this field; the live value lives in encrypted at-rest storage.
  ai_api_key: string
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
  isHideAmounts: boolean
  isOcrEnabled: boolean
  isAiProviderEnabled: boolean
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