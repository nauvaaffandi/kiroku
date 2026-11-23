import { BaseRepository } from '../crud/base';

export class NotificationRepository extends BaseRepository<any> {
  constructor() {
    super('notification');
  }
}

export const notificationRepository = new NotificationRepository();
