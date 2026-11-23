import { BaseRepository } from '../crud/base';

export class CategoryRepository extends BaseRepository<any> {
  constructor() {
    super('category');
  }
}

export const categoryRepository = new CategoryRepository();
