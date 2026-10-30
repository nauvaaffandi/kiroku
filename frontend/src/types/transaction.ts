// Transaction Types
import type { SortOrder, PaginationParams } from './common'

export type TransactionType = 'income' | 'expense' | 'transfer'
export type TransactionStatus = 'cleared' | 'pending'

export interface Transaction {
  id: string
  type: TransactionType
  date_time: string
  amount: number
  currency: string
  from_account_id: string | null
  to_account_id: string | null
  category_id: string | null
  labels: string | string[]
  label_ids?: string[]
  note: string
  payee_payer: string
  payment_type: string
  warranty_until: string | null
  status: TransactionStatus
  place_name: string
  place_lat: number | null
  place_lng: number | null
  attachments: string
  template_id: string | null
  created_at: string
  updated_at: string
}

export interface TransactionFilters {
  account_id: string
  category_id: string
  type: string
  status: string
  date_from: string
  date_to: string
  amount_min: string
  amount_max: string
  q: string
  label: string
  payee: string
  sort: string
  order: SortOrder
}

export interface TransactionSummary {
  [key: string]: {
    income: number
    expense: number
    net: number
  }
}

export interface CreateTransactionRequest {
  type: TransactionType
  date_time: string
  amount: number
  currency: string
  from_account_id?: string | null
  to_account_id?: string | null
  category_id?: string | null
  labels?: string | string[]
  label_ids?: string[]
  note?: string
  payee_payer?: string
  payment_type?: string
  warranty_until?: string | null
  status?: TransactionStatus
  place_name?: string
  place_lat?: number | null
  place_lng?: number | null
  attachments?: string
  template_id?: string | null
}

export interface UpdateTransactionRequest {
  type?: TransactionType
  date_time?: string
  amount?: number
  currency?: string
  from_account_id?: string | null
  to_account_id?: string | null
  category_id?: string | null
  labels?: string | string[]
  label_ids?: string[]
  note?: string
  payee_payer?: string
  payment_type?: string
  warranty_until?: string | null
  status?: TransactionStatus
  place_name?: string
  place_lat?: number | null
  place_lng?: number | null
  attachments?: string
  template_id?: string | null
}

export interface TransactionListParams extends PaginationParams {
  account_id?: string
  category_id?: string
  type?: string
  status?: string
  date_from?: string
  date_to?: string
  amount_min?: number
  amount_max?: number
  q?: string
  label?: string
  payee?: string
}