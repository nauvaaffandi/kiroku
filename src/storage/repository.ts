import { CapacitorSQLite, SQLiteConnection, type SQLiteDBConnection } from '@capacitor-community/sqlite';
import { MIGRATIONS, SCHEMA_VERSION } from './schema';

export type EntityKind = 'account' | 'category' | 'label' | 'notification' | 'transaction' | 'template' | 'scheduled' | 'budget' | 'goal' | 'setting';

export interface StoredEntity<T = Record<string, unknown>> {
  id: string;
  kind: EntityKind;
  payload: T;
  updated_at: number;
  deleted: boolean;
}

export class LocalDatabase {
  private db: SQLiteDBConnection | null = null;
  constructor(private readonly name = 'kiroku') {}

  async open(): Promise<void> {
    if (this.db) return;
    const connection = new SQLiteConnection(CapacitorSQLite);
    this.db = await connection.createConnection(this.name, false, 'no-encryption', SCHEMA_VERSION, false);
    await this.db.open();
    await this.migrate();
    // Always ensure the entities table exists, even if a previous migration left it missing.
    await this.db!.execute(`CREATE TABLE IF NOT EXISTS entities (
      id TEXT PRIMARY KEY,
      kind TEXT NOT NULL,
      payload TEXT NOT NULL,
      updated_at INTEGER NOT NULL,
      deleted INTEGER NOT NULL DEFAULT 0
    );`);
    await this.db!.execute(`CREATE INDEX IF NOT EXISTS idx_entities_kind_deleted ON entities(kind, deleted);`);
  }

  private async migrate(): Promise<void> {
    await this.db!.execute(`CREATE TABLE IF NOT EXISTS meta (key TEXT PRIMARY KEY, value TEXT NOT NULL);`);
    const result = await this.db!.query(`SELECT value FROM meta WHERE key = 'schema_version';`);
    const current = Number(result.values?.[0]?.value ?? 0);
    for (let version = current + 1; version <= SCHEMA_VERSION; version += 1) {
      await this.db!.execute('BEGIN IMMEDIATE;');
      try {
        for (const statement of MIGRATIONS[version]) await this.db!.execute(statement);
        await this.db!.execute('COMMIT;');
      } catch (error) {
        await this.db!.execute('ROLLBACK;');
        throw error;
      }
    }
  }

  async put<T>(kind: EntityKind, id: string, payload: T): Promise<void> {
    await this.open();
    await this.db!.run(
      `INSERT INTO entities(id, kind, payload, updated_at, deleted) VALUES (?, ?, ?, ?, 0)
       ON CONFLICT(id) DO UPDATE SET kind=excluded.kind, payload=excluded.payload, updated_at=excluded.updated_at, deleted=0;`,
      [id, kind, JSON.stringify(payload), Date.now()]
    );
  }

  async get<T>(kind: EntityKind, id: string): Promise<T | null> {
    await this.open();
    const result = await this.db!.query(
      `SELECT payload FROM entities WHERE kind = ? AND id = ? AND deleted = 0;`,
      [kind, id]
    );
    const payload = result.values?.[0]?.payload;
    if (!payload) return null;
    try {
      return JSON.parse(payload) as T;
    } catch {
      return null;
    }
  }

  async list<T>(kind: EntityKind): Promise<T[]> {
    await this.open();
    const result = await this.db!.query(
      `SELECT payload FROM entities WHERE kind = ? AND deleted = 0 ORDER BY updated_at DESC;`,
      [kind]
    );
    const rows = result.values ?? [];
    const items: T[] = [];
    for (const row of rows) {
      try {
        items.push(JSON.parse(row.payload) as T);
      } catch {
        // Skip baris data yang rusak/invalid JSON agar tidak gagal seluruhnya
        continue;
      }
    }
    return items;
  }

  async remove(kind: EntityKind, id: string): Promise<void> {
    await this.open();
    await this.db!.run(
      `UPDATE entities SET deleted = 1, updated_at = ? WHERE kind = ? AND id = ?;`,
      [Date.now(), kind, id]
    );
  }

  async close(): Promise<void> {
    if (!this.db) return;
    await this.db.close();
    await new SQLiteConnection(CapacitorSQLite).closeConnection(this.name, false);
    this.db = null;
  }
}

export const localDb = new LocalDatabase();
