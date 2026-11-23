import { BaseRepository } from '../crud/base';

export class TemplateRepository extends BaseRepository<any> {
  constructor() {
    super('template');
  }
}

export const templateRepository = new TemplateRepository();
