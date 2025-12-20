import { nanoid } from 'nanoid';

import { BaseClientService } from '@/services/baseClientService';
import { clientS3Storage } from '@/services/file/ClientS3';

import { IFileService } from './type';
import {
  DB,
  toNullString,
  toNullJSON,
  getNullableString,
  parseNullableJSON,
  currentTimestampMs,
  File as DBFile,
} from '@/types/database';
import { Service as DBService } from '@@/github.com/kawai-network/veridium/internal/database';

export class ClientService extends BaseClientService implements IFileService {
  createFile: IFileService['createFile'] = async (file) => {
    const { isExist } = await this.checkFileHash(file.hash!);

    // save to local storage
    // we may want to save to a remote server later

    const fileId = nanoid();
    const now = currentTimestampMs();

    await DBService.CreateFileWithLinks({
      File: {
        id: fileId,
        fileType: toNullString(file.fileType) as any,
        fileHash: toNullString(file.hash) as any,
        name: toNullString(file.name) as any,
        size: file.size || 0,
        url: toNullString(file.url) as any,
        source: toNullString('upload') as any, // Default source
        clientId: toNullString(this.userId) as any,
        metadata: toNullJSON(file.metadata) as any,
        chunkTaskId: toNullString('') as any,
        embeddingTaskId: toNullString('') as any,
        createdAt: now,
        updatedAt: now,
      },
      GlobalFile: !isExist ? {
        hashId: toNullString(file.hash) as any,
        fileType: toNullString(file.fileType) as any,
        size: file.size || 0,
        url: toNullString(file.url) as any,
        metadata: toNullJSON(file.metadata) as any,
        creator: toNullString(this.userId) as any,
        createdAt: now,
      } : null,
      KnowledgeBase: file.knowledgeBaseId || null,
    });

    // get file to base64 url
    let base64 = '';
    try {
      base64 = await this.getBase64ByFileHash(file.hash!);
    } catch (e) {
      console.warn('Failed to get base64 for file:', file.hash, e);
    }

    return {
      id: fileId,
      url: base64 ? `data:${file.fileType};base64,${base64}` : file.url || '',
    };
  };

  getFile: IFileService['getFile'] = async (id) => {
    const item = await DB.GetFile(id);
    if (!item) {
      throw new Error('file not found');
    }

    // arrayBuffer to url
    const fileHash = getNullableString(item.fileHash as any);
    if (!fileHash) throw new Error('file hash not found');

    const fileItem = await clientS3Storage.getObject(fileHash);
    if (!fileItem) throw new Error('file not found in storage');

    const url = URL.createObjectURL(fileItem);

    return {
      createdAt: new Date(item.createdAt),
      id,
      name: getNullableString(item.name as any) || '',
      size: item.size,
      type: getNullableString(item.fileType as any) || '',
      updatedAt: new Date(item.updatedAt),
      url,
    };
  };

  removeFile: IFileService['removeFile'] = async (id) => {
    const item = await DB.GetFile(id);
    if (!item) return;

    const fileHash = getNullableString(item.fileHash as any);

    await DBService.DeleteFileWithCascade({
      FileID: id,
      RemoveGlobalFile: false, // Original was false
      FileHash: fileHash || '',
    });
  };

  removeFiles: IFileService['removeFiles'] = async (ids) => {
    if (ids.length === 0) return;

    // 1. Get file list first
    const fileList: DBFile[] = [];
    for (const id of ids) {
      try {
        const file = await DB.GetFile(id);
        if (file) fileList.push(file);
      } catch (e) {
        // ignore
      }
    }

    if (fileList.length === 0) return;

    // 2. Delete chunks (simplified batch delete)
    for (const file of fileList) {
      const fileHash = getNullableString(file.fileHash as any);
      await DBService.DeleteFileWithCascade({
        FileID: file.id,
        RemoveGlobalFile: false,
        FileHash: fileHash || '',
      });
    }
  };

  removeAllFiles: IFileService['removeAllFiles'] = async () => {
    return DB.DeleteAllFiles();
  };

  checkFileHash: IFileService['checkFileHash'] = async (hash) => {
    const item = await DB.GetGlobalFile(hash);

    if (!item) return { isExist: false };

    return {
      fileType: getNullableString(item.fileType as any),
      isExist: true,
      metadata: parseNullableJSON(item.metadata as any),
      size: item.size,
      url: getNullableString(item.url as any),
    };
  };

  private getBase64ByFileHash = async (hash: string) => {
    const fileItem = await clientS3Storage.getObject(hash);
    if (!fileItem) throw new Error('file not found');

    const arrayBuffer = await fileItem.arrayBuffer();
    const bytes = new Uint8Array(arrayBuffer);
    let binary = '';
    for (let i = 0; i < bytes.length; i++) {
      binary += String.fromCharCode(bytes[i]);
    }
    return btoa(binary);
  };
}
