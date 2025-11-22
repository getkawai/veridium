import { TRPCError } from '@/types';
import * as GoFileService from 'bindings/github.com/kawai-network/veridium/internal/services/file/fileservice';

/**
 * FileService - Wrapper around Go backend file service
 * This class provides a TypeScript interface to the Go file service implementation
 */
export class FileService {
  constructor(_db: any, _userId: string) {
    // Constructor kept for backward compatibility
  }

  /**
   * 删除文件
   */
  public async deleteFile(key: string) {
    return GoFileService.DeleteFile(key);
  }

  /**
   * 批量删除文件
   */
  public async deleteFiles(keys: string[]) {
    return GoFileService.DeleteFiles(keys);
  }

  /**
   * 获取文件内容
   */
  public async getFileContent(key: string): Promise<string> {
    return GoFileService.GetFileContent(key);
  }

  /**
   * 获取文件字节数组
   */
  public async getFileByteArray(key: string): Promise<Uint8Array> {
    const base64 = await GoFileService.GetFileByteArray(key);
    // Convert base64 string to Uint8Array
    const binaryString = atob(base64);
    const bytes = new Uint8Array(binaryString.length);
    for (let i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i);
    }
    return bytes;
  }

  /**
   * 创建预签名上传URL
   */
  public async createPreSignedUrl(key: string): Promise<string> {
    return GoFileService.CreatePreSignedUrl(key);
  }

  /**
   * 创建预签名预览URL
   */
  public async createPreSignedUrlForPreview(key: string, _expiresIn?: number): Promise<string> {
    // Note: Go implementation doesn't support expiresIn parameter yet
    return GoFileService.CreatePreSignedUrlForPreview(key);
  }

  /**
   * 上传内容
   */
  public async uploadContent(path: string, content: string) {
    return GoFileService.UploadContent(path, content);
  }

  /**
   * 获取完整文件URL
   */
  public async getFullFileUrl(url?: string | null, _expiresIn?: number): Promise<string> {
    if (!url) return '';
    // Note: Go implementation doesn't support expiresIn parameter yet
    return GoFileService.GetFullFileUrl(url);
  }

  /**
   * 从完整 URL中 提取 key
   */
  public getKeyFromFullUrl(url: string): string {
    // This is synchronous in the original, but Go binding returns a promise
    // We'll need to make this async or handle it differently
    throw new Error('getKeyFromFullUrl is not yet migrated - use async version');
  }

  /**
   * 从完整 URL中 提取 key (async version)
   */
  public async getKeyFromFullUrlAsync(url: string): Promise<string> {
    return GoFileService.GetKeyFromFullUrl(url);
  }

  /**
   * 上传媒体文件
   */
  public async uploadMedia(key: string, buffer: Buffer): Promise<{ key: string }> {
    // Convert Buffer to base64 string for Go binding
    const base64 = buffer.toString('base64');
    const resultKey = await GoFileService.UploadMedia(key, base64);
    return { key: resultKey };
  }

  /**
   * 下载文件到本地临时存储
   */
  async downloadFileToLocal(
    fileId: string,
  ): Promise<{ cleanup: () => void; file: any; filePath: string }> {
    try {
      const [cleanup, file, filePath] = await GoFileService.DownloadFileToLocal(fileId);
      return {
        cleanup: cleanup || (() => { }),
        file,
        filePath,
      };
    } catch (e) {
      throw new TRPCError({ code: 'BAD_REQUEST', message: `Failed to download file: ${e}` });
    }
  }
}
