import { nanoid } from 'nanoid';

import { DocumentItem, NewDocument } from '../schemas';
import {
  DB,
  toNullString,
  toNullJSON,
  toNullInt,
  getNullableString,
  parseNullableJSON,
  currentTimestampMs,
} from '@/types/database';

export class DocumentModel {
  private userId: string;

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  create = async (params: Omit<NewDocument, 'userId'>) => {
    const now = currentTimestampMs();

    const result = await DB.CreateDocument({
      id: nanoid(),
      title: toNullString(params.title as any),
      content: toNullString(params.content as any),
      fileType: toNullString(params.fileType as any),
      filename: toNullString(params.filename as any),
      totalCharCount: toNullInt(params.totalCharCount as any),
      totalLineCount: toNullInt(params.totalLineCount as any),
      metadata: toNullJSON(params.metadata),
      pages: toNullInt(params.pages as any),
      sourceType: toNullString(params.sourceType as any),
      source: toNullString(params.source as any),
      fileId: toNullString(params.fileId as any),
      userId: this.userId,
      clientId: toNullString(params.clientId as any),
      editorData: toNullJSON(params.editorData),
      createdAt: now,
      updatedAt: now,
    });

    return this.mapDocument(result);
  };

  delete = async (id: string) => {
    await DB.DeleteDocument({
      id,
      userId: this.userId,
    });
  };

  deleteAll = async () => {
    await DB.DeleteAllDocuments(this.userId);
  };

  query = async () => {
    const results = await DB.ListDocuments({
      userId: this.userId,
      limit: 1000,
      offset: 0,
    });
    return results.map((r) => this.mapDocument(r));
  };

  findById = async (id: string) => {
    try {
      const result = await DB.GetDocument({
        id,
        userId: this.userId,
      });
      return this.mapDocument(result);
    } catch {
      return undefined;
    }
  };

  update = async (id: string, value: Partial<DocumentItem>) => {
    const now = currentTimestampMs();

    await DB.UpdateDocument({
      id,
      userId: this.userId,
      title: toNullString(value.title as any),
      content: toNullString(value.content as any),
      metadata: toNullJSON(value.metadata),
      editorData: toNullJSON(value.editorData),
      updatedAt: now,
    });
  };

  // **************** Helper *************** //

  private mapDocument = (doc: any): DocumentItem => {
    return {
      id: doc.id,
      title: getNullableString(doc.title as any),
      content: getNullableString(doc.content as any),
      fileType: getNullableString(doc.fileType as any),
      filename: getNullableString(doc.filename as any),
      totalCharCount: doc.totalCharCount,
      totalLineCount: doc.totalLineCount,
      metadata: parseNullableJSON(doc.metadata as any),
      pages: doc.pages,
      sourceType: getNullableString(doc.sourceType as any),
      source: getNullableString(doc.source as any),
      fileId: getNullableString(doc.fileId as any),
      userId: doc.userId,
      clientId: getNullableString(doc.clientId as any),
      editorData: parseNullableJSON(doc.editorData as any),
      createdAt: new Date(doc.createdAt),
      updatedAt: new Date(doc.updatedAt),
    } as DocumentItem;
  };
}
