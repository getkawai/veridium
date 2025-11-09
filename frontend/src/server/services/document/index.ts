import { DB } from '@@/database/sql/models';
import debug from 'debug';

import { LobeDocument } from '@/types/document';
import { GetDocument, DeleteDocument } from '@@/github.com/kawai-network/veridium/internal/database/generated/queries';

import { FileService } from '../file';

const log = debug('lobe-chat:service:document');

export class DocumentService {
  userId: string;
  private fileService: FileService;

  constructor(db: DB, userId: string) {
    this.userId = userId;
    this.fileService = new FileService(db, userId);
  }

  /**
   * Parse and save file content (handled in Go backend)
   */
  async parseFile(fileId: string): Promise<LobeDocument> {
    const { filePath, file, cleanup } = await this.fileService.downloadFileToLocal(fileId);

    const logPrefix = `[${file.name}]`;
    log(`${logPrefix} Processing file in Go backend: ${filePath}`);

    try {
      // Import FileProcessorService dynamically to avoid circular dependencies
      const { ProcessFileForStorage } = await import('@@/github.com/kawai-network/veridium/fileprocessorservice');

      // Single Go call: parse + save to database + RAG processing
      const result = await ProcessFileForStorage(
        filePath,
        file.name,
        file.fileType,
        this.userId,
        true, // enableRAG
      );

      if (!result) {
        throw new Error('ProcessFileForStorage returned null');
      }

      log(`${logPrefix} File processed successfully`, {
        fileId: result.fileId,
        documentId: result.documentId,
        chunks: result.chunkIds?.length || 0,
      });

      // Fetch document from database (already saved by Go)
      const document = await this.getDocument(result.documentId);

      return document;
    } catch (error) {
      console.error(`${logPrefix} File processing failed:`, error);
      throw error;
    } finally {
      cleanup();
    }
  }

  /**
   * Get document by ID (read-only)
   */
  async getDocument(documentId: string): Promise<LobeDocument> {
    const doc = await GetDocument({
      id: documentId,
      userId: this.userId,
    });

    // Helper to extract string from NullString
    const getNullableString = (ns: any): string | undefined => {
      if (!ns) return undefined;
      if (typeof ns === 'string') return ns;
      if (ns.String && ns.Valid) return ns.String;
      return undefined;
    };

    return {
      id: doc.id,
      title: getNullableString(doc.title),
      content: getNullableString(doc.content),
      fileType: doc.fileType,
      filename: getNullableString(doc.filename),
      totalCharCount: doc.totalCharCount || 0,
      totalLineCount: doc.totalLineCount || 0,
      metadata: doc.metadata ? JSON.parse(getNullableString(doc.metadata) || '{}') : undefined,
      pages: doc.pages ? JSON.parse(getNullableString(doc.pages) || '[]') : undefined,
      sourceType: doc.sourceType,
      source: doc.source,
      fileId: getNullableString(doc.fileId),
      userId: doc.userId,
      clientId: getNullableString(doc.clientId),
      createdAt: new Date(doc.createdAt),
      updatedAt: new Date(doc.updatedAt),
    } as LobeDocument;
  }

  /**
   * Delete document
   */
  async deleteDocument(documentId: string): Promise<void> {
    await DeleteDocument({
      id: documentId,
      userId: this.userId,
    });
  }
}
