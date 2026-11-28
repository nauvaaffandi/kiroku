// Currency Types
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
