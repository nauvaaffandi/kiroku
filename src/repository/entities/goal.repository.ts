import { BaseRepository } from '../crud/base';

export class GoalRepository extends BaseRepository<any> {
  constructor() {
    super('goal');
  }
}

export const goalRepository = new GoalRepository();
