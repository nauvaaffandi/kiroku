import { BaseRepository } from '../crud/base';
import type { EntityKind } from '../types';
import { filesystem } from '../core/filesystem';
import { cacheRegistry } from '../core/cache';
import type { Transaction } from '@/types/transaction';

function log(action: string, idOrPath: string, detail?: unknown) {
  console.log(`[TransactionRepository] ${action} target=${idOrPath}`, detail ?? '');
}

type TransactionEntity = Transaction;

export class TransactionRepository extends BaseRepository<Transaction> {
  private readonly index = new Map<string, string>();
  private indexBuilt = false;

  constructor() {
    super('transaction' as EntityKind);
    this.pathResolver = (id: string, entity?: TransactionEntity) => {
      if (entity?.date_time) {
        const d = new Date(entity.date_time);
        const year = d.getFullYear();
        const month = String(d.getMonth() + 1).padStart(2, '0');
        const day = String(d.getDate()).padStart(2, '0');
        return `${this.config.basePath}/${this.kind}/${year}/${month}/${day}.json`;
      }
      return `${this.config.basePath}/${this.kind}/${id}.json`;
    };
  }

  private async buildIndex(): Promise<void> {
    if (this.indexBuilt) return;
    const dir = `${this.config.basePath}/${this.kind}`;
    try {
      const years = await filesystem.readDir(dir);
      for (const y of years) {
        if (y.type !== 'directory') continue;
        const months = await filesystem.readDir(y.path);
        for (const m of months) {
          if (m.type !== 'directory') continue;
          const days = await filesystem.readDir(m.path);
          for (const dayFile of days) {
            if (dayFile.type !== 'file') continue;
            try {
              const raw = await filesystem.readFile(dayFile.path);
              const txns = JSON.parse(raw) as TransactionEntity[];
              for (const t of txns) {
                this.index.set(t.id, dayFile.path);
              }
            } catch {
              // skip corrupt files
            }
          }
        }
      }
    } catch {
      // directory may not exist yet
    }
    this.indexBuilt = true;
    log('buildIndex', '*', { indexedCount: this.index.size });
  }

  private invalidateIndex() {
    this.index.clear();
    this.indexBuilt = false;
  }

  async findById(id: string): Promise<Transaction | null> {
    const cacheKey = `transaction:${id}`;
    const cached = cacheRegistry.get<Transaction>('transaction', cacheKey);
    if (cached !== null) {
      log('findById', id, { source: 'cache' });
      return cached;
    }

    await this.buildIndex();
    const filePath = this.index.get(id);
    if (!filePath) {
      log('findById', id, { source: 'index', found: false });
      return null;
    }

    let raw: string;
    try {
      raw = await filesystem.readFile(filePath);
    } catch {
      return null;
    }

    let txns: TransactionEntity[];
    try {
      txns = JSON.parse(raw) as TransactionEntity[];
    } catch {
      return null;
    }

    const txn = (txns.find((t) => t.id === id && !(t as unknown as { deleted?: boolean }).deleted) ?? null) as Transaction | null;
    if (txn) {
      cacheRegistry.set('transaction', cacheKey, txn);
    }
    log('findById', id, { source: 'filesystem', found: !!txn });
    return txn;
  }

  async findAll(filters?: { [key: string]: unknown }): Promise<Transaction[]> {
    await this.buildIndex();

    const dir = `${this.config.basePath}/${this.kind}`;
    let years: { type: string; path: string; name: string }[];
    try {
      years = await filesystem.readDir(dir);
    } catch {
      return [];
    }

    const all: Transaction[] = [];

    for (const y of years) {
      if (y.type !== 'directory') continue;
      const months = await filesystem.readDir(y.path);
      for (const m of months) {
        if (m.type !== 'directory') continue;
        const days = await filesystem.readDir(m.path);
        for (const dayFile of days) {
          if (dayFile.type !== 'file') continue;
          try {
            const raw = await filesystem.readFile(dayFile.path);
            const txns = JSON.parse(raw) as TransactionEntity[];
            for (const t of txns) {
              if ((t as unknown as { deleted?: boolean }).deleted) continue;
              if (filters && !this.matchesFilters(t, filters)) continue;
              all.push(t as Transaction);
            }
          } catch {
            continue;
          }
        }
      }
    }

    log('findAll', dir, { count: all.length });
    return all;
  }

  override async create(data: Omit<Transaction, 'id'> & { id?: string }): Promise<string> {
    const id = data.id ?? crypto.randomUUID();
    const entity = { id, ...data } as Transaction;
    const path = this.pathResolver(id, entity);

    log('create', id, { path });

    await this.readModifyWrite(path, (current: unknown) => {
      const arr = Array.isArray(current) ? current : [current];
      arr.push(entity);
      return arr;
    });

    this.index.set(id, path);
    this.invalidateIndex();
    log('create', id, { path, status: 'committed' });
    return id;
  }

  override async update(id: string, data: Partial<Transaction>): Promise<void> {
    await this.buildIndex();
    const filePath = this.index.get(id);
    if (!filePath) {
      throw new Error(`Transaction not found: ${id}`);
    }

    log('update', id, { path: filePath });

    await this.readModifyWrite(filePath, (current: unknown) => {
      const arr = Array.isArray(current) ? current : [current];
      const idx = arr.findIndex((e: unknown) => (e as { id?: string })?.id === id);
      if (idx < 0) throw new Error(`Transaction not found in file: ${id}`);
      arr[idx] = { ...arr[idx], ...data, id } as Transaction;
      return arr;
    });

    if (data.date_time) {
      const newPath = this.pathResolver(id, { ...data, id, date_time: data.date_time } as Transaction);
      if (newPath !== filePath) {
        await this.readModifyWrite(filePath, (current: unknown) => {
          const arr = Array.isArray(current) ? current : [current];
          return arr.filter((e: unknown) => (e as { id?: string })?.id !== id);
        });
        await this.readModifyWrite(newPath, (current: unknown) => {
          const arr = Array.isArray(current) ? current : [current];
          const updated = { ...data, id, date_time: data.date_time } as Transaction;
          arr.push(updated);
          return arr;
        });
        this.index.set(id, newPath);
      }
    }

    this.invalidateIndex();
    log('update', id, { path: filePath, status: 'committed' });
  }

  override async delete(id: string): Promise<void> {
    await this.buildIndex();
    const filePath = this.index.get(id);
    if (!filePath) {
      throw new Error(`Transaction not found: ${id}`);
    }

    log('delete', id, { path: filePath });

    await this.readModifyWrite(filePath, (current: unknown) => {
      const arr = Array.isArray(current) ? current : [current];
      return arr.map((e: unknown) => {
        const item = e as { id?: string; deleted?: boolean };
        if (item.id === id) {
          return { ...item, deleted: true, deletedAt: Date.now() };
        }
        return e;
      });
    });

    this.index.delete(id);
    log('delete', id, { path: filePath, status: 'committed' });
  }

  private async readModifyWrite(path: string, mutate: (current: unknown) => unknown): Promise<void> {
    const dir = path.substring(0, path.lastIndexOf('/'));
    try {
      await filesystem.ensureDir(dir);
    } catch {
      // directory may already exist
    }

    let currentData: unknown = null;
    try {
      const exists = await filesystem.fileExists(path);
      if (exists) {
        const raw = await filesystem.readFile(path);
        currentData = JSON.parse(raw);
      }
    } catch {
      currentData = null;
    }

    const newData = mutate(currentData);
    await filesystem.writeFile(path, JSON.stringify(newData));
  }
}

export const transactionRepository = new TransactionRepository();
