// Compatibility layer for Node.js APIs using Wails bindings
// This provides synchronous-looking APIs that internally use async Wails bindings

import {
  ReadFileSync as WailsReadFile,
  StatSync as WailsStat,
  FileExistsSync as WailsFileExists,
} from '@@/github.com/kawai-network/veridium/nodefsservice';

import {
  Extname as WailsExtname,
  Basename as WailsBasename,
  Dirname as WailsDirname,
} from '@@/github.com/kawai-network/veridium/nodepathservice';

import {
  Alloc as WailsBufferAlloc,
  From as WailsBufferFrom,
  IsBuffer as WailsBufferIsBuffer,
} from '@@/github.com/kawai-network/veridium/nodebufferservice';

// File System APIs - Synchronous compatibility layer
export class fs {
  static promises = {
    // Read file synchronously (returns Promise<string> for text, Uint8Array for binary)
    async readFile(path: string, encoding?: string): Promise<string | Uint8Array> {
      // Get the base64 encoded data from Wails
      const base64Data = await WailsReadFile(path);

      if (encoding === 'utf8' || encoding === 'utf-8') {
        // For text files, decode the base64 to string
        return atob(base64Data);
      } else if (!encoding || encoding === 'buffer') {
        // For binary files, convert base64 to Uint8Array
        const binaryString = atob(base64Data);
        const bytes = new Uint8Array(binaryString.length);
        for (let i = 0; i < binaryString.length; i++) {
          bytes[i] = binaryString.charCodeAt(i);
        }
        return bytes;
      } else {
        // Default to string for backward compatibility
        return atob(base64Data);
      }
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
  static async extname(pathStr: string): Promise<string> {
    return await WailsExtname(pathStr);
  }

  // Get base name of path
  static async basename(pathStr: string, ext?: string): Promise<string> {
    if (ext) {
      // Remove extension if provided
      const base = await WailsBasename(pathStr);
      if (base.endsWith(ext)) {
        return base.slice(0, -ext.length);
      }
      return base;
    }
    return await WailsBasename(pathStr);
  }

  // Get directory name
  static async dirname(pathStr: string): Promise<string> {
    return await WailsDirname(pathStr);
  }
}

// Buffer APIs - Compatibility layer
export class Buffer {
  private data: string; // Base64 encoded data

  constructor(data?: any, encoding?: string) {
    // Constructor should be synchronous, so we can't await here
    // This is a limitation - we'll need to make constructors async or use factory methods
    throw new Error('Use Buffer.from() or Buffer.alloc() instead of new Buffer()');
  }

  // Static methods
  static async alloc(size: number): Promise<Buffer> {
    const buf = new (Buffer as any)();
    buf.data = await WailsBufferAlloc(size);
    return buf;
  }

  static async from(data: any, encoding?: string): Promise<Buffer> {
    const buf = new (Buffer as any)();
    buf.data = await WailsBufferFrom(data, encoding || 'utf8');
    return buf;
  }

  static async isBuffer(obj: any): Promise<boolean> {
    return await WailsBufferIsBuffer(obj);
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
