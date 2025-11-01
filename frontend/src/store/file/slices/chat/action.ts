import { t } from 'i18next';
import { StateCreator } from 'zustand/vanilla';

import { notification } from '@/components/AntdStaticMethods';
import { FILE_UPLOAD_BLACKLIST } from '@/const/file';
// import { ragService } from '@/services/rag';
import {
  UploadFileListDispatch,
  uploadFileListReducer,
} from '@/store/file/reducers/uploadFileList';
import { FileListItem } from '@/types/files';
import { UploadFileItem } from '@/types/files/upload';
import { isChunkingUnsupported } from '@/utils/isChunkingUnsupported';
import { setNamespace } from '@/utils/storeDebug';

import { FileStore } from '../../store';

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
      chunkingStatus: 'success',
      embeddingError: null,
      embeddingStatus: 'success',
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
    // 1. add files with base64
    const uploadFiles: UploadFileItem[] = await Promise.all(
      files.map(async (file) => {
        let previewUrl: string | undefined = undefined;
        let base64Url: string | undefined = undefined;

        // only image and video can be previewed, we create a previewUrl and base64Url for them
        if (file.type.startsWith('image') || file.type.startsWith('video')) {
          const data = await file.arrayBuffer();

          previewUrl = URL.createObjectURL(new Blob([data!], { type: file.type }));

          const bytes = new Uint8Array(data!);
          let binary = '';
          for (let i = 0; i < bytes.length; i++) {
            binary += String.fromCharCode(bytes[i]);
          }
          const base64 = btoa(binary);
          base64Url = `data:${file.type};base64,${base64}`;
        }

        return { base64Url, file, id: file.name, previewUrl, status: 'pending' } as UploadFileItem;
      }),
    );

    dispatchChatUploadFileList({ files: uploadFiles, type: 'addFiles' });

    // upload files with dummy implementation for UI focus
    const pools = files.map(async (file) => {
      console.log('Processing file:', file.name);

      // Simulate upload progress
      dispatchChatUploadFileList({
        id: file.name,
        type: 'updateFile',
        value: { status: 'uploading', uploadState: { progress: 50, restTime: 1, speed: 1000 } },
      });

      // Simulate async upload delay
      await new Promise(resolve => setTimeout(resolve, 800));

      const mockFileResult = {
        id: `mock-${file.name}-${Date.now()}`,
        url: `mock://file/${file.name}`,
      };

      // Update to success
      dispatchChatUploadFileList({
        id: file.name,
        type: 'updateFile',
        value: {
          status: 'success',
          uploadState: { progress: 100, restTime: 0, speed: 0 },
        },
      });

      // image don't need to be chunked and embedding
      if (isChunkingUnsupported(file.type)) return;

      // Dummy file processing simulation
      console.log('File processing completed for:', file.name);
    });

    await Promise.all(pools);
  },
});
