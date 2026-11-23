import { BaseRepository } from '../crud/base';

export class LabelRepository extends BaseRepository<any> {
  constructor() {
    super('label');
  }
}

export const labelRepository = new LabelRepository();
