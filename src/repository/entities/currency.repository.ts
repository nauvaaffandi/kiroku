import { BaseRepository } from '../crud/base'
import type { Currency } from '@/types/currency'

export class CurrencyRepository extends BaseRepository<Currency> {
	constructor() {
		super('currency')
	}

	async pick(code: string): Promise<void> {
		const normalized = code.toUpperCase()
		const exists = await this.findById(normalized)
		if (exists) {
			await this.update(normalized, { enabled: true })
			console.log('[CurrencyRepository] pick', normalized, {
				status: 'already-existed',
				enabled: true,
			})
			return
		}
		await this.create({
			id: normalized,
			code: normalized,
			name: normalized,
			symbol: normalized,
			enabled: true,
		})
		console.log('[CurrencyRepository] pick', normalized, {
			status: 'created',
			enabled: true,
		})
	}

	async unpick(code: string): Promise<void> {
		const normalized = code.toUpperCase()
		const exists = await this.findById(normalized)
		if (!exists) {
			console.log('[CurrencyRepository] unpick', normalized, {
				status: 'not-found',
			})
			return
		}
		await this.update(normalized, { enabled: false })
		console.log('[CurrencyRepository] unpick', normalized, {
			status: 'disabled',
		})
	}

	async findPicked(): Promise<Currency[]> {
		const all = await this.findAll()
		const picked = all.filter((c) => c.enabled !== false)
		console.log('[CurrencyRepository] findPicked', {
			total: all.length,
			picked: picked.length,
		})
		return picked
	}
}

export const currencyRepository = new CurrencyRepository()
