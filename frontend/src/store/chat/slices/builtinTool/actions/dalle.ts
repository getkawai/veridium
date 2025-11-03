import { produce } from 'immer';
import { SWRResponse } from 'swr';
import { StateCreator } from 'zustand/vanilla';

import { useClientDataSWR } from '@/libs/swr';
// import { imageGenerationService } from '@/services/textToImage';
// import { uploadService } from '@/services/upload';
import { chatSelectors } from '../../message/selectors';
import { ChatStore } from '@/store/chat/store';
import { DallEImageItem } from '@/types/tool/dalle';
import { setNamespace } from '@/utils/storeDebug';

const n = setNamespace('tool');

const SWR_FETCH_KEY = 'FetchImageItem';

export interface ChatDallEAction {
  generateImageFromPrompts: (items: DallEImageItem[], id: string) => Promise<void>;
  text2image: (id: string, data: DallEImageItem[]) => Promise<void>;
  toggleDallEImageLoading: (key: string, value: boolean) => void;
  updateImageItem: (id: string, updater: (data: DallEImageItem[]) => void) => Promise<void>;
  useFetchDalleImageItem: (id: string) => SWRResponse;
}

export const dalleSlice: StateCreator<
  ChatStore,
  [['zustand/devtools', never]],
  [],
  ChatDallEAction
> = (set, get) => ({
  generateImageFromPrompts: async (items, messageId) => {
    const { toggleDallEImageLoading, updateImageItem } = get();
    // eslint-disable-next-line unicorn/consistent-function-scoping
    const getMessageById = (id: string) => chatSelectors.getMessageById(id)(get());

    const message = getMessageById(messageId);
    if (!message) return;

    const parent = getMessageById(message!.parentId!);
    const originPrompt = parent?.content;

    // Dummy implementation for UI focus
    console.log('Generating images from prompts:', items.map(item => item.prompt));

    await Promise.all(items.map(async (params, index) => {
      toggleDallEImageLoading(messageId + params.prompt, true);

      // Simulate image generation delay
      await new Promise(resolve => setTimeout(resolve, 1500 + Math.random() * 1000));

      try {
        // Mock generated image URL (base64 placeholder)
        const mockImageUrl = 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==';

        await updateImageItem(messageId, (draft) => {
          draft[index].previewUrl = mockImageUrl;
        });

        toggleDallEImageLoading(messageId + params.prompt, false);

        // Simulate file creation and upload
        await new Promise(resolve => setTimeout(resolve, 500));

        const mockImageFile = new File([await (await fetch(mockImageUrl)).blob()], `${originPrompt || params.prompt}_${index}.png`, { type: 'image/png' });

        // Mock upload result
        const mockUploadResult = {
          id: `mock-dalle-${params.prompt}-${Date.now()}-${index}`,
          url: `mock://dalle-image/${params.prompt}-${index}`,
        };

        await updateImageItem(messageId, (draft) => {
          draft[index].imageId = mockUploadResult.id;
          draft[index].previewUrl = undefined;
        });

      } catch (e) {
        toggleDallEImageLoading(messageId + params.prompt, false);
        console.error('Mock image generation failed:', e);
      }
    }));
  },
  text2image: async (id, data) => {
    // const isAutoGen = settingsSelectors.isDalleAutoGenerating(useGlobalStore.getState());
    // if (!isAutoGen) return;

    await get().generateImageFromPrompts(data, id);
  },

  toggleDallEImageLoading: (key, value) => {
    set(
      { dalleImageLoading: { ...get().dalleImageLoading, [key]: value } },
      false,
      n('toggleDallEImageLoading'),
    );
  },

  updateImageItem: async (id, updater) => {
    const message = chatSelectors.getMessageById(id)(get());
    if (!message) return;

    const data: DallEImageItem[] = JSON.parse(message.content);

    const nextContent = produce(data, updater);
    await get().internal_updateMessageContent(id, JSON.stringify(nextContent));
  },

  useFetchDalleImageItem: (id) =>
    useClientDataSWR([SWR_FETCH_KEY, id], async () => {
      // Dummy implementation for UI focus
      console.log('Fetching DALL-E image item:', id);

      const mockItem = {
        id,
        name: `Mock DALL-E Image ${id.slice(0, 8)}`,
        type: 'image/png',
        size: 2048000,
        url: `mock://dalle-image/${id}`,
        createdAt: new Date(),
        updatedAt: new Date(),
      };

      set(
        produce((draft) => {
          if (draft.dalleImageMap[id]) return;

          draft.dalleImageMap[id] = mockItem;
        }),
        false,
        n('useFetchFile'),
      );

      return mockItem;
    }),
});
