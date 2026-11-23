import type { WriteOperation } from '../types';
import { filesystem } from '../core/filesystem';

function log(action: string, target: string, detail?: unknown) {
  console.log(`[WriteQueue] ${action} target=${target}`, detail ?? '');
}

export interface QueuedOperation extends WriteOperation {
  mutate: (current: unknown) => unknown;
  resolve: (value: unknown) => void;
  reject: (error: Error) => void;
}

interface FileBatch {
  target: string;
  ops: QueuedOperation[];
}

export class WriteQueue {
  private readonly queue: QueuedOperation[] = [];
  private processing = false;
  private readonly fileLocks = new Map<string, Promise<void>>();

  enqueue(op: Omit<QueuedOperation, 'resolve' | 'reject'>): Promise<unknown> {
    return new Promise((resolve, reject) => {
      const queued: QueuedOperation = { ...op, resolve, reject };
      this.queue.push(queued);
      log('enqueue', op.target, { type: op.type, entityId: op.entityId, queueSize: this.queue.length });
      void this.process();
    });
  }

  async flush(): Promise<void> {
    log('flush', '*', { queueSize: this.queue.length, processing: this.processing });
    while (this.queue.length > 0 || this.processing) {
      await new Promise((r) => setTimeout(r, 50));
    }
    log('flush', '*', { status: 'done' });
  }

  private async process(): Promise<void> {
    if (this.processing) return;
    this.processing = true;

    try {
      while (this.queue.length > 0) {
        const batches = this.groupByTarget();
        for (const batch of batches) {
          await this.processBatch(batch);
        }
      }
    } finally {
      this.processing = false;
    }
  }

  private groupByTarget(): FileBatch[] {
    const byTarget = new Map<string, QueuedOperation[]>();
    while (this.queue.length > 0) {
      const op = this.queue.shift()!;
      const group = byTarget.get(op.target) ?? [];
      group.push(op);
      byTarget.set(op.target, group);
    }
    return Array.from(byTarget.entries()).map(([target, ops]) => ({ target, ops }));
  }

  private async processBatch(batch: FileBatch): Promise<void> {
    const { target, ops } = batch;

    const existing = this.fileLocks.get(target);
    const promise = (async () => {
      if (existing) await existing;

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

      let result = currentData;
      for (const op of ops) {
        try {
          result = op.mutate(result);
        } catch (mutateErr) {
          const err = mutateErr as Error;
          for (const op2 of ops) {
            op2.reject(new Error(`Mutate failed for ${op2.entityId}: ${err.message}`));
          }
          log('processBatch', target, { status: 'mutate-failed', error: err.message });
          throw err;
        }
      }

      let payload: string;
      try {
        payload = JSON.stringify(result);
      } catch (serializeErr) {
        const err = serializeErr as Error;
        for (const op of ops) {
          op.reject(new Error(`Serialize failed: ${err.message}`));
        }
        log('processBatch', target, { status: 'serialize-failed', error: err.message });
        throw err;
      }

      try {
        await filesystem.writeFile(target, payload);
      } catch (writeErr) {
        const err = writeErr as Error;
        for (const op of ops) {
          op.reject(new Error(`Write failed for ${op.entityId}: ${err.message}`));
        }
        log('processBatch', target, { status: 'write-failed', error: err.message });
        throw err;
      }

      for (const op of ops) {
        op.resolve(result);
      }

      log('processBatch', target, {
        operations: ops.length,
        entityIds: ops.map((o) => o.entityId),
        status: 'committed',
      });
    })();

    this.fileLocks.set(target, promise);
    try {
      await promise;
    } finally {
      if (this.fileLocks.get(target) === promise) {
        this.fileLocks.delete(target);
      }
    }
  }
}

export const writeQueue = new WriteQueue();
