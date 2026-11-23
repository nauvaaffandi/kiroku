import type {
  EntityKind,
  FileEntry,
  Filter,
  RepositoryConfig,
} from '../types';
import { filesystem } from '../core/filesystem';
import { cacheRegistry } from '../core/cache';

function log(action: string, kind: string, idOrPath: string, detail?: unknown) {
  console.log(`[BaseRepository] ${action} kind=${kind} target=${idOrPath}`, detail ?? '');
}

const DEFAULT_CONFIG: RepositoryConfig = {
  basePath: 'data',
  maxCacheEntries: 500,
  writeDebounceMs: 500,
  maxWriteRetries: 3,
};

const fileLocks = new Map<string, Promise<unknown>>();

async function withFileLock<T>(path: string, fn: () => Promise<T>): Promise<T> {
  const existing = fileLocks.get(path);
  const promise = (async () => {
    if (existing) await existing;
    return fn();
  })();
  fileLocks.set(path, promise as unknown as Promise<unknown>);
  try {
    return await promise;
  } finally {
    if (fileLocks.get(path) === promise) {
      fileLocks.delete(path);
    }
  }
}

export abstract class BaseRepository<T extends { id: string }> {
  protected readonly kind: EntityKind;
  protected readonly config: RepositoryConfig;
  protected pathResolver: (id: string, entity?: T) => string;

  constructor(kind: EntityKind, config?: Partial<RepositoryConfig>) {
    this.kind = kind;
    this.config = { ...DEFAULT_CONFIG, ...config };
    this.pathResolver = (id: string) => `${this.config.basePath}/${this.kind}/${id}.json`;
  }

  private cacheKey(id: string): string {
    return `${this.kind}:${id}`;
  }

  private dirPath(): string {
    return `${this.config.basePath}/${this.kind}`;
  }

  async findById(id: string): Promise<T | null> {
    const key = this.cacheKey(id);
    const cached = cacheRegistry.get<T>(this.kind, key);
    if (cached !== null) {
      log('findById', this.kind, id, { source: 'cache' });
      return cached;
    }

    const path = this.pathResolver(id);
    let exists = false;
    try {
      exists = await filesystem.fileExists(path);
    } catch {
      exists = false;
    }

    if (!exists) {
      log('findById', this.kind, id, { source: 'filesystem', found: false });
      return null;
    }

    let raw: string;
    try {
      raw = await filesystem.readFile(path);
    } catch (error) {
      console.error(`[BaseRepository] findById readFile error kind=${this.kind} id=${id}`, error);
      return null;
    }

    let entity: T;
    try {
      entity = JSON.parse(raw) as T;
    } catch {
      console.error(`[BaseRepository] findById invalid JSON kind=${this.kind} id=${id} path=${path}`);
      return null;
    }

    if ((entity as unknown as { deleted?: boolean }).deleted) {
      log('findById', this.kind, id, { source: 'filesystem', found: false, reason: 'soft-deleted' });
      return null;
    }

    cacheRegistry.set(this.kind, key, entity);
    log('findById', this.kind, id, { source: 'filesystem', found: true });
    return entity;
  }

  async findAll(filters?: Filter): Promise<T[]> {
    const dir = this.dirPath();
    log('findAll', this.kind, dir, { hasFilters: !!filters });

    let entries: FileEntry[];
    try {
      entries = await filesystem.readDir(dir);
    } catch {
      return [];
    }

    const fileEntries = entries.filter((e) => e.type === 'file' && e.name.endsWith('.json'));
    const results: T[] = [];

    for (const entry of fileEntries) {
      let raw: string;
      try {
        raw = await filesystem.readFile(entry.path);
      } catch {
        continue;
      }

      let entity: T;
      try {
        entity = JSON.parse(raw) as T;
      } catch {
        continue;
      }

      if ((entity as unknown as { deleted?: boolean }).deleted) {
        continue;
      }

      if (filters && !this.matchesFilters(entity, filters)) {
        continue;
      }

      const key = this.cacheKey(entity.id);
      cacheRegistry.set(this.kind, key, entity);
      results.push(entity);
    }

    log('findAll', this.kind, dir, { count: results.length });
    return results;
  }

  protected matchesFilters(entity: T, filters: Filter): boolean {
    for (const [key, value] of Object.entries(filters)) {
      const entityValue = (entity as unknown as Record<string, unknown>)[key];
      if (entityValue === undefined || entityValue !== value) {
        return false;
      }
    }
    return true;
  }

  async create(data: Omit<T, 'id'> & { id?: string }): Promise<string> {
    const id = data.id ?? crypto.randomUUID();
    const entity = { id, ...data } as T;
    const path = this.pathResolver(id);

    log('create', this.kind, id, { path });

    await withFileLock(path, async () => {
      const dir = this.dirPath();
      try {
        await filesystem.ensureDir(dir);
      } catch {
        // directory may already exist
      }

      await filesystem.writeFile(path, JSON.stringify(entity));
      cacheRegistry.set(this.kind, this.cacheKey(id), entity);
    });

    log('create', this.kind, id, { path, status: 'committed' });
    return id;
  }

  async update(id: string, data: Partial<T>): Promise<void> {
    const path = this.pathResolver(id);

    log('update', this.kind, id, { path });

    await withFileLock(path, async () => {
      const current = await this.findById(id);
      if (!current) {
        throw new Error(`${this.kind} not found: ${id}`);
      }

      const updated = { ...current, ...data, id } as T;
      await filesystem.writeFile(path, JSON.stringify(updated));
      cacheRegistry.set(this.kind, this.cacheKey(id), updated);
    });

    log('update', this.kind, id, { path, status: 'committed' });
  }

  async delete(id: string): Promise<void> {
    const path = this.pathResolver(id);

    log('delete', this.kind, id, { path });

    await withFileLock(path, async () => {
      const current = await this.findById(id);
      if (!current) {
        throw new Error(`${this.kind} not found: ${id}`);
      }

      const softDeleted = { ...current, deleted: true, deletedAt: Date.now() } as T & { deleted: boolean; deletedAt: number };
      await filesystem.writeFile(path, JSON.stringify(softDeleted));
      cacheRegistry.delete(this.kind, this.cacheKey(id));
    });

    log('delete', this.kind, id, { path, status: 'committed' });
  }

  async count(filters?: Filter): Promise<number> {
    const all = await this.findAll(filters);
    log('count', this.kind, this.dirPath(), { count: all.length });
    return all.length;
  }

  async invalidateCache(id?: string): Promise<void> {
    if (id) {
      cacheRegistry.delete(this.kind, this.cacheKey(id));
      log('invalidateCache', this.kind, id);
    } else {
      cacheRegistry.clear(this.kind);
      log('invalidateCache', this.kind, 'all');
    }
  }
}
