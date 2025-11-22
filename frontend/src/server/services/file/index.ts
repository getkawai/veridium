import { TRPCError } from '@/types';
import { WriteTempFile, Cleanup } from 'bindings/github.com/kawai-network/veridium/tempfileservice';

import { FileServiceImpl, createFileServiceModule } from './impls';
import {
  DB,
  getNullableString,
} from '@/types/database';
import { Service as DBService } from '@@/github.com/kawai-network/veridium/internal/database';

/**
 * 文件服务类
 * 使用模块化实现方式，提供文件操作服务
 */
export class FileService {
  private userId: string;

  private impl: FileServiceImpl = createFileServiceModule();

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  /**
   * 删除文件
   */
  public async deleteFile(key: string) {
    return this.impl.deleteFile(key);
  }

  /**
   * 批量删除文件
   */
  public async deleteFiles(keys: string[]) {
    return this.impl.deleteFiles(keys);
  }

  /**
   * 获取文件内容
   */
  public async getFileContent(key: string): Promise<string> {
    return this.impl.getFileContent(key);
  }

  /**
   * 获取文件字节数组
   */
  public async getFileByteArray(key: string): Promise<Uint8Array> {
    return this.impl.getFileByteArray(key);
  }

  /**
   * 创建预签名上传URL
   */
  public async createPreSignedUrl(key: string): Promise<string> {
    return this.impl.createPreSignedUrl(key);
  }

  /**
   * 创建预签名预览URL
   */
  public async createPreSignedUrlForPreview(key: string, expiresIn?: number): Promise<string> {
    return this.impl.createPreSignedUrlForPreview(key, expiresIn);
  }

  /**
   * 上传内容
   */
  public async uploadContent(path: string, content: string) {
    return this.impl.uploadContent(path, content);
  }

  /**
   * 获取完整文件URL
   */
  public async getFullFileUrl(url?: string | null, expiresIn?: number): Promise<string> {
    return this.impl.getFullFileUrl(url, expiresIn);
  }

  /**
   * 从完整 URL中 提取 key
   */
  public getKeyFromFullUrl(url: string): string {
    return this.impl.getKeyFromFullUrl(url);
  }

  /**
   * 上传媒体文件
   */
  public async uploadMedia(key: string, buffer: Buffer): Promise<{ key: string }> {
    return this.impl.uploadMedia(key, buffer);
  }

  async downloadFileToLocal(
    fileId: string,
  ): Promise<{ cleanup: () => void; file: any; filePath: string }> {
    const file = await DB.GetFile({ id: fileId, userId: this.userId });
    if (!file) {
      throw new TRPCError({ code: 'BAD_REQUEST', message: 'File not found' });
    }

    const fileUrl = getNullableString(file.url as any);
    if (!fileUrl) {
      throw new TRPCError({ code: 'BAD_REQUEST', message: 'File URL not found' });
    }

    let content: Uint8Array | undefined;
    try {
      content = await this.getFileByteArray(fileUrl);
    } catch (e) {
      console.error(e);
      // if file not found, delete it from db
      if ((e as any).Code === 'NoSuchKey') {
        const fileHash = getNullableString(file.fileHash as any);
        await DBService.DeleteFileWithCascade({
          FileID: fileId,
          UserID: this.userId,
          RemoveGlobalFile: true,
          FileHash: fileHash || '',
        });
        throw new TRPCError({ code: 'BAD_REQUEST', message: 'File not found' });
      }
    }

    if (!content) throw new TRPCError({ code: 'BAD_REQUEST', message: 'File content is empty' });

    // Convert Uint8Array to base64 string for Wails binding
    const dataStr = btoa(String.fromCharCode(...content));
    const filePath = await WriteTempFile(dataStr, getNullableString(file.name as any) || 'unknown');
    return { cleanup: () => Cleanup(), file, filePath };
  }
}
