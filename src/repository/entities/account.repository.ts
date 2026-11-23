import { BaseRepository } from '../crud/base';

export class AccountRepository extends BaseRepository<any> {
  constructor() {
    super('account');
  }
}

export const accountRepository = new AccountRepository();
