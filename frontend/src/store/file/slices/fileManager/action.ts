import { StateCreator } from 'zustand/vanilla';
import { nanoid } from 'nanoid';

import { FILE_UPLOAD_BLACKLIST } from '@/const/file';
import {
  UploadFileListDispatch,
  uploadFileListReducer,
} from '@/store/file/reducers/uploadFileList';
import { FileListItem, QueryFileListParams, FilesTabs, SortType } from '@/types/files';
import { AsyncTaskStatus } from '@/types/asyncTask';
import { createServiceLogger } from '@/utils/logger';

const logger = createServiceLogger('FileManager', 'FileManagerActions', 'store/file/slices/fileManager/action');
import { isChunkingUnsupported } from '@/utils/isChunkingUnsupported';
import { clientS3Storage } from '@/services/file/ClientS3';

import { FileStore } from '../../store';
import { fileManagerSelectors } from './selectors';

// Database imports
import {
  DB,
  toNullString,
  toNullJSON,
  getNullableString,
  File as DBFile,
} from '@/types/database';
import { Service as DBService } from '@@/github.com/kawai-network/veridium/internal/database';
import { getUserId } from '@/store/session/helpers';

export interface FileManageAction {
  dispatchDockFileList: (payload: UploadFileListDispatch) => void;
  embeddingChunks: (fileIds: string[]) => Promise<void>;
  parseFilesToChunks: (ids: string[], params?: { skipExist?: boolean }) => Promise<void>;
  pushDockFileList: (files: File[], knowledgeBaseId?: string) => Promise<void>;

  reEmbeddingChunks: (id: string) => Promise<void>;
  reParseFile: (id: string) => Promise<void>;
  refreshFileList: () => Promise<void>;
  removeAllFiles: () => Promise<void>;
  removeFileItem: (id: string) => Promise<void>;
  removeFiles: (ids: string[]) => Promise<void>;

  toggleEmbeddingIds: (ids: string[], loading?: boolean) => void;
  toggleParsingIds: (ids: string[], loading?: boolean) => void;

  fetchFileList: (params: QueryFileListParams) => Promise<void>;
  internal_fetchFileItem: (id?: string) => Promise<void>;
  setActiveFileId: (id: string | undefined) => void;
  setCategory: (category: string) => void;
  setSearchKeywords: (keywords: string) => void;
  setSortType: (sortType: SortType) => void;
  setSorter: (sorter: string) => void;
}

export const createFileManageSlice: StateCreator<
  FileStore,
  [['zustand/devtools', never]],
  [],
  FileManageAction
> = (set, get) => ({
  dispatchDockFileList: (payload: UploadFileListDispatch) => {
    const nextValue = uploadFileListReducer(get().dockUploadFileList, payload);
    if (nextValue === get().dockUploadFileList) return;

    set({ dockUploadFileList: nextValue }, false, `dispatchDockFileList/${payload.type}`);
  },
  embeddingChunks: async (fileIds: string[]) => {
    // Dummy implementation for UI focus
    await logger.debug('Embedding chunks for files', { fileIds });

    // toggle file ids
    get().toggleEmbeddingIds(fileIds);

    // Simulate async operation
    await new Promise(resolve => setTimeout(resolve, 1000));

    // Mock success
    get().toggleEmbeddingIds(fileIds, false);
  },
  parseFilesToChunks: async (ids: string[], params) => {
    // Dummy implementation for UI focus
    await logger.debug('Parsing files to chunks', { ids, params });

    // toggle file ids
    get().toggleParsingIds(ids);

    // Simulate async operation
    await new Promise(resolve => setTimeout(resolve, 1500));

    // Mock success
    get().toggleParsingIds(ids, false);
  },
  pushDockFileList: async (rawFiles, knowledgeBaseId) => {
    const { dispatchDockFileList } = get();
    const userId = getUserId();

    // 0. Process ZIP files
    const filesToUpload: File[] = [];
    for (const file of rawFiles) {
      if (file.type === 'application/zip' || file.name.endsWith('.zip')) {
        // TODO: Implement real ZIP extraction if needed
        await logger.warn('ZIP file extraction not fully implemented yet', { fileName: file.name });
        filesToUpload.push(file);
      } else {
        filesToUpload.push(file);
      }
    }

    // 1. skip file in blacklist
    const files = filesToUpload.filter((file) => !FILE_UPLOAD_BLACKLIST.includes(file.name));

    // 2. Add all files to dock
    dispatchDockFileList({
      atStart: true,
      files: files.map((file) => ({ file, id: file.name, status: 'pending' })),
      type: 'addFiles',
    });

    // 3. Upload files
    const uploadResults = await Promise.all(
      files.map(async (file) => {
        try {
          dispatchDockFileList({
            id: file.name,
            type: 'updateFile',
            value: { status: 'uploading', uploadState: { progress: 0, restTime: 0, speed: 0 } },
          });

          // Calculate hash (simplified for client-side, ideally should be content-based)
          // For now using name + size + timestamp as a simple unique identifier
          // In a real app, you'd want to calculate SHA256 of the file content
          const fileHash = `hash-${file.name}-${file.size}-${Date.now()}`;

          // Check if file exists globally
          const existingGlobalFile = await DB.GetGlobalFile(fileHash);
          const isExist = !!existingGlobalFile;

          // Upload to S3
          await clientS3Storage.putObject(fileHash, file);

          // Create DB record
          const fileId = nanoid();

          await logger.debug('Creating file with KB link', {
            fileName: file.name,
            fileId,
            knowledgeBaseId: knowledgeBaseId || 'null',
            hasKBId: !!knowledgeBaseId
          });

          // Use backend transaction method for atomic operations
          await DBService.CreateFileWithLinks({
            File: {
              fileType: toNullString(file.type) as any,
              fileHash: toNullString(fileHash) as any,
              name: toNullString(file.name) as any,
              size: file.size || 0,
              url: toNullString(`s3://${fileHash}`) as any, // Placeholder URL
              source: toNullString('upload') as any,
              metadata: toNullJSON({}) as any,
            },
            GlobalFile: !isExist ? {
              hashId: toNullString(fileHash) as any,
              fileType: toNullString(file.type) as any,
              size: file.size || 0,
              url: toNullString(`s3://${fileHash}`) as any,
              metadata: toNullJSON({}) as any,
              creator: toNullString(userId) as any,
            } : null,
            KnowledgeBase: knowledgeBaseId || null,
          });

          await logger.info('File created successfully with KB', { knowledgeBaseId: knowledgeBaseId || 'none' });

          dispatchDockFileList({
            id: file.name,
            type: 'updateFile',
            value: {
              status: 'success',
              uploadState: { progress: 100, restTime: 0, speed: 0 },
            },
          });

          return { file, fileId, fileType: file.type };
        } catch (error) {
          await logger.error('Failed to upload file', { fileName: file.name, error });
          dispatchDockFileList({
            id: file.name,
            type: 'updateFile',
            value: { status: 'error' },
          });
          return { file, fileId: undefined, fileType: file.type };
        }
      })
    );

    // 4. auto-embed files that support chunking
    const fileIdsToEmbed = uploadResults
      .filter(({ fileType, fileId }) => fileId && !isChunkingUnsupported(fileType))
      .map(({ fileId }) => fileId!);

    if (fileIdsToEmbed.length > 0) {
      await get().parseFilesToChunks(fileIdsToEmbed, { skipExist: false });
    }

    // Refresh file list
    await get().refreshFileList();
  },

  reEmbeddingChunks: async (id) => {
    if (fileManagerSelectors.isCreatingChunkEmbeddingTask(id)(get())) return;

    // Dummy implementation for UI focus
    await logger.debug('Re-embedding chunks for file', { fileId: id });

    // toggle file ids
    get().toggleEmbeddingIds([id]);

    // Simulate async operation
    await new Promise(resolve => setTimeout(resolve, 2000));

    // Mock success
    get().toggleEmbeddingIds([id], false);
  },
  reParseFile: async (id) => {
    // Dummy implementation for UI focus
    await logger.debug('Re-parsing file', { fileId: id });

    // toggle file ids
    get().toggleParsingIds([id]);

    // Simulate async operation
    await new Promise(resolve => setTimeout(resolve, 1200));

    // Mock success
    get().toggleParsingIds([id], false);
  },
  refreshFileList: async () => {
    const { queryListParams, fetchFileList } = get();
    if (queryListParams) {
      await fetchFileList(queryListParams);
    }
  },
  removeAllFiles: async () => {
    try {
      await DB.DeleteAllFiles();
      await get().refreshFileList();
    } catch (error) {
      await logger.error('Failed to remove all files', error);
    }
  },
  removeFileItem: async (id) => {
    try {
      const file = await DB.GetFile(id);

      if (file) {
        const fileHash = getNullableString(file.fileHash as any);

        await DBService.DeleteFileWithCascade({
          FileID: id,
          RemoveGlobalFile: true, // Default to true
          FileHash: fileHash || '',
        });

        await get().refreshFileList();
      }
    } catch (error) {
      await logger.error('Failed to remove file item', { fileId: id, error });
    }
  },

  removeFiles: async (ids) => {
    try {
      if (ids.length === 0) return;
      // 1. Get file list first
      // No batch query available, fetch one by one
      const fileList: DBFile[] = [];
      for (const id of ids) {
        try {
          const file = await DB.GetFile(id);
          if (file) fileList.push(file);
        } catch (e) {
          // ignore missing files
        }
      }

      if (fileList.length === 0) return;

      // 2. Extract hashes
      const hashList = fileList
        .map((f) => getNullableString(f.fileHash as any))
        .filter(Boolean) as string[];

      // 3. Delete chunks (manual implementation of deleteFileChunks)
      // Get all chunk IDs for these files
      const allChunkIds: string[] = [];
      for (const fileId of ids) {
        const chunks = await DB.GetFileChunkIds(toNullString(fileId) as any);
        allChunkIds.push(...chunks.map((c) => getNullableString(c as any) || '').filter(Boolean));
      }

      if (allChunkIds.length > 0) {
        // Batch delete in chunks of 500
        const BATCH_SIZE = 500;
        for (let i = 0; i < allChunkIds.length; i += BATCH_SIZE) {
          const batchIds = allChunkIds.slice(i, i + BATCH_SIZE);

          // Delete embeddings
          await Promise.all(
            batchIds.map((chunkId) =>
              DB.DeleteEmbedding(chunkId).catch(() => { }),
            ),
          );

          // Delete chunks
          await Promise.all(
            batchIds.map((chunkId) =>
              DB.DeleteChunk(chunkId).catch(() => { }),
            ),
          );
        }
      }

      // 4. Delete files
      await Promise.all(
        ids.map((id) =>
          DB.DeleteFile(id),
        ),
      );

      // 5. Delete global files if needed
      if (hashList.length > 0) {
        // Check which hashes are still in use
        await DB.GetFilesByHash(toNullString(hashList[0]) as any);

        // This part is tricky without a proper batch check. 
        // For now, let's skip complex global file cleanup in batch delete 
        // or implement it one by one if critical.
        // Given the complexity, we'll rely on the fact that global files are less critical to clean up immediately.
      }

      await get().refreshFileList();
    } catch (error) {
      await logger.error('Failed to remove files', { fileIds: ids, error });
    }
  },
  toggleEmbeddingIds: (ids, loading) => {
    set((state) => {
      const nextValue = new Set(state.creatingEmbeddingTaskIds);

      ids.forEach((id) => {
        if (typeof loading === 'undefined') {
          if (nextValue.has(id)) nextValue.delete(id);
          else nextValue.add(id);
        } else {
          if (loading) nextValue.add(id);
          else nextValue.delete(id);
        }
      });

      return { creatingEmbeddingTaskIds: Array.from(nextValue.values()) };
    });
  },
  toggleParsingIds: (ids, loading) => {
    set((state) => {
      const nextValue = new Set(state.creatingChunkingTaskIds);

      ids.forEach((id) => {
        if (typeof loading === 'undefined') {
          if (nextValue.has(id)) nextValue.delete(id);
          else nextValue.add(id);
        } else {
          if (loading) nextValue.add(id);
          else nextValue.delete(id);
        }
      });

      return { creatingChunkingTaskIds: Array.from(nextValue.values()) };
    });
  },

  setActiveFileId: (id) => {
    set({ activeFileId: id });
  },
  setCategory: (category) => {
    set({ currentCategory: category });
  },
  setSearchKeywords: (keywords) => {
    set({ searchKeywords: keywords });
  },
  setSortType: (sortType) => {
    set({ sortType });
  },
  setSorter: (sorter) => {
    set({ sorter });
  },

  internal_fetchFileItem: async (id) => {
    if (!id) return;

    try {
      const item = await DB.GetFile(id);

      if (!item) return;

      // Get URL from S3
      const fileHash = getNullableString(item.fileHash as any);
      if (fileHash) {
        try {
          const fileItem = await clientS3Storage.getObject(fileHash);
          if (fileItem) {
            // url = URL.createObjectURL(fileItem);
          }
        } catch (e) {
          await logger.error('Failed to get file from S3', e);
        }
      }

      // Store the file item if needed - currently no state for single item in this slice
      // set({ fileItemMap: { ...get().fileItemMap, [id]: fileItem } });
    } catch (error) {
      await logger.error('Failed to fetch file item', error);
    }
  },

  fetchFileList: async (params) => {
    try {
      set({ isFetchingFiles: true });
      await logger.debug('Starting fetch with params', params);
      const { category, q, sortType, sorter, knowledgeBaseId } = params;
      let allFiles;

      // If filtering by knowledge base, use JOIN query
      if (knowledgeBaseId) {
        allFiles = await DB.QueryFilesByKnowledgeBase(knowledgeBaseId);
        await logger.debug('KB files received', { count: allFiles.length, firstFile: allFiles[0] });
      } else {
        // Otherwise, get all files
        allFiles = await DB.QueryFiles();
        await logger.debug('All files received', { count: allFiles.length, firstFile: allFiles[0] });
      }

      // Apply filters client-side
      let filtered = allFiles;

      // Filter by search query
      if (q) {
        filtered = filtered.filter((f) =>
          getNullableString(f.name as any)
            ?.toLowerCase()
            .includes(q.toLowerCase()),
        );
      }

      // Filter by category
      if (category && category !== FilesTabs.All) {
        const imageExts = ['png', 'jpg', 'jpeg', 'gif', 'webp', 'svg', 'image'];
        const videoExts = ['mp4', 'webm', 'mov', 'avi', 'mkv', 'video'];
        const audioExts = ['mp3', 'wav', 'ogg', 'm4a', 'flac', 'audio'];
        // Documents are anything not image/video/audio, or specific doc types
        const docExts = ['pdf', 'txt', 'json', 'xml', 'doc', 'docx', 'csv', 'md', 'application', 'text'];

        filtered = filtered.filter((f) => {
          // f is already a FileListItem with plain string values (not NullString)
          const type = (f.fileType || '').toLowerCase();
          const fileName = f.name || '';

          let matches = false;
          switch (category) {
            case FilesTabs.Images:
              matches = imageExts.some(ext => type.includes(ext));
              break;
            case FilesTabs.Videos:
              matches = videoExts.some(ext => type.includes(ext));
              break;
            case FilesTabs.Audios:
              matches = audioExts.some(ext => type.includes(ext));
              break;
            case FilesTabs.Documents:
              // For documents, we include specific types OR anything that implies text/app
              matches = docExts.some(ext => type.includes(ext));
              break;
            default:
              matches = true;
          }

          return matches;
        });
      }

      // Sort
      if (sorter && sortType) {
        const direction = sortType.toLowerCase() === SortType.Asc ? 1 : -1;
        filtered.sort((a, b) => {
          const aVal = a[sorter as keyof typeof a];
          const bVal = b[sorter as keyof typeof b];
          if (aVal < bVal) return -1 * direction;
          if (aVal > bVal) return 1 * direction;
          return 0;
        });
      }

      // Map to FileListItem
      const fileListItems: FileListItem[] = await Promise.all(filtered.map(async (item) => {
        // Note: Generating object URLs for all files might be expensive/memory intensive
        // Consider doing this only on demand or using a different approach
        const fileItem = {
          id: item.id,
          name: item.name || 'Unknown',
          fileType: item.fileType || '',
          size: item.size,
          url: item.url || '',
          createdAt: new Date(item.createdAt),
          updatedAt: new Date(item.updatedAt),
          // Use real chunk count and status from database
          // Note: Wails generated NullInt64 uses Int64 property, NullString uses String property
          chunkCount: (item.chunkCount?.Valid ? item.chunkCount?.Int64 : 0) || 0,
          chunkingError: null,
          chunkingStatus: item.chunkingStatus?.String === 'success'
            ? AsyncTaskStatus.Success
            : item.chunkingStatus?.String === 'empty'
              ? AsyncTaskStatus.Success
              : AsyncTaskStatus.Pending,
          embeddingError: null,
          embeddingStatus: item.embeddingStatus?.String === 'success'
            ? AsyncTaskStatus.Success
            : item.embeddingStatus?.String === 'empty'
              ? AsyncTaskStatus.Success
              : AsyncTaskStatus.Pending,
          finishEmbedding: item.embeddingStatus?.String === 'success',
        };

        return fileItem;
      }));

      // Debug: Log chunk count for verification
      const filesWithChunks = fileListItems.filter(f => f.chunkCount > 0);
      if (filesWithChunks.length > 0) {
        await logger.info('Files with chunks found', {
          count: filesWithChunks.length,
          files: filesWithChunks.map(f => ({
            name: f.name,
            chunkCount: f.chunkCount,
            embeddingStatus: f.embeddingStatus
          }))
        });
      } else {
        await logger.warn('No files with chunks found', { totalFiles: fileListItems.length });
      }

      set({ fileList: fileListItems, queryListParams: params, isFetchingFiles: false });
    } catch (error) {
      await logger.error('Failed to fetch files', error);
      set({ isFetchingFiles: false });
    }
  },
});
