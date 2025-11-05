import { shallow } from 'zustand/shallow';
import { createWithEqualityFn } from 'zustand/traditional';
import { StateCreator } from 'zustand/vanilla';

import { createDevtools } from '../middleware/createDevtools';
import { FilesStoreState, initialState } from './initialState';
import { FileAction, createFileSlice } from './slices/chat';
import { FileChunkAction, createFileChunkSlice } from './slices/chunk';
import { FileManageAction, createFileManageSlice } from './slices/fileManager';
import { FileUploadAction, createFileUploadSlice } from './slices/upload/action';
import { TTSFileAction, createTTSFileSlice } from './slices/tts/action';

//  ===============  聚合 createStoreFn ============ //

export type FileStore = FilesStoreState &
  FileAction &
  FileManageAction &
  FileChunkAction &
  FileUploadAction &
  TTSFileAction;

const createStore: StateCreator<FileStore, [['zustand/devtools', never]]> = (...parameters) => ({
  ...initialState,
  ...createFileSlice(...parameters),
  ...createFileManageSlice(...parameters),
  ...createFileChunkSlice(...parameters),
  ...createFileUploadSlice(...parameters),
  ...createTTSFileSlice(...parameters),
});

//  ===============  实装 useStore ============ //
const devtools = createDevtools('file');

export const useFileStore = createWithEqualityFn<FileStore>()(devtools(createStore), shallow);

export const getFileStoreState = () => useFileStore.getState();
