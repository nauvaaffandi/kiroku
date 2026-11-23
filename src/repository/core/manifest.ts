import { filesystem } from './filesystem';

const MANIFEST_PATH = 'manifest.json';

export interface Manifest {
  storageFormatVersion: number;
  schemaVersion: number;
  createdAt: string;
  updatedAt: string;
  fileRevisions: Record<string, { version: number; checksum?: string; updatedAt: string }>;
}

const DEFAULT_MANIFEST: Manifest = {
  storageFormatVersion: 1,
  schemaVersion: 1,
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString(),
  fileRevisions: {},
};

function log(action: string, detail?: unknown) {
  console.log(`[Manifest] ${action}`, detail ?? '');
}

export class ManifestManager {
  private manifest: Manifest = { ...DEFAULT_MANIFEST };

  async load(): Promise<Manifest> {
    try {
      const raw = await filesystem.readFile(MANIFEST_PATH);
      this.manifest = JSON.parse(raw) as Manifest;
      log('load', { storageFormatVersion: this.manifest.storageFormatVersion });
      return this.manifest;
    } catch {
      await this.save();
      log('load', { status: 'created-default' });
      return this.manifest;
    }
  }

  async save(): Promise<void> {
    this.manifest.updatedAt = new Date().toISOString();
    await filesystem.writeFile(MANIFEST_PATH, JSON.stringify(this.manifest));
    log('save', { updatedAt: this.manifest.updatedAt });
  }

  get(): Manifest {
    return { ...this.manifest };
  }

  updateFileRevision(path: string, checksum?: string): void {
    this.manifest.fileRevisions[path] = {
      version: (this.manifest.fileRevisions[path]?.version ?? 0) + 1,
      checksum,
      updatedAt: new Date().toISOString(),
    };
  }

  removeFileRevision(path: string): void {
    delete this.manifest.fileRevisions[path];
  }
}

export const manifestManager = new ManifestManager();
