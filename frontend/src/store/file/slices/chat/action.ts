// import { t } from 'i18next';
import { StateCreator } from 'zustand/vanilla';

// import { notification } from '@/components/AntdStaticMethods';
import { FILE_UPLOAD_BLACKLIST } from '@/const/file';
// import { ragService } from '@/services/rag';
import {
  UploadFileListDispatch,
  uploadFileListReducer,
} from '@/store/file/reducers/uploadFileList';
import { FileListItem } from '@/types/files';
import { UploadFileItem } from '@/types/files/upload';
import { setNamespace } from '@/utils/storeDebug';

// @ts-ignore - Wails binding
import * as FileService from '@@/github.com/kawai-network/veridium/internal/services/fileservice';
import { ProcessFileForStorage } from '@@/github.com/kawai-network/veridium/fileprocessorservice';

import { FileStore } from '../../store';
import { AsyncTaskStatus } from '@/types/asyncTask';

const n = setNamespace('chat');


export interface FileAction {
  clearChatUploadFileList: () => void;
  dispatchChatUploadFileList: (payload: UploadFileListDispatch) => void;

  removeChatUploadFile: (id: string) => Promise<void>;
  startAsyncTask: (
    fileId: string,
    runner: (id: string) => Promise<string>,
    onFileItemChange: (fileItem: FileListItem) => void,
  ) => Promise<void>;

  uploadChatFiles: (files: File[]) => Promise<void>;
}

export const createFileSlice: StateCreator<
  FileStore,
  [['zustand/devtools', never]],
  [],
  FileAction
> = (set, get) => ({
  clearChatUploadFileList: () => {
    set({ chatUploadFileList: [] }, false, n('clearChatUploadFileList'));
  },
  dispatchChatUploadFileList: (payload) => {
    const nextValue = uploadFileListReducer(get().chatUploadFileList, payload);
    if (nextValue === get().chatUploadFileList) return;

    set({ chatUploadFileList: nextValue }, false, `dispatchChatFileList/${payload.type}`);
  },
  removeChatUploadFile: async (id) => {
    const { dispatchChatUploadFileList } = get();

    // Dummy implementation for UI focus
    console.log('Removing chat upload file:', id);
    dispatchChatUploadFileList({ id, type: 'removeFile' });
  },

  startAsyncTask: async (id, runner, onFileItemUpdate) => {
    // Dummy implementation for UI focus
    console.log('Starting async task for file:', id);

    // Simulate task progress with mock data
    const mockFileItem: FileListItem = {
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
    };

    // Simulate async delay
    await new Promise(resolve => setTimeout(resolve, 1000));

    // Call callback with mock data
    onFileItemUpdate(mockFileItem);
  },

  uploadChatFiles: async (rawFiles) => {
    const { dispatchChatUploadFileList } = get();
    // 0. skip file in blacklist
    const files = rawFiles.filter((file) => !FILE_UPLOAD_BLACKLIST.includes(file.name));
    
    // 1. add files with local preview URLs (for desktop app)
    const uploadFiles: UploadFileItem[] = await Promise.all(
      files.map(async (file) => {
        let previewUrl: string | undefined = undefined;
        let base64Url: string | undefined = undefined;

        // only image and video can be previewed
        if (file.type.startsWith('image') || file.type.startsWith('video')) {
          // For small images (< 1MB), use base64 for instant preview
          // For larger images, we'll save to local storage and use file URL
          if (file.size < 1024 * 1024) {
            // Small file: use base64
            base64Url = await new Promise<string>((resolve) => {
              const reader = new FileReader();
              reader.onloadend = () => {
                resolve(reader.result as string);
              };
              reader.readAsDataURL(file);
            });
            previewUrl = base64Url;
          } else {
            // Large file: create blob URL temporarily
            // This will be replaced with local file URL after upload
            const data = await file.arrayBuffer();
            previewUrl = URL.createObjectURL(new Blob([data!], { type: file.type }));
          }
        }
        return { base64Url, file, id: file.name, previewUrl, status: 'pending' } as UploadFileItem;
      }),
    );

    dispatchChatUploadFileList({ files: uploadFiles, type: 'addFiles' });

    // upload files - save to local storage for desktop app
    const pools = files.map(async (file) => {

      try {
        // Update to uploading status
        dispatchChatUploadFileList({
          id: file.name,
          type: 'updateFile',
          value: { status: 'uploading', uploadState: { progress: 50, restTime: 1, speed: 1000 } },
        });

        // Read file content and convert to base64 for transfer
        const arrayBuffer = await file.arrayBuffer();
        const bytes = new Uint8Array(arrayBuffer);
        
        // Convert to base64 string for Wails binding transfer
        let binary = '';
        const chunkSize = 0x8000; // 32KB chunks
        for (let i = 0; i < bytes.length; i += chunkSize) {
          const chunk = bytes.subarray(i, Math.min(i + chunkSize, bytes.length));
          binary += String.fromCharCode.apply(null, Array.from(chunk));
        }
        const base64Data = btoa(binary);
        
        // Generate unique filename with timestamp
        const timestamp = Date.now();
        const sanitizedName = file.name.replace(/[^a-zA-Z0-9.-]/g, '_');
        const uniqueFileName = `${timestamp}-${sanitizedName}`;
        const uploadPath = `uploads/${uniqueFileName}`;
        
        // Save file to local storage using FileService with base64 data
        const savedKey = await FileService.UploadMedia(uploadPath, base64Data);
        
        // Process file for document storage (BLOCKING)
        // Note: File card already shows "uploading" status, so user sees progress
        // Backend will automatically skip RAG for images/videos
        await ProcessFileForStorage(
          savedKey,           // filePath
          file.name,          // filename
          file.type,          // fileType
          'system',           // userID
          true                // enableRAG (backend decides based on file type)
        );
        
        // Create local file URL for preview
        const localFileUrl = `/files/${savedKey}`;

        // Update to success with local file URL
        dispatchChatUploadFileList({
          id: file.name,
          type: 'updateFile',
          value: {
            status: 'success',
            uploadState: { progress: 100, restTime: 0, speed: 0 },
            // Update preview URL to use local file URL for large files
            ...(file.size >= 1024 * 1024 && (file.type.startsWith('image') || file.type.startsWith('video')) 
              ? { previewUrl: localFileUrl } 
              : {}),
          },
        });

        console.log('[uploadChatFiles] File processing completed for:', file.name);
      } catch (error) {
        console.error('[uploadChatFiles] Error processing file:', file.name, error);
        dispatchChatUploadFileList({
          id: file.name,
          type: 'updateFile',
          value: { status: 'error' },
        });
      }
    });

    await Promise.all(pools);
  },
});
