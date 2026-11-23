import { BaseRepository } from '../crud/base';

export class BudgetRepository extends BaseRepository<any> {
  constructor() {
    super('budget');
  }
}

export const budgetRepository = new BudgetRepository();
