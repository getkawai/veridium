import { StateCreator } from 'zustand/vanilla';

import {
  UploadFileListDispatch,
  uploadFileListReducer,
} from '@/store/file/reducers/uploadFileList';
import { FileListItem } from '@/types/files';
import { setNamespace } from '@/utils/storeDebug';

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
});
