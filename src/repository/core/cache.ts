import type { CacheEntry } from '../types';

function log(action: string, kind: string, key: string, detail?: unknown) {
  console.log(`[MemoryCache] ${action} kind=${kind} key=${key}`, detail ?? '');
}

export class MemoryCache<T> {
  private readonly map = new Map<string, CacheEntry<T>>();
  private readonly maxEntries: number;
  private readonly kind: string;

  constructor(kind: string, maxEntries = 500) {
    this.kind = kind;
    this.maxEntries = maxEntries;
  }

  get(key: string): T | null {
    const entry = this.map.get(key);
    if (!entry) {
      log('miss', this.kind, key);
      return null;
    }
    this.map.delete(key);
    this.map.set(key, { ...entry, accessCount: entry.accessCount + 1 });
    log('hit', this.kind, key, { accessCount: entry.accessCount + 1 });
    return entry.data;
  }

  set(key: string, data: T): void {
    if (this.map.size >= this.maxEntries && !this.map.has(key)) {
      const firstKey = this.map.keys().next().value;
      if (firstKey !== undefined) {
        log('evict', this.kind, firstKey);
        this.map.delete(firstKey);
      }
    }
    this.map.set(key, { data, timestamp: Date.now(), accessCount: 1 });
    log('set', this.kind, key);
  }

  has(key: string): boolean {
    return this.map.has(key);
  }

  delete(key: string): void {
    const existed = this.map.has(key);
    this.map.delete(key);
    log('delete', this.kind, key, { existed });
  }

  clear(): void {
    log('clear', this.kind, 'all', { count: this.map.size });
    this.map.clear();
  }

  size(): number {
    return this.map.size;
  }

  keys(): string[] {
    return Array.from(this.map.keys());
  }
}

export interface CacheRegistry {
  get<T>(kind: string, key: string): T | null;
  set<T>(kind: string, key: string, data: T): void;
  has(kind: string, key: string): boolean;
  delete(kind: string, key: string): void;
  clear(kind?: string): void;
  size(kind: string): number;
}

export class CacheRegistryImpl implements CacheRegistry {
  private readonly caches = new Map<string, MemoryCache<unknown>>();
  private readonly defaultMaxEntries: number;

  constructor(defaultMaxEntries = 500) {
    this.defaultMaxEntries = defaultMaxEntries;
  }

  private getCache(kind: string): MemoryCache<unknown> {
    let cache = this.caches.get(kind);
    if (!cache) {
      cache = new MemoryCache<unknown>(kind, this.defaultMaxEntries);
      this.caches.set(kind, cache);
    }
    return cache;
  }

  get<T>(kind: string, key: string): T | null {
    return this.getCache(kind).get(key) as T | null;
  }

  set<T>(kind: string, key: string, data: T): void {
    this.getCache(kind).set(key, data);
  }

  has(kind: string, key: string): boolean {
    return this.getCache(kind).has(key);
  }

  delete(kind: string, key: string): void {
    this.getCache(kind).delete(key);
  }

  clear(kind?: string): void {
    if (kind) {
      this.getCache(kind).clear();
    } else {
      for (const cache of this.caches.values()) {
        cache.clear();
      }
    }
  }

  size(kind: string): number {
    return this.getCache(kind).size();
  }
}

export const cacheRegistry = new CacheRegistryImpl();
