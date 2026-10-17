// Report Types
export interface SpendingGroup {
  category_id: string | null
  category_name: string | null
  category_icon: string | null
  category_color: string | null
  amount: number
  count: number
  percentage: number
  children: SpendingGroup[]
}

export interface SpendingReport {
  period: string
  groups: SpendingGroup[]
  total_income: number
  total_expense: number
  net: number
}

export interface TrendPoint {
  period: string
  income: number
  expense: number
  net: number
}

export interface TrendReport {
  series: TrendPoint[]
}

export interface BalancePoint {
  date: string
  balance: number
  account_id: string
  account_name: string
}

export interface BalanceReport {
  series: BalancePoint[]
}

export interface BudgetItem {
  category_id: string
  category_name: string
  category_icon: string
  category_color: string
  budget_amount: number
  spent_amount: number
  remaining: number
  percentage: number
}

export interface BudgetReport {
  items: BudgetItem[]
  total_budget: number
  total_spent: number
  total_remaining: number
}

export interface CashflowProjectionPoint {
  date: string
  projected_balance: number
  actual_balance: number | null
}

export interface CashflowProjectionReport {
  series: CashflowProjectionPoint[]
  accounts: {
    id: string
    name: string
    current_balance: number
  }[]
}

export interface ProjectionPoint {
  date: string
  income: number
  expense: number
  net: number
  running_balance: number
  currency: string
}

export interface CashFlowResponse {
  projected_balance: ProjectionPoint[]
  scheduled_transactions: any[]
  future_transactions: any[]
}

export interface BalanceTrendPoint {
  period: string
  opening_balance: number
  closing_balance: number
  net_change: number
  currency: string
}

export interface AccountBalanceSummary {
  id: number
  name: string
  currency: string
  opening_balance: number
  closing_balance: number
  net_change: number
  color: string
}

export interface CurrencyBalanceSummary {
  currency: string
  opening_balance: number
  closing_balance: number
  net_change: number
  fx_rate_to_base: number
  converted_total: number
}

export interface Goal {
  id: string
  name: string
  target_amount: number
  current_amount: number
  target_date: string
  category_id: string | null
  account_id: string | null
  created_at: string
  updated_at: string
}

export interface GoalItem {
  category_id: string
  category_name: string
  icon: string
  color: string
  target_amount: number
  target_period: string
  current_amount: number
  progress_pct: number
  remaining: number
  projected_completion: string
  status: 'on_track' | 'behind' | 'achieved'
}

export interface GoalSummary {
  total_targets: number
  total_current: number
  overall_progress: number
  goals_achieved: number
  goals_on_track: number
  goals_behind: number
}

export interface GoalsReport {
  goals: Goal[]
}

export interface OutlookItem {
  date: string
  type: 'income' | 'expense' | 'transfer'
  amount: number
  description: string
  category_id: string | null
  account_id: string | null
}

export interface OutlookReport {
  items: OutlookItem[]
  summary: {
    total_income: number
    total_expense: number
    net: number
  }
}