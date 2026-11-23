import { BaseRepository } from '../crud/base';

export class ScheduledRepository extends BaseRepository<any> {
  constructor() {
    super('scheduled');
  }
}

export const scheduledRepository = new ScheduledRepository();
