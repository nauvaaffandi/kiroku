// Exchange Types
export interface RateState {
	rate: number | null
	loading: boolean
	error: string | null
	lastFetched: number | null
}
