import type { JournalEntry } from '../types';
import { filesystem } from './filesystem';
import { journal } from './journal';

function log(action: string, opId: string, detail?: unknown) {
  console.log(`[AtomicWrite] ${action} opId=${opId}`, detail ?? '');
}

export interface AtomicWriteResult {
  opId: string;
  target: string;
  success: boolean;
}

export class AtomicWrite {
  async execute(
    opId: string,
    target: string,
    mutate: (current: unknown) => unknown,
  ): Promise<AtomicWriteResult> {
    log('execute', opId, { target });

    await journal.create({
      id: opId,
      target,
      state: 'prepared',
      createdAt: new Date().toISOString(),
    });

    let currentData: unknown = null;
    try {
      const exists = await filesystem.fileExists(target);
      if (exists) {
        const raw = await filesystem.readFile(target);
        currentData = JSON.parse(raw);
      }
    } catch {
      currentData = null;
    }

    let newData: unknown;
    try {
      newData = mutate(currentData);
    } catch (error) {
      await journal.updateState(opId, 'failed');
      log('execute', opId, { target, status: 'mutate-failed', error: (error as Error).message });
      return { opId, target, success: false };
    }

    let payload: string;
    try {
      payload = JSON.stringify(newData);
    } catch (error) {
      await journal.updateState(opId, 'failed');
      log('execute', opId, { target, status: 'serialize-failed', error: (error as Error).message });
      return { opId, target, success: false };
    }

    await journal.writeTemp(opId, payload);

    try {
      JSON.parse(payload);
    } catch (error) {
      await journal.updateState(opId, 'failed');
      await journal.removeTemp(opId);
      log('execute', opId, { target, status: 'validate-failed', error: (error as Error).message });
      return { opId, target, success: false };
    }

    await journal.updateState(opId, 'committing');
    try {
      await filesystem.writeFile(target, payload);
      await journal.removeTemp(opId);
    } catch (error) {
      await journal.updateState(opId, 'failed');
      await journal.removeTemp(opId);
      log('execute', opId, { target, status: 'replace-failed', error: (error as Error).message });
      return { opId, target, success: false };
    }

    await journal.updateState(opId, 'committed');
    log('execute', opId, { target, status: 'committed' });
    return { opId, target, success: true };
  }

  async recover(): Promise<void> {
    const pending = await journal.listPending();
    log('recover', '*', { pendingCount: pending.length });

    for (const entry of pending) {
      await this.recoverEntry(entry);
    }
  }

  private async recoverEntry(entry: JournalEntry): Promise<void> {
    const { id: opId, target, state } = entry;

    switch (state) {
      case 'prepared':
        await journal.removeTemp(opId);
        await journal.remove(opId);
        log('recover', opId, { target, state, action: 'discarded' });
        break;

      case 'committing':
        const tempData = await journal.readTemp(opId);
        if (tempData) {
          try {
            await filesystem.writeFile(target, tempData);
            await journal.removeTemp(opId);
            await journal.updateState(opId, 'committed');
            log('recover', opId, { target, state, action: 're-applied-and-committed' });
          } catch {
            await journal.updateState(opId, 'failed');
            log('recover', opId, { target, state, action: 're-apply-failed' });
          }
        } else {
          await journal.updateState(opId, 'committed');
          log('recover', opId, { target, state, action: 'optimistic-commit' });
        }
        break;

      case 'failed':
        await journal.removeTemp(opId);
        await journal.remove(opId);
        log('recover', opId, { target, state, action: 'cleaned-up' });
        break;

      case 'committed':
        await journal.removeTemp(opId);
        await journal.remove(opId);
        log('recover', opId, { target, state, action: 'cleaned-up' });
        break;
    }
  }
}

export const atomicWrite = new AtomicWrite();
