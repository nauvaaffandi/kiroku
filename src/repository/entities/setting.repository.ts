import { BaseRepository } from '../crud/base';

export class SettingRepository extends BaseRepository<any> {
  constructor() {
    super('setting');
  }
}

export const settingRepository = new SettingRepository();
