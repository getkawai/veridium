// Compatibility layer for Node.js APIs using Wails bindings
// This provides synchronous-looking APIs that internally use async Wails bindings

import {
  ReadFileSync as WailsReadFile,
  StatSync as WailsStat,
  FileExistsSync as WailsFileExists,
} from '@/bindings/github.com/kawai-network/veridium/nodefsservice';

import {
  Extname as WailsExtname,
  Basename as WailsBasename,
  Dirname as WailsDirname,
} from '@/bindings/github.com/kawai-network/veridium/nodepathservice';

import {
  Alloc as WailsBufferAlloc,
  From as WailsBufferFrom,
  IsBuffer as WailsBufferIsBuffer,
  ToBytes as WailsBufferToBytes,
} from '@/bindings/github.com/kawai-network/veridium/nodebufferservice';

// File System APIs - Synchronous compatibility layer
export class fs {
  static promises = {
    // Read file synchronously (returns Promise<string>)
    async readFile(path: string, encoding?: string): Promise<string> {
      // For now, we assume UTF-8 encoding
      return WailsReadFile(path);
    },

    // Get file stats synchronously (returns Promise with stats object)
    async stat(path: string): Promise<{
      size: number;
      mtime: Date;
      ctime: Date;
      isDirectory(): boolean;
      isFile(): boolean;
    }> {
      const stats = await WailsStat(path);

      // Convert Wails stats to Node.js-like stats object
      return {
        size: stats.size || 0,
        mtime: new Date((stats.modTime as number) * 1000), // Convert Unix timestamp to Date
        ctime: new Date((stats.modTime as number) * 1000), // Use mtime as ctime for simplicity
        isDirectory(): boolean {
          return stats.isDir || false;
        },
        isFile(): boolean {
          return !stats.isDir;
        },
      };
    },
  };

  // Synchronous file existence check
  static existsSync(path: string): boolean {
    try {
      // This will throw if file doesn't exist
      WailsFileExists(path);
      return true;
    } catch {
      return false;
    }
  }
}

// Path APIs - Synchronous compatibility layer
export class path {
  // Get file extension
  static extname(pathStr: string): string {
    return WailsExtname(pathStr);
  }

  // Get base name of path
  static basename(pathStr: string, ext?: string): string {
    if (ext) {
      // Remove extension if provided
      const base = WailsBasename(pathStr);
      if (base.endsWith(ext)) {
        return base.slice(0, -ext.length);
      }
      return base;
    }
    return WailsBasename(pathStr);
  }

  // Get directory name
  static dirname(pathStr: string): string {
    return WailsDirname(pathStr);
  }
}

// Buffer APIs - Compatibility layer
export class Buffer {
  private data: string; // Base64 encoded data

  constructor(data?: any, encoding?: string) {
    if (data) {
      this.data = WailsBufferFrom(data, encoding || 'utf8');
    } else {
      this.data = WailsBufferAlloc(0);
    }
  }

  // Static methods
  static alloc(size: number): Buffer {
    const buf = new Buffer();
    buf.data = WailsBufferAlloc(size);
    return buf;
  }

  static from(data: any, encoding?: string): Buffer {
    const buf = new Buffer();
    buf.data = WailsBufferFrom(data, encoding || 'utf8');
    return buf;
  }

  static isBuffer(obj: any): boolean {
    return WailsBufferIsBuffer(obj);
  }

  // Instance methods
  toString(encoding?: string): string {
    // For simplicity, return the base64 data as-is
    // In a real implementation, you'd decode it
    return this.data;
  }

  // Get underlying bytes (as base64 string for compatibility)
  toBytes(): string {
    return this.data;
  }

  // Get length (would need to decode and get actual length)
  get length(): number {
    // This is a simplified implementation
    // In reality, you'd need to decode the base64 and get the actual length
    try {
      const decoded = atob(this.data);
      return decoded.length;
    } catch {
      return 0;
    }
  }
}
