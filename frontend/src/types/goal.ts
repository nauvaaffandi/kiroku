export interface Goal {
  id: string
  name: string
  target_amount: number
  icon?: string
  desired_date: string
  color_theme?: string
  note?: string
  saved_already: number
  currency: string
  created_at: string
  updated_at: string
}

export interface GoalProjection {
  current_saved: number
  target_amount: number
  remaining: number
  weekly_rate: number
  estimated_weeks: number
  estimated_reach_date: string
  currency: string
}

export interface CreateGoalInput {
  name: string
  target_amount: number
  icon?: string
  desired_date?: string
  color_theme?: string
  note?: string
  currency: string
}

export interface UpdateGoalInput {
  name?: string
  target_amount?: number
  icon?: string
  desired_date?: string
  color_theme?: string
  note?: string
  currency?: string
}

export interface ContributeGoalInput {
  amount: number
  note?: string
}