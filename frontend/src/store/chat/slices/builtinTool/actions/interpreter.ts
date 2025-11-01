import {
  CodeInterpreterFileItem,
  CodeInterpreterParams,
  CodeInterpreterResponse,
} from '@/types';
import { produce } from 'immer';
import pMap from 'p-map';
import { SWRResponse } from 'swr';
import { StateCreator } from 'zustand/vanilla';

import { useClientDataSWR } from '@/libs/swr';
// import { pythonService } from '@/services/python';
import { chatSelectors } from '@/store/chat/selectors';
import { ChatStore } from '@/store/chat/store';
import { useFileStore } from '@/store/file';
import { CodeInterpreterIdentifier } from '@/tools/code-interpreter';
import { setNamespace } from '@/utils/storeDebug';

const n = setNamespace('codeInterpreter');

const SWR_FETCH_INTERPRETER_FILE_KEY = 'FetchCodeInterpreterFileItem';

export interface ChatCodeInterpreterAction {
  python: (id: string, params: CodeInterpreterParams) => Promise<boolean | undefined>;
  toggleInterpreterExecuting: (id: string, loading: boolean) => void;
  updateInterpreterFileItem: (
    id: string,
    updater: (data: CodeInterpreterResponse) => void,
  ) => Promise<void>;
  uploadInterpreterFiles: (id: string, files: CodeInterpreterFileItem[]) => Promise<void>;
  useFetchInterpreterFileItem: (id?: string) => SWRResponse;
}

export const codeInterpreterSlice: StateCreator<
  ChatStore,
  [['zustand/devtools', never]],
  [],
  ChatCodeInterpreterAction
> = (set, get) => ({
  python: async (id: string, params: CodeInterpreterParams) => {
    const {
      toggleInterpreterExecuting,
      updatePluginState,
      internal_updateMessageContent,
      uploadInterpreterFiles,
    } = get();

    // Dummy implementation for UI focus
    console.log('Running Python code:', params.code, 'with packages:', params.packages);

    toggleInterpreterExecuting(id, true);

    // Simulate async execution time
    await new Promise(resolve => setTimeout(resolve, 2000));

    try {
      // Mock Python execution result
      const mockResult: CodeInterpreterResponse = {
        output: `Mock execution result for: ${params.code.substring(0, 50)}...`,
        executionTime: 1.5,
        files: [
          {
            filename: 'result.png',
            data: new File(['mock image data'], 'result.png', { type: 'image/png' }),
            fileId: undefined,
            previewUrl: 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==',
          },
          {
            filename: 'data.csv',
            data: new File(['col1,col2\n1,2\n3,4'], 'data.csv', { type: 'text/csv' }),
            fileId: undefined,
            previewUrl: undefined,
          },
        ],
      };

      await internal_updateMessageContent(id, JSON.stringify(mockResult));
      await uploadInterpreterFiles(id, mockResult.files);

    } catch (error) {
      updatePluginState(id, { error });
      return;
    } finally {
      toggleInterpreterExecuting(id, false);
    }

    return true;
  },

  toggleInterpreterExecuting: (id: string, executing: boolean) => {
    set(
      { codeInterpreterExecuting: { ...get().codeInterpreterExecuting, [id]: executing } },
      false,
      n('toggleInterpreterExecuting'),
    );
  },

  updateInterpreterFileItem: async (
    id: string,
    updater: (data: CodeInterpreterResponse) => void,
  ) => {
    const message = chatSelectors.getMessageById(id)(get());
    if (!message) return;

    const result: CodeInterpreterResponse = JSON.parse(message.content);
    if (!result.files) return;

    const nextResult = produce(result, updater);

    await get().internal_updateMessageContent(id, JSON.stringify(nextResult));
  },

  uploadInterpreterFiles: async (id: string, files: CodeInterpreterFileItem[]) => {
    const { updateInterpreterFileItem } = get();

    // Dummy implementation for UI focus
    console.log('Uploading interpreter files:', files.map(f => f.filename));

    if (!files) return;

    // Simulate file uploads with mock IDs
    await Promise.all(files.map(async (file, index) => {
      if (!file.data) return;

      // Simulate async upload delay
      await new Promise(resolve => setTimeout(resolve, 500));

      const mockUploadResult = {
        id: `mock-interpreter-${file.filename}-${Date.now()}`,
        url: `mock://file/${file.filename}`,
      };

      await updateInterpreterFileItem(id, (draft) => {
        if (draft.files?.[index]) {
          draft.files[index].fileId = mockUploadResult.id;
          draft.files[index].previewUrl = undefined;
          draft.files[index].data = undefined;
        }
      });
    }));
  },

  useFetchInterpreterFileItem: (id) =>
    useClientDataSWR(id ? [SWR_FETCH_INTERPRETER_FILE_KEY, id] : null, async () => {
      if (!id) return null;

      // Dummy implementation for UI focus
      console.log('Fetching interpreter file item:', id);

      const mockItem = {
        id,
        name: `Mock Interpreter File ${id.slice(0, 8)}`,
        type: 'application/octet-stream',
        size: 1024000,
        url: `mock://interpreter-file/${id}`,
        createdAt: new Date(),
        updatedAt: new Date(),
      };

      set(
        produce((draft) => {
          if (!draft.codeInterpreterFileMap) {
            draft.codeInterpreterFileMap = {};
          }
          if (draft.codeInterpreterFileMap[id]) return;

          draft.codeInterpreterFileMap[id] = mockItem;
        }),
        false,
        n('useFetchInterpreterFileItem'),
      );

      return mockItem;
    }),
});
