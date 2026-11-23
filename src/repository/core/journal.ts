import type { JournalEntry } from '../types';
import { filesystem } from './filesystem';

const RECOVERY_DIR = 'tmp/recovery';
const WRITES_DIR = 'tmp/writes';

function log(action: string, opId: string, detail?: unknown) {
  console.log(`[Journal] ${action} opId=${opId}`, detail ?? '');
}

export class Journal {
  async ensureDirs(): Promise<void> {
    await filesystem.ensureDir(RECOVERY_DIR);
    await filesystem.ensureDir(WRITES_DIR);
  }

  private journalPath(opId: string): string {
    return `${RECOVERY_DIR}/${opId}.json`;
  }

  async create(entry: Omit<JournalEntry, 'completedAt'>): Promise<void> {
    await this.ensureDirs();
    const journal: JournalEntry = { ...entry, completedAt: undefined };
    await filesystem.writeFile(this.journalPath(entry.id), JSON.stringify(journal));
    log('create', entry.id, { target: entry.target, state: entry.state });
  }

  async get(opId: string): Promise<JournalEntry | null> {
    try {
      const raw = await filesystem.readFile(this.journalPath(opId));
      return JSON.parse(raw) as JournalEntry;
    } catch {
      return null;
    }
  }

  async updateState(opId: string, state: JournalEntry['state']): Promise<void> {
    const entry = await this.get(opId);
    if (!entry) {
      log('updateState', opId, { state, error: 'not-found' });
      return;
    }
    const updated: JournalEntry = {
      ...entry,
      state,
      ...(state === 'committed' || state === 'failed' ? { completedAt: new Date().toISOString() } : {}),
    };
    await filesystem.writeFile(this.journalPath(opId), JSON.stringify(updated));
    log('updateState', opId, { state, completedAt: updated.completedAt });
  }

  async listPending(): Promise<JournalEntry[]> {
    try {
      await this.ensureDirs();
      const entries = await filesystem.readDir(RECOVERY_DIR);
      const journals: JournalEntry[] = [];
      for (const e of entries) {
        if (e.type !== 'file' || !e.name.endsWith('.json')) continue;
        try {
          const raw = await filesystem.readFile(e.path);
          const journal = JSON.parse(raw) as JournalEntry;
          if (journal.state !== 'committed' && journal.state !== 'failed') {
            journals.push(journal);
          }
        } catch {
          // skip corrupt journal files
        }
      }
      return journals;
    } catch {
      return [];
    }
  }

  async remove(opId: string): Promise<void> {
    try {
      await filesystem.deleteFile(this.journalPath(opId));
      log('remove', opId);
    } catch {
      // ignore — may already be deleted
    }
  }

  async listAll(): Promise<JournalEntry[]> {
    try {
      await this.ensureDirs();
      const entries = await filesystem.readDir(RECOVERY_DIR);
      const journals: JournalEntry[] = [];
      for (const e of entries) {
        if (e.type !== 'file' || !e.name.endsWith('.json')) continue;
        try {
          const raw = await filesystem.readFile(e.path);
          journals.push(JSON.parse(raw) as JournalEntry);
        } catch {
          // skip corrupt
        }
      }
      return journals;
    } catch {
      return [];
    }
  }

  tempPath(opId: string): string {
    return `${WRITES_DIR}/${opId}.json`;
  }

  async writeTemp(opId: string, data: string): Promise<void> {
    await this.ensureDirs();
    await filesystem.writeFile(this.tempPath(opId), data);
    log('writeTemp', opId);
  }

  async readTemp(opId: string): Promise<string | null> {
    try {
      return await filesystem.readFile(this.tempPath(opId));
    } catch {
      return null;
    }
  }

  async removeTemp(opId: string): Promise<void> {
    try {
      await filesystem.deleteFile(this.tempPath(opId));
      log('removeTemp', opId);
    } catch {
      // ignore
    }
  }
}

export const journal = new Journal();
