// Currency Types
export type CurrencyCode = string

export interface Currency {
	id: string
	code: string
	name: string
	symbol: string
	enabled?: boolean
}

export interface CurrencyPick {
	code: string
}

export interface MasterCurrency {
	code: CurrencyCode
	name: string
	symbol: string
}
