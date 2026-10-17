// Scheduled Types
export type ScheduledType = 'income' | 'expense' | 'transfer'
export type ScheduledFrequency = 'daily' | 'weekly' | 'monthly' | 'yearly' | 'custom'

export interface ScheduledTransaction {
  id: string
  name: string
  type: ScheduledType
  frequency: ScheduledFrequency
  start_date: string
  end_date: string | null
  cron_expression: string | null
  amount: number
  currency: string
  from_account_id: string | null
  to_account_id: string | null
  category_id: string | null
  description: string
  active: boolean
  next_run: string | null
  occurrences_done: number
  created_at: string
  updated_at: string
}

export interface ScheduledFilters {
  q: string
  type: string
  frequency: string
  active: string
  sort: string
  order: string
}

export interface CreateScheduledRequest {
  name: string
  type: ScheduledType
  frequency: ScheduledFrequency
  start_date: string
  end_date?: string | null
  cron_expression?: string | null
  amount: number
  currency: string
  from_account_id: string | null
  to_account_id?: string | null
  category_id?: string | null
  description?: string
  active?: boolean
}

export interface UpdateScheduledRequest {
  name?: string
  type?: ScheduledType
  frequency?: ScheduledFrequency
  start_date?: string
  end_date?: string | null
  cron_expression?: string | null
  amount?: number
  currency?: string
  from_account_id?: string | null
  to_account_id?: string | null
  category_id?: string | null
  description?: string
  active?: boolean
}