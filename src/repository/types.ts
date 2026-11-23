// Storage-specific types (kept separate from domain types in src/types/**)

export type StorageDirectory = 'DATA' | 'DOCUMENTS' | 'CACHE' | 'LIBRARY' | 'EXTERNAL';

export interface FileEntry {
  path: string;
  uri: string;
  name: string;
  type: 'file' | 'directory';
  size: number;
  lastModified: number;
}

export interface WriteOperation {
  id: string;
  type: 'create' | 'update' | 'delete' | 'replace';
  target: string;
  entityId: string;
  payload?: unknown;
  previousPayload?: unknown;
  timestamp: number;
}

export interface JournalEntry {
  id: string;
  target: string;
  state: 'prepared' | 'committing' | 'committed' | 'failed';
  createdAt: string;
  completedAt?: string;
  operations?: WriteOperation[];
  checksum?: string;
}

export interface CacheEntry<T> {
  data: T;
  timestamp: number;
  accessCount: number;
}

export interface RepositoryConfig {
  basePath: string;
  maxCacheEntries: number;
  writeDebounceMs: number;
  maxWriteRetries: number;
}

export type EntityKind =
  | 'account'
  | 'category'
  | 'label'
  | 'notification'
  | 'transaction'
  | 'template'
  | 'scheduled'
  | 'budget'
  | 'goal'
  | 'setting';

export interface Filter {
  [key: string]: unknown;
}

export interface Aggregation {
  [key: string]: unknown;
}

export interface CreateInput {
  [key: string]: unknown;
}

export interface UpdateInput {
  [key: string]: unknown;
}

export interface RepositoryStats {
  totalReads: number;
  totalWrites: number;
  cacheHits: number;
  cacheMisses: number;
  pendingWrites: number;
}
