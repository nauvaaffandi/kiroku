import { Filesystem, Directory } from '@capacitor/filesystem';
import type { FileEntry } from '../types';

const DATA_DIR = Directory.Documents;

function log(action: string, path: string, detail?: unknown) {
  console.log(`[FilesystemAdapter] ${action} path=${path}`, detail ?? '');
}

export class FilesystemAdapter {
  private readonly directory = DATA_DIR;

  async ensureDir(path: string): Promise<void> {
    log('ensureDir', path);
    await Filesystem.mkdir({ path, directory: this.directory, recursive: true });
  }

  async writeFile(path: string, data: string): Promise<void> {
    log('writeFile', path, { dataLength: data.length });
    await Filesystem.writeFile({ path, data, directory: this.directory });
  }

  async readFile(path: string): Promise<string> {
    log('readFile', path);
    const result = await Filesystem.readFile({ path, directory: this.directory });
    return result.data as string;
  }

  async deleteFile(path: string): Promise<void> {
    log('deleteFile', path);
    await Filesystem.deleteFile({ path, directory: this.directory });
  }

  async fileExists(path: string): Promise<boolean> {
    log('fileExists', path);
    try {
      await Filesystem.stat({ path, directory: this.directory });
      return true;
    } catch {
      return false;
    }
  }

  async readDir(path: string): Promise<FileEntry[]> {
    log('readDir', path);
    const result = await Filesystem.readdir({ path, directory: this.directory });
    return (result.files ?? []).map((f) => ({
      path: `${path}/${f.name}`,
      uri: f.uri ?? '',
      name: f.name,
      type: f.type === 'file' ? 'file' : 'directory',
      size: f.size ?? 0,
      lastModified: f.mtime ?? 0,
    }));
  }

  async rename(from: string, to: string): Promise<void> {
    log('rename', from, { to });
    await Filesystem.rename({ from, to, directory: this.directory });
  }
}

export const filesystem = new FilesystemAdapter();
