import type { EntityKind } from '../types';
import { filesystem } from './filesystem';

const BASE_PATH = 'data';

const ENTITY_DIRS = {
  account: true,
  category: true,
  label: true,
  notification: true,
  transaction: true,
  template: true,
  scheduled: true,
  budget: true,
  goal: true,
  setting: true,
} satisfies Record<EntityKind, true>;

export async function seed(): Promise<void> {
  const dirs = Object.keys(ENTITY_DIRS).map((kind) => `${BASE_PATH}/${kind}`);
  console.log('[Seed] start', { dirs });

  let created = 0;
  for (const dir of dirs) {
    try {
      if (await filesystem.fileExists(dir)) {
        console.log('[Seed] exists, skip', { dir });
        continue;
      }
      await filesystem.ensureDir(dir);
      created++;
      console.log('[Seed] created', { dir });
    } catch (error) {
      console.error('[Seed] ensure failed', { dir, error });
    }
  }

  console.log('[Seed] done', { count: dirs.length, created });
}
