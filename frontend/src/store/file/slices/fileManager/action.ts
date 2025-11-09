import { useEffect } from 'react';
import { StateCreator } from 'zustand/vanilla';

import { FILE_UPLOAD_BLACKLIST, MAX_UPLOAD_FILE_COUNT } from '@/const/file';
// import { lambdaClient } from '@/libs/trpc/client';
// import { ragService } from '@/services/rag';
import {
  UploadFileListDispatch,
  uploadFileListReducer,
} from '@/store/file/reducers/uploadFileList';
import { FileListItem, QueryFileListParams } from '@/types/files';
import { AsyncTaskStatus } from '@/types/asyncTask';
import { isChunkingUnsupported } from '@/utils/isChunkingUnsupported';

import { FileStore } from '../../store';
import { fileManagerSelectors } from './selectors';


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

  useFetchFileItem: (id?: string) => void;
  useFetchFileManage: (params: QueryFileListParams) => void;
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
    console.log('Embedding chunks for files:', fileIds);

    // toggle file ids
    get().toggleEmbeddingIds(fileIds);

    // Simulate async operation
    await new Promise(resolve => setTimeout(resolve, 1000));

    // Mock success
    get().toggleEmbeddingIds(fileIds, false);
  },
  parseFilesToChunks: async (ids: string[], params) => {
    // Dummy implementation for UI focus
    console.log('Parsing files to chunks:', ids, params);

    // toggle file ids
    get().toggleParsingIds(ids);

    // Simulate async operation
    await new Promise(resolve => setTimeout(resolve, 1500));

    // Mock success
    get().toggleParsingIds(ids, false);
  },
  pushDockFileList: async (rawFiles, knowledgeBaseId) => {
    const { dispatchDockFileList } = get();

    // 0. Process ZIP files with dummy implementation for UI focus
    const filesToUpload: File[] = [];
    for (const file of rawFiles) {
      if (file.type === 'application/zip' || file.name.endsWith('.zip')) {
        // Dummy ZIP extraction for UI focus
        console.log('Mock extracting ZIP file:', file.name);
        // Simulate extracting 2-3 files from ZIP
        const mockExtractedFiles = [
          new File(['mock content 1'], `extracted1_${file.name}.txt`, { type: 'text/plain' }),
          new File(['mock content 2'], `extracted2_${file.name}.pdf`, { type: 'application/pdf' }),
        ];
        filesToUpload.push(...mockExtractedFiles);
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

    // 3. Upload files with dummy implementation for UI focus
    const uploadResults = await Promise.all(
      files.map(async (file) => {
        // Simulate upload progress
        dispatchDockFileList({
          id: file.name,
          type: 'updateFile',
          value: { status: 'uploading', uploadState: { progress: 50, restTime: 1, speed: 1000 } },
        });

        // Simulate completion
        await new Promise(resolve => setTimeout(resolve, 500));

        dispatchDockFileList({
          id: file.name,
          type: 'updateFile',
          value: {
            status: 'success',
            uploadState: { progress: 100, restTime: 0, speed: 0 },
          },
        });

        return { file, fileId: `mock-${file.name}`, fileType: file.type };
      })
    );

    // 4. auto-embed files that support chunking
    const fileIdsToEmbed = uploadResults
      .filter(({ fileType, fileId }) => fileId && !isChunkingUnsupported(fileType))
      .map(({ fileId }) => fileId!);

    if (fileIdsToEmbed.length > 0) {
      await get().parseFilesToChunks(fileIdsToEmbed, { skipExist: false });
    }
  },

  reEmbeddingChunks: async (id) => {
    if (fileManagerSelectors.isCreatingChunkEmbeddingTask(id)(get())) return;

    // Dummy implementation for UI focus
    console.log('Re-embedding chunks for file:', id);

    // toggle file ids
    get().toggleEmbeddingIds([id]);

    // Simulate async operation
    await new Promise(resolve => setTimeout(resolve, 2000));

    // Mock success
    get().toggleEmbeddingIds([id], false);
  },
  reParseFile: async (id) => {
    // Dummy implementation for UI focus
    console.log('Re-parsing file:', id);

    // toggle file ids
    get().toggleParsingIds([id]);

    // Simulate async operation
    await new Promise(resolve => setTimeout(resolve, 1200));

    // Mock success
    get().toggleParsingIds([id], false);
  },
  refreshFileList: async () => {
    // Dummy implementation for UI focus - avoid mutate to prevent loops
    console.log('Refreshing file list');
    // In real implementation, this would trigger SWR revalidation
  },
  removeAllFiles: async () => {
    // Dummy implementation for UI focus
    console.log('Removing all files');
    // Mock success
  },
  removeFileItem: async (id) => {
    // Dummy implementation for UI focus
    console.log('Removing file item:', id);
    // Mock success - file would be removed from UI state
  },

  removeFiles: async (ids) => {
    // Dummy implementation for UI focus
    console.log('Removing files:', ids);
    // Mock success - files would be removed from UI state
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

  useFetchFileItem: (id) => {
    useEffect(() => {
      if (!id) return;

      const fetchFileItem = async () => {
        try {
          // Dummy implementation for UI focus
          console.log('Fetching file item:', id);
          const fileItem: FileListItem | undefined = id ? {
            id,
            name: `Mock File ${id.slice(0, 8)}`,
            fileType: 'application/pdf',
            size: 1024000,
            url: `mock://file/${id}`,
            createdAt: new Date(),
            updatedAt: new Date(),
            chunkCount: 10,
            chunkingError: null,
            chunkingStatus: AsyncTaskStatus.Success,
            embeddingError: null,
            embeddingStatus: AsyncTaskStatus.Success,
            finishEmbedding: true,
          } : undefined;
          
          // Store the file item if needed
          // set({ fileItemMap: { ...get().fileItemMap, [id]: fileItem } });
        } catch (error) {
          console.error('[useFetchFileItem] Error:', error);
        }
      };

      fetchFileItem();
    }, [id]);
  },

  useFetchFileManage: (params) => {
    useEffect(() => {
      const fetchFileList = async () => {
        try {
          // Dummy implementation for UI focus
          console.log('Fetching file list:', params);
          const mockFiles: FileListItem[] = [
            {
              id: 'mock-1',
              name: 'Sample Document.pdf',
              fileType: 'application/pdf',
              size: 2048000,
              url: 'mock://file/mock-1',
              createdAt: new Date(Date.now() - 86400000),
              updatedAt: new Date(Date.now() - 86400000),
              chunkCount: 25,
              chunkingError: null,
              chunkingStatus: AsyncTaskStatus.Success,
              embeddingError: null,
              embeddingStatus: AsyncTaskStatus.Success,
              finishEmbedding: true,
            },
            {
              id: 'mock-2',
              name: 'Presentation.pptx',
              fileType: 'application/vnd.openxmlformats-officedocument.presentationml.presentation',
              size: 5120000,
              url: 'mock://file/mock-2',
              createdAt: new Date(Date.now() - 172800000),
              updatedAt: new Date(Date.now() - 172800000),
              chunkCount: 15,
              chunkingError: null,
              chunkingStatus: AsyncTaskStatus.Success,
              embeddingError: null,
              embeddingStatus: AsyncTaskStatus.Processing,
              finishEmbedding: false,
            },
            {
              id: 'mock-3',
              name: 'Spreadsheet.xlsx',
              fileType: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
              size: 1024000,
              url: 'mock://file/mock-3',
              createdAt: new Date(Date.now() - 259200000),
              updatedAt: new Date(Date.now() - 259200000),
              chunkCount: 8,
              chunkingError: null,
              chunkingStatus: AsyncTaskStatus.Success,
              embeddingError: null,
              embeddingStatus: AsyncTaskStatus.Success,
              finishEmbedding: true,
            },
          ];
          
          set({ fileList: mockFiles, queryListParams: params });
        } catch (error) {
          console.error('[useFetchFileManage] Error:', error);
        }
      };

      fetchFileList();
    }, [params]);
  },
});
