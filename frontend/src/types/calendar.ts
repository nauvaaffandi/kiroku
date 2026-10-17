export interface CalendarTxn {
  id: string
  name: string
  amount: number
  type: string
  category: string
  currency: string
}

export interface CalendarDayBackend {
  date: string
  transactions: CalendarTxn[]
  scheduled: CalendarTxn[]
  balance_delta: number
  projected_balance: number
}

export interface DisplayDay {
  date: string
  dayNumber: number
  isCurrentMonth: boolean
  isToday: boolean
  transactions: CalendarTxn[]
  scheduled: CalendarTxn[]
  projectedBalance: number
}