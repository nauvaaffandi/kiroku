export const SCHEMA_VERSION = 1;

export const MIGRATIONS: Record<number, string[]> = {
  1: [
    `CREATE TABLE IF NOT EXISTS meta (key TEXT PRIMARY KEY, value TEXT NOT NULL);`,
    `CREATE TABLE IF NOT EXISTS entities (
      id TEXT PRIMARY KEY,
      kind TEXT NOT NULL,
      payload TEXT NOT NULL,
      updated_at INTEGER NOT NULL,
      deleted INTEGER NOT NULL DEFAULT 0
    );`,
    `CREATE INDEX IF NOT EXISTS idx_entities_kind_deleted ON entities(kind, deleted);`,
    `INSERT INTO meta(key, value) VALUES ('schema_version', '1') ON CONFLICT(key) DO UPDATE SET value='1';`
  ]
};
