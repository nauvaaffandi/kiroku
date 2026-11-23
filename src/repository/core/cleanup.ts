import { filesystem } from './filesystem';
import { journal } from './journal';
import { atomicWrite } from './atomic-write';

const RETENTION_MS_FAILED = 24 * 60 * 60 * 1000;
const TMP_DIRS = ['tmp/writes', 'tmp/recovery', 'tmp/orphan'];

function log(action: string, detail?: unknown) {
  console.log(`[Cleanup] ${action}`, detail ?? '');
}

export class Cleanup {
  async runStartup(): Promise<void> {
    log('startup', { dirs: TMP_DIRS });

    for (const dir of TMP_DIRS) {
      try {
        await filesystem.ensureDir(dir);
      } catch {
        // ignore
      }
    }

    await atomicWrite.recover();
    await this.cleanOldJournals();
    await this.cleanOrphanTemps();

    log('startup', { status: 'done' });
  }

  async runPostCommit(opId: string): Promise<void> {
    log('postCommit', opId);
    await journal.removeTemp(opId);
    await journal.remove(opId);
  }

  private async cleanOldJournals(): Promise<void> {
    const all = await journal.listAll();
    const now = Date.now();
    let cleaned = 0;

    for (const entry of all) {
      if (entry.state === 'prepared' || entry.state === 'committing') {
        continue;
      }

      const completedAt = entry.completedAt ? new Date(entry.completedAt).getTime() : 0;
      const age = completedAt > 0 ? now - completedAt : 0;

      if (entry.state === 'failed' && age > RETENTION_MS_FAILED) {
        await journal.removeTemp(entry.id);
        await journal.remove(entry.id);
        cleaned++;
      } else if (entry.state === 'committed') {
        await journal.removeTemp(entry.id);
        await journal.remove(entry.id);
        cleaned++;
      }
    }

    log('cleanOldJournals', { cleaned, total: all.length });
  }

  private async cleanOrphanTemps(): Promise<void> {
    try {
      const temps = await filesystem.readDir('tmp/writes');
      const journals = await journal.listAll();
      const journalIds = new Set(journals.map((j) => j.id));
      let cleaned = 0;

      for (const t of temps) {
        if (t.type !== 'file') continue;
        const opId = t.name.replace(/\.json$/, '');
        if (!journalIds.has(opId)) {
          await filesystem.deleteFile(t.path);
          cleaned++;
          log('cleanOrphanTemps', { opId, path: t.path });
        }
      }

      log('cleanOrphanTemps', { cleaned, total: temps.length });
    } catch {
      // directory may not exist
    }
  }
}

export const cleanup = new Cleanup();
