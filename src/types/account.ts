// Account Types
export type AccountType = 'checking' | 'savings' | 'credit' | 'cash' | 'investment' | 'other'

export interface AccountFilters {
  q: string
  type: string
  currency: string
  archived: string
  sort: string
  order: string
}

export interface Account {
  id: string
  name: string
  type: AccountType
  currency: string
  balance: number
  icon: string
  color: string
  sort_order: number
  archived: boolean
}

export interface CreateAccountRequest {
  name: string
  type: AccountType
  currency: string
  balance?: number
  icon?: string
  color?: string
  sort_order?: number
}

export interface UpdateAccountRequest {
  name?: string
  type?: AccountType
  currency?: string
  balance?: number
  icon?: string
  color?: string
  sort_order?: number
  archived?: boolean
}